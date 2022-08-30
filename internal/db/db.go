/**
 * db provides set of utilities for working with cassandra.
 */
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
	"go.breu.io/ctrlplane/internal/cmn"
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
		Hosts              []string `env:"CASSANDRA_HOSTS" env-default:"cassandra"`
		Keyspace           string   `env:"CASSANDRA_KEYSPACE" env-default:"ctrlplane"`
		MigrationSourceURL string   `env:"CASSANDRA_MIGRATION_SOURCE_URL"`
	}
)

// Reads the environment variables
func (d *db) ReadEnv() {
	cleanenv.ReadEnv(d)
}

// Initializes the session with the configured hosts
func (d *db) InitSession() {
	cluster := gocql.NewCluster(d.Hosts...)
	cluster.Keyspace = d.Keyspace
	createSession := func() error {
		cmn.Logger.Info("db: connecting ...")
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			cmn.Logger.Error("db: failed to connect", "error", err)
			return err
		}

		d.Session = session
		cmn.Logger.Info("db: connected")
		return nil
	}

	if err := retry.Do(
		createSession,
		retry.Attempts(15),
		retry.Delay(6*time.Second),
	); err != nil {
		cmn.Logger.Error("db: aborting ....", "error", err)
	}
}

// Runs the migrations
func (d *db) RunMigrations() {
	cmn.Logger.Info("db: running migrations ...", "source", d.MigrationSourceURL)
	driver, err := cassandra.WithInstance(d.Session.Session, &cassandra.Config{KeyspaceName: d.Keyspace})
	if err != nil {
		cmn.Logger.Error("db: failed to initialize driver for migrations ...", "error", err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(d.MigrationSourceURL, "cassandra", driver)
	if err != nil {
		cmn.Logger.Error("db: failed to initialize migrations ...", "error", err)
	}

	err = migrations.Up()

	if err == migrate.ErrNoChange {
		cmn.Logger.Info("db: no migrations to run")
	}

	if err != nil && err != migrate.ErrNoChange {
		cmn.Logger.Error("db: failed to run migrations ...", "error", err)
	}
	cmn.Logger.Info("db: migrations done")
}

// Shorthand for initializing the database along with running migrations
func (d *db) InitSessionWithMigrations() {
	d.InitSession()
	d.RunMigrations()
}

// Register DB related validators
func (d *db) RegisterValidations() {
	cmn.Validate.RegisterValidation("db_unique", UniqueField)
}
