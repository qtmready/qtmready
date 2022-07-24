package common

import (
	"reflect"
	"strings"

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

// Reads the environment variables and initializes the service.
func (s *serviceconf) ReadEnv() {
	cleanenv.ReadEnv(s)
}

// Sets up global validator.
func (s *serviceconf) InitValidator() {
	Validator = validator.New()
	// by default, the validator will try to get json tag.
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Sets up global logger.
func (s *serviceconf) InitLogger() {
	if s.Debug {
		Logger, _ = zap.NewDevelopment()
	} else {
		Logger, _ = zap.NewProduction()
	}
}
