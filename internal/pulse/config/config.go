package config

import (
	"context"
	"sync"
)

type (
	Config struct {
		Clickhouse *Clickhouse `koanf:"CH"`
		QuestDB    *QuestDB    `koanf:"QDB"`

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

		if c.QuestDB == nil {
			c.QuestDB = &DefaultQuestDBConfig
		}

		if cerr := c.Clickhouse.Start(ctx); cerr != nil {
			err = cerr
			return
		}

		if qerr := c.QuestDB.Start(ctx); qerr != nil {
			err = qerr
			return
		}
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
		c.QuestDB = cfg.QuestDB
	}
}

func New(opts ...Option) *Config {
	cfg := &Config{once: &sync.Once{}}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
