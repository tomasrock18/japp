package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tomasrock18/japp/backend/internal/database"
	"github.com/tomasrock18/japp/backend/internal/handler"
	"github.com/tomasrock18/japp/backend/internal/repository/postgres"
)

func main() {
	_ = godotenv.Load("../.env")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.NewPostgresPool(
		ctx,
		fmt.Sprintf(
			"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASSWORD"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_NAME"),
		),
	)
	if err != nil {
		os.Exit(1)
	}
	defer pool.Close()

	productRepo := postgres.NewProductRepository(pool)

	productHandler := handler.NewProductHandler(productRepo)

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("%v | %v", r.Method, time.Now()))
			next.ServeHTTP(w, r)
		})
	}

	r := chi.NewRouter()
	r.Use(loggingMiddleware)

	// Product endpoints
	r.Get("/products/{barcode}", productHandler.GetProduct)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products", productHandler.GetAllProducts)
	r.Delete("/products/{barcode}", productHandler.DeleteProduct)
	r.Put("/products/{barcode}", productHandler.UpdateProduct)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		os.Exit(1)
	}

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		os.Exit(1)
	}
}
