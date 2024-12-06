package server

import (
	"connectrpc.com/connect"

	"go.breu.io/quantm/internal/auth"
	reposnmd "go.breu.io/quantm/internal/core/repos/nomad"
	githubnmd "go.breu.io/quantm/internal/hooks/github/nomad"
	slacknmd "go.breu.io/quantm/internal/hooks/slack/nomad"
	"go.breu.io/quantm/internal/observe"
)

// DefaultServer creates a new Nomad server instance with the provided options.
//
// FIXME: create an insecure handler for user registration and login. AccountServiceHandler?
func DefaultServer(opts ...Option) *Server {
	srv := New(opts...)

	// -- config/interceptors --

	interceptors := []connect.Interceptor{
		observe.NomadRequestLogger(),
	}

	// -- config/handlers --
	options := []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
	}

	// - insecure handlers -
	// -- auth --
	srv.add(auth.NomadAccountServiceHandler(options...))
	srv.add(auth.NomadOrgServiceHandler(options...))
	srv.add(auth.NomadUserServiceHandler(options...))

	// - secure handlers -

	options = append(options, connect.WithInterceptors(auth.NomadInterceptor()))

	// -- core/repos --
	srv.add(reposnmd.NewRepoServiceHandler(options...))

	// -- hooks/github --
	srv.add(githubnmd.NewGithubServiceHandler(options...))

	// -- hooks/slack --
	srv.add(slacknmd.NewSlackServiceHandler(options...))

	return srv
}
