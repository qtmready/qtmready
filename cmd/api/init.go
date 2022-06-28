package main

import (
	_cleanenv "github.com/ilyakaznacheev/cleanenv"
	_temporalClient "go.temporal.io/sdk/client"
	_zap "go.uber.org/zap"

	"go.breu.io/ctrlplane/internal/conf"
)

func init() {
	_cleanenv.ReadEnv(&conf.Github)
	_cleanenv.ReadEnv(&conf.Kratos)
	_cleanenv.ReadEnv(&conf.Service)

	if conf.Service.Debug {
		conf.Logger, _ = _zap.NewDevelopment()
	} else {
		conf.Logger, _ = _zap.NewProduction()
	}

	client, err := _temporalClient.NewClient(_temporalClient.Options{})

	if err != nil {
		conf.Logger.Fatal(err.Error())
	}

	conf.Temporal.Client = client
}
