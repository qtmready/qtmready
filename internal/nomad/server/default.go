package server

import (
	"connectrpc.com/connect"

	"go.breu.io/quantm/internal/auth"
	authnmd "go.breu.io/quantm/internal/auth/nomad"
	reposnmd "go.breu.io/quantm/internal/core/repos/nomad"
	githubnmd "go.breu.io/quantm/internal/hooks/github/nomad"
	"go.breu.io/quantm/internal/observe/logs"
)

// DefaultServer creates a new Nomad server instance with the provided options.
//
// FIXME: create an insecure handler for user registration and login. AccountServiceHandler?
func DefaultServer(opts ...Option) *Server {
	srv := New(opts...)

	// -- config/interceptors --

	interceptors := []connect.Interceptor{
		logs.NomadRequestLogger(),
	}

	// -- config/handlers --
	options := []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
	}

	// -- auth --
	srv.add(authnmd.NewAccountSericeServiceHandler(options...))
	srv.add(authnmd.NewOrgServiceServiceHandler(options...))
	srv.add(authnmd.NewUserSericeServiceHandler(options...))

	options = append(options, connect.WithInterceptors(auth.NomadAuthenticator()))

	// -- core/repos --
	srv.add(reposnmd.NewRepoServiceHandler(options...))

	// -- hooks/github --
	srv.add(githubnmd.NewGithubServiceHandler(options...))

	return srv
}
