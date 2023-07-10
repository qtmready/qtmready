package cli

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Cli interface {
		GetURL() string
	}

	Config struct {
		BaseURL string `env:"BASE_URL" env-default:"http://api.ctrlplane.ai"`
		APIKEY  string `env:"API_KEY" env-default:""`
	}

	ConfigOption func(*Config)
)

var ()

// Temporal returns the global temporal instance.

func NewCLI(opts ...ConfigOption) Cli {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// FromEnvironment reads the environment variables.
func FromEnvironment() ConfigOption {
	return func(t *Config) {
		if err := cleanenv.ReadEnv(t); err != nil {
			panic(fmt.Errorf("failed to read environment variables: %w", err))
		}

		// if shared.Service().GetDebug() == true {
		t.BaseURL = "http://localhost:8000"
		// }

	}
}

func (c *Config) GetURL() string {
	return c.BaseURL
}
