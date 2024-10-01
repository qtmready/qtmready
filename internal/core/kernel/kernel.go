// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

// Package kernel provides a central registry for various I/O providers in the application.
//
// The Kernel pattern implemented here serves several important purposes:
//
//  1. Centralized Configuration: It provides a single point of configuration for all I/O providers
//     (e.g., repository access, messaging systems) used throughout the application. This centralization makes it
//     easier to manage and modify the application's external dependencies.
//
//  2. Dependency Injection: By registering providers in the Kernel, we implement a form of dependency injection. This
//     allows for easier testing and more flexible architecture, as providers can be swapped out without changing the
//     core application logic.
//
//  3. Abstraction: The Kernel abstracts away the details of how different I/O operations are performed. This allows
//     the rest of the application to work with a consistent interface, regardless of the underlying implementation.
//
//  4. Singleton Pattern: The Kernel is implemented as a singleton, ensuring that there's only one instance managing
//     all providers across the application. This prevents duplication and ensures consistency.
//
//  5. Lazy Initialization: Providers are only initialized when first requested, which can help improve application
//     startup time and resource usage.
//
// Usage:
//
// Initialize the Kernel with providers:
//
//	gitRepoIO := &GitRepoIO{}
//	slackMessageIO := &SlackMessageIO{}
//
//	k := kernel.Instance(
//	  kernel.WithRepoProvider(defs.RepoProvider, gitRepoIO),
//	  kernel.WithMessageProvider(defs.MessageProviderSlack, slackMessageIO),
//	)
//
// Retrieve a RepoIO:
//
//	github := k.RepoIO(defs.RepoProviderGithub)
//	// Use gitIO to interact with Git repositories
//
// Retrieve a MessageIO:
//
//	slack := k.MessageIO(defs.MessageProviderSlack)
//	// Use slack to send messages via Slack
//
// Register a new provider after initialization:
//
//	k.RegisterRepoProvider(defs.SVNRepoProvider, &SVNRepoIO{})
//
// Access the Kernel instance from anywhere in the application:
//
//	k := kernel.Instance()
//	// Use k to access registered providers
package kernel

import (
	"log/slog"
	"sync"

	"go.breu.io/quantm/internal/core/defs"
)

type (
	// Kernel defines the interface for managing repo and message providers.
	Kernel interface {
		// RegisterRepoProvider registers a RepoIO for a given RepoProvider.
		//
		// Usage:
		//  kernel.Instance().RegisterRepoProvider(defs.RepoProvider, &github.RepoIO{})
		RegisterRepoProvider(defs.RepoProvider, RepoIO)

		// RepoIO retrieves the RepoIO for a given RepoProvider.
		//
		// It panics with a ProviderNotFoundError if the provider is not registered.
		//
		// Panicking is used here because a missing provider indicates a critical configuration error. This approach ensures
		// that such errors are caught early, typically during application startup or in tests, rather than manifesting as
		// nil pointer exceptions later in the program execution.
		//
		// Usage:
		//  io := kernel.Instance().RepoIO(defs.RepoProviderGithub)
		//  io.SetEarlyWarning(ctx, repo.ID.String(), true)
		RepoIO(defs.RepoProvider) RepoIO

		// RegisterMessageProvider registers a MessageIO for a given MessageProvider.
		//
		// Usage:
		//  kernel.Instance().RegisterMessageProvider(defs.MessageProviderSlack, &slack.MessageIO{})
		RegisterMessageProvider(defs.MessageProvider, MessageIO)

		// MessageIO retrieves the MessageIO for a given MessageProvider.
		//
		// It panics with a ProviderNotFoundError if the provider is not registered.
		//
		// As with RepoIO, panicking is used to immediately surface critical configuration errors. This fail-fast approach
		// helps identify missing or misconfigured providers during application initialization or testing, rather than
		// allowing the application to continue running in an invalid state.
		//
		// Usage:
		//  io := kernel.Instance().MessageIO(defs.MessageProviderSlack)
		//  // Use slack to send messages via Slack
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

	slog.Info("kernel: registered repo provider", "provider", provider)
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

	slog.Info("kernel: registered message provider", "provider", provider)
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

// Instance returns the singleton instance of the Kernel. It initializes the Kernel on the first call and applies the
// provided options.
//
// It's recommended to instantiate the Kernel with all necessary providers in the main function of your application.
// This ensures that all required providers are registered before any part of the application attempts to use them.
//
// Please note that providers cannot be registred after the Kernel has been initialized.
//
// The kernel panics if a provider is not found. This is by design to ensure that critical configuration errors are
// caught early, typically during application startup or in tests, rather than manifesting as nil pointer exceptions
// later.
//
// Please note that providers cannot be registred after the Kernel has been initialized.
//
// The kernel panics if a provider is not found. This is by design to ensure that
// critical configuration errors are caught early, typically during application
// startup or in tests, rather than manifesting as nil pointer exceptions later.
//
// Example:
//
//	func main() {}
//	  kernel.Instance(
//	    kernel.WithRepoProvider(defs.RepoProvider, &github.RepoIO{}),
//	    kernel.WithMessageProvider(defs.MessageProviderSlack, slack.MessageIO{}),
//	  )
//
//	  // Rest of your application setup...
//	  // ...
//
//	  // Run your application
//	  app.Run()
//	}
//
// After this initial setup, you can retrieve the Kernel instance from anywhere in your application using
// kernel.Instance() without any arguments.
func Instance(opts ...Option) Kernel {
	once.Do(func() {
		slog.Info("kernel: init ...")

		instance = &_k{
			providers: Providers{
				repos:   make(map[defs.RepoProvider]RepoIO),
				message: make(map[defs.MessageProvider]MessageIO),
			},
		}

		for _, opt := range opts {
			opt(instance)
		}

		slog.Info("kernel: initialized")
	})

	return instance
}
