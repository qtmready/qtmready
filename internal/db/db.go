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
	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db/validations"
	"go.uber.org/zap"
)

var DB db

type db struct {
	gocqlx.Session
	Hosts              []string `env:"CASSANDRA_HOSTS" env-default:"cassandra"`
	Keyspace           string   `env:"CASSANDRA_KEYSPACE" env-default:"ctrlplane"`
	MigrationSourceURL string   `env:"CASSANDRA_MIGRATION_SOURCE_URL"`
}

func (d *db) ReadEnv() {
	cleanenv.ReadEnv(d)
}

func (d *db) InitSession() {
	cluster := gocql.NewCluster(d.Hosts...)
	cluster.Keyspace = d.Keyspace

	retryCassandra := func() error {
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return err
		}

		d.Session = session
		return nil
	}

	if err := retry.Do(
		retryCassandra,
		retry.Attempts(10),
		retry.Delay(6*time.Second),
	); err != nil {
		common.Logger.Fatal("Failed to initialize Cassandra Session", zap.Error(err))
	}
}

func (d *db) RunMigrations() {
	common.Logger.Info("Running Migrations ...", zap.String("source", d.MigrationSourceURL))
	driver, err := cassandra.WithInstance(d.Session.Session, &cassandra.Config{KeyspaceName: d.Keyspace})
	if err != nil {
		common.Logger.Fatal("Failed to initialize DB Driver", zap.Error(err))
	}

	migrations, err := migrate.NewWithDatabaseInstance(d.MigrationSourceURL, "cassandra", driver)
	if err != nil {
		common.Logger.Fatal("Failed to initialize DB Migrations", zap.Error(err))
	}

	err = migrations.Up()

	if err == migrate.ErrNoChange {
		common.Logger.Info("Running Migrations ... No Changes")
	}

	if err != nil && err != migrate.ErrNoChange {
		common.Logger.Fatal("Failed to run DB Migrations", zap.Error(err))
	}
	common.Logger.Info("Running Migrations ... Done")
}

func (d *db) InitSessionWithMigrations() {
	d.InitSession()
	d.RunMigrations()
}

func (d *db) RegisterValidations() {
	common.Validator.RegisterValidation("db_unique", validations.Unique)
}
