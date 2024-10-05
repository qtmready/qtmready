package ws

import (
	"sync"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/shared"
)

var (
	_q     queues.Queue
	_qonce sync.Once
)

// Queue returns the default queue for the websockets. This queue is used to
// manage the connections hub.
func Queue() queues.Queue {
	_qonce.Do(func() {
		_q = queues.New(
			queues.WithName("websockets"),
			queues.WithClient(shared.Temporal().Client()),
		)

		_q.RegisterWorkflow(ConnectionsHubWorkflow)
	})

	return _q
}
