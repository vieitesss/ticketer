package logger

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func Init() {
	// Set UTC+2 timezone
	loc, err := time.LoadLocation("Europe/Madrid") // UTC+2 (CET/CEST)
	if err != nil {
		// Fallback to fixed UTC+2 offset
		loc = time.FixedZone("UTC+2", 2*60*60)
	}
	time.Local = loc

	// Configure logger with vibrant colors
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("[DEBUG]").
		Foreground(lipgloss.Color("63")). // Bright purple
		Bold(true)
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("[INFO]").
		Foreground(lipgloss.Color("86")). // Bright cyan
		Bold(true)
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("[WARN]").
		Foreground(lipgloss.Color("226")). // Bright yellow
		Bold(true)
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("[ERROR]").
		Foreground(lipgloss.Color("196")). // Bright red
		Bold(true)
	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("[FATAL]")

	styles.Key = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))       // Bright blue
	styles.Value = lipgloss.NewStyle().Foreground(lipgloss.Color("87"))     // Light cyan
	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // Gray

	log.SetStyles(styles)
	log.SetTimeFormat(time.Stamp)
	log.SetReportTimestamp(true)
	log.SetReportCaller(false)

	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	log.SetLevel(log.DebugLevel)

	// Set custom formatter for values to pretty print JSON
	log.SetFormatter(log.TextFormatter)
}

// FormatJSON formats a JSON string with indentation for logging
func FormatJSON(jsonStr string) string {
	var obj any
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return jsonStr // Return as-is if not valid JSON
	}
	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return jsonStr
	}
	return "\n" + string(formatted)
}
