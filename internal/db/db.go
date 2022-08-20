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
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v2"
	"go.breu.io/ctrlplane/internal/cmn"
	"go.uber.org/zap"
)

var (
	DB         = &db{}
	NullUUID   = "00000000-0000-0000-0000-000000000000"
	NullString = ""
)

type (
	// Defines the query params required for DB lookup queries
	QueryParams map[string]interface{}

	// An Entity defines the interface for a database entity
	Entity interface {
		GetTable() *table.Table
		PreCreate() error
		PreUpdate() error
	}

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
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return err
		}

		d.Session = session
		return nil
	}

	if err := retry.Do(
		createSession,
		retry.Attempts(10),
		retry.Delay(6*time.Second),
	); err != nil {
		cmn.Log.Fatal("Failed to initialize Cassandra Session", zap.Error(err))
	}
}

// Runs the migrations
func (d *db) RunMigrations() {
	cmn.Log.Info("Running Migrations ...", zap.String("source", d.MigrationSourceURL))
	driver, err := cassandra.WithInstance(d.Session.Session, &cassandra.Config{KeyspaceName: d.Keyspace})
	if err != nil {
		cmn.Log.Fatal("Failed to initialize DB Driver", zap.Error(err))
	}

	migrations, err := migrate.NewWithDatabaseInstance(d.MigrationSourceURL, "cassandra", driver)
	if err != nil {
		cmn.Log.Fatal("Failed to initialize DB Migrations", zap.Error(err))
	}

	err = migrations.Up()

	if err == migrate.ErrNoChange {
		cmn.Log.Info("Running Migrations ... No Changes")
	}

	if err != nil && err != migrate.ErrNoChange {
		cmn.Log.Fatal("Failed to run DB Migrations", zap.Error(err))
	}
	cmn.Log.Info("Running Migrations ... Done")
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
