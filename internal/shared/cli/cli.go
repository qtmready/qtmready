package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Cli interface {
		GetURL() string
		GetConfigFile() string
	}

	Config struct {
		BaseURL    string `env:"BASE_URL" env-default:"http://api.ctrlplane.ai"`
		APIKEY     string `env:"API_KEY" env-default:""`
		CONFIGFILE string
	}

	ConfigOption func(*Config)
)

var ()

// NewCLI returns the global temporal instance.
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
		// }
		t.BaseURL = "http://localhost:8000"

		// set location for quantum's data file, it will contain the logged in user's access token etc

		path := ""
		op := runtime.GOOS

		switch op {
		case "windows":
			path = os.Getenv("APPDATA") + `\quantum\`
			t.CONFIGFILE = path + `access_token`
		case "darwin":
		case "linux":
			path = `~/.config/quantum/`
			t.CONFIGFILE = path + `access_token`
		default:
			fmt.Printf("%s OS is not supported by quantum yet\n", op)
		}

		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("Unable to create/locate path: %s", path)
			os.Exit(1)
		}
	}
}

func (c *Config) GetURL() string {
	return c.BaseURL
}

func (c *Config) GetConfigFile() string {
	return c.CONFIGFILE
}
