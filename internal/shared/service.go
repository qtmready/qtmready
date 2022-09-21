// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package shared

import (
	"os"
	"path"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"

	"go.breu.io/ctrlplane/internal/shared/logger"
)

type (
	service struct {
		Name    string `env:"SERVICE_NAME" env-default:"service"`
		Debug   bool   `env:"DEBUG" env-default:"false"`
		Secret  string `env:"SECRET" env-default:""`
		CLI     cli    `env-prefix:"CLI_" env-allow-empty:"true"`
		version string `env:"VERSION" env-default:""`
	}

	cli struct {
		BaseUrl      string `env:"BASE_URL" env-default:"http://localhost:8000"`
		AccessToken  string `env:"ACCESS_TOKEN" env-default:""`
		RefreshToken string `env:"REFRESH_TOKEN" env-default:""`
	}
)

var (
	Logger   *logger.ZapAdapter
	Service  = &service{}
	Validate *validator.Validate
)

func (s *service) Version() string {
	if s.version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			var revision string
			var modified string
			var timestamp time.Time
			for _, s := range info.Settings {
				if s.Key == "vcs.revision" {
					revision = s.Value
				}

				if s.Key == "vcs.modified" {
					modified = s.Value
				}

				if s.Key == "vcs.time" {
					timestamp, _ = time.Parse(time.RFC3339, s.Value)
				}
			}

			version := timestamp.Format("060102") + "." + revision[:8]
			if modified == "true" {
				version += "-dev"
			}
			s.version = version
		}
	}

	return s.version
}

// ReadEnv reads the environment variables and initializes the service.
func (s *service) ReadEnv() {
	if err := cleanenv.ReadEnv(s); err != nil {
		panic("Failed to read environment variables")
	}
}

func (s *service) GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".ctrlplane", "config.json"), nil
}

func (s *service) ReadFile() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	conf := path.Join(home, ".ctrlplane", "config.json")
	if err := cleanenv.ReadConfig(conf, s); err != nil {
		return err
	}

	return nil
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
