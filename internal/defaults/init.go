package defaults

import (
	_cleanenv "github.com/ilyakaznacheev/cleanenv"
	_tclient "go.temporal.io/sdk/client"
	_zap "go.uber.org/zap"
)

var Conf conf
var Logger *_zap.Logger

var TemporalClient _tclient.Client

var err error = nil

func init() {
	_cleanenv.ReadEnv(&Conf)

	if Conf.Debug {
		Logger, _ = _zap.NewDevelopment()
	} else {
		Logger, _ = _zap.NewProduction()
	}

	// TODO: ysf - handle this for production
	TemporalClient, err = _tclient.Dial(_tclient.Options{})

	if err != nil {
		Logger.Fatal(err.Error())
	}

	defer TemporalClient.Close()
}
