package kernel

import (
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
		_k = New()
	})

	return _k
}
