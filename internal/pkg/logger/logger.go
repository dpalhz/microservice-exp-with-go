package logger

import (
	"log/slog"
	"os"
)

// New mengembalikan logger slog terstruktur baru yang menulis ke stdout dalam format JSON.
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}