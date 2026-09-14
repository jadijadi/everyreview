package review

import "time"

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
