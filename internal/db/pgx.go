package db

import (
	"embed"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"go.breu.io/quantm/internal/db/config"
)

var (
	//go:embed migrations/postgres/*.sql
	sql embed.FS
)

// WithPostgresMigrations configures PostgreSQL database migrations.
// TODO - move to function return.
func WithPostgresMigrations() {
	// TODO: read from .env
	c := &config.DefaultConn

	dir, err := iofs.New(sql, "migrations/postgres")
	if err != nil {
		slog.Error("db: failed to initialize migrations", "error", err.Error())
		return
	}

	migrations, err := migrate.NewWithSourceInstance(
		"iofs",
		dir,
		c.ConnectionUri(),
	)
	if err != nil {
		slog.Error("db: failed to create migrations instance", "error", err.Error())
		return
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Warn("db: failed to run migrations", "error", err.Error())
		return
	}

	if err == migrate.ErrNoChange {
		slog.Info("db: no new migrations to run")
	}

	slog.Info("db: migrations done successfully")
}
