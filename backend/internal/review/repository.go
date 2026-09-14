package review

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	ListByProduct(ctx context.Context, productID string, limit int) ([]Review, error)
	Create(ctx context.Context, r NewReview) (Review, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (repo *postgresRepository) ListByProduct(ctx context.Context, productID string, limit int) ([]Review, error) {
	const query = `
		SELECT id, product_id, author_name, body, rating, status, created_at, updated_at
		FROM reviews
		WHERE product_id = $1 AND status = 'published' AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := repo.db.QueryContext(ctx, query, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("review: list by product: %w", err)
	}
	defer rows.Close()

	reviews := []Review{}
	for rows.Next() {
		var r Review
		if err := rows.Scan(&r.ID, &r.ProductID, &r.AuthorName, &r.Body, &r.Rating, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("review: list by product: scan: %w", err)
		}
		reviews = append(reviews, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("review: list by product: rows: %w", err)
	}
	return reviews, nil
}

func (repo *postgresRepository) Create(ctx context.Context, newReview NewReview) (Review, error) {
	const query = `
		INSERT INTO reviews (product_id, author_name, body, rating)
		VALUES ($1, $2, $3, $4)
		RETURNING id, product_id, author_name, body, rating, status, created_at, updated_at
	`
	var r Review
	err := repo.db.QueryRowContext(ctx, query, newReview.ProductID, newReview.AuthorName, newReview.Body, newReview.Rating).Scan(
		&r.ID, &r.ProductID, &r.AuthorName, &r.Body, &r.Rating, &r.Status, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return Review{}, fmt.Errorf("review: create: %w", err)
	}
	return r, nil
}
