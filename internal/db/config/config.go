package config

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5"
)

var (
	_c     *connection // Global connection instance.
	_conce sync.Once   // Ensures connection initialization occurs only once.
)

type (
	// connection struct holds database connection parameters and the established connection.
	connection struct {
		Host      string `env:"DB__HOST" env-default:"db"`          // Database host.
		Name      string `env:"DB__NAME" env-default:"ctrlplane"`   // Database name.
		Port      int    `env:"DB__PORT" env-default:"5432"`        // Database port.
		User      string `env:"DB__USER" env-default:"postgres"`    // Database user.
		Password  string `env:"DB__PASS" env-default:"postgres"`    // Database password.
		EnableSSL bool   `env:"DB__ENABLE_SSL" env-default:"false"` // Enable SSL.

		conn *pgx.Conn // Database connection.
	}

	// Option defines functional options for connection.
	Option func(*connection)
)

// ConnectionString builds a connection string from connection parameters.
func (c *connection) ConnectionString() string {
	ssl := "disable"
	if c.EnableSSL {
		ssl = "require"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, ssl,
	)
}

// IsConnected checks if a database connection exists.
func (c *connection) IsConnected() bool {
	return c.conn != nil
}

// Connect establishes a database connection using retry logic.
//
// Panics if a connection cannot be established after multiple retries.
func (c *connection) Connect(ctx context.Context) {
	if c.conn != nil {
		slog.Warn("db: already connected")

		return
	}

	if c.Host == "" || c.Name == "" || c.User == "" {
		slog.Error("db: invalid configuration", "host", c.Host, "name", c.Name, "user", c.User)

		panic("db: invalid configuration")
	}

	slog.Info("db: connecting ...", "host", c.Host, "port", c.Port, "name", c.Name, "user", c.User, "ssl", c.EnableSSL)

	err := retry.Do(
		c.retryfn(ctx),
		retry.Attempts(10),
		retry.Delay(500*time.Millisecond),
		retry.OnRetry(func(count uint, err error) {
			slog.Warn(
				"db: error connecting, retrying ...",
				"remaining_attempts", 10-count,
				"host", c.Host,
				"port", c.Port,
				"name", c.Name,
				"user", c.User,
				"ssl", c.EnableSSL,
				"error", err.Error(),
			)
		}),
	)

	if err != nil {
		panic(fmt.Errorf("db: unable to connect, %v", err))
	}
}

// Ping checks the database connection health by sending a ping.
//
// Returns an error if the ping fails.
func (c *connection) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

// Close closes the database connection.
func (c *connection) Close(ctx context.Context) {
	c.conn.Close(ctx)
}

// retryfn returns a function that attempts to establish a database connection.
//
// This function is used internally by the `Connect` method for retry logic. The returned function returns an error if the connection fails.
func (c *connection) retryfn(ctx context.Context) func() error {
	return func() error {
		conn, err := pgx.Connect(ctx, c.ConnectionString())
		if err != nil {
			return err
		}

		c.conn = conn

		return nil
	}
}

// WithHost sets the database host.
func WithHost(host string) Option {
	return func(c *connection) {
		c.Host = host
	}
}

// WithPort sets the database port.
func WithPort(port int) Option {
	return func(c *connection) {
		c.Port = port
	}
}

// WithName sets the database name.
func WithName(name string) Option {
	return func(c *connection) {
		c.Name = name
	}
}

// WithUser sets the database user.
func WithUser(user string) Option {
	return func(c *connection) {
		c.User = user
	}
}

// WithPassword sets the database password.
func WithPassword(password string) Option {
	return func(c *connection) {
		c.Password = password
	}
}

// WithConfigFromEnvironment reads connection parameters from environment variables.
//
// Panics if environment variables cannot be read.
func WithConfigFromEnvironment() Option {
	return func(c *connection) {
		if err := cleanenv.ReadEnv(c); err != nil {
			panic(fmt.Errorf("db: unable to read environment variables, %v", err))
		}
	}
}

// Connection creates a new global connection instance with functional options.
//
// Uses `sync.Once` to ensure the connection is initialized only once.
func Connection(opts ...Option) *connection {
	_conce.Do(func() {
		_c = &connection{}

		for _, opt := range opts {
			opt(_c)
		}
	})

	return _c
}

// ConnectionFromEnvironment creates a new connection instance with values from environment variables.
func ConnectionFromEnvironment() *connection {
	return Connection(
		WithConfigFromEnvironment(),
	)
}
