package pulse

import (
	"log/slog"
	"sync"

	"go.breu.io/quantm/internal/pulse/config"
)

type (
	Config = config.Config
	Option = config.Option
)

var (
	DefaultConfig = config.DefaultConfig

	_c   *Config
	once sync.Once
)

func WithConfig(cfg *Config) Option {
	return config.WithConfig(cfg)
}

func Instance(opts ...Option) *Config {
	once.Do(func() {
		slog.Info("pulse: configuring ...")

		_c = config.New(opts...)
	})

	return _c
}
