// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.  

package db

import (
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/scylladb/gocqlx/v2"
	"go.breu.io/ctrlplane/internal/shared"
)

const (
	NullUUID   = "00000000-0000-0000-0000-000000000000"
	NullString = ""
)

var (
	DB = &db{}
)

type (
	// Holds the information about the database
	db struct {
		gocqlx.Session
		Hosts              []string `env:"CASSANDRA_HOSTS" env-default:"ctrlplane-database"`
		Keyspace           string   `env:"CASSANDRA_KEYSPACE" env-default:"ctrlplane"`
		MigrationSourceURL string   `env:"CASSANDRA_MIGRATION_SOURCE_URL"`
	}
)

// ReadEnv reads the environment variables
func (d *db) ReadEnv() {
	_ = cleanenv.ReadEnv(d)
}

// InitSession initializes the session with the configured hosts
func (d *db) InitSession() {
	cluster := gocql.NewCluster(d.Hosts...)
	cluster.Keyspace = d.Keyspace
	createSession := func() error {
		shared.Logger.Info("db: connecting ...")
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			shared.Logger.Error("db: failed to connect", "error", err)
			return err
		}

		d.Session = session
		shared.Logger.Info("db: connected")
		return nil
	}

	if err := retry.Do(
		createSession,
		retry.Attempts(15),
		retry.Delay(6*time.Second),
	); err != nil {
		shared.Logger.Error("db: aborting ....", "error", err)
	}
}

// RunMigrations runs database migrations if any.
func (d *db) RunMigrations() {
	shared.Logger.Info("db: running migrations ...", "source", d.MigrationSourceURL)

	config := &cassandra.Config{KeyspaceName: d.Keyspace, MultiStatementEnabled: true}
	driver, err := cassandra.WithInstance(d.Session.Session, config)
	if err != nil {
		shared.Logger.Error("db: failed to initialize driver for migrations ...", "error", err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(d.MigrationSourceURL, "cassandra", driver)
	if err != nil {
		shared.Logger.Error("db: failed to initialize migrations ...", "error", err)
	}

	err = migrations.Up()

	if err == migrate.ErrNoChange {
		shared.Logger.Info("db: no migrations to run")
	}

	if err != nil && err != migrate.ErrNoChange {
		shared.Logger.Error("db: failed to run migrations ...", "error", err)
	}
	shared.Logger.Info("db: migrations done")
}

// InitSessionWithMigrations is a shorthand for initializing the database along with running migrations
func (d *db) InitSessionWithMigrations() {
	d.InitSession()
	d.RunMigrations()
}

// RegisterValidations registers any field or entity related validators
func (d *db) RegisterValidations() {
	_ = shared.Validate.RegisterValidation("db_unique", UniqueField)
}
