package config

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-playground/validator/v10"
)

type (
	// Config encapsulates configuration and connection management for a ClickHouse database.
	Config struct {
		Host     string `json:"host" koanf:"HOST" validate:"required"` // Database host address.
		Port     int    `json:"port" koanf:"PORT" validate:"required"` // Database port number.
		User     string `json:"user" koanf:"USER" validate:"required"` // Database username.
		Password string `json:"pass" koanf:"PASS" validate:"required"` // Database password.
		Name     string `json:"name" koanf:"NAME" validate:"required"` // Database name.

		conn driver.Conn // Established database connection.
		once *sync.Once  // Ensures single connection initialization.
	}

	// Option provides a functional option for customizing Clickhouse configurations.
	Option func(*Config)
)

var (
	// DefaultConfig defines the default configuration for connecting to a ClickHouse database.
	DefaultConfig = Config{
		Host:     "localhost", // Default host is localhost.
		Port:     6666,        // Default port is 9000.  Native ClickHouse port.
		User:     "ctrlplane", // Default username.
		Password: "ctrlplane", // Default password.
		Name:     "ctrlplane", // Default database name.

		once: &sync.Once{}, // Guarantees single connection attempt.
	}
)

func (c *Config) Validate() error {
	v := validator.New()

	return v.Struct(c)
}

// connect establishes a connection to the ClickHouse database.  The function attempts to connect to ClickHouse using the
// instance's configuration parameters.  Includes a ping to verify the connection's health.  Returns an error if the connection
// cannot be established or the ping fails. The context allows for connection timeout and cancellation.
func (c *Config) connect(ctx context.Context) error {
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

// GetAddress formats the ClickHouse server address as "host:port".
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Connection returns the established ClickHouse database connection.
func (c *Config) Connection() driver.Conn {
	return c.conn
}

// Start initiates a connection to the ClickHouse database.  Uses a sync.Once to ensure the connection is established only
// once, even with concurrent calls.  The provided context allows for cancellation or timeout during connection establishment.
// Returns an error from the connect function.
func (c *Config) Start(ctx context.Context) error {
	var err error

	c.once.Do(func() {
		err = c.connect(ctx)
	})

	return err
}

// Stop closes the existing ClickHouse database connection gracefully.  Checks for a nil connection to avoid potential
// panics. Returns any error encountered while closing the connection. The context is not utilized in the current
// implementation, but remains for potential future enhancements (e.g., connection draining).
func (c *Config) Stop(_ context.Context) error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

// WithHost sets the host address for the ClickHouse connection.
func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

// WithPort sets the port number for the ClickHouse connection.
func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithUser sets the username for the ClickHouse connection.
func WithUser(user string) Option {
	return func(c *Config) {
		c.User = user
	}
}

// WithPassword sets the password for the ClickHouse connection.
func WithPassword(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

// WithName sets the database name for the ClickHouse connection.
func WithName(name string) Option {
	return func(c *Config) {
		c.Name = name
	}
}

// WithConfig applies a given Clickhouse configuration.
func WithConfig(cfg *Config) Option {
	return func(c *Config) {
		c.Host = cfg.Host
		c.Port = cfg.Port
		c.User = cfg.User
		c.Password = cfg.Password
		c.Name = cfg.Name
	}
}

// New creates a new Clickhouse instance with the provided options.  Applies the functional options to
// customize the default configuration. Returns a pointer to the newly created Clickhouse instance.
func New(opts ...Option) *Config {
	cfg := &DefaultConfig

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
