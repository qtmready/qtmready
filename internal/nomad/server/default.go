package server

import (
	authnmd "go.breu.io/quantm/internal/auth/nomad"
	githubnmd "go.breu.io/quantm/internal/hooks/github/nomad"
)

// DefaultServer creates a new Nomad server instance with the provided options.
func DefaultServer(opts ...Option) *Server {
	srv := New(opts...)

	srv.add(authnmd.NewAccountSericeServiceHandler())
	srv.add(authnmd.NewOrgServiceServiceHandler())
	srv.add(authnmd.NewUserSericeServiceHandler())

	srv.add(githubnmd.NewGithubServiceHandler())

	return srv
}
