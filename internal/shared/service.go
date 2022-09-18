// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.  

package shared

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"

	"go.breu.io/ctrlplane/internal/shared/logger"
)

type (
	service struct {
		Name    string `env:"SERVICE_NAME" env-default:"service"`
		Debug   bool   `env:"DEBUG" env-default:"false"`
		Version string `env:"VERSION" env-default:"0.0.0-dev"`
		Secret  string `env:"SECRET" env-default:""`
	}
)

var (
	Logger   *logger.ZapAdapter
	Service  = &service{}
	Validate *validator.Validate
)

// ReadEnv reads the environment variables and initializes the service.
func (s *service) ReadEnv() {
	if err := cleanenv.ReadEnv(s); err != nil {
		Logger.Error("Failed to read environment variables", "error", err)
	}
}

// InitValidator sets up global validator.
func (s *service) InitValidator() {
	Validate = validator.New()
	// by default, the validator will try to get json tag.
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// InitLogger sets up global logger.
func (s *service) InitLogger() {
	var zl *zap.Logger

	if s.Debug {
		zl, _ = zap.NewDevelopment()
	} else {
		zl, _ = zap.NewProduction()
	}

	Logger = logger.NewZapAdapter(zl)
}
