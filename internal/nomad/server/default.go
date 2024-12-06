package server

import (
	"connectrpc.com/connect"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/hooks/slack"
	"go.breu.io/quantm/internal/nomad/intercepts"
)

// DefaultServer creates a new Nomad server instance with the provided options.
//
// FIXME: create an insecure handler for user registration and login. AccountServiceHandler?
func DefaultServer(opts ...Option) *Server {
	srv := New(opts...)

	// -- config/interceptors --

	interceptors := []connect.Interceptor{
		intercepts.RequestLogger(),
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
	srv.add(repos.NomadHandler(options...))

	// -- hooks/github --
	srv.add(github.NomadHandler(options...))

	// -- hooks/slack --
	srv.add(slack.NomadHandler(options...))

	return srv
}
