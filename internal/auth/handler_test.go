package auth_test

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
	"go.breu.io/ctrlplane/internal/testutils"
)

type (
	TestLogConsumer struct {
		Msgs []string
	}
)

func (t *TestLogConsumer) Accept(content testcontainers.Log) {
	t.Msgs = append(t.Msgs, string(content.Content))
}

func TestHandler(t *testing.T) {
	ctx := context.Background()
	shared.InitForTest()
	dbcon, err := testutils.StartDBContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err = dbcon.CreateKeyspace(db.TestKeyspace); err != nil {
		t.Fatal(err)
	}

	port, err := dbcon.Container.MappedPort(context.Background(), "9042")
	if err != nil {
		t.Fatal(err)
	}

	err = db.DB.InitSessionForTests(port.Int(), "file://../db/migrations")
	shared.Logger.Warn("session gets initiated, but if we catch the error and do t.Fatal(err), the test panics!")
	if db.DB.Session.Session().S == nil {
		t.Fatal("session is nil")
	}

	db.DB.RunMigrations()

	t.Cleanup(func() {
		db.DB.Session.Close()
		_ = dbcon.ShutdownCassandra()
	})
}
