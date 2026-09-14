package product

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"everyreview/backend/internal/media"
)

var (
	ErrInvalidBarcode  = errors.New("product: barcode must not be empty")
	ErrEmptyName       = errors.New("product: name must not be empty")
	ErrNameTooLong     = errors.New("product: name must be at most 200 characters")
	ErrFieldTooLong    = errors.New("product: brand and manufacturer must be at most 120 characters")
	ErrDescTooLong     = errors.New("product: description must be at most 2000 characters")
	ErrInvalidPrice    = errors.New("product: price must be a non-negative number")
	ErrInvalidCurrency = errors.New("product: currency must be a 3-letter ISO 4217 code")
)

const (
	maxNameLength        = 200
	maxShortFieldLength  = 120
	maxDescriptionLength = 2000
	maxPrice             = 9_999_999_999.99 // numeric(12,2)
)

type Service struct {
	repo  Repository
	media media.Store
}

func NewService(repo Repository, mediaStore media.Store) *Service {
	return &Service{repo: repo, media: mediaStore}
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
	return s.withSummary(ctx, p)
}

func (s *Service) GetByID(ctx context.Context, id string) (Product, error) {
	return s.repo.GetByID(ctx, id)
}

// DetailsInput carries the raw form values; validation and normalisation happen in SubmitDetails.
type DetailsInput struct {
	Name         string
	Brand        string
	Manufacturer string
	Description  string
	Price        string
	Currency     string
}

// SubmitDetails turns a placeholder into a described product. image may be nil.
// The image is stored before the row is updated and removed again if the update
// fails, so a lost race never leaves an orphaned file behind.
func (s *Service) SubmitDetails(ctx context.Context, id string, in DetailsInput, image io.Reader) (WithSummary, error) {
	details, err := validateDetails(in)
	if err != nil {
		return WithSummary{}, err
	}

	if image != nil {
		key, err := s.media.SaveProductImage(ctx, image)
		if err != nil {
			return WithSummary{}, err
		}
		details.ImageObjectKey = &key
	}

	p, err := s.repo.UpdateDetails(ctx, id, details)
	if err != nil {
		if details.ImageObjectKey != nil {
			_ = s.media.Delete(ctx, *details.ImageObjectKey)
		}
		if errors.Is(err, ErrNotFound) || errors.Is(err, ErrAlreadyDetailed) {
			return WithSummary{}, err
		}
		return WithSummary{}, fmt.Errorf("product: submit details: %w", err)
	}
	return s.withSummary(ctx, p)
}

func (s *Service) withSummary(ctx context.Context, p Product) (WithSummary, error) {
	summary, err := s.repo.RatingSummary(ctx, p.ID)
	if err != nil {
		return WithSummary{}, fmt.Errorf("product: rating summary: %w", err)
	}
	return WithSummary{Product: p, RatingSummary: summary}, nil
}

func validateDetails(in DetailsInput) (Details, error) {
	name := strings.TrimSpace(in.Name)
	switch {
	case name == "":
		return Details{}, ErrEmptyName
	case utf8.RuneCountInString(name) > maxNameLength:
		return Details{}, ErrNameTooLong
	}
	d := Details{Name: name}

	if d.Brand = optionalText(in.Brand); d.Brand != nil && utf8.RuneCountInString(*d.Brand) > maxShortFieldLength {
		return Details{}, ErrFieldTooLong
	}
	if d.Manufacturer = optionalText(in.Manufacturer); d.Manufacturer != nil && utf8.RuneCountInString(*d.Manufacturer) > maxShortFieldLength {
		return Details{}, ErrFieldTooLong
	}
	if d.Description = optionalText(in.Description); d.Description != nil && utf8.RuneCountInString(*d.Description) > maxDescriptionLength {
		return Details{}, ErrDescTooLong
	}

	if raw := strings.TrimSpace(in.Price); raw != "" {
		// Accept "12,50" as well as "12.50": comma decimal separators are common on phone keyboards.
		price, err := strconv.ParseFloat(strings.Replace(raw, ",", ".", 1), 64)
		if err != nil || math.IsNaN(price) || math.IsInf(price, 0) || price < 0 || price > maxPrice {
			return Details{}, ErrInvalidPrice
		}
		price = math.Round(price*100) / 100
		d.Price = &price
	}
	if raw := strings.ToUpper(strings.TrimSpace(in.Currency)); raw != "" {
		if len(raw) != 3 || strings.Trim(raw, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") != "" {
			return Details{}, ErrInvalidCurrency
		}
		d.Currency = &raw
	}
	return d, nil
}

func optionalText(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
