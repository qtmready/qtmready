package cli

import (
	"fmt"
	"os"
	"path"
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

		// set location for quantm's data file, it will contain the logged in user's access token etc

		configpath := ""
		hostos := runtime.GOOS

		switch hostos {
		case "windows":
			configpath = path.Join(os.Getenv("APPDATA"), "quantm")
		case "darwin", "linux":
			home, err := os.UserHomeDir()
			if err != nil {
				os.Exit(1)
			}

			configpath = path.Join(home, ".config", "quantm")
		default:
			fmt.Printf("%s OS is not supported by quantm yet\n", hostos)
		}

		t.CONFIGFILE = path.Join(configpath, "access_token")

		err := os.MkdirAll(configpath, os.ModeDir)
		if err != nil {
			fmt.Printf(err.Error())
			fmt.Printf("Unable to create/locate path: %s", configpath)
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
