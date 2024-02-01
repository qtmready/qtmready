// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

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
