package product

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/products/{barcode}", h.GetByBarcode)
}

type ratingSummaryResponse struct {
	ReviewCount   int      `json:"review_count"`
	AverageRating *float64 `json:"average_rating"`
}

type productResponse struct {
	ID            string                `json:"id"`
	Barcode       string                `json:"barcode"`
	Name          *string               `json:"name"`
	Brand         *string               `json:"brand"`
	ImageURL      *string               `json:"image_url"`
	IsPlaceholder bool                  `json:"is_placeholder"`
	RatingSummary ratingSummaryResponse `json:"rating_summary"`
}

func toProductResponse(p WithSummary) productResponse {
	return productResponse{
		ID:            p.ID,
		Barcode:       p.Barcode,
		Name:          p.Name,
		Brand:         p.Brand,
		ImageURL:      nil, // no CDN/media handling in this MVP slice
		IsPlaceholder: p.IsPlaceholder(),
		RatingSummary: ratingSummaryResponse{
			ReviewCount:   p.RatingSummary.ReviewCount,
			AverageRating: p.RatingSummary.AverageRating,
		},
	}
}

func (h *Handler) GetByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := chi.URLParam(r, "barcode")
	result, err := h.service.LookupByBarcode(r.Context(), barcode)
	if err != nil {
		if errors.Is(err, ErrInvalidBarcode) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(result))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
