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
	"time"

	"github.com/Guilospanck/gocqlxmock"
	"github.com/Guilospanck/igocqlx"
	"github.com/avast/retry-go/v4"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	"github.com/ilyakaznacheev/cleanenv"

	"go.breu.io/ctrlplane/internal/shared"

	_ "github.com/golang-migrate/migrate/v4/source/file" // required for file:// migrations
)

const (
	NullUUID     = "00000000-0000-0000-0000-000000000000"
	NullString   = ""
	TestKeyspace = "ctrlplane_test"
)

var (
	DB = &db{} // represents the initialized gocql session.
)

type (
	// Holds the information about the database.
	db struct {
		Session            igocqlx.ISessionx
		Hosts              []string `env:"CASSANDRA_HOSTS" env-default:"database"`
		Keyspace           string   `env:"CASSANDRA_KEYSPACE" env-default:"ctrlplane"`
		MigrationSourceURL string   `env:"CASSANDRA_MIGRATION_SOURCE_URL"`
	}

	ms struct {
		*gocqlxmock.SessionxMock
	}
)

func (m *ms) Session() *igocqlx.Session {
	return nil
}

// ReadEnv reads the environment variables.
func (d *db) ReadEnv() {
	_ = cleanenv.ReadEnv(d)
}

// InitSession initializes the session with the configured hosts.
func (d *db) InitSession() {
	cluster := gocql.NewCluster(d.Hosts...)
	cluster.Keyspace = d.Keyspace
	createSession := func() error {
		shared.Logger.Info("db: connecting ...", "hosts", d.Hosts, "keyspace", d.Keyspace)

		session, err := igocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			shared.Logger.Error("db: failed to connect", "error", err)
			return err
		}

		d.Session = session

		shared.Logger.Info("db: connected")

		return nil
	}

	if err := retry.Do(
		createSession,
		retry.Attempts(15),
		retry.Delay(6*time.Second),
	); err != nil {
		shared.Logger.Error("db: aborting ....", "error", err)
	}
}

// RunMigrations runs database migrations if any.
func (d *db) RunMigrations() {
	shared.Logger.Info("db: running migrations ...", "source", d.MigrationSourceURL)

	config := &cassandra.Config{KeyspaceName: d.Keyspace, MultiStatementEnabled: true}
	driver, err := cassandra.WithInstance(d.Session.Session().S.Session, config)

	if err != nil {
		shared.Logger.Error("db: failed to initialize driver for migrations ...", "error", err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(d.MigrationSourceURL, "cassandra", driver)
	if err != nil {
		shared.Logger.Error("db: failed to initialize migrations ...", "error", err)
	}

	version, dirty, err := migrations.Version()
	if err == migrate.ErrNilVersion {
		shared.Logger.Info("db: no migrations have been run ...")
	}

	if dirty {
		shared.Logger.Warn("db: migration is dirty, forcing fix ...", "version", version)
		if err = migrations.Force(int(version) - 1); err != nil {
			shared.Logger.Error("db: failed to force migration ...", "error", err)
		}
	}
	err = migrations.Up()

	if err == migrate.ErrNoChange {
		shared.Logger.Info("db: no migrations to run")
	}

	shared.Logger.Info("db: migrations done")
}

// InitSessionWithMigrations is a shorthand for initializing the database along with running migrations.
func (d *db) InitSessionWithMigrations() {
	d.InitSession()
	d.RunMigrations()
}

// RegisterValidations registers any field or entity related validators.
func (d *db) RegisterValidations() {
	_ = shared.Validator.RegisterValidation("db_unique", UniqueField)
}

// InitMockSession initializes the session with the provided mock session.
func (d *db) InitMockSession(session *gocqlxmock.SessionxMock) {
	ms := &ms{session}
	d.Session = ms
}

// InitSessionForTests initializes the session with the configured hosts.
//
// NOTE: It might appear that the client throws error as explained at [issue], which will eventially point to [gocql github],
// but IRL, it will work. This is a known issue with gocql and it's not a problem for us.
//
// [issue]: https://app.shortcut.com/ctrlplane/story/2509/migrate-testing-to-use-test-containers-instead-of-mocks#activity-2749
// [gocql github]: https://github.com/gocql/gocql/issues/575
func (d *db) InitSessionForTests(port int, migrationsPath string) error {
	d.Hosts = []string{"localhost"}
	d.MigrationSourceURL = migrationsPath
	d.Keyspace = TestKeyspace
	cluster := gocql.NewCluster(d.Hosts...)
	// cluster.ProtoVersion = 4
	cluster.Keyspace = d.Keyspace
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Port = port
	// NOTE: Workaround for https://github.com/gocql/gocql/issues/575#issuecomment-172124342
	cluster.IgnorePeerAddr = true
	cluster.DisableInitialHostLookup = true
	cluster.Events.DisableTopologyEvents = true
	cluster.Events.DisableNodeStatusEvents = true
	cluster.Events.DisableSchemaEvents = true
	session, err := igocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return err
	}

	d.Session = session
	return nil
}
