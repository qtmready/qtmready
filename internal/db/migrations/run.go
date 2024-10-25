package migrations

import (
	"context"
	"embed"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"go.breu.io/quantm/internal/db/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

var (
	//go:embed postgres/*.sql
	sql embed.FS
)

// Run runs the migrations for the PostgreSQL database.
func Run(ctx context.Context, connection *config.Connection) {
	if !connection.IsConnected() {
		connection.Start(ctx)
	}

	dir, err := iofs.New(sql, "migrations/postgres")
	if err != nil {
		slog.Error("db: unable to read migrations ...", "error", err.Error())

		return
	}

	migrations, err := migrate.NewWithSourceInstance(
		"iofs",
		dir,
		connection.ConnectionURI(),
	)

	if err != nil {
		slog.Error("db: failed to create migrations instance ...", "error", err.Error())
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
