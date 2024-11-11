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
		QuestDB:    &DefaultQuestDBConfig,

		once: &sync.Once{},
	}
)

// Start starts the Clickhouse and QuestDB clients
//
// It starts both clients in a goroutine, and returns an error
// if either client fails to start.
//
//	cfg := config.New(config.WithClickhouse(clickhouse.New("localhost:8123")))
//	if err := cfg.Start(context.Background()); err != nil {
//	    log.Fatal(err)
//	}
//	defer cfg.Stop(context.Background())
//
//	// use the Clickhouse and QuestDB clients
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

// Stop stops the Clickhouse client
//
// It returns an error if the Clickhouse client fails to stop.
func (c *Config) Stop(ctx context.Context) error {
	if c.Clickhouse == nil {
		return nil
	}

	return c.Clickhouse.Stop(ctx)
}

// WithClickhouse returns an Option that sets the Clickhouse client
//
// Example:
//
//	cfg := config.New(config.WithClickhouse(clickhouse.New("localhost:8123")))
func WithClickhouse(ch *Clickhouse) Option {
	return func(c *Config) {
		c.Clickhouse = ch
	}
}

// WithConfig returns an Option that sets the Clickhouse and QuestDB clients
//
// Example:
//
//	cfg := config.New(config.WithConfig(&config.Config{
//	    Clickhouse: clickhouse.New("localhost:8123"),
//	    QuestDB:    questdb.New("localhost:8888"),
//	}))
func WithConfig(cfg *Config) Option {
	return func(c *Config) {
		c.Clickhouse = cfg.Clickhouse
		c.QuestDB = cfg.QuestDB
	}
}

// New returns a new Config
//
// Example:
//
//	cfg := config.New()
//
//	cfg := config.New(config.WithClickhouse(clickhouse.New("localhost:8123")))
//
//	cfg := config.New(config.WithConfig(&config.Config{
//	    Clickhouse: clickhouse.New("localhost:8123"),
//	    QuestDB:    questdb.New("localhost:8888"),
//	}))
func New(opts ...Option) *Config {
	cfg := &Config{once: &sync.Once{}}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
