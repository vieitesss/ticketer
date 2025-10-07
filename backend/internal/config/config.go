package config

import (
	"os"
)

type Config struct {
	GeminiAPIKey string
	LogLevel     string
	ServerPort   string
	DatabaseURL  string
}

func Load() *Config {
	return &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
		ServerPort:   getEnvOrDefault("PORT", "8080"),
		DatabaseURL:  getEnvOrDefault("DATABASE_URL", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
