package service

import (
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/fatih/color"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	config struct {
		Name            string `env:"SERVICE_NAME" env-default:"service"`
		Debug           bool   `env:"DEBUG" env-default:"false"`
		Secret          string `env:"SECRET" env-default:""`
		Version         string `env:"VERSION" env-default:"dev"`
		LogSkipper      int    `env:"LOG_SKIPPER" env-default:"1"`
		CloudRunService string `env:"K_SERVICE" env-default:"unset"`
		CloudRunJob     string `env:"CLOUD_RUN_JOB" env-default:"unset"`
	}

	Service interface {
		GetName() string
		SetName(name string)
		GetVersion() string
		GetSecret() string
		GetDebug() bool
		GetLogSkipper() int
		Banner()
		GetCloudRunService() string
		GetCloudRunJob() string
	}

	ServiceOption func(Service)
)

func (s *config) GetName() string {
	return s.Name
}

func (s *config) SetName(name string) {
	s.Name = name
}

func (s *config) GetVersion() string {
	return s.Version
}

func (s *config) GetSecret() string {
	return s.Secret
}

func (s *config) GetDebug() bool {
	return s.Debug
}

func (s *config) GetLogSkipper() int {
	return s.LogSkipper
}

func (s *config) Banner() {
	banner := `
                           __          
  ____  __  ______  ____  / /_____ ___ 
 / __ \/ / / / __ \/ __ \/ __/ __ ˇ__ \
/ /_/ / /_/ / /_/ / / / / /_  / / / / /
\__, /\__▲_/\__▲_/_/ /_/\__/_/ /_/ /_/ 
  /_/ durable delivery for distributed systems.

component: %s
version: %s

%s
  `
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf(banner, green(s.Name), blue(s.Version), yellow("https://quantm.io"))
}

func (s *config) GetCloudRunService() string {
	return s.CloudRunService
}

func (s *config) GetCloudRunJob() string {
	return s.CloudRunJob
}

// WithName sets the service name.
func WithName(name string) ServiceOption {
	return func(s Service) { s.(*config).Name = name }
}

// WithDebug sets the debug flag.
func WithDebug(debug bool) ServiceOption {
	return func(s Service) { s.(*config).Debug = debug }
}

// WithSecret sets the secret. Secret is used to sign JWT and API keys.
func WithSecret(secret string) ServiceOption {
	return func(s Service) { s.(*config).Secret = secret }
}

// WithVersion sets the version.
func WithVersion(version string) ServiceOption {
	return func(s Service) { s.(*config).Version = version }
}

func WithLogSkipper(skipper int) ServiceOption {
	return func(s Service) { s.(*config).LogSkipper = skipper }
}

// WithVersionFromBuildInfo sets the version from the build info.
func WithVersionFromBuildInfo() ServiceOption {
	return func(s Service) {
		if info, ok := debug.ReadBuildInfo(); ok {
			var (
				revision  string
				modified  string
				timestamp time.Time
				version   string
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

			if len(revision) > 0 && len(modified) > 0 && timestamp.Unix() > 0 {
				version = timestamp.Format("2006.01.02") + "." + revision[:7]
			} else {
				version = "debug"
			}

			if modified == "true" {
				version += "-dev"
			}

			s.(*config).Version = version
		}
	}
}

// FromEnvironment reads the environment variables and sets the config.
func FromEnvironment() ServiceOption {
	return func(s Service) {
		if err := cleanenv.ReadEnv(s.(*config)); err != nil {
			panic(fmt.Errorf("failed to read environment variables: %w", err))
		}
	}
}

// FromFile reads the config from the given path.
func FromFile(path string) ServiceOption {
	return func(s Service) {
		if err := cleanenv.ReadConfig(path, s.(*config)); err != nil {
			panic(fmt.Errorf("failed to read config: %w", err))
		}
	}
}

// WithConfigFromDefault reads the config from the default path.
func DefaultConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home dir: %w", err))
	}

	return path.Join(home, ".ctrlplane", "config.json")
}

// New creates a new instance of the service.
func New(opts ...ServiceOption) Service {
	s := &config{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}
