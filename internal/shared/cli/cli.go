// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2023, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.


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

// New returns the global temporal instance.
func New(opts ...ConfigOption) Cli {
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

		err := os.MkdirAll(configpath, 0744) // os.ModeSticky|os.ModePerm
		if err != nil {
			fmt.Printf("error: %s", err.Error())
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
