package media

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store *DiskStore
}

func NewHandler(store *DiskStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/media/*", h.Serve)
}

func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "*")
	if !keyPattern.MatchString(key) {
		http.NotFound(w, r)
		return
	}
	// Keys are random and never reused, so the bytes behind a URL never change.
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeFile(w, r, h.store.path(key))
}

// PublicURL builds the absolute URL a client fetches an object from. The API sits
// behind Apache/Caddy on the VPS, so scheme/host come from the forwarded headers
// unless PUBLIC_BASE_URL pins them explicitly.
func PublicURL(r *http.Request, key string) string {
	base := strings.TrimSuffix(os.Getenv("PUBLIC_BASE_URL"), "/")
	if base == "" {
		scheme := r.Header.Get("X-Forwarded-Proto")
		if scheme == "" {
			scheme = "http"
			if r.TLS != nil {
				scheme = "https"
			}
		}
		base = scheme + "://" + r.Host
	}
	return base + "/v1/media/" + key
}
