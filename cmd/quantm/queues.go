package main

import (
	"log/slog"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/durable"
	githubacts "go.breu.io/quantm/internal/hooks/github/activities"
	githubwfs "go.breu.io/quantm/internal/hooks/github/workflows"
)

func q_prefix() {
	queues.SetDefaultPrefix("ai.ctrlplane.")
}

func q_hooks(q queues.Queues) {
	slog.Info("main: configuring hooks queue ...")

	durable.OnHooks().CreateWorker()

	durable.OnHooks().RegisterWorkflow(githubwfs.Install)
	durable.OnHooks().RegisterActivity(&githubacts.Install{})

	q[durable.OnHooks().Name()] = durable.OnHooks()

	slog.Info("main: hooks queue configured")
}
