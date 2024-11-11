package kernel

import (
	"log/slog"
	"sync"
)

var (
	_k   Kernel
	once sync.Once
)

func Configure(opts ...Option) {
	once.Do(func() {
		_k = New(opts...)
	})
}

func Get() Kernel {
	once.Do(func() {
		slog.Warn("kernel: Get() called before Configure()")

		_k = New()
	})

	return _k
}
