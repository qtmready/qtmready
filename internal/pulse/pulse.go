package pulse

import (
	"context"
	"sync"

	"go.breu.io/quantm/internal/pulse/config"
)

type (
	Config = config.Config
	Option = config.Option
)

var (
	_c   *Config
	once sync.Once
)

func Configure(opts ...Option) *Config {
	once.Do(func() {
		_c = config.New(opts...)
	})

	return _c
}

func Add(ctx context.Context) error {
	return nil
}
