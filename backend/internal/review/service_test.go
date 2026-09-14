package review

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	created []NewReview
	toList  []Review
}

func (f *fakeRepository) ListByProduct(ctx context.Context, productID string, limit int) ([]Review, error) {
	return f.toList, nil
}

func (f *fakeRepository) Create(ctx context.Context, r NewReview) (Review, error) {
	f.created = append(f.created, r)
	return Review{
		ID:         "review-1",
		ProductID:  r.ProductID,
		AuthorName: r.AuthorName,
		Body:       r.Body,
		Rating:     r.Rating,
		Status:     "published",
	}, nil
}

func TestSubmit_RejectsEmptyBody(t *testing.T) {
	svc := NewService(&fakeRepository{})

	_, err := svc.Submit(context.Background(), "product-1", "Alice", "   ", nil)

	require.ErrorIs(t, err, ErrEmptyBody)
}

func TestSubmit_RejectsOutOfRangeRating(t *testing.T) {
	svc := NewService(&fakeRepository{})
	rating := 6

	_, err := svc.Submit(context.Background(), "product-1", "Alice", "Great product", &rating)

	require.ErrorIs(t, err, ErrInvalidRating)
}

func TestSubmit_DefaultsAuthorNameWhenBlank(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewService(repo)

	created, err := svc.Submit(context.Background(), "product-1", "  ", "Great product", nil)

	require.NoError(t, err)
	require.Equal(t, "Anonymous", created.AuthorName)
	require.Len(t, repo.created, 1)
	require.Equal(t, "Anonymous", repo.created[0].AuthorName)
}

func TestListByProduct_ReturnsRepositoryResults(t *testing.T) {
	repo := &fakeRepository{toList: []Review{{ID: "r1"}, {ID: "r2"}}}
	svc := NewService(repo)

	reviews, err := svc.ListByProduct(context.Background(), "product-1")

	require.NoError(t, err)
	require.Len(t, reviews, 2)
}
