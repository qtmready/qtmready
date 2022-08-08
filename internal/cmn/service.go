package cmn

import (
	"reflect"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

var (
	Log       *zap.Logger
	Service   = &service{}
	Validator *validator.Validate
	JWT       *jwtauth.JWTAuth
)

type service struct {
	Name    string `env:"SERVICE_NAME" env-default:"service"`
	Debug   bool   `env:"DEBUG" env-default:"false"`
	Version string `env:"VERSION" env-default:"0.0.0-dev"`
	Secret  string `env:"SECRET" env-default:""`
}

// ReadEnv reads the environment variables and initializes the service.
func (s *service) ReadEnv() {
	if err := cleanenv.ReadEnv(s); err != nil {
		Log.Fatal("Failed to read environment variables", zap.Error(err))
	}
}

// InitValidator sets up global validator.
func (s *service) InitValidator() {
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

// InitJWT sets up global JWT.
func (s *service) InitJWT() {
	JWT = jwtauth.New("HS256", []byte(s.Secret), nil)
}

// InitLogger sets up global logger.
func (s *service) InitLogger() {
	if s.Debug {
		Log, _ = zap.NewDevelopment()
	} else {
		Log, _ = zap.NewProduction()
	}
}
