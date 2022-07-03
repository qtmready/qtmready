package conf

import (
	"github.com/gocql/gocql"
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
func InitService(name string) {
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

func InitDB() {
	cleanenv.ReadEnv(&DB)
}

func InitDBSession() {
	Logger.Info("Initializing DB Session ...", zap.Strings("hosts", DB.Hosts), zap.String("keyspace", DB.KeySpace))
	cluster := gocql.NewCluster(DB.Hosts...)
	cluster.Keyspace = DB.KeySpace
	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		Logger.Fatal(err.Error())
	}

	DB.Session = session
	Logger.Info("Initializing DB Session ... Done")
}

// Initialize Kratos (https://ory.sh)
func InitKratos() {
	cleanenv.ReadEnv(&Kratos)
}

// Initialize Github App
func InitGithub() {
	cleanenv.ReadEnv(&Github)
}

// Initialize Temporal
func InitTemporal() {
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

	client, err := tc.Dial(options)

	if err != nil {
		Logger.Fatal(err.Error())
	}

	Temporal.Client = client
	Logger.Info("Initializing Temporal Client ... Done")
}
