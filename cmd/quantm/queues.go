package main

import (
	"log/slog"

	"go.breu.io/durex/queues"

	pkg_durable "go.breu.io/quantm/internal/durable"
	githubwfs "go.breu.io/quantm/internal/hooks/github/workflows"
)

func q_prefix() {
	queues.SetDefaultPrefix("ai.ctrlplane.")
}

func q_hooks(q queues.Queues) {
	slog.Info("main: configuring hooks queue ...")

	pkg_durable.OnHooks().CreateWorker()

	pkg_durable.OnHooks().RegisterWorkflow(githubwfs.Install)

	q[pkg_durable.OnHooks().Name()] = pkg_durable.OnHooks()

	slog.Info("main: hooks queue configured")
}
