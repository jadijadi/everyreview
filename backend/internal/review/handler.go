package review

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
	r.Get("/products/{productId}/reviews", h.List)
	r.Post("/products/{productId}/reviews", h.Create)
}

type reviewResponse struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Author    string `json:"author_name"`
	Body      string `json:"body"`
	Rating    *int   `json:"rating"`
	CreatedAt string `json:"created_at"`
}

func toReviewResponse(r Review) reviewResponse {
	return reviewResponse{
		ID:        r.ID,
		ProductID: r.ProductID,
		Author:    r.AuthorName,
		Body:      r.Body,
		Rating:    r.Rating,
		CreatedAt: r.CreatedAt.Format(timeFormat),
	}
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

type reviewPageResponse struct {
	Items []reviewResponse `json:"items"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")
	reviews, err := h.service.ListByProduct(r.Context(), productID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	items := make([]reviewResponse, 0, len(reviews))
	for _, rev := range reviews {
		items = append(items, toReviewResponse(rev))
	}
	writeJSON(w, http.StatusOK, reviewPageResponse{Items: items})
}

type createReviewRequest struct {
	AuthorName string `json:"author_name"`
	Body       string `json:"body"`
	Rating     *int   `json:"rating"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")

	var req createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.service.Submit(r.Context(), productID, req.AuthorName, req.Body, req.Rating)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyBody), errors.Is(err, ErrBodyTooLong), errors.Is(err, ErrInvalidRating):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	writeJSON(w, http.StatusCreated, toReviewResponse(created))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
