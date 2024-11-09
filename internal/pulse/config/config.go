package config

import (
	"context"
	"sync"
)

type (
	Config struct {
		Clickhouse *Clickhouse `koanf:"CH"`

		once *sync.Once
	}

	Option func(*Config)
)

var (
	DefaultConfig = Config{
		Clickhouse: &DefaultClickhouseConfig,

		once: &sync.Once{},
	}
)

func (c *Config) Start(ctx context.Context) error {
	var err error

	c.once.Do(func() {
		if c.Clickhouse == nil {
			c.Clickhouse = &DefaultClickhouseConfig
		}

		err = c.Clickhouse.Start(ctx)
	})

	return err
}

func (c *Config) Stop(ctx context.Context) error {
	if c.Clickhouse == nil {
		return nil
	}

	return c.Clickhouse.Stop(ctx)
}

func WithClickhouse(ch *Clickhouse) Option {
	return func(c *Config) {
		c.Clickhouse = ch
	}
}

func WithConfig(cfg *Config) Option {
	return func(c *Config) {
		c.Clickhouse = cfg.Clickhouse
	}
}

func Default() *Config {
	return &DefaultConfig
}

func New(opts ...Option) *Config {
	cfg := &Config{once: &sync.Once{}}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
