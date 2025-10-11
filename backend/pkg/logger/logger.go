package logger

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// Init initializes the global Charmbracelet logger with configuration matching the Fiber middleware
func Init() {
	// Set timezone to Europe/Madrid (UTC+2)
	loc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		// Fallback to fixed UTC+2 offset
		loc = time.FixedZone("UTC+2", 2*60*60)
	}
	time.Local = loc

	// Configure time format to match middleware (HH:MM:SS.mmm)
	log.SetTimeFormat("15:04:05.000")
	log.SetReportTimestamp(true)
	log.SetReportCaller(false)

	// Configure color styles to match Fiber middleware
	styles := log.DefaultStyles()

	// Standard log level colors (matching ANSI colors)
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBU").
		Foreground(lipgloss.Color("13")) // Bright Magenta (ANSI)

	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Foreground(lipgloss.Color("14")) // Bright Cyan (ANSI)

	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		Foreground(lipgloss.Color("11")) // Bright Yellow (ANSI)

	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERRO").
		Foreground(lipgloss.Color("9")) // Bright Red (ANSI)

	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FTAL").
		Bold(true).
		Foreground(lipgloss.Color("9")) // Bright Red Bold (ANSI)

	// Keys in bright green (matches HTTP label in middleware)
	styles.Key = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")) // Bright Green (ANSI)

	// Values in white (matches latency in middleware)
	styles.Value = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")) // Bright White (ANSI)

	// Message text in white
	styles.Message = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	log.SetStyles(styles)

	// Set log level from environment variable
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

// DebugJSON logs a value as formatted JSON at debug level
func DebugJSON(msg string, key string, v any) {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Error("Failed to marshal JSON", "error", err)
		return
	}
	log.Debug(msg, key, "\n"+string(jsonBytes))
}
