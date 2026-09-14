package product

import "time"

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
