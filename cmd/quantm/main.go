package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.breu.io/durex/queues"
	"go.breu.io/graceful"

	"go.breu.io/quantm/internal/db"
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

	ctx := context.Background()
	quit := make(chan os.Signal, 1)

	githubcfg.Configure(githubcfg.WithConfig(cfg.Github))

	if err := durable.Configure(durable.WithConfig(cfg.Durable)); err != nil {
		slog.Error("unable to configure durable layer", "error", err.Error())
		os.Exit(1)
	}

	queues.SetDefaultPrefix("io.ctrlplane.")
	configure_qhooks()

	nmd := nomad.New(nomad.WithConfig(cfg.Nomad))
	app := graceful.New()

	app.Add(DB, db.Connection(db.WithConfig(cfg.DB)))
	app.Add(Durable, durable.Instance())
	app.Add(Github, githubcfg.Instance())
	app.Add(Nomad, nmd, DB, Durable, Github)
	app.Add(HooksQ, durable.OnHooks(), DB, Durable, Github)
	app.Add(Webhook, NewWebhookServer(), Durable, Github)

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
}

func configure_qhooks() {
	q := durable.OnHooks()

	q.CreateWorker()

	if q != nil {
		q.RegisterWorkflow(githubwfs.Install)
		q.RegisterActivity(&githubacts.Install{})
	}
}
