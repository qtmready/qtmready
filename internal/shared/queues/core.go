package queues

import (
	"log/slog"
	"sync"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/shared"
)

var (
	core     queues.Queue
	coreonce sync.Once
)

// Core is a singleton instance of the core queue.
func Core() queues.Queue {
	coreonce.Do(func() {
		slog.Info("queues/core: init ...")

		core = queues.New(
			queues.WithName("core"),
			queues.WithClient(shared.Temporal().Client()),
		)
	})

	return core
}
