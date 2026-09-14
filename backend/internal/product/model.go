package product

import (
	"time"

	"everyreview/backend/internal/platform/timeseries"
)

type Product struct {
	ID             string
	Barcode        string
	Name           *string
	Brand          *string
	Manufacturer   *string
	Description    *string
	Price          *float64
	Currency       *string
	ImageObjectKey *string
	Source         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (p Product) IsPlaceholder() bool {
	return p.Source == "placeholder"
}

// Details is what a user fills in for a product nobody has described yet.
// Only Name is mandatory; ImageObjectKey is set by the service after the upload is stored.
type Details struct {
	Name           string
	Brand          *string
	Manufacturer   *string
	Description    *string
	Price          *float64
	Currency       *string
	ImageObjectKey *string
}

type RatingSummary struct {
	ReviewCount   int
	AverageRating *float64
}

// Stats feeds the admin dashboard; counts are over all products regardless of source.
type Stats struct {
	Total        int
	Placeholders int // scanned but nobody has described them yet
	WithPhoto    int
	Last24h      int
	Last7d       int
	Recent       []Product             // newest first
	PerDay       []timeseries.DayCount // oldest first, one entry per day, gaps filled with zero
}
