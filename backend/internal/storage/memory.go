package storage

import (
	"errors"
	"sync"

	"github.com/tomasrock18/japp/backend/internal/model"
)

type MemoryStorage struct {
	mu       sync.RWMutex
	Products map[string]model.Product
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Products: make(map[string]model.Product),
	}
}

func (s *MemoryStorage) GetProduct(barcode string) (model.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, exists := s.Products[barcode]
	if !exists {
		return model.Product{}, errors.New("product not found")
	}
	return product, nil
}

func (s *MemoryStorage) CreateProduct(product model.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Products[product.Barcode] = product
	return nil
}
