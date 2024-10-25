package config

import (
	"context"
	"log/slog"
	"sync"

	"go.breu.io/quantm/internal/db/entities"
)

var (
	_qry       *entities.Queries // Global database queries instance.
	_queryonce sync.Once         // Ensures queries initialization occurs only once.
)

// Queries returns a singleton instance of SQLC-generated queries, initialized with the Connection singleton's database connection.
//
// If no connection exists, Queries establishes one using the default environment-based configuration.  For more predictable behavior,
// explicitly establishing a connection with `Connection()` and `Connect()` is recommended.
func Queries() *entities.Queries {
	_queryonce.Do(func() {
		slog.Info("db: initializing queries ...")

		if _c == nil {
			slog.Warn("db: no connection, using environment variables to connect ...")

			ConnectionFromEnvironment().Connect(context.Background())
		}

		_qry = entities.New(Conn().conn)
	})

	return _qry
}
