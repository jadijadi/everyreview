package product

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"everyreview/backend/internal/media"
)

// Multipart form size cap: the image limit plus room for the text fields.
const maxDetailsRequestBytes = media.MaxImageBytes + 64<<10

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/products/{barcode}", h.GetByBarcode)
	r.Post("/products/{productId}/details", h.SubmitDetails)
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
	Manufacturer  *string               `json:"manufacturer"`
	Description   *string               `json:"description"`
	Price         *float64              `json:"price"`
	Currency      *string               `json:"currency"`
	ImageURL      *string               `json:"image_url"`
	IsPlaceholder bool                  `json:"is_placeholder"`
	RatingSummary ratingSummaryResponse `json:"rating_summary"`
}

func toProductResponse(r *http.Request, p WithSummary) productResponse {
	var imageURL *string
	if p.ImageObjectKey != nil {
		u := media.PublicURL(r, *p.ImageObjectKey)
		imageURL = &u
	}
	return productResponse{
		ID:            p.ID,
		Barcode:       p.Barcode,
		Name:          p.Name,
		Brand:         p.Brand,
		Manufacturer:  p.Manufacturer,
		Description:   p.Description,
		Price:         p.Price,
		Currency:      p.Currency,
		ImageURL:      imageURL,
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
	writeJSON(w, http.StatusOK, toProductResponse(r, result))
}

// SubmitDetails accepts multipart/form-data (text fields + optional "image" file)
// rather than JSON so the photo travels in the same request as the details.
func (h *Handler) SubmitDetails(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")

	r.Body = http.MaxBytesReader(w, r.Body, maxDetailsRequestBytes)
	// Small memory threshold: the image spills to a temp file instead of the heap.
	if err := r.ParseMultipartForm(256 << 10); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, media.ErrTooLarge.Error())
			return
		}
		writeError(w, http.StatusBadRequest, "expected multipart/form-data")
		return
	}
	defer func() { _ = r.MultipartForm.RemoveAll() }()

	input := DetailsInput{
		Name:         r.FormValue("name"),
		Brand:        r.FormValue("brand"),
		Manufacturer: r.FormValue("manufacturer"),
		Description:  r.FormValue("description"),
		Price:        r.FormValue("price"),
		Currency:     r.FormValue("currency"),
	}

	var image io.Reader
	if file, _, err := r.FormFile("image"); err == nil {
		defer file.Close()
		image = file
	} else if !errors.Is(err, http.ErrMissingFile) {
		writeError(w, http.StatusBadRequest, "invalid image upload")
		return
	}

	result, err := h.service.SubmitDetails(r.Context(), productID, input, image)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyName), errors.Is(err, ErrNameTooLong), errors.Is(err, ErrFieldTooLong),
			errors.Is(err, ErrDescTooLong), errors.Is(err, ErrInvalidPrice), errors.Is(err, ErrInvalidCurrency),
			errors.Is(err, media.ErrUnsupportedType):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, media.ErrTooLarge):
			writeError(w, http.StatusRequestEntityTooLarge, err.Error())
		case errors.Is(err, ErrNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrAlreadyDetailed):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(r, result))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
