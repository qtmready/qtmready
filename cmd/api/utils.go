package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/core/ws"
)

// _run runs a function in a goroutine.
func _run(fn func() error, ch chan error) {
	if err := fn(); err != nil {
		ch <- err
	}
}

// _serve starts the echo server in a goroutine.
func _serve(e *echo.Echo, port string) func() error {
	return func() error { return e.Start(":" + port) }
}

func _hub() error {
	worker := ws.ConnectionsHubWorker()

	return worker.Start()
}

// _graceful shuts down each goroutine gracefully.
func _graceful(ctx context.Context, fns []shutdownfn, signals []chan any, code int) {
	for _, signal := range signals {
		signal <- true
	}

	for _, fn := range fns {
		if err := fn(ctx); err != nil {
			code = 1
		}
	}

	slog.Info("shutdown complete, exiting.")

	os.Exit(code)
}
