package queue

import (
	"log/slog"

	"go.breu.io/durex/queues"
)

func init() {
	slog.Info("queues: init ...")
	queues.SetDefaultPrefix("ai.ctrlplane.")
}
