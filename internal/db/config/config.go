package config

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"

	"go.breu.io/quantm/internal/erratic"
)

var (
	_c     *Config   // Global connection instance.
	_conce sync.Once // Ensures connection initialization occurs only once.
)

type (
	// connection struct holds database connection parameters and the established connection.
	Config struct {
		Host      string `json:"host" koanf:"HOST" validate:"required"` // Database host.
		Name      string `json:"name" koanf:"NAME" validate:"required"` // Database name.
		Port      int    `json:"port" koanf:"PORT" validate:"required"` // Database port.
		User      string `json:"user" koanf:"USER" validate:"required"` // Database user.
		Password  string `json:"pass" koanf:"PASS" validate:"required"` // Database password.
		EnableSSL bool   `json:"enable_ssl" koanf:"ENABLE_SSL"`         // Enable SSL.

		conn *pgx.Conn // Database connection.
	}

	// ConfigOption defines functional options for connection.
	ConfigOption func(*Config)
)

var (
	// Default is the default database connection configuration.
	Default = Config{
		Host:      "localhost",
		Name:      "ctrlplane",
		Port:      5432,
		User:      "postgres",
		Password:  "postgres",
		EnableSSL: false,
	}
)

func (c *Config) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}

// ConnectionString builds a connection string from connection parameters.
func (c *Config) ConnectionString() string {
	ssl := "disable"
	if c.EnableSSL {
		ssl = "require"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, ssl,
	)
}

func (c *Config) ConnectionURI() string {
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
func (c *Config) IsConnected() bool {
	return c.conn != nil
}

// Start establishes a database connection using retry logic.
//
// Panics if a connection cannot be established after multiple retries.
func (c *Config) Start(ctx context.Context) error {
	if c.IsConnected() {
		slog.Warn("db: already connected.")

		return nil
	}

	if err := c.Validate(); err != nil {
		slog.Error("db: invalid configuration.", "error", err.Error())

		return erratic.NewConfigError(erratic.CommonModule, err.Error())
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
		return erratic.NewSystemError(erratic.CommonModule).Wrap(err)
	}

	slog.Info("db: connected.")

	return nil
}

// Ping checks the database connection health by sending a ping.
//
// Returns an error if the ping fails.
func (c *Config) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

func (c *Config) Get() *pgx.Conn {
	return c.conn
}

// Stop closes the database connection.
func (c *Config) Stop(ctx context.Context) error {
	if c.IsConnected() {
		c.conn.Close(ctx)
	} else {
		slog.Warn("db: already closed.")
	}

	return nil
}

// retryfn returns a function that attempts to establish a database connection.
//
// This function is used internally by the `Connect` method for retry logic. The returned function returns an error if the connection fails.
func (c *Config) retryfn(ctx context.Context) func() error {
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
func WithHost(host string) ConfigOption {
	return func(c *Config) {
		c.Host = host
	}
}

// WithPort sets the database port.
func WithPort(port int) ConfigOption {
	return func(c *Config) {
		c.Port = port
	}
}

// WithName sets the database name.
func WithName(name string) ConfigOption {
	return func(c *Config) {
		c.Name = name
	}
}

// WithUser sets the database user.
func WithUser(user string) ConfigOption {
	return func(c *Config) {
		c.User = user
	}
}

// WithPassword sets the database password.
func WithPassword(password string) ConfigOption {
	return func(c *Config) {
		c.Password = password
	}
}

func WithConfig(config *Config) ConfigOption {
	return func(c *Config) {
		if c.IsConnected() {
			slog.Warn("db: already connected.")

			return
		}

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
func WithConfigFromEnvironment(opts ...string) ConfigOption {
	return func(c *Config) {
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
		_ = k.Load(structs.Provider(Default, "__"), nil)

		if err := k.Load(env.Provider(prefix, "__", nil), nil); err != nil {
			panic(err)
		}

		if err := k.Unmarshal("", k); err != nil {
			panic(err)
		}
	}
}

// Instance creates a new global connection instance with functional Configoptions.
//
// Uses `sync.Once` to ensure the connection is initialized only once.
func Instance(opts ...ConfigOption) *Config {
	_conce.Do(func() {
		slog.Info("db: configuring connection ...")

		_c = &Config{}

		for _, opt := range opts {
			opt(_c)
		}

		slog.Info("db: connection configured.")
	})

	return _c
}
