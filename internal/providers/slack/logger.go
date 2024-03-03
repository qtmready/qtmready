package slack

import (
	"log"
	"log/slog"
	"os"

	"go.breu.io/quantm/internal/shared"
)

// returns a logger with slog handler.
// TODO: figure out a way to add prefix. slog doesn't support it.
func logger() *log.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
	}

	if shared.Service().GetDebug() {
		return slog.NewLogLogger(slog.NewTextHandler(os.Stdout, opts), slog.LevelDebug)
	}

	return slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, opts), slog.LevelInfo)
}
