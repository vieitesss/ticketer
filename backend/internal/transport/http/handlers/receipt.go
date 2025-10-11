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

func (h *ReceiptHandler) HealthCheck(c fiber.Ctx) error {
	c.Response().Header.Add("Content-Type", "application/json")

	return c.SendString(`{"status": "ok"}`)
}
