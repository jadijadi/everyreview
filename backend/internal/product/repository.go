package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound        = errors.New("product: not found")
	ErrAlreadyDetailed = errors.New("product: details were already submitted")
)

type Repository interface {
	GetOrCreateByBarcode(ctx context.Context, barcode string) (Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	// UpdateDetails fills in a placeholder product. Returns ErrAlreadyDetailed if
	// someone else got there first, ErrNotFound if the id is unknown.
	UpdateDetails(ctx context.Context, id string, d Details) (Product, error)
	RatingSummary(ctx context.Context, productID string) (RatingSummary, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

const productColumns = `id, barcode, name, brand, manufacturer, description, price, currency, image_object_key, source, created_at, updated_at`

func scanProduct(row interface{ Scan(dest ...any) error }) (Product, error) {
	var p Product
	var price sql.NullFloat64
	err := row.Scan(
		&p.ID, &p.Barcode, &p.Name, &p.Brand, &p.Manufacturer, &p.Description, &price, &p.Currency,
		&p.ImageObjectKey, &p.Source, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return Product{}, err
	}
	if price.Valid {
		p.Price = &price.Float64
	}
	return p, nil
}

func (r *postgresRepository) GetOrCreateByBarcode(ctx context.Context, barcode string) (Product, error) {
	query := `
		INSERT INTO products (barcode, source)
		VALUES ($1, 'placeholder')
		ON CONFLICT (barcode) DO UPDATE SET updated_at = products.updated_at
		RETURNING ` + productColumns
	p, err := scanProduct(r.db.QueryRowContext(ctx, query, barcode))
	if err != nil {
		return Product{}, fmt.Errorf("product: get or create by barcode: %w", err)
	}
	return p, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (Product, error) {
	query := `SELECT ` + productColumns + ` FROM products WHERE id = $1`
	p, err := scanProduct(r.db.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, fmt.Errorf("product: get by id: %w", err)
	}
	return p, nil
}

func (r *postgresRepository) UpdateDetails(ctx context.Context, id string, d Details) (Product, error) {
	// The source guard makes this first-writer-wins without a separate lock;
	// a lost race surfaces as zero rows, which we disambiguate below.
	query := `
		UPDATE products
		SET name = $2, brand = $3, manufacturer = $4, description = $5, price = $6, currency = $7,
		    image_object_key = COALESCE($8, image_object_key), source = 'user', updated_at = now()
		WHERE id = $1 AND source = 'placeholder'
		RETURNING ` + productColumns
	p, err := scanProduct(r.db.QueryRowContext(ctx, query,
		id, d.Name, d.Brand, d.Manufacturer, d.Description, d.Price, d.Currency, d.ImageObjectKey,
	))
	if errors.Is(err, sql.ErrNoRows) {
		if _, getErr := r.GetByID(ctx, id); getErr != nil {
			return Product{}, getErr
		}
		return Product{}, ErrAlreadyDetailed
	}
	if err != nil {
		return Product{}, fmt.Errorf("product: update details: %w", err)
	}
	return p, nil
}

func (r *postgresRepository) RatingSummary(ctx context.Context, productID string) (RatingSummary, error) {
	const query = `
		SELECT COUNT(*), AVG(rating)
		FROM reviews
		WHERE product_id = $1 AND status = 'published' AND deleted_at IS NULL
	`
	var summary RatingSummary
	var avg sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, query, productID).Scan(&summary.ReviewCount, &avg); err != nil {
		return RatingSummary{}, fmt.Errorf("product: rating summary: %w", err)
	}
	if avg.Valid {
		summary.AverageRating = &avg.Float64
	}
	return summary, nil
}
