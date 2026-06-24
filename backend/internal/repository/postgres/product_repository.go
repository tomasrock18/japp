package postgres

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tomasrock18/japp/backend/internal/model"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) GetProduct(ctx context.Context, barcode string) (model.Product, error) {
	var product model.Product

	query := `SELECT * FROM products WHERE barcode = $1`

	err := r.pool.QueryRow(ctx, query, barcode).Scan(
		&product.Barcode,
		&product.Name,
		&product.KcalPer100g,
		&product.ProteinPer100g,
		&product.FatPer100g,
		&product.CarbsPer100g,
		&product.CreatedBy,
		&product.CreatedAt,
	)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	query := `
        INSERT INTO products (barcode, name, kcal_per_100g, protein_per_100g, fat_per_100g, carbs_per_100g, created_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	telegramId, _ := strconv.Atoi(product.CreatedBy)

	_, err := r.pool.Exec(
		ctx,
		query,
		product.Barcode,
		product.Name,
		product.KcalPer100g,
		product.ProteinPer100g,
		product.FatPer100g,
		product.CarbsPer100g,
		telegramId,
	)

	return product, err
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	var products []model.Product

	query := `SELECT * FROM products ORDER BY created_at`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		if err := rows.Scan(
			&product.Barcode,
			&product.Name,
			&product.KcalPer100g,
			&product.ProteinPer100g,
			&product.FatPer100g,
			&product.CarbsPer100g,
			&product.CreatedBy,
			&product.CreatedAt,
		); err != nil {
			return products, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, barcode string) error {
	query := `DELETE FROM products WHERE barcode = $1`

	tag, err := r.pool.Exec(ctx, query, barcode)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product model.Product) error {
	query := `
        UPDATE products
        SET name = $2,
            kcal_per_100g = $3,
            protein_per_100g = $4,
            fat_per_100g = $5,
            carbs_per_100g = $6
        WHERE barcode = $1
    `
	tag, err := r.pool.Exec(
		ctx,
		query,
		product.Barcode,
		product.Name,
		product.KcalPer100g,
		product.ProteinPer100g,
		product.FatPer100g,
		product.CarbsPer100g,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("product not found")
	}

	return nil
}
