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
		Port      int  `json:"port" koanf:"PORT"`
		EnableSSL bool `json:"enable_ssl" koanf:"ENABLE_SSL"`
	}

	ConfigOption func(*Config)
)

var (
	DefaultConfig = Config{
		Port:      7070,
		EnableSSL: false,
	}
)

func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

func WithPortConfig(port int) ConfigOption {
	return func(c *Config) {
		c.Port = port
	}
}

func WithSSLConfig(enableSSL bool) ConfigOption {
	return func(c *Config) {
		c.EnableSSL = enableSSL
	}
}

func WithEnvironmentConfig(opts ...string) ConfigOption {
	return func(c *Config) {
		var prefix string

		if len(opts) > 0 {
			prefix = strings.ToUpper(opts[0])

			if !strings.HasSuffix(prefix, "__") {
				prefix += "__"
			}
		} else {
			prefix = "NOMAD__"
		}

		k := koanf.New("__")
		_ = k.Load(structs.Provider(DefaultConfig, "__"), nil)

		if err := k.Load(env.Provider(prefix, "__", nil), nil); err != nil {
			panic(err)
		}

		if err := k.Unmarshal("", k); err != nil {
			panic(err)
		}
	}
}

func NewConfig(opts ...ConfigOption) *Config {
	c := &Config{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
