package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tomasrock18/japp/backend/internal/handler"
	"github.com/tomasrock18/japp/backend/internal/storage"
)

func main() {
	_ = godotenv.Load("../.env")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("%v | %v", r.Method, time.Now()))
			next.ServeHTTP(w, r)
		})
	}

	store := storage.NewMemoryStorage()

	productHandler := handler.NewProductHandler(store)

	r := chi.NewRouter()
	r.Use(loggingMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Health check requested")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			slog.Error("Health check error", "error", err)
		}
	})
	r.Get("/products/{barcode}", productHandler.GetProduct)

	r.Post("/products", productHandler.CreateProduct)

	r.Get("/products", productHandler.GetAllProducts)

	r.Delete("/products/{barcode}", productHandler.DeleteProduct)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		slog.Warn("Miss BACKEND_PORT environment variable")
		os.Exit(1)
	}
	slog.Info("Server is starting", "port", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
