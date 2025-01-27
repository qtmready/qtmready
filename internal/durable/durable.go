package durable

import (
	"log/slog"
	"sync"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/durable/config"
)

// -- Internal --

var (
	// configured is the instantiated configuration.
	configured *Config
	configonce sync.Once

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

	WithConfig = config.WithConfig
)

// Get is a singleton that holds the temporal client.
//
// Please note that the actual client is not created until the first call Get().Client().
func Get(opts ...ConfigOption) *Config {
	configonce.Do(func() {
		configured = config.New(opts...)
	})

	return configured
}

// -- Queues --

// OnCore returns the core queue.
//
// All workflows on this queue will have the ID prefix of
//
//	io.ctrlpane.core.{block}.{block_id}.{element}.{element_id}.{modifier}.{modifier_id}....
func OnCore() queues.Queue {
	coreqonce.Do(func() {
		client, err := Get().Client()
		if err != nil {
			slog.Error("durable: unable to connect to durable server ...", "error", err.Error())
			panic(err)
		}

		coreq = queues.New(queues.WithName("core"), queues.WithClient(client))
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
		client, err := Get().Client()
		if err != nil {
			slog.Error("durable: unable to connect to durable server ...", "error", err.Error())
			panic(err)
		}

		hooksq = queues.New(queues.WithName("hooks"), queues.WithClient(client))
	})

	return hooksq
}
