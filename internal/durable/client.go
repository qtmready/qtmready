package durable

import (
	"log/slog"
	"sync"

	sdk "go.temporal.io/sdk/client"

	"go.breu.io/quantm/internal/durable/config"
)

type (
	Config       = config.Config
	ConfigOption = config.Option
)

var (
	DefaultConfig = config.Default

	configured *Config
	configonce sync.Once

	client sdk.Client
)

func WithConfig(conf *Config) ConfigOption {
	return config.WithConfig(conf)
}

// Configure returns the configured Temporal client.
func Configure(opts ...ConfigOption) error {
	var err error

	configonce.Do(func() {
		configured = config.New(opts...)
		client, err = configured.Client()
	})

	return err
}

// Client returns the configured Temporal client. If not yet initialized, it will use the default
// configuration. If initialization fails, the program will panic. For predictable behavior,
// initialize the client using Configure prior to usage, typically within the main function.
func Client() sdk.Client {
	if client == nil {
		slog.Warn("durable: client not configured, configuring using default configuration")

		if err := Configure(); err != nil {
			panic(err)
		}
	}

	return client
}
