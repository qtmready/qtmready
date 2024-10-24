package db

import (
	"context"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db/config"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	Config = config.Connection
)

var (
	DefaultConfig = config.DefaultConnection
)

func WithConfig(conf *Config) config.Option {
	return config.WithConfig(conf)
}

// Connection is a wrapper around the config.Instance singleton.
func Connection(opts ...config.Option) *config.Connection {
	return config.Instance(opts...)
}

// Queries is a wrapper around the config.Queries singleton.
func Queries() *entities.Queries {
	return config.Queries()
}

// Transaction begins the transaction and wraps the queries in a transaction.
//
// Example:
//
//	tx, qtx, err := db.Transaction(ctx)
//	if err != nil { return err }
//
//	defer func() { _ = tx.Rollback(ctx) }()
//
//	// Do something with qtx. Any time you return on error, the transaction will be rolled back.
//	...
//
//	// Commit the transaction.
//	err = tx.Commit(ctx)
//	if err != nil { return err }
//
//	return nil
func Transaction(ctx context.Context) (pgx.Tx, *entities.Queries, error) {
	tx, err := Connection().Get().Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	return tx, Queries().WithTx(tx), nil
}
