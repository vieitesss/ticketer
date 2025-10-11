package database

import (
	"context"

	"github.com/vieitesss/ticketer/internal/models"
)

// ReceiptRepository defines the interface for receipt data access operations.
// This interface allows for easy swapping of database implementations (e.g., PostgreSQL, MySQL, MongoDB).
type ReceiptRepository interface {
	// CreateReceipt inserts a new receipt and its items into the database
	CreateReceipt(ctx context.Context, receipt *models.Receipt) (string, error)

	// GetReceipt retrieves a receipt by ID with all its items
	GetReceipt(ctx context.Context, id string) (*models.Receipt, error)

	// ListReceipts retrieves all receipts (without items, but with calculated totals)
	ListReceipts(ctx context.Context, limit, offset int) ([]models.Receipt, error)

	// Close closes the database connection
	Close()
}
