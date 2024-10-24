package server

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type (
	// Config represents the Nomad server configuration.
	Config struct {
		Port      int  `json:"port" koanf:"PORT"`             // Server port.
		EnableSSL bool `json:"enable_ssl" koanf:"ENABLE_SSL"` // Enables TLS/SSL.
	}

	ConfigOption func(*Config) // ConfigOption is a function that modifies the Config.
)

var (
	// DefaultConfig contains the default configuration values.
	DefaultConfig = Config{
		Port:      7070,  // Default port.
		EnableSSL: false, // SSL disabled by default.
	}
)

// Address returns the formatted address for the Nomad server.
func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.Port) // Formatted address string.
}

// WithPortConfig returns a ConfigOption that sets the server port.
func WithPortConfig(port int) ConfigOption {
	return func(c *Config) {
		c.Port = port // Set the server port.
	}
}

// WithSSLConfig returns a ConfigOption that enables or disables TLS/SSL.
func WithSSLConfig(enableSSL bool) ConfigOption {
	return func(c *Config) {
		c.EnableSSL = enableSSL // Enable or disable SSL.
	}
}

// WithEnvironmentConfig returns a ConfigOption that loads configuration from environment variables.
//
// It reads environment variables prefixed with the specified prefix, or "NOMAD__" if no prefix is provided.
// The environment variable names are mapped to the corresponding struct fields using `koanf` and `structs`.
//
// For example, with the prefix "APP__", the environment variable "APP__PORT" will be used to set the `Port` field.
func WithEnvironmentConfig(opts ...string) ConfigOption {
	return func(c *Config) {
		var prefix string

		if len(opts) > 0 {
			prefix = strings.ToUpper(opts[0]) // Prefix for environment variables.

			if !strings.HasSuffix(prefix, "__") {
				prefix += "__" // Ensure prefix ends with "__".
			}
		} else {
			prefix = "NOMAD__" // Default prefix.
		}

		k := koanf.New("__")                                   // Create a new `koanf` instance.
		_ = k.Load(structs.Provider(DefaultConfig, "__"), nil) // Load default values from the struct.

		if err := k.Load(env.Provider(prefix, "__", nil), nil); err != nil {
			panic(err) // Panic if an error occurs while loading environment variables.
		}

		if err := k.Unmarshal("", k); err != nil {
			panic(err) // Panic if an error occurs while unmarshaling values.
		}
	}
}

// NewConfig creates a new Config instance with the specified options.
func NewConfig(opts ...ConfigOption) *Config {
	c := &Config{} // Initialize the Config.

	for _, opt := range opts {
		opt(c) // Apply each option to the Config.
	}

	return c // Return the configured Config.
}
