package nomad

import (
	"go.breu.io/quantm/internal/nomad/server"
)

type (
	Config = server.Config
)

var (
	DefaultConfig = server.DefaultConfig
)

// WithConfig sets the server configuration.
func WithConfig(config *Config) server.Option {
	return server.WithConfig(config)
}

// New creates a new Nomad server instance with the provided options.
func New(opts ...server.Option) *server.Server {
	return server.DefaultServer(opts...)
}
