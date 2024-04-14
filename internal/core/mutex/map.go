package mutex

import (
	"sync"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (

	// Map is a thread-safe map that holds the lock durations for each resource.
	// NOTE: This must be used in the context of a Temporal workflow.
	Map struct {
		sync.Mutex                          // Locks the pool for concurrent access
		Internal   map[string]time.Duration // Holds the lock durations for each resource
	}

	Wrap struct {
		Map     *Map
		Timeout time.Duration
		Ok      bool
	}
)

// Add adds a resource to the pool.
func (p *Map) Add(ctx workflow.Context, id string, timeout time.Duration) {
	fn := workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		p.Lock()
		defer p.Unlock()

		p.Internal[id] = timeout

		return p
	})

	_ = fn.Get(p)
}

// Remove removes the resource from the pool.
func (p *Map) Remove(ctx workflow.Context, id string) {
	fn := workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		p.Lock()
		defer p.Unlock()

		delete(p.Internal, id)

		return p
	})

	_ = fn.Get(p)
}

// Get returns the timeout for the requested resource.
func (p *Map) Get(ctx workflow.Context, id string) (time.Duration, bool) {
	fn := workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		p.Lock()
		defer p.Unlock()

		timeout, ok := p.Internal[id]

		return Wrap{Map: p, Timeout: timeout, Ok: ok}
	})

	w := &Wrap{}
	_ = fn.Get(w)

	return w.Timeout, w.Ok
}
