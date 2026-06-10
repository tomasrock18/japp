package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomasrock18/japp/backend/internal/model"
	"github.com/tomasrock18/japp/backend/internal/storage"
)

type ProductHandler struct {
	storage *storage.MemoryStorage
}

func NewProductHandler(storage *storage.MemoryStorage) *ProductHandler {
	return &ProductHandler{storage: storage}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	barcode := chi.URLParam(r, "barcode")
	slog.Info("Getting product", "barcode", barcode)

	product, err := h.storage.GetProduct(barcode)
	if err != nil {
		slog.Warn("Product not found", "barcode", barcode)
		http.Error(w, `{"error": "product not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		slog.Warn("Error encoding product", "error", err)
		http.Error(w, `{"error": "error encoding product"}`, http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		slog.Error("Failed to create product", "error", err)
		http.Error(w, `{"error": "failed to create product"}`, http.StatusInternalServerError)
		return
	}

	if err := h.storage.CreateProduct(product); err != nil {
		slog.Error("Failed to create product", "error", err)
		http.Error(w, `{"error": "failed to create product"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(product)
	if err != nil {
		slog.Warn("Error encoding product", "error", err)
		http.Error(w, `{"error": "error encoding product"}`, http.StatusInternalServerError)
		return
	}
}
