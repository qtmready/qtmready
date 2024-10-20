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

	"go.breu.io/quantm/internal/db/entities"
)

var (
	_conf *config
	_qry  *entities.Queries
	once  sync.Once
)

type (
	config struct {
		Host      string `env:"DB__HOST" env-default:"db"`
		Name      string `env:"DB__NAME" env-default:"ctrlplane"`
		Port      int    `env:"DB__PORT" env-default:"5432"`
		User      string `env:"DB__USER" env-default:"postgres"`
		Password  string `env:"DB__PASS" env-default:"postgres"`
		EnableSSL bool   `env:"DB__ENABLE_SSL" env-default:"false"`

		conn *pgx.Conn
	}

	Option func(*config)
)

func (c *config) ConnectionString() string {
	ssl := "disable"
	if c.EnableSSL {
		ssl = "require"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, ssl,
	)
}

func (c *config) retryfn() error {
	conn, err := pgx.Connect(context.Background(), c.ConnectionString())

	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *config) connect() {
	if c.conn != nil {
		slog.Warn("db: already connected")
	}

	if c.Host == "" || c.Name == "" || c.User == "" {
		slog.Error("db: invalid configuration", "host", c.Host, "name", c.Name, "user", c.User)

		panic("db: invalid configuration")
	}

	slog.Info("db: connecting ...", "host", c.Host, "port", c.Port, "name", c.Name, "user", c.User, "ssl", c.EnableSSL)

	err := retry.Do(
		c.retryfn,
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

func WithConfigFromEnvironment() Option {
	return func(c *config) {
		if err := cleanenv.ReadEnv(c); err != nil {
			panic(fmt.Errorf("db: unable to read environment variables, %v", err))
		}
	}
}

func WithConnect() Option {
	return func(c *config) {
		c.connect()
	}
}

func Close() {
	if _conf != nil && _conf.conn != nil {
		slog.Info("db: closing connection ...")

		_conf.conn.Close(context.Background())
	}
}

func Queries() *entities.Queries {
	once.Do(func() {
		slog.Info("db: initializing queries ...")

		_conf, _ = NewConfig(
			WithConfigFromEnvironment(),
			WithConnect(),
		)

		_qry = entities.New(_conf.conn)
	})

	return _qry
}

func NewConfig(opts ...Option) (*config, error) {
	_conf = &config{}
	for _, opt := range opts {
		opt(_conf)
	}

	return _conf, nil
}
