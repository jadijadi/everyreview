package product

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	products map[string]Product
	summary  RatingSummary
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{products: map[string]Product{}}
}

func (f *fakeRepository) GetOrCreateByBarcode(ctx context.Context, barcode string) (Product, error) {
	if p, ok := f.products[barcode]; ok {
		return p, nil
	}
	p := Product{ID: "id-" + barcode, Barcode: barcode, Source: "placeholder"}
	f.products[barcode] = p
	return p, nil
}

func (f *fakeRepository) GetByID(ctx context.Context, id string) (Product, error) {
	for _, p := range f.products {
		if p.ID == id {
			return p, nil
		}
	}
	return Product{}, ErrNotFound
}

func (f *fakeRepository) RatingSummary(ctx context.Context, productID string) (RatingSummary, error) {
	return f.summary, nil
}

func TestLookupByBarcode_RejectsEmptyBarcode(t *testing.T) {
	svc := NewService(newFakeRepository())

	_, err := svc.LookupByBarcode(context.Background(), "   ")

	require.ErrorIs(t, err, ErrInvalidBarcode)
}

func TestLookupByBarcode_CreatesPlaceholderForUnknownBarcode(t *testing.T) {
	repo := newFakeRepository()
	avg := 4.5
	repo.summary = RatingSummary{ReviewCount: 2, AverageRating: &avg}
	svc := NewService(repo)

	result, err := svc.LookupByBarcode(context.Background(), "0000000000000")

	require.NoError(t, err)
	require.Equal(t, "0000000000000", result.Barcode)
	require.True(t, result.IsPlaceholder())
	require.Equal(t, 2, result.RatingSummary.ReviewCount)
	require.NotNil(t, result.RatingSummary.AverageRating)
	require.Equal(t, 4.5, *result.RatingSummary.AverageRating)
}
