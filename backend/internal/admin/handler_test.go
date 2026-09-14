package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"everyreview/backend/internal/platform/timeseries"
	"everyreview/backend/internal/product"
	"everyreview/backend/internal/review"
)

type fakeProducts struct{ byID map[string]product.Product }

func (f fakeProducts) Stats(ctx context.Context, recentLimit, days int) (product.Stats, error) {
	recent := make([]product.Product, 0, len(f.byID))
	for _, p := range f.byID {
		recent = append(recent, p)
	}
	return product.Stats{
		Total: len(f.byID), Placeholders: 1, WithPhoto: 1, Last24h: 2, Last7d: 3,
		Recent: recent,
		PerDay: timeseries.FillDays(map[string]int{time.Now().UTC().Format(time.DateOnly): 2}, time.Now().UTC(), days),
	}, nil
}

func (f fakeProducts) GetByIDs(ctx context.Context, ids []string) (map[string]product.Product, error) {
	out := map[string]product.Product{}
	for _, id := range ids {
		if p, ok := f.byID[id]; ok {
			out[id] = p
		}
	}
	return out, nil
}

type fakeReviews struct{ stats review.Stats }

func (f fakeReviews) Stats(ctx context.Context, recentLimit, topLimit, days int) (review.Stats, error) {
	s := f.stats
	s.PerDay = timeseries.FillDays(map[string]int{time.Now().UTC().Format(time.DateOnly): 5}, time.Now().UTC(), days)
	return s, nil
}

func fixture() (*Service, string) {
	name := "Choco Bar"
	five := 5
	products := fakeProducts{byID: map[string]product.Product{
		"p1": {ID: "p1", Barcode: "111", Name: &name, Source: "user", CreatedAt: time.Now()},
		"p2": {ID: "p2", Barcode: "222", Source: "placeholder", CreatedAt: time.Now()},
	}}
	avg := 4.3333
	reviews := fakeReviews{stats: review.Stats{
		Total: 3, AverageRating: &avg, RatingCounts: [5]int{0, 0, 0, 1, 2}, DistinctAuthor: 2,
		Recent: []review.Review{
			{ID: "r1", ProductID: "p1", AuthorName: "Alice", Body: "Lovely <b>stuff</b>", Rating: &five, CreatedAt: time.Now()},
			{ID: "r2", ProductID: "p2", AuthorName: "Bob", Body: "meh", CreatedAt: time.Now().Add(-3 * time.Hour)},
		},
		TopProducts: []review.ProductCount{{ProductID: "p1", Count: 2}, {ProductID: "p2", Count: 1}},
	}}
	return NewService(products, reviews), "s3cret"
}

func router(password string) http.Handler {
	svc, _ := fixture()
	r := chi.NewRouter()
	NewHandler(svc, password).Routes(r)
	return r
}

func get(t *testing.T, h http.Handler, path, user, password string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if password != "" {
		req.SetBasicAuth(user, password)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestDashboard_DisabledWithoutPassword(t *testing.T) {
	rec := get(t, router(""), "/admin", "admin", "anything")

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDashboard_RequiresPassword(t *testing.T) {
	h := router("s3cret")

	rec := get(t, h, "/admin", "", "")
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Header().Get("WWW-Authenticate"), "Basic")

	rec = get(t, h, "/admin", "admin", "wrong")
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDashboard_RendersStats(t *testing.T) {
	rec := get(t, router("s3cret"), "/admin", "whoever", "s3cret")

	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	require.Contains(t, body, "Choco Bar")
	require.Contains(t, body, "Unnamed · 222", "reviews on placeholders fall back to the barcode")
	require.Contains(t, body, "Lovely &lt;b&gt;stuff&lt;/b&gt;", "review text must be escaped")
	require.Contains(t, body, "★★★★★")
	require.Contains(t, body, ">4.33<", "average rating is formatted, not printed as a pointer")
	require.Contains(t, body, "50% of products are still placeholders")
}

func TestStatsJSON(t *testing.T) {
	rec := get(t, router("s3cret"), "/admin/stats.json", "x", "s3cret")

	require.Equal(t, http.StatusOK, rec.Code)
	var d Dashboard
	require.NoError(t, json.NewDecoder(strings.NewReader(rec.Body.String())).Decode(&d))
	require.Equal(t, 2, d.Products.Total)
	require.Equal(t, 3, d.Reviews.Total)
	require.Len(t, d.Series, seriesDays)
	require.Equal(t, 5, d.SeriesMax, "tallest bar sets the chart scale")
	require.Equal(t, "Choco Bar", d.TopProducts[0].DisplayName())
}
