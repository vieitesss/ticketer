package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vieitesss/ticketer/internal/models"
)

// CreateReceipt inserts a new receipt and its items into the database
func (db *DB) CreateReceipt(ctx context.Context, receipt *models.Receipt) (string, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	receiptID := uuid.New().String()
	now := time.Now()

	// Insert receipt (without total_amount - it's calculated)
	_, err = tx.Exec(ctx, `
		INSERT INTO receipts (id, store_name, discounts, created_at)
		VALUES ($1, $2, $3, $4)
	`, receiptID, receipt.StoreName, receipt.Discounts, now)
	if err != nil {
		return "", fmt.Errorf("failed to insert receipt: %w", err)
	}

	// Insert items
	for _, item := range receipt.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO items (id, receipt_id, name, quantity, price)
			VALUES ($1, $2, $3, $4, $5)
		`, uuid.New().String(), receiptID, item.Name, item.Quantity, item.Price)
		if err != nil {
			return "", fmt.Errorf("failed to insert item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return receiptID, nil
}

// GetReceipt retrieves a receipt by ID with all its items
func (db *DB) GetReceipt(ctx context.Context, id string) (*models.Receipt, error) {
	// Get receipt
	var receipt models.Receipt
	var discounts *float64
	err := db.Pool.QueryRow(ctx, `
		SELECT id, store_name, discounts, created_at
		FROM receipts
		WHERE id = $1
	`, id).Scan(&receipt.ID, &receipt.StoreName, &discounts, &receipt.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("receipt not found")
		}
		return nil, fmt.Errorf("failed to get receipt: %w", err)
	}

	if discounts != nil {
		receipt.Discounts = *discounts
	} else {
		receipt.Discounts = 0
	}

	// Get items
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, quantity, price
		FROM items
		WHERE receipt_id = $1
		ORDER BY name
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	receipt.Items = []models.Item{}
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Price); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		receipt.Items = append(receipt.Items, item)
	}

	return &receipt, nil
}

// ListReceipts retrieves all receipts (without items, but with calculated totals)
func (db *DB) ListReceipts(ctx context.Context, limit, offset int) ([]models.Receipt, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT
			r.id,
			r.store_name,
			r.discounts,
			COALESCE(SUM(i.quantity * i.price), 0) as subtotal
		FROM receipts r
		LEFT JOIN items i ON r.id = i.receipt_id
		GROUP BY r.id, r.store_name, r.discounts
		ORDER BY r.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list receipts: %w", err)
	}
	defer rows.Close()

	receipts := []models.Receipt{}
	for rows.Next() {
		var receipt models.Receipt
		var discounts *float64
		var subtotal float64
		if err := rows.Scan(&receipt.ID, &receipt.StoreName, &discounts, &subtotal); err != nil {
			return nil, fmt.Errorf("failed to scan receipt: %w", err)
		}
		if discounts != nil {
			receipt.Discounts = *discounts
		} else {
			receipt.Discounts = 0
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}
