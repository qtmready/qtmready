package mutex

import (
	"time"
)

type (
	Pool map[string]time.Duration // Pool holds the timeout against the resource ID.
)

func (p Pool) Add(id string, timeout time.Duration) {
	p[id] = timeout
}

func (p Pool) Remove(id string) {
	delete(p, id)
}

func (p Pool) Get(id string) (time.Duration, bool) {
	timeout, ok := p[id]
	return timeout, ok
}

func (p Pool) Size() int {
	return len(p)
}

func NewSimpleMap() Pool {
	return make(Pool)
}
