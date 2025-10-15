package services

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/vieitesss/ticketer/internal/database"
	"github.com/vieitesss/ticketer/internal/models"
	"github.com/vieitesss/ticketer/internal/services/ai"
	"github.com/vieitesss/ticketer/internal/transport/dto"
)

type ReceiptService struct {
	aiService *ai.GeminiService
	db        database.ReceiptRepository
}

func NewReceiptService(ctx context.Context, db database.ReceiptRepository) (*ReceiptService, error) {
	// Initialize AI service
	geminiService, err := ai.NewGeminiService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini service: %w", err)
	}

	return &ReceiptService{
		aiService: geminiService,
		db:        db,
	}, nil
}

func (s *ReceiptService) ProcessReceipt(ctx context.Context, imagePath string) (*dto.ReceiptResponse, error) {
	// Process receipt using AI service
	receipt, err := s.aiService.ProcessReceipt(ctx, imagePath)
	if err != nil {
		return nil, err
	}

	// Save to database if available
	if s.db != nil {
		receiptID, err := s.db.CreateReceipt(ctx, receipt)
		if err != nil {
			log.Warn("Failed to save receipt to database", "error", err)
			// Don't fail the request if database save fails
		} else {
			log.Info("Receipt saved to database", "id", receiptID)
		}
	}

	// Convert to DTO with calculated fields
	return s.modelToDTO(receipt), nil
}

func (s *ReceiptService) GetReceipt(ctx context.Context, id string) (*dto.ReceiptResponse, error) {
	receipt, err := s.db.GetReceipt(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.modelToDTO(receipt), nil
}

func (s *ReceiptService) ListReceipts(ctx context.Context, limit, offset int) ([]dto.ReceiptListItem, error) {
	receipts, err := s.db.ListReceipts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	listItems := make([]dto.ReceiptListItem, len(receipts))
	for i, receipt := range receipts {
		listItems[i] = dto.ReceiptListItem{
			ID:          receipt.ID,
			StoreName:   receipt.StoreName,
			ItemCount:   receipt.ItemCount,
			BoughtDate:  receipt.BoughtDate,
			TotalAmount: receipt.TotalAmount,
		}
	}

	return listItems, nil
}

// DeleteReceipt deletes a receipt by ID
func (s *ReceiptService) DeleteReceipt(ctx context.Context, id string) error {
	return s.db.DeleteReceipt(ctx, id)
}

// UpdateItem updates an item's quantity and price
func (s *ReceiptService) UpdateItem(ctx context.Context, itemID string, quantity, pricePaid float64) error {
	return s.db.UpdateItem(ctx, itemID, quantity, pricePaid)
}

// modelToDTO converts a receipt model to a DTO with calculated fields
func (s *ReceiptService) modelToDTO(receipt *models.Receipt) *dto.ReceiptResponse {
	// Calculate subtotal from items
	subtotal := 0.0
	items := make([]dto.ItemResponse, len(receipt.Items))
	for i, item := range receipt.Items {
		itemSubtotal := item.Quantity * item.Price
		subtotal += itemSubtotal
		items[i] = dto.ItemResponse{
			ID:          item.ID,
			ProductID:   item.ID, // TODO: Need to get actual product_id from database
			ProductName: item.Name,
			Quantity:    item.Quantity,
			PricePaid:   item.Price,
			Subtotal:    itemSubtotal,
		}
	}

	// Calculate total amount
	totalAmount := subtotal - receipt.Discounts

	// Extract store info (for now, we only have store name)
	// TODO: Need to get actual store ID from database
	storeResponse := dto.StoreResponse{
		ID:   "",
		Name: receipt.StoreName,
	}

	return &dto.ReceiptResponse{
		ID:          receipt.ID,
		Store:       storeResponse,
		BoughtDate:  receipt.BoughtDate,
		Items:       items,
		Subtotal:    subtotal,
		Discounts:   receipt.Discounts,
		TotalAmount: totalAmount,
	}
}
