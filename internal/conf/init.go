package conf

import (
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	c "github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/scylladb/gocqlx/v2"
	tc "go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

var DB cassandra
var Github github
var Kratos kratos
var Service service
var Logger *zap.Logger
var Temporal temporal

// Initialize the service
func ReadSvcConfig(name string) {
	cleanenv.ReadEnv(&Service)

	if Service.Name == "" {
		Service.Name = name
	}

	if Service.Debug {
		Logger, _ = zap.NewDevelopment()
	} else {
		Logger, _ = zap.NewProduction()
	}

	Logger.Info("Initializing Service ...", zap.String("name", Service.Name), zap.String("version", Service.Version))
}

func ReadDBConfig() {
	cleanenv.ReadEnv(&DB)
}

func InitDBSession() {
	Logger.Info("Initializing DB Session ...", zap.Strings("hosts", DB.Hosts), zap.String("keyspace", DB.KeySpace))
	cluster := gocql.NewCluster(DB.Hosts...)
	cluster.Keyspace = DB.KeySpace

	retryCassandra := func() error {
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			return err
		}

		DB.Session = session
		Logger.Info("Initializing DB Session ... Done")
		return nil
	}

	if err := retry.Do(
		retryCassandra,
		retry.Attempts(10),
		retry.Delay(1*time.Second),
	); err != nil {
		Logger.Fatal("Failed to initialize DB Session", zap.Error(err))
	}

	if DB.MigrationFilesPath != "" {
		Logger.Info("Running DB Migration ...", zap.String("path", DB.MigrationFilesPath))

		config := &c.Config{
			KeyspaceName: DB.KeySpace,
		}
		driver, err := c.WithInstance(DB.Session.Session, config)
		if err != nil {
			Logger.Fatal("Failed to initialize DB Migration", zap.Error(err))
		}

		m, err := migrate.NewWithDatabaseInstance("file://"+DB.MigrationFilesPath, "cassandra", driver)
		if err != nil {
			Logger.Fatal("Failed to initialize DB Migration", zap.Error(err))
		}
		m.Up()

		Logger.Info("Running DB Migration ... Done")
	}
}

// Initialize Kratos (https://ory.sh)
func ReadKratosConfig() {
	cleanenv.ReadEnv(&Kratos)
}

// Initialize Github App
func ReadGithubConfig() {
	cleanenv.ReadEnv(&Github)
}

// Initialize Temporal
func ReadTemporalConfig() {
	cleanenv.ReadEnv(&Temporal)
}

// Initalize Temporal Client.
//
// Must be called after `InitService()` & `InitTemporal()`
//
// Must do `defer conf.TemporalClient.Close()` after calling `conf.InitTemporalClient()`
func InitTemporalClient() {
	Logger.Info("Initializing Temporal Client ...", zap.String("host", Temporal.ServerHost), zap.String("port", Temporal.ServerPort))
	options := tc.Options{
		HostPort: Temporal.GetConnectionString(),
	}

	retryTemporal := func() error {
		client, err := tc.Dial(options)
		if err != nil {
			return err
		}

		Temporal.Client = client
		Logger.Info("Initializing Temporal Client ... Done")
		return nil
	}

	if err := retry.Do(
		retryTemporal,
		retry.Attempts(10),
		retry.Delay(1*time.Second),
	); err != nil {
		Logger.Fatal("Failed to initialize Temporal Client", zap.Error(err))
	}
}
