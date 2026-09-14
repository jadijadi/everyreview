package admin

import (
	"crypto/sha256"
	"crypto/subtle"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

//go:embed templates/dashboard.html
var templateFS embed.FS

// Delay on a wrong password: cheap brute-force damping without any state.
const wrongPasswordDelay = time.Second

type Handler struct {
	service  *Service
	password string
	tmpl     *template.Template
}

// NewHandler wires the dashboard; an empty password leaves it disabled (404) so a
// deployment never exposes stats by accident.
func NewHandler(service *Service, password string) *Handler {
	tmpl := template.Must(template.New("dashboard.html").Funcs(template.FuncMap{
		"ago":     ago,
		"avg":     formatAverage,
		"percent": percent,
		"stars":   stars,
		"snippet": snippet,
	}).ParseFS(templateFS, "templates/dashboard.html"))
	return &Handler{service: service, password: password, tmpl: tmpl}
}

func (h *Handler) Routes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(h.requirePassword)
		r.Get("/admin", h.Dashboard)
		r.Get("/admin/stats.json", h.StatsJSON)
	})
}

func (h *Handler) requirePassword(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.password == "" {
			http.NotFound(w, r)
			return
		}
		_, given, ok := r.BasicAuth()
		if !ok || !equalPasswords(given, h.password) {
			if ok {
				time.Sleep(wrongPasswordDelay)
				log.Printf("admin: rejected login from %s", r.RemoteAddr)
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="EveryReview admin", charset="UTF-8"`)
			http.Error(w, "password required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Hash both sides so the comparison is constant-time regardless of length.
func equalPasswords(given, want string) bool {
	g, w := sha256.Sum256([]byte(given)), sha256.Sum256([]byte(want))
	return subtle.ConstantTimeCompare(g[:], w[:]) == 1
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Dashboard(r.Context())
	if err != nil {
		log.Printf("admin: dashboard: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if err := h.tmpl.Execute(w, d); err != nil {
		log.Printf("admin: render: %v", err)
	}
}

func (h *Handler) StatsJSON(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Dashboard(r.Context())
	if err != nil {
		log.Printf("admin: stats: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(d)
}

func ago(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return plural(int(d.Minutes()), "minute")
	case d < 48*time.Hour:
		return plural(int(d.Hours()), "hour")
	default:
		return plural(int(d.Hours()/24), "day")
	}
}

func plural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit + " ago"
	}
	return fmt.Sprintf("%d %ss ago", n, unit)
}

func formatAverage(v *float64) string {
	if v == nil {
		return "—"
	}
	return fmt.Sprintf("%.2f", *v)
}

func percent(part, total int) int {
	if total == 0 {
		return 0
	}
	return part * 100 / total
}

func stars(rating *int) string {
	if rating == nil {
		return "—"
	}
	return strings.Repeat("★", *rating) + strings.Repeat("☆", 5-*rating)
}

func snippet(s string) string {
	const max = 120
	runes := []rune(strings.Join(strings.Fields(s), " "))
	if len(runes) <= max {
		return string(runes)
	}
	return string(runes[:max]) + "…"
}
