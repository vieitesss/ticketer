package main

import (
	"context"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/vieitesss/ticketer/internal/handlers"
	"github.com/vieitesss/ticketer/internal/logger"
	"github.com/vieitesss/ticketer/internal/services/ai"
)

func main() {
	// Initialize logger
	logger.Init()

	ctx := context.Background()

	// Initialize Gemini service
	geminiService, err := ai.NewGeminiService(ctx)
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v", err)
	}

	// Initialize handlers
	receiptHandler := handlers.NewReceiptHandler(geminiService)

	// Setup routes
	http.HandleFunc("/api/receipts/upload", receiptHandler.UploadAndProcess)
	http.HandleFunc("/api/health", receiptHandler.HealthCheck)

	// Enable CORS for development
	handler := enableCORS(http.DefaultServeMux)

	// Start server
	port := ":8080"
	log.Info("Server starting", "port", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal("Server failed to start", "error", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
