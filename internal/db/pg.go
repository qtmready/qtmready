package db

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
	_pg     *pg
	queries *entities.Queries
	once    sync.Once
)

type (
	pg struct {
		Host      string `env:"DB__HOST" env-default:"pg"`
		Name      string `env:"DB__NAME" env-default:"ctrlplane"`
		Port      int    `env:"DB__PORT" env-default:"5432"`
		User      string `env:"DB__USER" env-default:"postgres"`
		Password  string `env:"DB__PASS" env-default:"postgres"`
		EnableSSL bool   `env:"DB__ENABLE_SSL" env-default:"false"`

		conn *pgx.Conn
	}

	Option func(*pg)
)

func (c *pg) ConnectionString() string {
	ssl := "disable"
	if c.EnableSSL {
		ssl = "require"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, ssl,
	)
}

func (c *pg) retryfn() error {
	conn, err := pgx.Connect(context.Background(), c.ConnectionString())

	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *pg) connect() {
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
	return func(c *pg) {
		if err := cleanenv.ReadEnv(c); err != nil {
			panic(fmt.Errorf("db: unable to read environment variables, %v", err))
		}
	}
}

func WithConnect() Option {
	return func(c *pg) {
		c.connect()
	}
}

func Close() {
	if _pg != nil && _pg.conn != nil {
		slog.Info("db: closing connection ...")

		_pg.conn.Close(context.Background())
	}
}

func Queries() *entities.Queries {
	once.Do(func() {
		slog.Info("db: initializing queries ...")

		_pg, _ = NewConfig(
			WithConfigFromEnvironment(),
			WithConnect(),
		)

		queries = entities.New(_pg.conn)
	})

	return queries
}

func NewConfig(opts ...Option) (*pg, error) {
	_pg = &pg{}
	for _, opt := range opts {
		opt(_pg)
	}

	return _pg, nil
}
