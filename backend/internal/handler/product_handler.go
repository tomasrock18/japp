package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomasrock18/japp/backend/internal/model"
	"github.com/tomasrock18/japp/backend/internal/repository"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	barcode := chi.URLParam(r, "barcode")

	product, err := h.repo.GetProduct(ctx, barcode)
	if err != nil {
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.repo.GetAllProducts(ctx)
	if err != nil {
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var product model.Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := product.IsValid(); err != nil {
		http.Error(w, `{"error": "product validation failed"}`, http.StatusBadRequest)
		return
	}

	if _, err := h.repo.CreateProduct(ctx, product); err != nil {
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	barcode := chi.URLParam(r, "barcode")

	if err := h.repo.DeleteProduct(ctx, barcode); err != nil {
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	barcode := chi.URLParam(r, "barcode")

	prodcut, err := h.repo.GetProduct(ctx, barcode)
	if err != nil {
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}

	var bodyMap map[string]any
	if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	for field, value := range bodyMap {
		if err := prodcut.UpdateField(field, value); err != nil {
			http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
			return
		}
	}

	if err := h.repo.UpdateProduct(ctx, prodcut); err != nil {
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, prodcut)
}
