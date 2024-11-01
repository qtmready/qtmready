package durable

import (
	"sync"

	"go.breu.io/durex/queues"
)

var (
	coreq     queues.Queue
	coreqonce sync.Once

	hooksq     queues.Queue
	hooksqonce sync.Once
)

// OnCore returns the core queue.
//
// All workflows on this queue will have the ID prefix of
//
//	io.ctrlpane.core.{block}.{block_id}.{element}.{element_id}.{modifier}.{modifier_id}....
func OnCore() queues.Queue {
	coreqonce.Do(func() {
		coreq = queues.New(queues.WithName("core"), queues.WithClient(Client()))
	})

	return coreq
}

// OnHooks returns the hooks queue.
//
// All workflows on this queue will have the ID prefix of
//
//	io.ctrlpane.hooks.{block}.{block_id}.{element}.{element_id}.{modifier}.{modifier_id}....
func OnHooks() queues.Queue {
	hooksqonce.Do(func() {
		hooksq = queues.New(queues.WithName("hooks"), queues.WithClient(Client()))
	})

	return hooksq
}
