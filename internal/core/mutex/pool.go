package mutex

import (
	"sync"
	"time"
)

type (

	// Pool of locks waiting to be acquired.
	Pool struct {
		sync.RWMutex
		internal map[string]time.Duration // Holds the lock durations for each resource
	}
)

/**
 * Receivers for Pool
 **/

func (p *Pool) Add(id string, timeout time.Duration) {
	p.Lock()
	defer p.Unlock()

	p.internal[id] = timeout
}

func (p *Pool) Remove(id string) {
	p.Lock()
	defer p.Unlock()

	delete(p.internal, id)
}

func (p *Pool) Read(id string) (time.Duration, bool) {
	p.RLock()
	defer p.RUnlock()

	timeout, ok := p.internal[id]

	return timeout, ok
}
