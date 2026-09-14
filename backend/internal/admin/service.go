// Package admin serves the password-protected operator dashboard: counts and
// recent activity across products and reviews. It reads through the other
// modules' repositories rather than issuing SQL of its own (docs/coding-conventions.md).
package admin

import (
	"context"
	"fmt"
	"time"

	"everyreview/backend/internal/platform/timeseries"
	"everyreview/backend/internal/product"
	"everyreview/backend/internal/review"
)

const (
	recentLimit = 10
	topLimit    = 5
	seriesDays  = 14
)

type ProductSource interface {
	Stats(ctx context.Context, recentLimit, days int) (product.Stats, error)
	GetByIDs(ctx context.Context, ids []string) (map[string]product.Product, error)
}

type ReviewSource interface {
	Stats(ctx context.Context, recentLimit, topLimit, days int) (review.Stats, error)
}

type Service struct {
	products ProductSource
	reviews  ReviewSource
	now      func() time.Time
}

func NewService(products ProductSource, reviews ReviewSource) *Service {
	return &Service{products: products, reviews: reviews, now: time.Now}
}

// Dashboard is the fully resolved view model: review rows already carry their
// product's name so the template stays free of lookups.
type Dashboard struct {
	GeneratedAt time.Time
	Products    product.Stats
	Reviews     review.Stats
	Recent      []RecentReview
	TopProducts []TopProduct
	Series      []DaySeries
	SeriesMax   int
}

type RecentReview struct {
	review.Review
	ProductName string
	Barcode     string
}

type TopProduct struct {
	product.Product
	ReviewCount int
}

type DaySeries struct {
	Day      time.Time
	Products int
	Reviews  int
}

func (s *Service) Dashboard(ctx context.Context) (Dashboard, error) {
	ps, err := s.products.Stats(ctx, recentLimit, seriesDays)
	if err != nil {
		return Dashboard{}, fmt.Errorf("admin: product stats: %w", err)
	}
	rs, err := s.reviews.Stats(ctx, recentLimit, topLimit, seriesDays)
	if err != nil {
		return Dashboard{}, fmt.Errorf("admin: review stats: %w", err)
	}

	ids := make([]string, 0, len(rs.Recent)+len(rs.TopProducts))
	for _, r := range rs.Recent {
		ids = append(ids, r.ProductID)
	}
	for _, t := range rs.TopProducts {
		ids = append(ids, t.ProductID)
	}
	named, err := s.products.GetByIDs(ctx, ids)
	if err != nil {
		return Dashboard{}, fmt.Errorf("admin: resolve product names: %w", err)
	}

	d := Dashboard{GeneratedAt: s.now().UTC(), Products: ps, Reviews: rs}
	for _, r := range rs.Recent {
		p := named[r.ProductID]
		d.Recent = append(d.Recent, RecentReview{Review: r, ProductName: displayName(p), Barcode: p.Barcode})
	}
	for _, t := range rs.TopProducts {
		d.TopProducts = append(d.TopProducts, TopProduct{Product: named[t.ProductID], ReviewCount: t.Count})
	}
	d.Series, d.SeriesMax = mergeSeries(ps.PerDay, rs.PerDay)
	return d, nil
}

func (t TopProduct) DisplayName() string { return displayName(t.Product) }

func displayName(p product.Product) string {
	if p.Name != nil && *p.Name != "" {
		return *p.Name
	}
	if p.Barcode == "" {
		return "(deleted product)"
	}
	return "Unnamed · " + p.Barcode
}

// mergeSeries zips the two per-day series (same length and days, both produced by
// timeseries.FillDays) and reports the tallest bar so the template can scale.
func mergeSeries(products, reviews []timeseries.DayCount) ([]DaySeries, int) {
	n := len(products)
	if len(reviews) < n {
		n = len(reviews)
	}
	series := make([]DaySeries, 0, n)
	maxCount := 1
	for i := 0; i < n; i++ {
		day := DaySeries{Day: products[i].Day, Products: products[i].Count, Reviews: reviews[i].Count}
		series = append(series, day)
		if day.Products > maxCount {
			maxCount = day.Products
		}
		if day.Reviews > maxCount {
			maxCount = day.Reviews
		}
	}
	return series, maxCount
}
