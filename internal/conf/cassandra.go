package conf

import (
	"github.com/avast/retry-go/v4"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	c "github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/scylladb/gocqlx/v2"
	"go.uber.org/zap"
)

var DB cassandra

type cassandra struct {
	gocqlx.Session
	Hosts              []string `env:"CASSANDRA_HOSTS" env-default:"cassandra"`
	Keyspace           string   `env:"CASSANDRA_KEYSPACE" env-default:"ctrlplane"`
	MigrationSourceURL string   `env:"CASSANDRA_MIGRATION_SOURCE_URL"`
}

func (cd *cassandra) ReadConf() {
	cleanenv.ReadEnv(cd)
}

func (cd *cassandra) InitSession() {
	Logger.Info("Initializing DB Session ...", zap.Strings("hosts", cd.Hosts), zap.String("keyspace", cd.Keyspace))
	cluster := gocql.NewCluster(cd.Hosts...)
	cluster.Keyspace = cd.Keyspace

	retryCassandra := func() error {
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return err
		}

		cd.Session = session
		Logger.Info("Initializing DB Session ... Done")
		return nil
	}

	if err := retry.Do(
		retryCassandra,
		retry.Attempts(10),
	); err != nil {
		Logger.Fatal("Failed to initialize DB Session", zap.Error(err))
	}
}

func (cd *cassandra) RunMigrations() {
	Logger.Info("Running Migrations ...", zap.String("source", cd.MigrationSourceURL))
	driver, err := c.WithInstance(cd.Session.Session, &c.Config{KeyspaceName: cd.Keyspace})
	if err != nil {
		Logger.Fatal("Failed to initialize DB Driver", zap.Error(err))
	}

	migrations, err := migrate.NewWithDatabaseInstance(DB.MigrationSourceURL, "cassandra", driver)
	if err != nil {
		Logger.Fatal("Failed to initialize DB Migrations", zap.Error(err))
	}

	err = migrations.Up()

	if err == migrate.ErrNoChange {
		Logger.Info("Running Migrations ... No Changes")
	}

	if err != nil && err != migrate.ErrNoChange {
		Logger.Fatal("Failed to run DB Migrations", zap.Error(err))
	}
	Logger.Info("Running Migrations ... Done")
}

func (cd *cassandra) InitSessionWithRunMigrations() {
	cd.InitSession()
	cd.RunMigrations()
}
