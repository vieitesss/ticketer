package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vieitesss/ticketer/internal/models"
)

// UpsertProduct inserts or updates a product by name and store_id, returns the product ID
func (r *PostgresRepository) UpsertProduct(ctx context.Context, name, storeID string) (string, error) {
	var productID string

	// UPSERT: Insert if not exists (unique on name + store_id), otherwise return existing ID
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO products (id, name, store_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (name, store_id) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, uuid.New().String(), name, storeID).Scan(&productID)

	if err != nil {
		return "", fmt.Errorf("failed to upsert product: %w", err)
	}

	return productID, nil
}

// GetProduct retrieves a product by ID
func (r *PostgresRepository) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product

	err := r.Pool.QueryRow(ctx, `
		SELECT id, name, store_id
		FROM products
		WHERE id = $1
	`, id).Scan(&product.ID, &product.Name, &product.StoreID)

	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// ListProductsByStore retrieves all products for a specific store
func (r *PostgresRepository) ListProductsByStore(ctx context.Context, storeID string) ([]models.Product, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, name, store_id
		FROM products
		WHERE store_id = $1
		ORDER BY name
	`, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.StoreID); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}
