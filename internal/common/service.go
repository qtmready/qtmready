package common

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

var (
	Logger    *zap.Logger
	Service   serviceconf
	Validator *validator.Validate
)

type serviceconf struct {
	Name    string `env:"SERVICE_NAME" env-default:"service"`
	Debug   bool   `env:"DEBUG" env-default:"false"`
	Version string `env:"VERSION" env-default:"0.0.0-dev"`
	Secret  string `env:"SECRET" env-default:""`
}

func (s *serviceconf) ReadEnv() {
	cleanenv.ReadEnv(s)
}

func (s *serviceconf) InitValidator() {
	Validator = validator.New()
}

func (s *serviceconf) InitLogger() {
	if s.Debug {
		Logger, _ = zap.NewDevelopment()
	} else {
		Logger, _ = zap.NewProduction()
	}
}
