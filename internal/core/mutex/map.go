package mutex

import (
	"sync"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (

	// SafeMap is a deterministic thread-safe map that holds the lock durations for each resource.
	// NOTE: This must be used in the context of a Temporal workflow.
	SafeMap struct {
		sync.RWMutex                          // Locks the pool for concurrent access
		internal     map[string]time.Duration // Holds the lock durations for each resource
	}
)

// Add adds a resource to the pool.
func (p *SafeMap) Add(ctx workflow.Context, id string, timeout time.Duration) {
	p.Lock()
	defer p.Unlock()

	p.internal[id] = timeout
}

// Remove removes the resource from the pool.
func (p *SafeMap) Remove(ctx workflow.Context, id string) {
	p.Lock()
	defer p.Unlock()

	delete(p.internal, id)
}

// Get returns the timeout for the requested resource.
func (p *SafeMap) Get(ctx workflow.Context, id string) (time.Duration, bool) {
	p.RLock()
	defer p.RUnlock()

	timeout, ok := p.internal[id]

	return timeout, ok
}
