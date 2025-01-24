package durable

import (
	"log/slog"
	"sync"

	"go.breu.io/durex/queues"
	sdk "go.temporal.io/sdk/client"

	"go.breu.io/quantm/internal/durable/config"
)

// -- Internal --

var (
	// configured is the instantiated configuration.
	configured *Config
	configonce sync.Once

	// client is the configured Temporal client.
	client sdk.Client

	// coreq is the core queue.
	coreq     queues.Queue
	coreqonce sync.Once

	// hooksq is the hooks queue.
	hooksq     queues.Queue
	hooksqonce sync.Once
)

// -- Types --

type (
	// Config represents the configuration for the durable layer.
	Config = config.Config

	// ConfigOption is an option for configuring the durable layer.
	ConfigOption = config.Option
)

// -- Configuration --

var (
	// DefaultConfig is the default configuration for the durable layer.
	DefaultConfig = config.Default
)

// Configure initializes the durable layer and instantiates the Temporal client.
//
// Configuration is applied once; subsequent calls are no-ops. An error is returned if Temporal client
// initialization fails.
func Configure(opts ...ConfigOption) error {
	var err error

	configonce.Do(func() {
		configured = config.New(opts...)
		client, err = configured.Client()
	})

	return err
}

func Instance(opts ...ConfigOption) *Config {
	if configured == nil {
		slog.Warn("durable: instance not configured, configuring using default configuration")

		if err := Configure(opts...); err != nil {
			panic(err)
		}
	}

	return configured
}

// Client returns the configured Temporal client.
//
// If the client is not yet initialized, it will be initialized using the default configuration.
// If initialization fails, the program will panic.
//
// For predictable behavior, initialize the client using Configure prior to usage, typically within the main function.
func Client() sdk.Client {
	if client == nil {
		slog.Warn("durable: client not configured, configuring using default configuration")

		if err := Configure(); err != nil {
			panic(err)
		}
	}

	return client
}

// -- Queues --

// OnCore returns the core queue.
//
// All workflows on this queue will have the ID prefix of
//
//	io.ctrlpane.core.{block}.{block_id}.{element}.{element_id}.{modifier}.{modifier_id}....
func OnCore() queues.Queue {
	coreqonce.Do(func() {
		coreq = queues.New(queues.WithName("core"), queues.WithClient(Client()))
	})

	return coreq
}

// OnHooks returns the hooks queue.
//
// All workflows on this queue will have the ID prefix of
//
//	io.ctrlpane.hooks.{block}.{block_id}.{element}.{element_id}.{modifier}.{modifier_id}....
func OnHooks() queues.Queue {
	hooksqonce.Do(func() {
		hooksq = queues.New(queues.WithName("hooks"), queues.WithClient(Client()))
	})

	return hooksq
}

// -- Helpers --

// WithConfig returns a ConfigOption that sets the configuration.
func WithConfig(conf *Config) ConfigOption {
	return config.WithConfig(conf)
}
