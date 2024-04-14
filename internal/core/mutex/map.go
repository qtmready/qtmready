package mutex

import (
	"sync"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (

	// SafeMap is a thread-safe map that holds the lock durations for each resource.
	// NOTE - This must be used in the context of a Temporal workflow.
	// TODO - The go standard library says "A Mutex must not be copied after first use". The way side effect works is, it stores the result
	// after the execution. That means the mutex is copied after the first use. This is a potential bug. We need to investigate this
	// further.
	SafeMap struct {
		*sync.Mutex                          // Locks the pool for concurrent access
		Internal    map[string]time.Duration // Holds the lock durations for each resource
	}

	Wrap struct {
		Map     *SafeMap
		Timeout time.Duration
		Ok      bool
	}
)

// Add adds a resource to the pool.
func (m *SafeMap) Add(ctx workflow.Context, id string, timeout time.Duration) {
	fn := workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		m.Lock()
		defer m.Unlock()

		m.Internal[id] = timeout

		return m
	})

	_ = fn.Get(m)
}

// Remove removes the resource from the pool.
func (m *SafeMap) Remove(ctx workflow.Context, id string) {
	fn := workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		m.Lock()
		defer m.Unlock()

		delete(m.Internal, id)

		return m
	})

	_ = fn.Get(m)
}

// Get returns the timeout for the requested resource.
func (m *SafeMap) Get(ctx workflow.Context, id string) (time.Duration, bool) {
	fn := workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		m.Lock()
		defer m.Unlock()

		timeout, ok := m.Internal[id]

		return Wrap{Map: m, Timeout: timeout, Ok: ok}
	})

	w := &Wrap{}
	_ = fn.Get(w)

	return w.Timeout, w.Ok
}
