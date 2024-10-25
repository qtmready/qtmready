package db

import (
	"go.breu.io/quantm/internal/db/config"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	Config = config.Connection
)

var (
	DefaultConfig = config.DefaultConnection
)

func WithConfig(conf *Config) config.Option {
	return config.WithConfig(conf)
}

// Connection is a wrapper around the config.Instance singleton.
func Connection(opts ...config.Option) *config.Connection {
	return config.Instance(opts...)
}

// Queries is a wrapper around the config.Queries singleton.
func Queries() *entities.Queries {
	return config.Queries()
}
