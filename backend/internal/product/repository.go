package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"everyreview/backend/internal/platform/timeseries"
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
	// GetByIDs returns the products it finds, keyed by id; unknown ids are simply absent.
	GetByIDs(ctx context.Context, ids []string) (map[string]Product, error)
	Stats(ctx context.Context, recentLimit, days int) (Stats, error)
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

func (r *postgresRepository) GetByIDs(ctx context.Context, ids []string) (map[string]Product, error) {
	result := make(map[string]Product, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	query := `SELECT ` + productColumns + ` FROM products WHERE id = ANY($1::uuid[])`
	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("product: get by ids: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("product: get by ids: scan: %w", err)
		}
		result[p.ID] = p
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("product: get by ids: rows: %w", err)
	}
	return result, nil
}

func (r *postgresRepository) Stats(ctx context.Context, recentLimit, days int) (Stats, error) {
	var s Stats
	const totals = `
		SELECT COUNT(*),
		       COUNT(*) FILTER (WHERE source = 'placeholder'),
		       COUNT(*) FILTER (WHERE image_object_key IS NOT NULL),
		       COUNT(*) FILTER (WHERE created_at >= now() - interval '1 day'),
		       COUNT(*) FILTER (WHERE created_at >= now() - interval '7 days')
		FROM products
	`
	if err := r.db.QueryRowContext(ctx, totals).Scan(&s.Total, &s.Placeholders, &s.WithPhoto, &s.Last24h, &s.Last7d); err != nil {
		return Stats{}, fmt.Errorf("product: stats totals: %w", err)
	}

	recent := `SELECT ` + productColumns + ` FROM products ORDER BY created_at DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, recent, recentLimit)
	if err != nil {
		return Stats{}, fmt.Errorf("product: stats recent: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return Stats{}, fmt.Errorf("product: stats recent: scan: %w", err)
		}
		s.Recent = append(s.Recent, p)
	}
	if err := rows.Err(); err != nil {
		return Stats{}, fmt.Errorf("product: stats recent: rows: %w", err)
	}

	const perDay = `
		SELECT date_trunc('day', created_at)::date, COUNT(*)
		FROM products
		WHERE created_at >= date_trunc('day', now()) - make_interval(days => $1 - 1)
		GROUP BY 1
	`
	dayRows, err := r.db.QueryContext(ctx, perDay, days)
	if err != nil {
		return Stats{}, fmt.Errorf("product: stats per day: %w", err)
	}
	defer dayRows.Close()
	if s.PerDay, err = timeseries.ScanPerDay(dayRows, days); err != nil {
		return Stats{}, fmt.Errorf("product: stats per day: %w", err)
	}
	return s, nil
}
