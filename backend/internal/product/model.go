package product

import "time"

type Product struct {
	ID             string
	Barcode        string
	Name           *string
	Brand          *string
	ImageObjectKey *string
	Source         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (p Product) IsPlaceholder() bool {
	return p.Source == "placeholder"
}

type RatingSummary struct {
	ReviewCount   int
	AverageRating *float64
}
