package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vieitesss/ticketer/internal/models"
)

// UpsertStore inserts or updates a store by name, returns the store ID
func (r *PostgresRepository) UpsertStore(ctx context.Context, name string) (string, error) {
	var storeID string

	// UPSERT: Insert if not exists, otherwise return existing ID
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO stores (id, name)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, uuid.New().String(), name).Scan(&storeID)

	if err != nil {
		return "", fmt.Errorf("failed to upsert store: %w", err)
	}

	return storeID, nil
}

// GetStore retrieves a store by ID
func (r *PostgresRepository) GetStore(ctx context.Context, id string) (*models.Store, error) {
	var store models.Store

	err := r.Pool.QueryRow(ctx, `
		SELECT id, name
		FROM stores
		WHERE id = $1
	`, id).Scan(&store.ID, &store.Name)

	if err != nil {
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	return &store, nil
}

// ListStores retrieves all stores
func (r *PostgresRepository) ListStores(ctx context.Context) ([]models.Store, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, name
		FROM stores
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list stores: %w", err)
	}
	defer rows.Close()

	stores := []models.Store{}
	for rows.Next() {
		var store models.Store
		if err := rows.Scan(&store.ID, &store.Name); err != nil {
			return nil, fmt.Errorf("failed to scan store: %w", err)
		}
		stores = append(stores, store)
	}

	return stores, nil
}
