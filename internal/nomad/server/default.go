package server

import (
	"go.breu.io/quantm/internal/nomad/handlers"
)

// DefaultServer creates a new Nomad server instance with the provided options.
func DefaultServer(opts ...Option) *Server {
	srv := New(opts...)

	srv.add(handlers.NewHealthCheckServiceHandler())
	srv.add(handlers.NewAccountSericeServiceHandler())
	srv.add(handlers.NewUserSericeServiceHandler())
	srv.add(handlers.NewGithubServiceHandler())

	return srv
}
