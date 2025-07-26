package logger

import (
	"log/slog"
	"os"
)

// New creates a new slog.Logger instance with a JSON handler that writes to standard output.
// This logger can be used throughout the application for structured logging.
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}