package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tomasrock18/japp/backend/internal/handler"
	"github.com/tomasrock18/japp/backend/internal/storage"
)

func main() {
	_ = godotenv.Load("../.env")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	store := storage.NewMemoryStorage()

	productHandler := handler.NewProductHandler(store)

	r := chi.NewRouter()

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
