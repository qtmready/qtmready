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
		CLI     *cli   `env-prefix:"CTRLPLANE_" env-allow-empty:"true"`
		version string `env:"VERSION" env-default:""`
	}

	cli struct {
		BaseURL string `env:"BASE_URL" env-default:"http://api.ctrlplane.ai"`
		APIKEY  string `env:"API_KEY" env-default:""`
	}
)

var (
	Logger   *logger.ZapAdapter
	Service  = &service{}
	Validate *validator.Validate
)

// Version creates the version string as per [calver].
//
// The scheme currently being followed is YYYY.0M.0D.<git commit hash>-<channel> where:
//   - YYYY.0M.0D is the date of the commit
//   - <git commit hash> is the first 8 characters of the git commit hash
//   - <channel> is the channel of the build (e.g. dev, alpha, beta, rc, stable).
//
// For out purposes, -<channel> is optional and will be set to "dev" if the git is dirty.
//
// [calver]: https://calver.org/
func (s *service) Version() string {
	if s.version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			var (
				revision  string
				modified  string
				timestamp time.Time
			)

			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					revision = setting.Value
				}

				if setting.Key == "vcs.modified" {
					modified = setting.Value
				}

				if setting.Key == "vcs.time" {
					timestamp, _ = time.Parse(time.RFC3339, setting.Value)
				}
			}

			version := timestamp.Format("2006.01.02") + "." + revision[:8]

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
	conf, err := s.GetConfigPath()

	if err != nil {
		return err
	}

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

func (s *service) InitCLI() {
	s.Name = "ctrlplane-cli"
	s.Debug = false // FIXME: this should be set to true if the CLI is run in debug mode

	if err := s.ReadFile(); err != nil {
		panic(err)
	}
}
