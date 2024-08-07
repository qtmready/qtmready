// Package kernel provides a central registry for various I/O providers in the application.
//
// The Kernel pattern implemented here serves several important purposes:
//
//  1. Centralized Configuration: It provides a single point of configuration for
//     all I/O providers (e.g., repository access, messaging systems) used throughout
//     the application. This centralization makes it easier to manage and modify the
//     application's external dependencies.
//
//  2. Dependency Injection: By registering providers in the Kernel, we implement a
//     form of dependency injection. This allows for easier testing and more flexible
//     architecture, as providers can be swapped out without changing the core application logic.
//
//  3. Abstraction: The Kernel abstracts away the details of how different I/O operations
//     are performed. This allows the rest of the application to work with a consistent
//     interface, regardless of the underlying implementation.
//
//  4. Singleton Pattern: The Kernel is implemented as a singleton, ensuring that there's
//     only one instance managing all providers across the application. This prevents
//     duplication and ensures consistency.
//
//  5. Lazy Initialization: Providers are only initialized when first requested, which
//     can help improve application startup time and resource usage.
//
// Usage:
//
// 1. Initialize the Kernel with providers:
//
//	gitRepoIO := &GitRepoIO{}
//	slackMessageIO := &SlackMessageIO{}
//
//	k := kernel.Instance(
//		kernel.WithRepoProvider(defs.GitRepoProvider, gitRepoIO),
//		kernel.WithMessageProvider(defs.SlackMessageProvider, slackMessageIO),
//	)
//
// 2. Retrieve a RepoIO:
//
//	gitIO := k.RepoIO(defs.GitRepoProvider)
//	// Use gitIO to interact with Git repositories
//
// 3. Retrieve a MessageIO:
//
//	slackIO := k.MessageIO(defs.SlackMessageProvider)
//	// Use slackIO to send messages via Slack
//
// 4. Register a new provider after initialization:
//
//	k.RegisterRepoProvider(defs.SVNRepoProvider, &SVNRepoIO{})
//
// 5. Access the Kernel instance from anywhere in the application:
//
//	k := kernel.Instance()
//	// Use k to access registered providers
package kernel

import (
	"sync"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/shared"
)

type (
	// Kernel defines the interface for managing repo and message providers.
	Kernel interface {
		// RegisterRepoProvider registers a RepoIO for a given RepoProvider.
		//
		// Usage:
		//  gitRepoIO := &GitRepoIO{}
		//  k.RegisterRepoProvider(defs.GitRepoProvider, gitRepoIO)
		RegisterRepoProvider(defs.RepoProvider, RepoIO)

		// RepoIO retrieves the RepoIO for a given RepoProvider.
		// It panics with a ProviderNotFoundError if the provider is not registered.
		//
		// Panicking is used here because a missing provider indicates a critical
		// configuration error. This approach ensures that such errors are caught
		// early, typically during application startup or in tests, rather than
		// manifesting as nil pointer exceptions later in the program execution.
		//
		// Usage:
		//  gitIO := k.RepoIO(defs.GitRepoProvider)
		//  // Use gitIO to interact with Git repositories
		RepoIO(defs.RepoProvider) RepoIO

		// RegisterMessageProvider registers a MessageIO for a given MessageProvider.
		//
		// Usage:
		//  slackMessageIO := &SlackMessageIO{}
		//  k.RegisterMessageProvider(defs.SlackMessageProvider, slackMessageIO)
		RegisterMessageProvider(defs.MessageProvider, MessageIO)

		// MessageIO retrieves the MessageIO for a given MessageProvider.
		// It panics with a ProviderNotFoundError if the provider is not registered.
		//
		// As with RepoIO, panicking is used to immediately surface critical
		// configuration errors. This fail-fast approach helps identify missing
		// or misconfigured providers during application initialization or testing,
		// rather than allowing the application to continue running in an invalid state.
		//
		// Usage:
		//  slackIO := k.MessageIO(defs.SlackMessageProvider)
		//  // Use slackIO to send messages via Slack
		MessageIO(defs.MessageProvider) MessageIO
	}

	// Option is a function type used for configuring the Kernel.
	Option func(Kernel)

	// Providers struct holds maps of registered providers.
	Providers struct {
		repos   map[defs.RepoProvider]RepoIO
		message map[defs.MessageProvider]MessageIO
	}

	// _k is the internal implementation of the Kernel interface.
	_k struct {
		providers Providers
	}
)

var (
	once     sync.Once // Ensures that the Kernel is initialized only once
	instance Kernel    // The singleton instance of the Kernel
)

// RegisterRepoProvider registers a RepoIO for a given RepoProvider.
func (k *_k) RegisterRepoProvider(provider defs.RepoProvider, io RepoIO) {
	k.providers.repos[provider] = io

	shared.Logger().Info("kernel: registered repo provider", "provider", provider)
}

// RepoIO retrieves the RepoIO for a given RepoProvider.
func (k *_k) RepoIO(provider defs.RepoProvider) RepoIO {
	io, ok := k.providers.repos[provider]
	if !ok {
		panic(defs.NewProviderNotFoundError(provider.String()).Error())
	}

	return io
}

// RegisterMessageProvider registers a MessageIO for a given MessageProvider.
func (k *_k) RegisterMessageProvider(provider defs.MessageProvider, io MessageIO) {
	k.providers.message[provider] = io

	shared.Logger().Info("kernel: registered message provider", "provider", provider)
}

// MessageIO retrieves the MessageIO for a given MessageProvider.
func (k *_k) MessageIO(provider defs.MessageProvider) MessageIO {
	io, ok := k.providers.message[provider]
	if !ok {
		panic(defs.NewProviderNotFoundError(provider.String()).Error())
	}

	return io
}

// WithRepoProvider creates an Option to register a RepoProvider.
func WithRepoProvider(provider defs.RepoProvider, io RepoIO) Option {
	return func(k Kernel) {
		k.RegisterRepoProvider(provider, io)
	}
}

// WithMessageProvider creates an Option to register a MessageProvider.
func WithMessageProvider(provider defs.MessageProvider, io MessageIO) Option {
	return func(k Kernel) {
		k.RegisterMessageProvider(provider, io)
	}
}

// Instance returns the singleton instance of the Kernel.
// It initializes the Kernel on the first call and applies the provided options.
//
// It's recommended to instantiate the Kernel with all necessary providers
// in the main function of your application. This ensures that all required
// providers are registered before any part of the application attempts to use them.
//
// Example:
//
//	func main() {
//		gitRepoIO := &GitRepoIO{}
//		slackMessageIO := &SlackMessageIO{}
//
//		kernel.Instance(
//			kernel.WithRepoProvider(defs.GitRepoProvider, gitRepoIO),
//			kernel.WithMessageProvider(defs.SlackMessageProvider, slackMessageIO),
//		)
//
//		// Rest of your application setup...
//		// ...
//
//		// Run your application
//		app.Run()
//	}
//
// After this initial setup, you can retrieve the Kernel instance
// from anywhere in your application using kernel.Instance() without any arguments.
func Instance(opts ...Option) Kernel {
	once.Do(func() {
		shared.Logger().Info("kernel: init ...")

		instance = &_k{
			providers: Providers{
				repos:   make(map[defs.RepoProvider]RepoIO),
				message: make(map[defs.MessageProvider]MessageIO),
			},
		}

		for _, opt := range opts {
			opt(instance)
		}

		shared.Logger().Info("kernel: initialized")
	})

	return instance
}
