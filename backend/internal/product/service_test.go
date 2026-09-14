package product

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	products  map[string]Product
	summary   RatingSummary
	updateErr error
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

func (f *fakeRepository) UpdateDetails(ctx context.Context, id string, d Details) (Product, error) {
	if f.updateErr != nil {
		return Product{}, f.updateErr
	}
	p, err := f.GetByID(ctx, id)
	if err != nil {
		return Product{}, err
	}
	if !p.IsPlaceholder() {
		return Product{}, ErrAlreadyDetailed
	}
	p.Name = &d.Name
	p.Brand, p.Manufacturer, p.Description = d.Brand, d.Manufacturer, d.Description
	p.Price, p.Currency = d.Price, d.Currency
	if d.ImageObjectKey != nil {
		p.ImageObjectKey = d.ImageObjectKey
	}
	p.Source = "user"
	f.products[p.Barcode] = p
	return p, nil
}

func (f *fakeRepository) RatingSummary(ctx context.Context, productID string) (RatingSummary, error) {
	return f.summary, nil
}

func (f *fakeRepository) GetByIDs(ctx context.Context, ids []string) (map[string]Product, error) {
	result := map[string]Product{}
	for _, id := range ids {
		if p, err := f.GetByID(ctx, id); err == nil {
			result[id] = p
		}
	}
	return result, nil
}

func (f *fakeRepository) Stats(ctx context.Context, recentLimit, days int) (Stats, error) {
	return Stats{Total: len(f.products)}, nil
}

type fakeMediaStore struct {
	saved   []string
	deleted []string
	saveErr error
}

func (f *fakeMediaStore) SaveProductImage(ctx context.Context, r io.Reader) (string, error) {
	if f.saveErr != nil {
		return "", f.saveErr
	}
	key := "products/" + strings.Repeat("a", 32) + ".jpg"
	f.saved = append(f.saved, key)
	return key, nil
}

func (f *fakeMediaStore) Delete(ctx context.Context, key string) error {
	f.deleted = append(f.deleted, key)
	return nil
}

func newService(repo *fakeRepository) (*Service, *fakeMediaStore) {
	store := &fakeMediaStore{}
	return NewService(repo, store), store
}

func placeholder(t *testing.T, repo *fakeRepository, barcode string) Product {
	t.Helper()
	p, err := repo.GetOrCreateByBarcode(context.Background(), barcode)
	require.NoError(t, err)
	return p
}

func TestLookupByBarcode_RejectsEmptyBarcode(t *testing.T) {
	svc, _ := newService(newFakeRepository())

	_, err := svc.LookupByBarcode(context.Background(), "   ")

	require.ErrorIs(t, err, ErrInvalidBarcode)
}

func TestLookupByBarcode_CreatesPlaceholderForUnknownBarcode(t *testing.T) {
	repo := newFakeRepository()
	avg := 4.5
	repo.summary = RatingSummary{ReviewCount: 2, AverageRating: &avg}
	svc, _ := newService(repo)

	result, err := svc.LookupByBarcode(context.Background(), "0000000000000")

	require.NoError(t, err)
	require.Equal(t, "0000000000000", result.Barcode)
	require.True(t, result.IsPlaceholder())
	require.Equal(t, 2, result.RatingSummary.ReviewCount)
	require.NotNil(t, result.RatingSummary.AverageRating)
	require.Equal(t, 4.5, *result.RatingSummary.AverageRating)
}

func TestSubmitDetails_FillsPlaceholder(t *testing.T) {
	repo := newFakeRepository()
	p := placeholder(t, repo, "111")
	svc, store := newService(repo)

	result, err := svc.SubmitDetails(context.Background(), p.ID, DetailsInput{
		Name:         "  Choco Bar ",
		Brand:        "Choco",
		Manufacturer: "",
		Description:  "Tasty",
		Price:        "12,505",
		Currency:     "usd",
	}, strings.NewReader("fake-image-bytes"))

	require.NoError(t, err)
	require.False(t, result.IsPlaceholder())
	require.Equal(t, "Choco Bar", *result.Name)
	require.Equal(t, "Choco", *result.Brand)
	require.Nil(t, result.Manufacturer)
	require.Equal(t, "Tasty", *result.Description)
	require.Equal(t, 12.51, *result.Price)
	require.Equal(t, "USD", *result.Currency)
	require.Len(t, store.saved, 1)
	require.Equal(t, store.saved[0], *result.ImageObjectKey)
	require.Empty(t, store.deleted)
}

func TestSubmitDetails_WithoutImageOrOptionalFields(t *testing.T) {
	repo := newFakeRepository()
	p := placeholder(t, repo, "222")
	svc, store := newService(repo)

	result, err := svc.SubmitDetails(context.Background(), p.ID, DetailsInput{Name: "Just a name"}, nil)

	require.NoError(t, err)
	require.Equal(t, "Just a name", *result.Name)
	require.Nil(t, result.Price)
	require.Nil(t, result.Currency)
	require.Nil(t, result.ImageObjectKey)
	require.Empty(t, store.saved)
}

func TestSubmitDetails_Validation(t *testing.T) {
	cases := map[string]struct {
		in   DetailsInput
		want error
	}{
		"empty name":       {DetailsInput{Name: "   "}, ErrEmptyName},
		"name too long":    {DetailsInput{Name: strings.Repeat("n", 201)}, ErrNameTooLong},
		"brand too long":   {DetailsInput{Name: "x", Brand: strings.Repeat("b", 121)}, ErrFieldTooLong},
		"desc too long":    {DetailsInput{Name: "x", Description: strings.Repeat("d", 2001)}, ErrDescTooLong},
		"negative price":   {DetailsInput{Name: "x", Price: "-1"}, ErrInvalidPrice},
		"garbage price":    {DetailsInput{Name: "x", Price: "ten"}, ErrInvalidPrice},
		"bad currency":     {DetailsInput{Name: "x", Currency: "dollars"}, ErrInvalidCurrency},
		"numeric currency": {DetailsInput{Name: "x", Currency: "123"}, ErrInvalidCurrency},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			repo := newFakeRepository()
			p := placeholder(t, repo, "333")
			svc, store := newService(repo)

			_, err := svc.SubmitDetails(context.Background(), p.ID, tc.in, strings.NewReader("img"))

			require.ErrorIs(t, err, tc.want)
			require.Empty(t, store.saved, "invalid input must be rejected before the image is stored")
		})
	}
}

func TestSubmitDetails_SecondSubmissionIsRejectedAndImageRemoved(t *testing.T) {
	repo := newFakeRepository()
	p := placeholder(t, repo, "444")
	svc, store := newService(repo)
	_, err := svc.SubmitDetails(context.Background(), p.ID, DetailsInput{Name: "First"}, nil)
	require.NoError(t, err)

	_, err = svc.SubmitDetails(context.Background(), p.ID, DetailsInput{Name: "Second"}, strings.NewReader("img"))

	require.ErrorIs(t, err, ErrAlreadyDetailed)
	require.Equal(t, store.saved, store.deleted)
	got, _ := repo.GetByID(context.Background(), p.ID)
	require.Equal(t, "First", *got.Name)
}

func TestSubmitDetails_UnknownProduct(t *testing.T) {
	svc, _ := newService(newFakeRepository())

	_, err := svc.SubmitDetails(context.Background(), "nope", DetailsInput{Name: "x"}, nil)

	require.ErrorIs(t, err, ErrNotFound)
}

func TestSubmitDetails_ImageStoreFailureIsSurfaced(t *testing.T) {
	repo := newFakeRepository()
	p := placeholder(t, repo, "555")
	svc, store := newService(repo)
	store.saveErr = errors.New("disk full")

	_, err := svc.SubmitDetails(context.Background(), p.ID, DetailsInput{Name: "x"}, strings.NewReader("img"))

	require.ErrorContains(t, err, "disk full")
	got, _ := repo.GetByID(context.Background(), p.ID)
	require.True(t, got.IsPlaceholder(), "row must not change when the image could not be stored")
}
