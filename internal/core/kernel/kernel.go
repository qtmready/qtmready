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
package kernel

import (
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
)

type (
	// Kernel provides a central registry for various I/O providers in the application by exposing methods to register
	// and retrieve implementations for different hooks. This pattern allows for DRY implementation of I/O operations
	// and provides a single point of configuration for all I/O providers.
	Kernel interface {
		// RegisterRepoHook registers the given Repo implementation for the specified RepoHook.
		RegisterRepoHook(enum commonv1.RepoHook, hook Repo)

		// RegisterMessagingHook registers the given Messaging implementation for the specified MessagingHook.
		RegisterMessagingHook(enum commonv1.MessagingHook, hook Messaging)

		// RepoHook returns the Repo implementation registered for the specified RepoHook.
		//
		// It panics if no implementation is registered for the given hook.
		// It is the caller's responsibility to ensure that an implementation is registered before calling this method.
		// By panicking, we ensure that the application fails fast during development if a required implementation is missing.
		RepoHook(enum commonv1.RepoHook) Repo

		// MessagingHook returns the Messaging implementation registered for the specified MessagingHook.
		//
		// It panics if no implementation is registered for the given hook.
		// It is the caller's responsibility to ensure that an implementation is registered before calling this method.
		// By panicking, we ensure that the application fails fast during development if a required implementation is missing.
		MessagingHook(enum commonv1.MessagingHook) Messaging
	}

	Option func(k Kernel)

	kernel struct {
		hooks_repo      map[commonv1.RepoHook]Repo
		hooks_messaging map[commonv1.MessagingHook]Messaging
	}
)

func (k *kernel) RegisterRepoHook(hook commonv1.RepoHook, repo Repo) {
	if k.hooks_repo == nil {
		k.hooks_repo = make(map[commonv1.RepoHook]Repo)
	}

	k.hooks_repo[hook] = repo
}

func (k *kernel) RepoHook(enum commonv1.RepoHook) Repo {
	return k.hooks_repo[enum]
}

func (k *kernel) RegisterMessagingHook(hook commonv1.MessagingHook, messaging Messaging) {
	if k.hooks_messaging == nil {
		k.hooks_messaging = make(map[commonv1.MessagingHook]Messaging)
	}

	k.hooks_messaging[hook] = messaging
}

func (k *kernel) MessagingHook(enum commonv1.MessagingHook) Messaging {
	return k.hooks_messaging[enum]
}

func WithRepo(hook commonv1.RepoHook, repo Repo) Option {
	return func(k Kernel) {
		k.RegisterRepoHook(hook, repo)
	}
}

func WithMessaging(hook commonv1.MessagingHook, messaging Messaging) Option {
	return func(k Kernel) {
		k.RegisterMessagingHook(hook, messaging)
	}
}

func New(opts ...Option) Kernel {
	k := &kernel{
		hooks_repo:      make(map[commonv1.RepoHook]Repo),
		hooks_messaging: make(map[commonv1.MessagingHook]Messaging),
	}

	for _, opt := range opts {
		opt(k)
	}

	return k
}
