package server

import (
	"go.breu.io/quantm/internal/nomad/handler"
)

// DefaultServer creates a new Nomad server instance with the provided options.
func DefaultServer(opts ...Option) *Server {
	srv := New(opts...)
	srv.add(handler.NewHealthCheckServiceHandler())

	return srv
}
