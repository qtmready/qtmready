package config

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	qdb "github.com/questdb/go-questdb-client/v3"
)

type (
	// QuestDB represents the configuration for a QuestDB instance.
	QuestDB struct {
		Host     string `json:"host" koanf:"HOST"` // Database host.
		Port     int    `json:"port" koanf:"PORT"` // Database port.
		User     string `json:"user" koanf:"USER"` // Database user.
		Password string `json:"pass" koanf:"PASS"` // Database password.

		pool *qdb.LineSenderPool
		once *sync.Once
	}
)

var (
	// DefaultQuestDBConfig holds the default configuration for a QuestDB instance.
	DefaultQuestDBConfig = QuestDB{
		Host:     "localhost",
		Port:     6060, // Default HTTP port for QuestDB
		User:     "ctrlplane",
		Password: "ctrlplane",

		once: &sync.Once{},
	}
)

// connect establishes a connection pool to the QuestDB database.
func (q *QuestDB) connect(_ context.Context) error {
	slog.Info("pusle/questdb: connecting questdb ...", "host", q.Host, "port", q.Port, "user", q.User)

	pool, err := qdb.PoolFromOptions(
		qdb.WithAddress(q.GetAddress()),
		qdb.WithTcp(),
	)
	if err != nil {
		return err
	}

	q.pool = pool

	slog.Info("pusle/questdb: questdb connection pool created.")

	return nil
}

func (q *QuestDB) GetAddress() string {
	return fmt.Sprintf("http::addr=%s:%d;username=%s;password=%s;", q.Host, q.Port, q.User, q.Password)
}

// Start establishes a connection to the QuestDB database.
func (q *QuestDB) Start(ctx context.Context) error {
	var err error

	q.once.Do(func() {
		err = q.connect(ctx)
	})

	return err
}

// Stop closes the connection to the QuestDB database.
func (q *QuestDB) Stop(ctx context.Context) error {
	if q.pool == nil {
		return nil
	}

	return q.pool.Close(ctx)
}

// Sender retrieves a sender from the connection pool.
func (q *QuestDB) Sender(ctx context.Context) (qdb.LineSender, error) {
	if q.pool == nil {
		return nil, fmt.Errorf("questdb connection pool not initialized")
	}

	return q.pool.Sender(ctx)
}

// WithQuestDBHost sets the host for the QuestDB instance.
func WithQuestDBHost(host string) func(*QuestDB) {
	return func(q *QuestDB) {
		q.Host = host
	}
}

// WithQuestDBPort sets the port for the QuestDB instance.
func WithQuestDBPort(port int) func(*QuestDB) {
	return func(q *QuestDB) {
		q.Port = port
	}
}

// WithQuestDBUser sets the user for the QuestDB instance.
func WithQuestDBUser(user string) func(*QuestDB) {
	return func(q *QuestDB) {
		q.User = user
	}
}

// WithQuestDBPassword sets the password for the QuestDB instance.
func WithQuestDBPassword(password string) func(*QuestDB) {
	return func(q *QuestDB) {
		q.Password = password
	}
}

// WithQuestDBConfig sets the configuration for the QuestDB instance.
func WithQuestDBConfig(cfg *QuestDB) func(*QuestDB) {
	return func(q *QuestDB) {
		q.Host = cfg.Host
		q.Port = cfg.Port
		q.User = cfg.User
		q.Password = cfg.Password
	}
}

// NewQuestDB creates a new QuestDB configuration with the given options.
func NewQuestDB(opts ...func(*QuestDB)) *QuestDB {
	cfg := &DefaultQuestDBConfig

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
