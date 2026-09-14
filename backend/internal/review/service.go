package review

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyBody     = errors.New("review: body must not be empty")
	ErrBodyTooLong   = errors.New("review: body must be at most 5000 characters")
	ErrInvalidRating = errors.New("review: rating must be between 1 and 5")
)

const maxBodyLength = 5000
const defaultListLimit = 50

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListByProduct(ctx context.Context, productID string) ([]Review, error) {
	reviews, err := s.repo.ListByProduct(ctx, productID, defaultListLimit)
	if err != nil {
		return nil, fmt.Errorf("review: list by product: %w", err)
	}
	return reviews, nil
}

func (s *Service) Submit(ctx context.Context, productID, authorName, body string, rating *int) (Review, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return Review{}, ErrEmptyBody
	}
	if len(body) > maxBodyLength {
		return Review{}, ErrBodyTooLong
	}
	if rating != nil && (*rating < 1 || *rating > 5) {
		return Review{}, ErrInvalidRating
	}
	authorName = strings.TrimSpace(authorName)
	if authorName == "" {
		authorName = "Anonymous"
	}

	created, err := s.repo.Create(ctx, NewReview{
		ProductID:  productID,
		AuthorName: authorName,
		Body:       body,
		Rating:     rating,
	})
	if err != nil {
		return Review{}, fmt.Errorf("review: submit: %w", err)
	}
	return created, nil
}
