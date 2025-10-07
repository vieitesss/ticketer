package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/vieitesss/ticketer/internal/config"
	"github.com/vieitesss/ticketer/internal/database"
	"github.com/vieitesss/ticketer/internal/services"
	"github.com/vieitesss/ticketer/internal/services/ai"
	httpTransport "github.com/vieitesss/ticketer/internal/transport/http"
	"github.com/vieitesss/ticketer/pkg/logger"
)

type App struct {
	config *config.Config
	db     *database.DB
	server *http.Server
}

func New() (*App, error) {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg := config.Load()

	// Initialize context
	ctx := context.Background()

	// Initialize database (optional - only if DATABASE_URL is set)
	var db *database.DB
	if cfg.DatabaseURL != "" {
		var err error
		db, err = database.NewPostgres(ctx, cfg.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create database connection: %w", err)
		}
		log.Info("Database connected successfully")
	} else {
		log.Warn("No DATABASE_URL provided, running without database")
	}

	// Initialize AI service
	geminiService, err := ai.NewGeminiService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini service: %w", err)
	}

	// Initialize services
	receiptService := services.NewReceiptService(geminiService, db)

	// Initialize HTTP handler
	handler := httpTransport.NewHandler(receiptService)

	// Setup routes
	router := httpTransport.SetupRoutes(handler)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	return &App{
		config: cfg,
		db:     db,
		server: server,
	}, nil
}

func (a *App) Run() error {
	log.Info("Server starting", "port", a.config.ServerPort)
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}

func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}
