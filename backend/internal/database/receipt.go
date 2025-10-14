package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vieitesss/ticketer/internal/models"
)

// calculateReceiptHash generates a unique hash for a receipt based on store, date, and items
func calculateReceiptHash(storeName, boughtDate string, items []models.Item) string {
	// Sort items to ensure consistent hash regardless of order
	sortedItems := make([]models.Item, len(items))
	copy(sortedItems, items)
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].Name < sortedItems[j].Name
	})

	// Build hash string
	hashInput := fmt.Sprintf("%s|%s", storeName, boughtDate)
	for _, item := range sortedItems {
		hashInput += fmt.Sprintf("|%s:%.3f:%.2f", item.Name, item.Quantity, item.Price)
	}

	// Log hash input for debugging
	fmt.Printf("DEBUG: Hash input: %s\n", hashInput)

	// Calculate SHA-256 hash
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// CreateReceipt inserts a new receipt and its items into the database
func (r *PostgresRepository) CreateReceipt(ctx context.Context, receipt *models.Receipt) (string, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Calculate receipt hash for duplicate detection
	receiptHash := calculateReceiptHash(receipt.StoreName, receipt.BoughtDate, receipt.Items)

	// Log hash for debugging
	fmt.Printf("DEBUG: Receipt hash: %s (Store: %s, Date: %s, Items: %d)\n",
		receiptHash, receipt.StoreName, receipt.BoughtDate, len(receipt.Items))

	// Check if receipt already exists
	var existingID string
	err = tx.QueryRow(ctx, `
		SELECT id FROM receipts WHERE receipt_hash = $1
	`, receiptHash).Scan(&existingID)
	if err == nil {
		// Receipt already exists
		fmt.Printf("DEBUG: Duplicate detected! Existing ID: %s\n", existingID)
		return "", fmt.Errorf("duplicate receipt: this receipt has already been uploaded (ID: %s)", existingID)
	} else if err != pgx.ErrNoRows {
		return "", fmt.Errorf("failed to check for duplicate receipt: %w", err)
	}

	fmt.Printf("DEBUG: No duplicate found, creating new receipt\n")

	// UPSERT store and get store ID
	var storeID string
	err = tx.QueryRow(ctx, `
		INSERT INTO stores (id, name)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, uuid.New().String(), receipt.StoreName).Scan(&storeID)
	if err != nil {
		return "", fmt.Errorf("failed to upsert store: %w", err)
	}

	// Parse bought_date
	boughtDate, err := time.Parse("2006-01-02", receipt.BoughtDate)
	if err != nil {
		return "", fmt.Errorf("invalid bought_date format (expected YYYY-MM-DD): %w", err)
	}

	// Insert receipt
	receiptID := uuid.New().String()
	_, err = tx.Exec(ctx, `
		INSERT INTO receipts (id, store_id, discounts, receipt_hash, bought_date)
		VALUES ($1, $2, $3, $4, $5)
	`, receiptID, storeID, receipt.Discounts, receiptHash, boughtDate)
	if err != nil {
		return "", fmt.Errorf("failed to insert receipt: %w", err)
	}

	// Insert items with product UPSERT
	for _, item := range receipt.Items {
		// UPSERT product and get product ID
		var productID string
		err = tx.QueryRow(ctx, `
			INSERT INTO products (id, name, store_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (name, store_id) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, uuid.New().String(), item.Name, storeID).Scan(&productID)
		if err != nil {
			return "", fmt.Errorf("failed to upsert product: %w", err)
		}

		// Insert item
		_, err = tx.Exec(ctx, `
			INSERT INTO items (id, receipt_id, product_id, quantity, price_paid)
			VALUES ($1, $2, $3, $4, $5)
		`, uuid.New().String(), receiptID, productID, item.Quantity, item.Price)
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
func (r *PostgresRepository) GetReceipt(ctx context.Context, id string) (*models.Receipt, error) {
	// Get receipt with store information
	var receipt models.Receipt
	var discounts *float64
	var boughtDate time.Time
	err := r.Pool.QueryRow(ctx, `
		SELECT r.id, s.name, r.discounts, r.bought_date
		FROM receipts r
		JOIN stores s ON r.store_id = s.id
		WHERE r.id = $1
	`, id).Scan(&receipt.ID, &receipt.StoreName, &discounts, &boughtDate)
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
	receipt.BoughtDate = boughtDate.Format("2006-01-02")

	// Get items with product information
	rows, err := r.Pool.Query(ctx, `
		SELECT i.id, p.name, i.quantity, i.price_paid
		FROM items i
		JOIN products p ON i.product_id = p.id
		WHERE i.receipt_id = $1
		ORDER BY p.name
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

// ReceiptListItem is a lightweight receipt representation for list views
type ReceiptListItem struct {
	ID          string
	StoreName   string
	ItemCount   int
	BoughtDate  string
	Subtotal    float64
	Discounts   float64
	TotalAmount float64
}

// ListReceipts retrieves all receipts (without items, but with calculated totals)
func (r *PostgresRepository) ListReceipts(ctx context.Context, limit, offset int) ([]ReceiptListItem, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT
			r.id,
			s.name,
			COALESCE(COUNT(i.id), 0) as item_count,
			r.bought_date,
			COALESCE(SUM(i.quantity * i.price_paid), 0) as subtotal,
			COALESCE(r.discounts, 0) as discounts
		FROM receipts r
		JOIN stores s ON r.store_id = s.id
		LEFT JOIN items i ON r.id = i.receipt_id
		GROUP BY r.id, s.name, r.discounts, r.bought_date
		ORDER BY r.bought_date DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list receipts: %w", err)
	}
	defer rows.Close()

	receipts := []ReceiptListItem{}
	for rows.Next() {
		var receipt ReceiptListItem
		var boughtDate time.Time
		if err := rows.Scan(&receipt.ID, &receipt.StoreName, &receipt.ItemCount, &boughtDate, &receipt.Subtotal, &receipt.Discounts); err != nil {
			return nil, fmt.Errorf("failed to scan receipt: %w", err)
		}
		receipt.BoughtDate = boughtDate.Format("2006-01-02")
		receipt.TotalAmount = receipt.Subtotal - receipt.Discounts
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

// ListReceiptsByDateRange retrieves receipts within a date range
func (r *PostgresRepository) ListReceiptsByDateRange(ctx context.Context, startDate, endDate string, limit, offset int) ([]models.Receipt, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT
			r.id,
			s.name,
			r.discounts,
			r.bought_date,
			COALESCE(SUM(i.quantity * i.price_paid), 0) as subtotal
		FROM receipts r
		JOIN stores s ON r.store_id = s.id
		LEFT JOIN items i ON r.id = i.receipt_id
		WHERE r.bought_date >= $1 AND r.bought_date <= $2
		GROUP BY r.id, s.name, r.discounts, r.bought_date
		ORDER BY r.bought_date DESC
		LIMIT $3 OFFSET $4
	`, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list receipts by date range: %w", err)
	}
	defer rows.Close()

	receipts := []models.Receipt{}
	for rows.Next() {
		var receipt models.Receipt
		var discounts *float64
		var boughtDate time.Time
		var subtotal float64
		if err := rows.Scan(&receipt.ID, &receipt.StoreName, &discounts, &boughtDate, &subtotal); err != nil {
			return nil, fmt.Errorf("failed to scan receipt: %w", err)
		}
		if discounts != nil {
			receipt.Discounts = *discounts
		} else {
			receipt.Discounts = 0
		}
		receipt.BoughtDate = boughtDate.Format("2006-01-02")
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

// GetReceiptsByStore retrieves all receipts from a specific store
func (r *PostgresRepository) GetReceiptsByStore(ctx context.Context, storeID string, limit, offset int) ([]models.Receipt, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT
			r.id,
			s.name,
			r.discounts,
			r.bought_date,
			COALESCE(SUM(i.quantity * i.price_paid), 0) as subtotal
		FROM receipts r
		JOIN stores s ON r.store_id = s.id
		LEFT JOIN items i ON r.id = i.receipt_id
		WHERE r.store_id = $1
		GROUP BY r.id, s.name, r.discounts, r.bought_date
		ORDER BY r.bought_date DESC
		LIMIT $2 OFFSET $3
	`, storeID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipts by store: %w", err)
	}
	defer rows.Close()

	receipts := []models.Receipt{}
	for rows.Next() {
		var receipt models.Receipt
		var discounts *float64
		var boughtDate time.Time
		var subtotal float64
		if err := rows.Scan(&receipt.ID, &receipt.StoreName, &discounts, &boughtDate, &subtotal); err != nil {
			return nil, fmt.Errorf("failed to scan receipt: %w", err)
		}
		if discounts != nil {
			receipt.Discounts = *discounts
		} else {
			receipt.Discounts = 0
		}
		receipt.BoughtDate = boughtDate.Format("2006-01-02")
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

// DeleteReceipt deletes a receipt and all its items (CASCADE)
func (r *PostgresRepository) DeleteReceipt(ctx context.Context, id string) error {
	result, err := r.Pool.Exec(ctx, `
		DELETE FROM receipts WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete receipt: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("receipt not found")
	}

	return nil
}

// UpdateItem updates an item's quantity and price_paid
func (r *PostgresRepository) UpdateItem(ctx context.Context, itemID string, quantity, pricePaid float64) error {
	result, err := r.Pool.Exec(ctx, `
		UPDATE items
		SET quantity = $1, price_paid = $2
		WHERE id = $3
	`, quantity, pricePaid, itemID)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}
