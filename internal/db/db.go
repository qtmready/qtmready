// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/Guilospanck/gocqlxmock"
	"github.com/Guilospanck/igocqlx"
	"github.com/avast/retry-go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	"github.com/ilyakaznacheev/cleanenv"

	"go.breu.io/quantm/internal/shared"

	_ "github.com/golang-migrate/migrate/v4/source/file" // required for file:// migrations
)

const (
	NullUUID     = "00000000-0000-0000-0000-000000000000"
	NullString   = ""
	TestKeyspace = "ctrlplane_test"
)

var (
	db   *Config
	once sync.Once
)

type (
	// Config holds the information about the database.
	Config struct {
		Session            igocqlx.ISessionx
		Hosts              []string      `env:"CASSANDRA_HOSTS" env-default:"database"`
		Port               int           `env:"CASSANDRA_PORT" env-default:"9042"`
		User               string        `env:"CASSANDRA_USER" env-default:""`
		Password           string        `env:"CASSANDRA_PASSWORD" env-default:""`
		Keyspace           string        `env:"CASSANDRA_KEYSPACE" env-default:"ctrlplane"`
		MigrationSourceURL string        `env:"CASSANDRA_MIGRATION_SOURCE_URL"`
		Timeout            time.Duration `env:"CASSANDRA_TIMEOUT" env-default:"1m"`
	}

	// MockConfig represents the mock session.
	MockConfig struct {
		*gocqlxmock.SessionxMock
	}

	ConfigOption func(*Config)
)

func (mc *MockConfig) Session() *igocqlx.Session {
	return nil
}

// WithHosts sets the hosts.
func WithHosts(hosts []string) ConfigOption {
	return func(c *Config) { c.Hosts = hosts }
}

// WithPort sets the port.
func WithPort(port int) ConfigOption {
	return func(c *Config) { c.Port = port }
}

// WithKeyspace sets the keyspace.
func WithKeyspace(keyspace string) ConfigOption {
	return func(c *Config) { c.Keyspace = keyspace }
}

// WithMigrationSourceURL sets the migration source URL.
func WithMigrationSourceURL(url string) ConfigOption {
	return func(c *Config) { c.MigrationSourceURL = url }
}

// WithTimeout sets the timeout.
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) { c.Timeout = timeout }
}

// FromEnvironment reads the configuration from the environment.
func FromEnvironment() ConfigOption {
	return func(c *Config) {
		if err := cleanenv.ReadEnv(c); err != nil {
			panic(fmt.Errorf("db: unable to read environment variables, %v", err))
		}
	}
}

// WithSessionCreation initializes the session.
func WithSessionCreation() ConfigOption {
	return func(c *Config) {
		if c.Hosts == nil || c.Keyspace == "" {
			panic(fmt.Errorf("db: hosts & keyspace not set, please set them before initializing session"))
		}

		cluster := gocql.NewCluster(c.Hosts...)
		cluster.Keyspace = c.Keyspace
		cluster.Timeout = c.Timeout
		cluster.ConnectTimeout = c.Timeout

		if c.User != "" {
			shared.Logger().Debug("db: authenticating ...", "user", c.User, "password", c.Password)
			cluster.Authenticator = gocql.PasswordAuthenticator{
				Username: c.User,
				Password: c.Password,
			}
		}

		createSession := func() error {
			shared.Logger().Info("db: connecting ...", "hosts", c.Hosts, "keyspace", c.Keyspace)

			session, err := igocqlx.WrapSession(cluster.CreateSession())
			if err != nil {
				shared.Logger().Error("db: failed to connect", "error", err)
				return err
			}

			c.Session = session

			shared.Logger().Info("db: connected")

			return nil
		}

		if err := retry.Do(
			createSession,
			retry.Attempts(15),
			retry.Delay(6*time.Second),
		); err != nil {
			shared.Logger().Error("db: aborting ....", "error", err)
		}
	}
}

// WithE2ESession initializes the session for end-to-end tests.
//
// NOTE: It might appear that the client throws error as explained at [issue], which will eventially point to [gocql github],
// but IRL, it will work. This is a known issue with gocql and it's not a problem for us.
//
// [issue]: https://app.shortcut.com/ctrlplane/story/2509/migrate-testing-to-use-test-containers-instead-of-mocks#activity-2749
// [gocql github]: https://github.com/gocql/gocql/issues/575
func WithE2ESession() ConfigOption {
	return func(c *Config) {
		cluster := gocql.NewCluster(c.Hosts...)
		cluster.Keyspace = c.Keyspace
		cluster.Timeout = 10 * time.Minute
		cluster.ConnectTimeout = 10 * time.Minute
		cluster.Port = c.Port
		// NOTE: Workaround for https://github.com/gocql/gocql/issues/575#issuecomment-172124342
		cluster.IgnorePeerAddr = true
		cluster.DisableInitialHostLookup = true
		cluster.Events.DisableTopologyEvents = true
		cluster.Events.DisableNodeStatusEvents = true
		cluster.Events.DisableSchemaEvents = true

		session, err := igocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			panic(fmt.Errorf("db: failed to connect to test database, %v", err))
		}

		c.Session = session
	}
}

// WithMockSession sets the mock session.
func WithMockSession(session *gocqlxmock.SessionxMock) ConfigOption {
	return func(c *Config) {
		c.Session = &MockConfig{session}
	}
}

// WithMigrations runs the migrations after the session has been initialized.
func WithMigrations() ConfigOption {
	return func(c *Config) {
		if c.MigrationSourceURL == "" {
			panic(fmt.Errorf("db: migration source url not set, please set it before running migrations"))
		}

		shared.Logger().Info("db: running migrations ...", "source", c.MigrationSourceURL)

		logger := shared.Logger()
		config := &cassandra.Config{KeyspaceName: c.Keyspace, MultiStatementEnabled: true}

		driver, err := cassandra.WithInstance(c.Session.Session().S.Session, config)
		if err != nil {
			shared.Logger().Error("db: failed to initialize driver for migrations ...", "error", err)
		}

		migrations, err := migrate.NewWithDatabaseInstance(c.MigrationSourceURL, "cassandra", driver)
		if err != nil {
			shared.Logger().Error("db: failed to initialize migrations ...", "error", err)
		}

		migrations.Log = logger

		version, dirty, err := migrations.Version()
		if err == migrate.ErrNilVersion {
			shared.Logger().Info("db: no migrations have been run, initializing keyspace ...")
		} else if err != nil {
			shared.Logger().Error("db: failed to get migration version ...", "error", err)
			return
		}

		shared.Logger().Info("db: migrations version ...", "version", version, "dirty", dirty)

		err = migrations.Up()
		if err != nil && err != migrate.ErrNoChange {
			shared.Logger().Warn("db: failed to run migrations ...", "error", err)
		}

		if err == migrate.ErrNoChange {
			shared.Logger().Info("db: no new migrations to run")
		}

		shared.Logger().Info("db: migrations done")
	}
}

func WithValidator(name string, validator validator.Func) ConfigOption {
	return func(c *Config) {
		if err := shared.Validator().RegisterValidation(name, validator); err != nil {
			panic(fmt.Errorf("db: failed to register validator %s, %v", name, err))
		}
	}
}

// NewSession creates a new session with the given options.
func NewSession(opts ...ConfigOption) *Config {
	db = &Config{}
	for _, opt := range opts {
		opt(db)
	}

	return db
}

// NewE2ESession creates a new session for end-to-end tests.
func NewE2ESession(port int, migrationPath string) {
	db = NewSession(
		WithHosts([]string{"localhost"}),
		WithKeyspace("ctrlplane_test"),
		WithPort(port),
		WithMigrationSourceURL(migrationPath),
		WithTimeout(10*time.Minute),
		WithE2ESession(),
		WithMigrations(),
		WithValidator("db_unique", UniqueField),
	)
}

// NewMockSession creates a new mock session.
func NewMockSession(session *gocqlxmock.SessionxMock) {
	db = NewSession(WithMockSession(session))
}

// DB returns the singleton database session.
func DB() *Config {
	if db == nil {
		shared.Logger().Info("db: creating new session")
		once.Do(func() {
			db = NewSession(
				FromEnvironment(),
				WithSessionCreation(),
				WithValidator("db_unique", UniqueField),
			)
		})
	}

	return db
}
