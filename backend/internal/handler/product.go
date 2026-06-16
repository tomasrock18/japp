package handler

import (
	"encoding/json"
	"fmt"
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

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(h.storage.Products)
	if err != nil {
		slog.Warn("Error encoding products", "error", err)
		http.Error(w, `{"error": "error encoding products"}`, http.StatusInternalServerError)
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

	if err := product.IsValid(); err != nil {
		slog.Error("Invalid product", "error", err)
		errorMsg := fmt.Sprintf(`{"error": "validation failed", "details": "%v"}`, err)
		http.Error(w, errorMsg, http.StatusBadRequest)
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

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	barcode := chi.URLParam(r, "barcode")
	slog.Info("Deleting product", "barcode", barcode)

	err := h.storage.DeleteProduct(barcode)
	if err != nil {
		http.Error(w, `{"error": product not found`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	barcode := chi.URLParam(r, "barcode")
	var bodyMap map[string]any

	product, err := h.storage.GetProduct(barcode)
	if err != nil {
		http.Error(w, `{"error": "product not found"}`, http.StatusNotFound)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
		http.Error(w, `{"error": "failed to parse product"}`, http.StatusBadRequest)
		return
	}

	for parameter := range bodyMap {
		err = product.UpdateField(parameter, bodyMap[parameter])
		if err != nil {
			errorMsg := fmt.Sprintf(`{"error": %v}`, err)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
	}

	err = h.storage.DeleteProduct(barcode)
	if err != nil {
		http.Error(w, `{"error": "failed to update product"}`, http.StatusNotFound)
		return
	}
	err = h.storage.CreateProduct(product)
	if err != nil {
		http.Error(w, `{"error": "failed to update product"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		slog.Warn("Error encoding product", "error", err)
		http.Error(w, `{"error": "error encoding product"}`, http.StatusInternalServerError)
		return
	}
}
