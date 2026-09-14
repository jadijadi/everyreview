package product

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidBarcode = errors.New("product: barcode must not be empty")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type WithSummary struct {
	Product
	RatingSummary RatingSummary
}

func (s *Service) LookupByBarcode(ctx context.Context, barcode string) (WithSummary, error) {
	barcode = strings.TrimSpace(barcode)
	if barcode == "" {
		return WithSummary{}, ErrInvalidBarcode
	}
	p, err := s.repo.GetOrCreateByBarcode(ctx, barcode)
	if err != nil {
		return WithSummary{}, fmt.Errorf("product: lookup by barcode: %w", err)
	}
	summary, err := s.repo.RatingSummary(ctx, p.ID)
	if err != nil {
		return WithSummary{}, fmt.Errorf("product: lookup by barcode: %w", err)
	}
	return WithSummary{Product: p, RatingSummary: summary}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (Product, error) {
	return s.repo.GetByID(ctx, id)
}
