package main

import (
	"github.com/charmbracelet/log"
	"github.com/vieitesss/ticketer/internal/app"
	"github.com/vieitesss/ticketer/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to initialize application", "error", err)
	}

	if err := application.Run(); err != nil {
		log.Fatal("Application error", "error", err)
	}
}
