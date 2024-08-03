package mutex

import (
	"time"
)

type (
	Pool map[string]time.Duration // Pool holds the timeout against the resource ID.
)

func (p Pool) add(id string, timeout time.Duration) {
	p[id] = timeout
}

func (p Pool) remove(id string) {
	delete(p, id)
}

func (p Pool) get(id string) (time.Duration, bool) {
	timeout, ok := p[id]
	return timeout, ok
}

func (p Pool) size() int {
	return len(p)
}

func NewPool() Pool {
	return make(Pool)
}
