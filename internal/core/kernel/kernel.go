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
	"context"
	"log/slog"

	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Kernel provides a central registry for various I/O providers in the application by exposing methods to register
	// and retrieve implementations for different hooks. This pattern allows for DRY implementation of I/O operations
	// and provides a single point of configuration for all I/O providers.
	Kernel interface {
		// Hooks returns a list of all hooks registered in the Kernel.
		Hooks() []string

		// RegisterRepoHook registers the given Repo implementation for the specified RepoHook.
		RegisterRepoHook(enum eventsv1.RepoHook, hook Repo)

		// RegisterChatHook registers the given Chat implementation for the specified ChatHook.
		RegisterChatHook(enum eventsv1.ChatHook, hook Chat)

		// RepoHook returns the Repo implementation registered for the specified RepoHook.
		//
		// It panics if no implementation is registered for the given hook.
		// It is the caller's responsibility to ensure that an implementation is registered before calling this method.
		// By panicking, we ensure that the application fails fast during development if a required implementation is missing.
		RepoHook(enum eventsv1.RepoHook) Repo

		// ChatHook returns the chat platform implementation registered for the specified ChatHook.
		//
		// It panics if no implementation is registered for the given hook.
		// It is the caller's responsibility to ensure that an implementation is registered before calling this method.
		// By panicking, we ensure that the application fails fast during development if a required implementation is missing.
		ChatHook(enum eventsv1.ChatHook) Chat

		// Start is a noop method that conforms to graceful.Service interface.
		Start(ctx context.Context) error

		// Stop is a noop method that conforms to graceful.Service interface.
		Stop(ctx context.Context) error
	}

	Option func(k Kernel)

	kernel struct {
		hooks_repo map[eventsv1.RepoHook]Repo
		hooks_chat map[eventsv1.ChatHook]Chat
	}
)

func (k *kernel) Hooks() []string {
	hooks := make([]string, 0)

	for hook := range k.hooks_repo {
		hooks = append(hooks, hook.String())
	}

	for hook := range k.hooks_chat {
		hooks = append(hooks, hook.String())
	}

	return hooks
}

func (k *kernel) RegisterRepoHook(hook eventsv1.RepoHook, repo Repo) {
	if k.hooks_repo == nil {
		k.hooks_repo = make(map[eventsv1.RepoHook]Repo)
	}

	slog.Info("kernel: registering repo hook", "hook", hook.String())

	k.hooks_repo[hook] = repo
}

func (k *kernel) RepoHook(enum eventsv1.RepoHook) Repo {
	return k.hooks_repo[enum]
}

func (k *kernel) RegisterChatHook(hook eventsv1.ChatHook, chat Chat) {
	if k.hooks_chat == nil {
		k.hooks_chat = make(map[eventsv1.ChatHook]Chat)
	}

	slog.Info("kernel: registering chat hook", "hook", hook.String())

	k.hooks_chat[hook] = chat
}

func (k *kernel) ChatHook(enum eventsv1.ChatHook) Chat {
	return k.hooks_chat[enum]
}

func (k *kernel) Start(ctx context.Context) error {
	slog.Info("kernel: starting ...", "hooks", k.Hooks())

	return nil
}

func (k *kernel) Stop(ctx context.Context) error { return nil }

func WithRepoHook(hook eventsv1.RepoHook, repo Repo) Option {
	return func(k Kernel) {
		k.RegisterRepoHook(hook, repo)
	}
}

func WithChatHook(hook eventsv1.ChatHook, chat Chat) Option {
	return func(k Kernel) {
		k.RegisterChatHook(hook, chat)
	}
}

func New(opts ...Option) Kernel {
	k := &kernel{
		hooks_repo: make(map[eventsv1.RepoHook]Repo),
		hooks_chat: make(map[eventsv1.ChatHook]Chat),
	}

	for _, opt := range opts {
		opt(k)
	}

	return k
}
