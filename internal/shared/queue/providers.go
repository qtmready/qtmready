package queue

import (
	"log/slog"
	"sync"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/shared"
)

var (
	provider     queues.Queue
	provideronce sync.Once
)

// Providers is a singleton instance of the providers queue.
func Providers() queues.Queue {
	provideronce.Do(func() {
		slog.Info("queues/providers: init ...")

		provider = queues.New(
			queues.WithName("providers"),
			queues.WithClient(shared.Temporal().Client()),
		)
	})

	return provider
}
