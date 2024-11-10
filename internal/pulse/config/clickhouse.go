package config

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type (
	Clickhouse struct {
		Host     string `json:"host" koanf:"HOST"` // Database host.
		Port     int    `json:"port" koanf:"PORT"` // Database port.
		User     string `json:"user" koanf:"USER"` // Database user.
		Password string `json:"pass" koanf:"PASS"` // Database password.
		Name     string `json:"name" koanf:"NAME"` // Database name.

		conn driver.Conn
		once *sync.Once
	}

	ClickhouseOption func(*Clickhouse)
)

var (
	DefaultClickhouseConfig = Clickhouse{
		Host:     "localhost",
		Port:     6666,
		User:     "ctrlplane",
		Password: "ctrlplane",
		Name:     "ctrlplane",

		once: &sync.Once{},
	}
)

func (c *Clickhouse) connect(ctx context.Context) error {
	slog.Info("pulse/clickhose: connecting clickhouse ...", "host", c.Host, "port", c.Port, "user", c.User, "name", c.Name)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{c.GetAddress()},
		Auth: clickhouse.Auth{
			Username: c.User,
			Password: c.Password,
			Database: c.Name,
		},
	})
	if err != nil {
		return err
	}

	if err := conn.Ping(ctx); err != nil {
		return err
	}

	c.conn = conn

	slog.Info("pulse/clickhose: clickhouse connected.")

	return nil
}

func (c *Clickhouse) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Start establishes a connection to the ClickHouse database.
func (c *Clickhouse) Start(ctx context.Context) error {
	var err error

	c.once.Do(func() {
		err = c.connect(ctx)
	})

	return err
}

// Stop closes the connection to the ClickHouse database.
func (c *Clickhouse) Stop(_ context.Context) error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func WithClickhouseHost(host string) ClickhouseOption {
	return func(c *Clickhouse) {
		c.Host = host
	}
}

func WithClickhousePort(port int) ClickhouseOption {
	return func(c *Clickhouse) {
		c.Port = port
	}
}

func WithClickhouseUser(user string) ClickhouseOption {
	return func(c *Clickhouse) {
		c.User = user
	}
}

func WithClickhousePassword(password string) ClickhouseOption {
	return func(c *Clickhouse) {
		c.Password = password
	}
}

func WithClickhouseName(name string) ClickhouseOption {
	return func(c *Clickhouse) {
		c.Name = name
	}
}

func WithClickhouseConfig(cfg *Clickhouse) ClickhouseOption {
	return func(c *Clickhouse) {
		c.Host = cfg.Host
		c.Port = cfg.Port
		c.User = cfg.User
		c.Password = cfg.Password
		c.Name = cfg.Name
	}
}

func NewClickhouse(opts ...ClickhouseOption) *Clickhouse {
	cfg := &DefaultClickhouseConfig

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
