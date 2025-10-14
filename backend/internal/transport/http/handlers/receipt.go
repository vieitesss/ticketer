package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/vieitesss/ticketer/internal/services"
)

type ReceiptHandler struct {
	receiptService *services.ReceiptService
}

func NewReceiptHandler(receiptService *services.ReceiptService) ReceiptHandler {
	return ReceiptHandler{
		receiptService: receiptService,
	}
}

func (h *ReceiptHandler) UploadAndProcess(c fiber.Ctx) error {
	// Parse multipart form (max 10MB)
	form, err := c.MultipartForm()

	if err != nil {
		log.Error("Failed to parse multipart form", "error", err)
		return c.Status(http.StatusBadRequest).SendString("Failed to parse form")
	}

	files := form.File["receipt"]
	if len(files) == 0 {
		return c.Status(http.StatusBadRequest).SendString("No file uploaded")
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		log.Error("Failed to open uploaded file", "error", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to open file")
	}
	defer file.Close()

	// Validate file extension
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return c.Status(http.StatusBadRequest).SendString("Invalid file format. Only JPG, JPEG, and PNG are allowed")
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "/app/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Error("Failed to create uploads directory", "error", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal server error")
	}

	// Save file temporarily
	tempPath := filepath.Join(uploadsDir, fileHeader.Filename)
	dst, err := os.Create(tempPath)
	if err != nil {
		log.Error("Failed to create temp file", "path", tempPath, "error", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to save file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Error("Failed to save file", "path", tempPath, "error", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to save file")
	}

	// Process receipt through service layer
	receipt, err := h.receiptService.ProcessReceipt(c.Context(), tempPath)
	if err != nil {
		log.Error("Failed to process receipt", "path", tempPath, "error", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Failed to process receipt: %v", err))
	}

	// Clean up temp file
	os.Remove(tempPath)

	// Return JSON response
	return c.JSON(receipt)
}

// GetReceipt retrieves a single receipt with full details
func (h *ReceiptHandler) GetReceipt(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("Receipt ID is required")
	}

	receipt, err := h.receiptService.GetReceipt(c.Context(), id)
	if err != nil {
		log.Error("Failed to get receipt", "id", id, "error", err)
		return c.Status(http.StatusNotFound).SendString("Receipt not found")
	}

	return c.JSON(receipt)
}

// ListReceipts retrieves all receipts (for the sidebar list)
func (h *ReceiptHandler) ListReceipts(c fiber.Ctx) error {
	limit := 50
	offset := 0

	if limitQuery := c.Query("limit"); limitQuery != "" {
		fmt.Sscanf(limitQuery, "%d", &limit)
	}
	if offsetQuery := c.Query("offset"); offsetQuery != "" {
		fmt.Sscanf(offsetQuery, "%d", &offset)
	}

	receipts, err := h.receiptService.ListReceipts(c.Context(), limit, offset)
	if err != nil {
		log.Error("Failed to list receipts", "error", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to list receipts")
	}

	return c.JSON(receipts)
}

// DeleteReceipt deletes a receipt by ID
func (h *ReceiptHandler) DeleteReceipt(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("Receipt ID is required")
	}

	err := h.receiptService.DeleteReceipt(c.Context(), id)
	if err != nil {
		log.Error("Failed to delete receipt", "id", id, "error", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to delete receipt")
	}

	return c.SendStatus(http.StatusNoContent)
}

// UpdateItem updates an item's quantity and price
func (h *ReceiptHandler) UpdateItem(c fiber.Ctx) error {
	itemID := c.Params("itemId")
	if itemID == "" {
		return c.Status(http.StatusBadRequest).SendString("Item ID is required")
	}

	var req struct {
		Quantity  float64 `json:"quantity"`
		PricePaid float64 `json:"price_paid"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}

	if req.Quantity <= 0 {
		return c.Status(http.StatusBadRequest).SendString("Quantity must be greater than 0")
	}

	if req.PricePaid < 0 {
		return c.Status(http.StatusBadRequest).SendString("Price must be non-negative")
	}

	err := h.receiptService.UpdateItem(c.Context(), itemID, req.Quantity, req.PricePaid)
	if err != nil {
		log.Error("Failed to update item", "itemId", itemID, "error", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to update item")
	}

	return c.SendStatus(http.StatusNoContent)
}

func (h *ReceiptHandler) HealthCheck(c fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "application/json")

	return c.SendString(`{"status": "ok"}`)
}
