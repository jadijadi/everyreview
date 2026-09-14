package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"everyreview/backend/internal/media"
	"everyreview/backend/internal/platform/db"
	"everyreview/backend/internal/product"
	"everyreview/backend/internal/review"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer conn.Close()

	mediaDir := os.Getenv("MEDIA_DIR")
	if mediaDir == "" {
		mediaDir = "./data/media"
	}
	mediaStore, err := media.NewDiskStore(mediaDir)
	if err != nil {
		log.Fatalf("preparing media dir: %v", err)
	}

	productHandler := product.NewHandler(product.NewService(product.NewPostgresRepository(conn), mediaStore))
	mediaHandler := media.NewHandler(mediaStore)
	reviewHandler := review.NewHandler(review.NewService(review.NewPostgresRepository(conn)))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(v1 chi.Router) {
		productHandler.Routes(v1)
		reviewHandler.Routes(v1)
		mediaHandler.Routes(v1)
	})

	addr := ":" + port
	log.Printf("everyreview backend listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
