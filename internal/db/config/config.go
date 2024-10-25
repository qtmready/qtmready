package config

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgx/v5"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"

	"go.breu.io/quantm/internal/db/status"
	"go.breu.io/quantm/internal/erratic"
)

var (
	_c     *Connection // Global connection instance.
	_conce sync.Once   // Ensures connection initialization occurs only once.
)

type (
	// connection struct holds database connection parameters and the established connection.
	Connection struct {
		Host      string `json:"host" koanf:"HOST"`             // Database host.
		Name      string `json:"name" koanf:"NAME"`             // Database name.
		Port      int    `json:"port" koanf:"PORT"`             // Database port.
		User      string `json:"user" koanf:"USER"`             // Database user.
		Password  string `json:"pass" koanf:"PASS"`             // Database password.
		EnableSSL bool   `json:"enable_ssl" koanf:"ENABLE_SSL"` // Enable SSL.

		conn *pgx.Conn // Database connection.
	}

	// Option defines functional options for connection.
	Option func(*Connection)
)

var (
	// DefaultConnection is the default database connection configuration.
	DefaultConnection = Connection{
		Host:      "localhost",
		Name:      "ctrlplane",
		Port:      5432,
		User:      "postgres",
		Password:  "postgres",
		EnableSSL: false,
	}
)

// ConnectionString builds a connection string from connection parameters.
func (c *Connection) ConnectionString() string {
	ssl := "disable"
	if c.EnableSSL {
		ssl = "require"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, ssl,
	)
}

func (c *Connection) ConnectionURI() string {
	ssl := "disable"
	if c.EnableSSL {
		ssl = "require"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		ssl,
	)
}

// IsConnected checks if a database connection exists.
func (c *Connection) IsConnected() bool {
	return c.conn != nil
}

// Start establishes a database connection using retry logic.
//
// Panics if a connection cannot be established after multiple retries.
func (c *Connection) Start(ctx context.Context) error {
	if c.IsConnected() {
		slog.Warn("db: already connected")

		return nil
	}

	if c.Host == "" || c.Name == "" || c.User == "" {
		slog.Error("db: invalid configuration", "host", c.Host, "name", c.Name, "user", c.User)

		return erratic.NewValidationError("reason", "database configuration is invalid", "host", c.Host, "name", c.Name, "user", c.User)
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
		return status.NewConnectionError().AddHint("error", err.Error())
	}

	slog.Info("db: connected")

	return nil
}

// Ping checks the database connection health by sending a ping.
//
// Returns an error if the ping fails.
func (c *Connection) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

// Stop closes the database connection.
func (c *Connection) Stop(ctx context.Context) error {
	if c.IsConnected() {
		c.conn.Close(ctx)
	} else {
		slog.Warn("db: already closed")
	}

	return nil
}

// retryfn returns a function that attempts to establish a database connection.
//
// This function is used internally by the `Connect` method for retry logic. The returned function returns an error if the connection fails.
func (c *Connection) retryfn(ctx context.Context) func() error {
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
	return func(c *Connection) {
		c.Host = host
	}
}

// WithPort sets the database port.
func WithPort(port int) Option {
	return func(c *Connection) {
		c.Port = port
	}
}

// WithName sets the database name.
func WithName(name string) Option {
	return func(c *Connection) {
		c.Name = name
	}
}

// WithUser sets the database user.
func WithUser(user string) Option {
	return func(c *Connection) {
		c.User = user
	}
}

// WithPassword sets the database password.
func WithPassword(password string) Option {
	return func(c *Connection) {
		c.Password = password
	}
}

func WithConfig(config *Connection) Option {
	return func(c *Connection) {
		c.Host = config.Host
		c.Port = config.Port
		c.Name = config.Name
		c.User = config.User
		c.Password = config.Password
		c.EnableSSL = config.EnableSSL
	}
}

// WithConfigFromEnvironment reads connection parameters from environment variables using koanf.
//
// Panics if environment variables cannot be read.
func WithConfigFromEnvironment(opts ...string) Option {
	return func(c *Connection) {
		var prefix string

		if len(opts) > 0 {
			prefix = strings.ToUpper(opts[0])

			if !strings.HasSuffix(prefix, "__") {
				prefix += "__"
			}
		} else {
			prefix = "DB__"
		}

		k := koanf.New("__")
		_ = k.Load(structs.Provider(DefaultConnection, "__"), nil)

		if err := k.Load(env.Provider(prefix, "__", nil), nil); err != nil {
			panic(err)
		}

		if err := k.Unmarshal("", k); err != nil {
			panic(err)
		}
	}
}

// Instance creates a new global connection instance with functional options.
//
// Uses `sync.Once` to ensure the connection is initialized only once.
func Instance(opts ...Option) *Connection {
	_conce.Do(func() {
		_c = &Connection{}

		for _, opt := range opts {
			opt(_c)
		}
	})

	slog.Info("db: instance created", "host", _c.Host, "port", _c.Port, "name", _c.Name, "user", _c.User, "ssl", _c.EnableSSL)
	return _c
}
