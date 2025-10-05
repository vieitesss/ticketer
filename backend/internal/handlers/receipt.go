package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/vieitesss/ticketer/internal/services/ai"
)

type ReceiptHandler struct {
	aiService *ai.GeminiService
}

func NewReceiptHandler(aiService *ai.GeminiService) *ReceiptHandler {
	return &ReceiptHandler{
		aiService: aiService,
	}
}

func (h *ReceiptHandler) UploadAndProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Error("Failed to parse form", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("receipt")
	if err != nil {
		log.Error("Failed to get file from form", "error", err)
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		http.Error(w, "Invalid file format. Only JPG, JPEG, and PNG are allowed", http.StatusBadRequest)
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "/app/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Error("Failed to create uploads directory", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Save file temporarily
	tempPath := filepath.Join(uploadsDir, header.Filename)
	dst, err := os.Create(tempPath)
	if err != nil {
		log.Error("Failed to create temp file", "path", tempPath, "error", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Error("Failed to save file", "path", tempPath, "error", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Process receipt
	receipt, err := h.aiService.ProcessReceipt(r.Context(), tempPath)
	if err != nil {
		log.Error("Failed to process receipt", "path", tempPath, "error", err)
		http.Error(w, fmt.Sprintf("Failed to process receipt: %v", err), http.StatusInternalServerError)
		return
	}

	// Clean up temp file
	os.Remove(tempPath)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func (h *ReceiptHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
