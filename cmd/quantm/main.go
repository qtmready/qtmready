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
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/migrations"
	"go.breu.io/quantm/internal/durable"
	githubacts "go.breu.io/quantm/internal/hooks/github/activities"
	githubcfg "go.breu.io/quantm/internal/hooks/github/config"
	githubwfs "go.breu.io/quantm/internal/hooks/github/workflows"
	"go.breu.io/quantm/internal/nomad"
)

const (
	DB      = "db"
	Durable = "durable"
	Github  = "github"
	Nomad   = "nomad"
	HooksQ  = "hooks_q"
	Webhook = "webhook"
)

func main() {
	cfg := &Config{}
	cfg.Load()

	configure_logger(cfg.Debug)
	auth.SetSecret(cfg.Secret)

	ctx := context.Background()
	quit := make(chan os.Signal, 1)

	githubcfg.Configure(githubcfg.WithConfig(cfg.Github))

	if err := durable.Configure(durable.WithConfig(cfg.Durable)); err != nil {
		slog.Error("unable to configure durable layer", "error", err.Error())
		os.Exit(1)
	}

	queues.SetDefaultPrefix("ai.ctrlplane.")
	configure_qhooks()

	nmd := nomad.New(nomad.WithConfig(cfg.Nomad))
	app := graceful.New()

	app.Add(DB, db.Connection(db.WithConfig(cfg.DB)))
	app.Add(Durable, durable.Instance())
	app.Add(Github, githubcfg.Instance())
	app.Add(Nomad, nmd, DB, Durable, Github)
	app.Add(HooksQ, durable.OnHooks(), DB, Durable, Github)
	app.Add(Webhook, NewWebhookServer(), Durable, Github)

	if cfg.Migrate {
		if err := migrations.Run(ctx, cfg.DB); err != nil {
			slog.Error("unable to migrate database", "error", err.Error())
		}

		os.Exit(0)
	}

	if err := app.Start(ctx); err != nil {
		slog.Error("unable to start service", "error", err.Error())
		os.Exit(1)
	}

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	<-quit

	if err := app.Stop(ctx); err != nil {
		slog.Error("unable to stop service", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("service stopped, exiting...")

	os.Exit(0)
}

func configure_qhooks() {
	q := durable.OnHooks()

	q.CreateWorker()

	if q != nil {
		q.RegisterWorkflow(githubwfs.Install)
		q.RegisterActivity(&githubacts.Install{})

		q.RegisterWorkflow(githubwfs.SyncRepos)
		q.RegisterActivity(&githubacts.InstallRepos{})
	}
}
