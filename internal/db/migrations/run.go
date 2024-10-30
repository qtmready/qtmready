package migrations

import (
	"context"
	"embed"
	"fmt"
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
func Run(ctx context.Context, connection *config.Config) error {
	slog.Info("migrations: running ...")

	if !connection.IsConnected() {
		_ = connection.Start(ctx)
		defer func() { _ = connection.Stop(ctx) }()
	}

	dir, err := iofs.New(sql, "postgres")
	if err != nil {
		slog.Error("migrations: unable to read ...", "error", err.Error())

		return err
	}

	migrations, err := migrate.NewWithSourceInstance(
		"iofs",
		dir,
		connection.ConnectionURI(),
	)

	if err != nil {
		slog.Error("migrations: unable to read data ...", "error", err.Error())
		return err
	}

	version, dirty, _ := migrations.Version()

	if dirty {
		return fmt.Errorf("migrations:  unapplied migrations at version, %d", version)
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		slog.Info("migrations: nothing new since ...", "version", version)
	}

	slog.Info("migrations: migrations done successfully")

	return nil
}
