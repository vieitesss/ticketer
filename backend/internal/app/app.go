package app

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v3"
	"github.com/vieitesss/ticketer/internal/config"
	"github.com/vieitesss/ticketer/internal/database"
	"github.com/vieitesss/ticketer/internal/services"
	"github.com/vieitesss/ticketer/internal/transport/http"
	"github.com/vieitesss/ticketer/internal/transport/http/handlers"
	"github.com/vieitesss/ticketer/internal/transport/http/routers"
)

type App struct {
	config *config.Config
	db     database.ReceiptRepository
	server *fiber.App
}

func New() (*App, error) {
	// Load configuration
	cfg := config.Load()

	// Initialize context
	ctx := context.Background()

	// Initialize database (optional - only if DATABASE_URL is set)
	var db database.ReceiptRepository
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

	// Initialize services
	receiptService, err := services.NewReceiptService(ctx, db)
	if err != nil {
		log.Error("Failed to initialize receipt service", "error", err)
		return nil, fmt.Errorf("Failed to initialize receipt service: %w", err)
	}

	// Create HTTP server
	server := http.NewServer()

	// Initialize HTTP handler
	receiptHandler := handlers.NewReceiptHandler(receiptService)

	// Setup routes
	api := server.Group("/api")
	routers.NewReceiptRouter(api, receiptHandler)

	return &App{
		config: cfg,
		db:     db,
		server: server,
	}, nil
}

func (a *App) Run() error {
	log.Info("Server starting", "port", a.config.ServerPort)
	if err := a.server.Listen(fmt.Sprint(":" + a.config.ServerPort)); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}

func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}
