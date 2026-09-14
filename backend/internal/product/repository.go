package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("product: not found")

type Repository interface {
	GetOrCreateByBarcode(ctx context.Context, barcode string) (Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	RatingSummary(ctx context.Context, productID string) (RatingSummary, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetOrCreateByBarcode(ctx context.Context, barcode string) (Product, error) {
	const query = `
		INSERT INTO products (barcode, source)
		VALUES ($1, 'placeholder')
		ON CONFLICT (barcode) DO UPDATE SET updated_at = products.updated_at
		RETURNING id, barcode, name, brand, image_object_key, source, created_at, updated_at
	`
	var p Product
	err := r.db.QueryRowContext(ctx, query, barcode).Scan(
		&p.ID, &p.Barcode, &p.Name, &p.Brand, &p.ImageObjectKey, &p.Source, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return Product{}, fmt.Errorf("product: get or create by barcode: %w", err)
	}
	return p, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (Product, error) {
	const query = `
		SELECT id, barcode, name, brand, image_object_key, source, created_at, updated_at
		FROM products WHERE id = $1
	`
	var p Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Barcode, &p.Name, &p.Brand, &p.ImageObjectKey, &p.Source, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, fmt.Errorf("product: get by id: %w", err)
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
