package defaults

import (
	_cleanenv "github.com/ilyakaznacheev/cleanenv"
	_sdkClient "go.temporal.io/sdk/client"
	_zap "go.uber.org/zap"
)

var Conf conf
var Logger *_zap.Logger

func init() {
	_cleanenv.ReadEnv(&Conf)

	if Conf.Debug {
		Logger, _ = _zap.NewDevelopment()
	} else {
		Logger, _ = _zap.NewProduction()
	}

	// TODO: ysf - handle this for production
	client, err := _sdkClient.Dial(_sdkClient.Options{})

	if err != nil {
		Logger.Fatal(err.Error())
	}

	Conf.Temporal.Client = client
}
