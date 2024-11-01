package main

import (
	"log/slog"

	"go.breu.io/durex/queues"

	pkg_durable "go.breu.io/quantm/internal/durable"
	github "go.breu.io/quantm/internal/hooks/github/workflows"
)

func configure_q_hooks(q queues.Queues) {
	slog.Info("main: configuring hooks queue ...")

	pkg_durable.OnHooks().CreateWorker()

	pkg_durable.OnHooks().RegisterWorkflow(github.Install)

	q[pkg_durable.OnHooks().Name()] = pkg_durable.OnHooks()

	slog.Info("main: hooks queue configured")
}
