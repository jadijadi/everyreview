package review

import (
	"time"

	"everyreview/backend/internal/platform/timeseries"
)

type Review struct {
	ID         string
	ProductID  string
	AuthorName string
	Body       string
	Rating     *int
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NewReview struct {
	ProductID  string
	AuthorName string
	Body       string
	Rating     *int
}

// Stats feeds the admin dashboard; only published, non-deleted reviews count.
type Stats struct {
	Total          int
	AverageRating  *float64
	RatingCounts   [5]int // index 0 = one star
	Last24h        int
	Last7d         int
	DistinctAuthor int
	Recent         []Review              // newest first
	TopProducts    []ProductCount        // most reviewed first
	PerDay         []timeseries.DayCount // oldest first, gaps zero-filled
}

type ProductCount struct {
	ProductID string
	Count     int
}
