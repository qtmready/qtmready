package config

import (
	"log/slog"
	"sync"
)

var (
	_c    *Config   // Global connection instance.
	_once sync.Once // Ensures connection initialization occurs only once.
)

// WithAppID sets the AppID field of the Config.
func WithAppID(id int64) ConfigOption {
	return func(config *Config) {
		config.AppID = id
	}
}

// WithClientID sets the ClientID field of the Config.
func WithClientID(id string) ConfigOption {
	return func(config *Config) {
		config.ClientID = id
	}
}

// WithWebhookSecret sets the WebhookSecret field of the Config.
func WithWebhookSecret(secret string) ConfigOption {
	return func(config *Config) {
		config.WebhookSecret = secret
	}
}

// WithPrivateKey sets the PrivateKey field of the Config.
func WithPrivateKey(key string) ConfigOption {
	return func(config *Config) {
		config.PrivateKey = key
	}
}

// WithConfig copies the values from the given Config into the target Config.
func WithConfig(cfg *Config) ConfigOption {
	return func(config *Config) {
		config.AppID = cfg.AppID
		config.ClientID = cfg.ClientID
		config.WebhookSecret = cfg.WebhookSecret
		config.PrivateKey = cfg.PrivateKey
	}
}

// Configure returns the singleton instance of the GitHub configuration.
//
// The function uses a `sync.Once` to ensure that the configuration is initialized only once. It initializes the
// instance with the `WithConfigFromEnv` option, which reads the configuration from environment variables.
func Configure(opts ...ConfigOption) *Config {
	_once.Do(func() {
		_c = &Config{}

		for _, opt := range opts {
			opt(_c)
		}
	})

	return _c
}

func Instance(opts ...ConfigOption) *Config {
	_once.Do(func() {
		slog.Warn("github: instance not initialized, this should not happen. Make sure that the configuration is loaded before calling this function.") // nolint

		_c = &Config{}

		for _, opt := range opts {
			opt(_c)
		}
	})

	return _c
}
