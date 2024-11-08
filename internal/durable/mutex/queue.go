package mutex

import (
	"sync"

	"go.breu.io/durex/queues"
)

var (
	_q     queues.Queue
	_qonce sync.Once
)

func Queue(opts ...queues.QueueOption) queues.Queue {
	_qonce.Do(func() {
		_q = queues.New(opts...)
	})

	return _q
}
