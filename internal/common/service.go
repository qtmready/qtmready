package common

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type serviceconf struct {
	Name    string `env:"SERVICE_NAME" env-default:"service"`
	Debug   bool   `env:"DEBUG" env-default:"false"`
	Version string `env:"VERSION" env-default:"0.0.0-dev"`
}

func (s *serviceconf) ReadConf() {
	cleanenv.ReadEnv(s)
}

func (s *serviceconf) InitLogger() {
	if s.Debug {
		Logger, _ = zap.NewDevelopment()
	} else {
		Logger, _ = zap.NewProduction()
	}

	Logger.Info("Initializing Service ...", zap.String("name", s.Name), zap.String("version", s.Version))
}
