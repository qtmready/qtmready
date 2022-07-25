package common

import (
	"reflect"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

var (
	Logger    *zap.Logger
	Service   serviceconf
	Validator *validator.Validate
	JWT       *jwtauth.JWTAuth
)

type serviceconf struct {
	Name    string `env:"SERVICE_NAME" env-default:"service"`
	Debug   bool   `env:"DEBUG" env-default:"false"`
	Version string `env:"VERSION" env-default:"0.0.0-dev"`
	Secret  string `env:"SECRET" env-default:""`
}

// Reads the environment variables and initializes the service.
func (s *serviceconf) ReadEnv() {
	if err := cleanenv.ReadEnv(s); err != nil {
		Logger.Fatal("Failed to read environment variables", zap.Error(err))
	}
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

func (s *serviceconf) InitJWT() {
	JWT = jwtauth.New("HS256", []byte(s.Secret), nil)
}

// Sets up global logger.
func (s *serviceconf) InitLogger() {
	if s.Debug {
		Logger, _ = zap.NewDevelopment()
	} else {
		Logger, _ = zap.NewProduction()
	}
}
