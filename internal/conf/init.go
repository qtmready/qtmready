package conf

import (
	"github.com/ilyakaznacheev/cleanenv"
	tclient "go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

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
	Logger.Info("Starting ...", zap.String("name", Service.Name), zap.String("version", Service.Version))
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
	Logger.Info("Initializing Temporal Client", zap.String("host", Temporal.ServerHost), zap.String("port", Temporal.ServerPort))
	options := tclient.Options{
		HostPort: Temporal.GetConnectionString(),
	}

	client, err := tclient.Dial(options)

	if err != nil {
		Logger.Fatal(err.Error())
	}

	Temporal.Client = client
}
