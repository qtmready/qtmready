package server

import (
	"connectrpc.com/connect"

	authnmd "go.breu.io/quantm/internal/auth/nomad"
	githubnmd "go.breu.io/quantm/internal/hooks/github/nomad"
	"go.breu.io/quantm/internal/observe/intercept"
)

// DefaultServer creates a new Nomad server instance with the provided options.
func DefaultServer(opts ...Option) *Server {
	srv := New(opts...)

	// -- config/interceptors --

	interceptors := []connect.Interceptor{
		intercept.RequestLogger(),
	}

	// -- config/handlers --
	options := []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
	}

	// -- auth --
	srv.add(authnmd.NewAccountSericeServiceHandler(options...))
	srv.add(authnmd.NewOrgServiceServiceHandler(options...))
	srv.add(authnmd.NewUserSericeServiceHandler(options...))

	// -- hooks/github --
	srv.add(githubnmd.NewGithubServiceHandler(options...))

	return srv
}
