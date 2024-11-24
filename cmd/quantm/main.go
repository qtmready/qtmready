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
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/migrations"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/nomad"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	"go.breu.io/quantm/internal/pulse"
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

	if err := durable.Configure(durable.WithConfig(cfg.Durable)); err != nil {
		slog.Error("unable to configure durable layer", "error", err.Error())
		os.Exit(1)
	}

	queues.SetDefaultPrefix("ai.ctrlplane.")

	q_core()
	q_hooks()

	nmd := nomad.New(nomad.WithConfig(cfg.Nomad))

	// - configure dependency graph for services

	app := graceful.New()

	app.Add(DB, db.Connection(db.WithConfig(cfg.DB)))
	app.Add(Pulse, pulse.Instance(pulse.WithConfig(cfg.Pulse)))
	app.Add(Durable, durable.Instance())
	app.Add(Github, github.Instance())
	app.Add(CoreQ, durable.OnCore(), DB, Durable, Pulse, Github)
	app.Add(HooksQ, durable.OnHooks(), DB, Durable, Pulse, Github)
	app.Add(Nomad, nmd, DB, Durable, Pulse, Github)
	app.Add(Webhook, NewWebhookServer(), DB, Durable, Github)
	app.Add(Kernel, kernel.New(kernel.WithRepoHook(eventsv1.RepoHook_REPO_HOOK_GITHUB, &github.KernelImpl{})), Github)

	// - if the migrate flag is set, run migrations and exit

	if cfg.Migrate {
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
