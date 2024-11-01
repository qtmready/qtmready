package db

import (
	"context"

	"github.com/jackc/pgx/v5"

	dbcfg "go.breu.io/quantm/internal/db/config"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	Config = dbcfg.Config
)

var (
	DefaultConfig = dbcfg.Default
)

func WithConfig(conf *Config) dbcfg.ConfigOption {
	return dbcfg.WithConfig(conf)
}

// Connection is a wrapper around the dbcfg.Instance singleton.
func Connection(opts ...dbcfg.ConfigOption) *dbcfg.Config {
	return dbcfg.Instance(opts...)
}

// Queries is a wrapper around the dbcfg.Queries singleton.
func Queries() *entities.Queries {
	return dbcfg.Queries()
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
