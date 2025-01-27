package pulse

import (
	"log/slog"
	"sync"

	"go.breu.io/quantm/internal/pulse/config"
)

type (
	Config = config.Config // Config represents the configuration for the Pulse package.
	Option = config.Option // Option is a functional option to configure Pulse.
)

var (
	DefaultConfig = config.DefaultConfig // DefaultConfig holds the default configuration values.

	_c   *Config   // _c stores the configured instance of Pulse.
	once sync.Once // once ensures the initialization happens only once.
)

// WithConfig allows customizing the Pulse configuration using functional options.
//
// It takes a Config pointer as input and returns an Option. This Option can then be passed to the Instance function.
func WithConfig(cfg *Config) Option {
	return config.WithConfig(cfg)
}

// Get returns the singleton instance of the Pulse configuration.
//
// It initializes Pulse with the provided options if it hasn't been initialized yet. The initialization is thread-safe,
// guaranteed by the sync.Once usage.  Get returns a pointer to the initialized Config instance.
func Get(opts ...Option) *Config {
	once.Do(func() {
		slog.Info("pulse: configuring ...")

		_c = config.New(opts...)
	})

	return _c
}
