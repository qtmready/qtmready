package defs

import (
	"sync"

	"go.breu.io/quantm/internal/shared"
)

var (
	instance Core
	once     sync.Once
)

type (
	// Core is the interface that defines the core of the application. It is the main entry point for the application.
	// It is responsible for registering different providers and exposing them to the rest of the application.
	//
	// NOTE: This is not an ideal design, because it only registers providers for the providers. It does not register
	// workflows. We may need to revisit this design in the future.
	Core interface {
		RegisterRepoProvider(RepoProvider, RepoIO)
		ResgisterMessageProvider(MessageProvider, MessageIO)

		RepoIO(RepoProvider) RepoIO
		MessageIO(MessageProvider) MessageIO
	}

	Option func(Core)

	// Providers is a struct that holds the different providers that are registered with the core.
	Providers struct {
		repos   map[RepoProvider]RepoIO
		message map[MessageProvider]MessageIO
	}

	core struct {
		providers Providers
		once      sync.Once // Do we really need this?
	}
)

func (c *core) RegisterRepoProvider(provider RepoProvider, activities RepoIO) {
	c.providers.repos[provider] = activities
}

func (c *core) RepoIO(name RepoProvider) RepoIO {
	if p, ok := c.providers.repos[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()).Error())
}

func (c *core) ResgisterMessageProvider(provider MessageProvider, activities MessageIO) {
	c.providers.message[provider] = activities
}

func (c *core) MessageIO(name MessageProvider) MessageIO {
	if p, ok := c.providers.message[name]; ok {
		return p
	}

	panic(NewProviderNotFoundError(name.String()).Error())
}

// WithMessageProvider registers a repo provider with the core.
func WithMessageProvider(provider MessageProvider, io MessageIO) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering message provider", "name", provider.String())
		c.ResgisterMessageProvider(provider, io)
	}
}

// WithRepoProvider registers a repo provider with the core.
func WithRepoProvider(provider RepoProvider, io RepoIO) Option {
	return func(c Core) {
		shared.Logger().Info("core: registering repo provider", "name", provider.String())
		c.RegisterRepoProvider(provider, io)
	}
}

// Instance returns a singleton instance of the core. It is best to call this function in the main() function to
// register the providers available to the service. This is because the core uses workflow and providers implementations
// to access the providers.
func Instance(opts ...Option) Core {
	if instance == nil {
		shared.Logger().Info("core: instance not initialized, initializing now ...")
		once.Do(func() {
			instance = &core{
				providers: Providers{
					repos:   make(map[RepoProvider]RepoIO),
					message: make(map[MessageProvider]MessageIO),
				},
			}

			for _, opt := range opts {
				opt(instance)
			}
		})
	}

	return instance
}
