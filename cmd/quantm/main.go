package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.breu.io/durex/queues"
	"go.breu.io/graceful"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db/migrations"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/hooks/slack"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

const (
	DB      = "db"
	Durable = "durable"
	Pulse   = "pulse"
	Kernel  = "kernel"
	Github  = "github"
	Nomad   = "nomad"
	CoreQ   = "core_q"
	HooksQ  = "hooks_q"
	Webhook = "webhook"
)

func main() {
	cfg := &Config{}
	cfg.Load()

	// - global configuration

	configure_logger(cfg.Debug)
	auth.SetSecret(cfg.Secret)

	ctx := context.Background()
	quit := make(chan os.Signal, 1)

	// - configure services
	github.Configure(github.WithConfig(cfg.Github))
	slack.Configure(slack.WithConfig(cfg.Slack))
	kernel.Configure(
		kernel.WithRepoHook(eventsv1.RepoHook_REPO_HOOK_GITHUB, &github.KernelImpl{}),
		kernel.WithChatHook(eventsv1.ChatHook_CHAT_HOOK_SLACK, &slack.KernelImpl{}),
	)

	if err := durable.Configure(durable.WithConfig(cfg.Durable)); err != nil {
		slog.Error("unable to configure durable layer", "error", err.Error())
		os.Exit(1)
	}

	queues.SetDefaultPrefix("ai.ctrlplane.")

	q_core()
	q_hooks()

	// - configure dependency graph for services

	app := graceful.New()
	cfg.Parse(app)
	cfg.Dependencies(app)

	// - if the migrate flag is set, run migrations and exit

	if cfg.Mode == ModeMigrate {
		if err := migrations.Run(ctx, cfg.DB); err != nil {
			slog.Error("unable to migrate database", "error", err.Error())
		}

		os.Exit(0)
	}

	// - start the services as defined in the dependency graph

	if err := app.Start(ctx); err != nil {
		slog.Error("unable to start service", "error", err.Error())
		os.Exit(1)
	}

	// - wait for a signal to stop the services

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	<-quit

	// - gracefully stop the services

	if err := app.Stop(ctx); err != nil {
		slog.Error("unable to stop service", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("service stopped, exiting...")

	os.Exit(0)
}
