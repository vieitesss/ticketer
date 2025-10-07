package services

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/vieitesss/ticketer/internal/database"
	"github.com/vieitesss/ticketer/internal/models"
	"github.com/vieitesss/ticketer/internal/services/ai"
	"github.com/vieitesss/ticketer/internal/transport/dto"
)

type ReceiptService struct {
	aiService *ai.GeminiService
	db        *database.DB
}

func NewReceiptService(aiService *ai.GeminiService, db *database.DB) *ReceiptService {
	return &ReceiptService{
		aiService: aiService,
		db:        db,
	}
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
			log.Error("Failed to save receipt to database", "error", err)
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
		// For list view, we need to calculate total from items if available
		// or return 0 if items weren't loaded
		subtotal := 0.0
		for _, item := range receipt.Items {
			subtotal += item.Quantity * item.Price
		}
		totalAmount := subtotal - receipt.Discounts

		listItems[i] = dto.ReceiptListItem{
			ID:          receipt.ID,
			StoreName:   receipt.StoreName,
			TotalAmount: totalAmount,
			Discounts:   receipt.Discounts,
		}
	}

	return listItems, nil
}

// modelToDTO converts a receipt model to a DTO with calculated fields
func (s *ReceiptService) modelToDTO(receipt *models.Receipt) *dto.ReceiptResponse {
	// Calculate subtotal from items
	subtotal := 0.0
	items := make([]dto.ItemResponse, len(receipt.Items))
	for i, item := range receipt.Items {
		subtotal += item.Quantity * item.Price
		items[i] = dto.ItemResponse{
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	// Calculate total amount
	totalAmount := subtotal - receipt.Discounts

	return &dto.ReceiptResponse{
		StoreName:   receipt.StoreName,
		Items:       items,
		Subtotal:    subtotal,
		Discounts:   receipt.Discounts,
		TotalAmount: totalAmount,
	}
}
