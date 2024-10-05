package queue

import (
	"log/slog"
	"sync"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/shared"
)

var (
	mutex     queues.Queue
	mutexonce sync.Once
)

// Mutex is a singleton instance of the mutex queue.
func Mutex() queues.Queue {
	mutexonce.Do(func() {
		slog.Info("queues/mutex: init ...")

		mutex = queues.New(
			queues.WithName("mutex"),
			queues.WithClient(shared.Temporal().Client()),
		)
	})

	return mutex
}
