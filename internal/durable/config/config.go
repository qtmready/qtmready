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

package config

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

type (
	// Config represents the Temporal configuration.
	Config struct {
		Namespace string `json:"namespace" koanf:"NAMESPACE"` // Temporal namespace.
		Host      string `json:"host" koanf:"HOST"`           // Temporal host.
		Port      int    `json:"port" koanf:"PORT"`           // Temporal port.
		Skip      int    `json:"skip" koanf:"LOG_SKIP"`       // Skip frames for logging.

		client client.Client // Temporal client.
		once   *sync.Once    // We can have only one Temporal client per configuration.
	}

	Option func(*Config) // ConfigOption is a function that modifies the Config.
)

var (
	// Default contains the default configuration values.
	Default = Config{
		Namespace: "default",
		Host:      "localhost",
		Port:      7233,
		Skip:      0,

		once: &sync.Once{},
	}
)

// Start is a no-op function to satisfy the graceful.Service interface.
func (c *Config) Start(ctx context.Context) error { return nil }

// Stop is a no-op function to satisfy the graceful.Service interface.
func (c *Config) Stop(ctx context.Context) error { return nil }

// Address returns the formatted address for the Temporal server.
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port) // Formatted address string.
}

// Client returns the Temporal client.
func (c *Config) Client() (client.Client, error) {
	var err error

	c.once.Do(func() {
		slog.Info("durable: establishing temporal connection ...", "host", c.Host, "port", c.Port)

		err = retry.Do(
			c.dial,
			retry.Attempts(10),
			retry.Delay(1*time.Second),
			retry.OnRetry(func(attempt uint, err error) {
				slog.Warn(
					"durable: unable to establish temporal conntection, retrying ...",
					"host", c.Host, "port", c.Port,
					"attempt", attempt, "remaining", 10-attempt,
					"error", err.Error(),
				)
			}),
		)

		if err != nil {
			slog.Error("durable: unable to establish temporal connection.", "host", c.Host, "port", c.Port, "error", err.Error())

			return
		}

		slog.Info("durable: temporal connection established", "host", c.Host, "port", c.Port)
	})

	return c.client, err
}

func (c *Config) dial() error {
	_c, err := client.Dial(c.options())
	if err != nil {
		return err
	}

	c.client = _c

	return nil
}

func (c *Config) options() client.Options {
	return client.Options{
		HostPort:  c.Address(),
		Namespace: c.Namespace,
		Logger:    log.Skip(log.NewStructuredLogger(slog.Default()), c.Skip),
	}
}

// WithNamespaceConfig returns a ConfigOption that sets the Temporal namespace.
func WithNamespaceConfig(namespace string) Option {
	return func(c *Config) {
		c.Namespace = namespace // Set the Temporal namespace.
	}
}

// WithHostConfig returns a ConfigOption that sets the Temporal host.
func WithHostConfig(host string) Option {
	return func(c *Config) {
		c.Host = host // Set the Temporal host.
	}
}

// WithPortConfig returns a ConfigOption that sets the Temporal port.
func WithPortConfig(port int) Option {
	return func(c *Config) {
		c.Port = port // Set the Temporal port.
	}
}

// WithEnvironmentConfig returns a ConfigOption that loads configuration from environment variables.
//
// It reads environment variables prefixed with the specified prefix, or "TEMPORAL__" if no prefix is provided.
// The environment variable names are mapped to the corresponding struct fields using `koanf` and `structs`.
//
// For example, with the prefix "APP__", the environment variable "APP__NAMESPACE" will be used to set the `Namespace`
// field.
func WithEnvironmentConfig(opts ...string) Option {
	return func(c *Config) {
		var prefix string

		if len(opts) > 0 {
			prefix = strings.ToUpper(opts[0]) // Prefix for environment variables.

			if !strings.HasSuffix(prefix, "__") {
				prefix += "__" // Ensure prefix ends with "__".
			}
		} else {
			prefix = "TEMPORAL__" // Default prefix.
		}

		k := koanf.New("__")                             // Create a new `koanf` instance.
		_ = k.Load(structs.Provider(Default, "__"), nil) // Load default values from the struct.

		if err := k.Load(env.Provider(prefix, "__", nil), nil); err != nil {
			panic(err) // Panic if an error occurs while loading environment variables.
		}

		if err := k.Unmarshal("", k); err != nil {
			panic(err) // Panic if an error occurs while unmarshaling values.
		}
	}
}

func WithConfig(conf *Config) Option {
	return func(c *Config) {
		c.Namespace = conf.Namespace
		c.Host = conf.Host
		c.Port = conf.Port
		c.Skip = conf.Skip
	}
}

// New creates a new Config instance with the specified options.
func New(opts ...Option) *Config {
	c := &Config{once: &sync.Once{}} // Initialize the Config.

	for _, opt := range opts {
		opt(c) // Apply each option to the Config.
	}

	return c // Return the configured Config.
}
