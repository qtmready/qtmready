package conf

import (
	_cleanenv "github.com/ilyakaznacheev/cleanenv"
	_tmprlclient "go.temporal.io/sdk/client"
	_zap "go.uber.org/zap"
)

// Initialize the service
func InitService(name string) {
	_cleanenv.ReadEnv(&Service)

	if Service.Name == "" {
		Service.Name = name
	}

	if Service.Debug {
		Logger, _ = _zap.NewDevelopment()
	} else {
		Logger, _ = _zap.NewProduction()
	}
}

// Initialize Kratos (https://ory.sh)
func InitKratos() {
	_cleanenv.ReadEnv(&Kratos)
}

// Initialize Github App
func InitGithub() {
	_cleanenv.ReadEnv(&Github)
}

// Initialize Temporal
func InitTemporal() {
	_cleanenv.ReadEnv(&Temporal)
}

// Initalize Temporal Client.
//
// Must be called after `InitService()` & `InitTemporal()`
//
// Must do `defer conf.TemporalClient.Close()` after calling `conf.InitTemporalClient()`
func InitTemporalClient() {
	client, err := _tmprlclient.Dial(_tmprlclient.Options{})

	if err != nil {
		Logger.Fatal(err.Error())
	}

	Temporal.Client = client
}
