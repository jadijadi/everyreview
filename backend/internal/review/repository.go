package review

import (
	"context"
	"database/sql"
	"fmt"

	"everyreview/backend/internal/platform/timeseries"
)

type Repository interface {
	ListByProduct(ctx context.Context, productID string, limit int) ([]Review, error)
	Create(ctx context.Context, r NewReview) (Review, error)
	Stats(ctx context.Context, recentLimit, topLimit, days int) (Stats, error)
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

const publishedFilter = `status = 'published' AND deleted_at IS NULL`

func (repo *postgresRepository) Stats(ctx context.Context, recentLimit, topLimit, days int) (Stats, error) {
	var s Stats
	totals := `
		SELECT COUNT(*), AVG(rating),
		       COUNT(*) FILTER (WHERE rating = 1), COUNT(*) FILTER (WHERE rating = 2),
		       COUNT(*) FILTER (WHERE rating = 3), COUNT(*) FILTER (WHERE rating = 4),
		       COUNT(*) FILTER (WHERE rating = 5),
		       COUNT(*) FILTER (WHERE created_at >= now() - interval '1 day'),
		       COUNT(*) FILTER (WHERE created_at >= now() - interval '7 days'),
		       COUNT(DISTINCT author_name)
		FROM reviews WHERE ` + publishedFilter
	var avg sql.NullFloat64
	err := repo.db.QueryRowContext(ctx, totals).Scan(
		&s.Total, &avg,
		&s.RatingCounts[0], &s.RatingCounts[1], &s.RatingCounts[2], &s.RatingCounts[3], &s.RatingCounts[4],
		&s.Last24h, &s.Last7d, &s.DistinctAuthor,
	)
	if err != nil {
		return Stats{}, fmt.Errorf("review: stats totals: %w", err)
	}
	if avg.Valid {
		s.AverageRating = &avg.Float64
	}

	recent := `
		SELECT id, product_id, author_name, body, rating, status, created_at, updated_at
		FROM reviews WHERE ` + publishedFilter + `
		ORDER BY created_at DESC LIMIT $1`
	rows, err := repo.db.QueryContext(ctx, recent, recentLimit)
	if err != nil {
		return Stats{}, fmt.Errorf("review: stats recent: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var r Review
		if err := rows.Scan(&r.ID, &r.ProductID, &r.AuthorName, &r.Body, &r.Rating, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return Stats{}, fmt.Errorf("review: stats recent: scan: %w", err)
		}
		s.Recent = append(s.Recent, r)
	}
	if err := rows.Err(); err != nil {
		return Stats{}, fmt.Errorf("review: stats recent: rows: %w", err)
	}

	top := `
		SELECT product_id, COUNT(*) FROM reviews WHERE ` + publishedFilter + `
		GROUP BY product_id ORDER BY 2 DESC, MAX(created_at) DESC LIMIT $1`
	topRows, err := repo.db.QueryContext(ctx, top, topLimit)
	if err != nil {
		return Stats{}, fmt.Errorf("review: stats top: %w", err)
	}
	defer topRows.Close()
	for topRows.Next() {
		var pc ProductCount
		if err := topRows.Scan(&pc.ProductID, &pc.Count); err != nil {
			return Stats{}, fmt.Errorf("review: stats top: scan: %w", err)
		}
		s.TopProducts = append(s.TopProducts, pc)
	}
	if err := topRows.Err(); err != nil {
		return Stats{}, fmt.Errorf("review: stats top: rows: %w", err)
	}

	perDay := `
		SELECT date_trunc('day', created_at)::date, COUNT(*)
		FROM reviews
		WHERE ` + publishedFilter + ` AND created_at >= date_trunc('day', now()) - make_interval(days => $1 - 1)
		GROUP BY 1`
	dayRows, err := repo.db.QueryContext(ctx, perDay, days)
	if err != nil {
		return Stats{}, fmt.Errorf("review: stats per day: %w", err)
	}
	defer dayRows.Close()
	if s.PerDay, err = timeseries.ScanPerDay(dayRows, days); err != nil {
		return Stats{}, fmt.Errorf("review: stats per day: %w", err)
	}
	return s, nil
}
