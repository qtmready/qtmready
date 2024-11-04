package main

import (
	"log/slog"
	"os"
)

func configure_logger(debug bool) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler = slog.NewJSONHandler(os.Stdout, opts)

	if debug {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
