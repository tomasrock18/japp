package repository

import (
	"context"

	"github.com/tomasrock18/japp/backend/internal/model"
)

type ProductRepository interface {
	GetProduct(ctx context.Context, barcode string) (model.Product, error)
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	DeleteProduct(ctx context.Context, barcode string) error
	UpdateProduct(ctx context.Context, product model.Product) error
}
