package main

import (
	"log/slog"
	"os"
)

func configure_logger(debug bool) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	handler = slog.NewJSONHandler(os.Stdout, opts)

	if debug {
		opts.Level = slog.LevelDebug
		opts.AddSource = false
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
