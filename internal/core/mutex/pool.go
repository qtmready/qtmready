package mutex

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	Pool struct {
		data  map[string]time.Duration
		mutex workflow.Mutex
	}
)

func (p *Pool) add(ctx workflow.Context, id string, timeout time.Duration) {
	_ = p.mutex.Lock(ctx)
	defer p.mutex.Unlock()
	p.data[id] = timeout
}

func (p *Pool) remove(ctx workflow.Context, id string) {
	_ = p.mutex.Lock(ctx)
	defer p.mutex.Unlock()
	delete(p.data, id)
}

func (p *Pool) get(id string) (time.Duration, bool) {
	timeout, ok := p.data[id]
	return timeout, ok
}

func (p *Pool) size() int {
	return len(p.data)
}

func NewPool(ctx workflow.Context) *Pool {
	return &Pool{
		data:  make(map[string]time.Duration),
		mutex: workflow.NewMutex(ctx),
	}
}
