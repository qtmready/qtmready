package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.breu.io/graceful"

	"go.breu.io/quantm/cmd/quantm/config"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/migrations"
)

func main() {
	conf := config.New()
	ctx := context.Background()

	conf.Load()
	conf.Parse()

	if conf.Mode == config.ModeMigrate {
		if err := migrations.Run(ctx, db.Get(db.WithConfig(conf.DB))); err != nil {
			slog.Error("unable to run migrations", "error", err.Error())

			os.Exit(1)
		}

		os.Exit(0)
	}

	quit := make(chan os.Signal, 1)
	app := graceful.New()

	if err := conf.Setup(app); err != nil {
		slog.Error("unable to setup ...", "error", err.Error())
		os.Exit(1)
	}

	if err := app.Start(ctx); err != nil {
		slog.Error("unable to start ...", "error", err.Error())

		os.Exit(1)
	}

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	<-quit

	// - gracefully stop the services

	if err := app.Stop(ctx); err != nil {
		slog.Error("unable to stop service", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("service stopped, exiting...")

	os.Exit(0)
}
