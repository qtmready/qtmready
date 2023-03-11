package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
	"go.breu.io/ctrlplane/internal/testutils"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	network, dbctr, temporalctr, natsctr, err := setup(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		shared.Logger.Info("shutting down ...")
		time.Sleep(5 * time.Second) // TODO: remove this
		db.DB.Session.Close()
		_ = temporalctr.Shutdown()
		_ = natsctr.Shutdown()
		_ = dbctr.ShutdownCassandra()
		_ = network.Remove(ctx)
		shared.Logger.Info("Test done. Exiting...")
	})
}

func setup(ctx context.Context, t *testing.T) (testcontainers.Network, *testutils.Container, *testutils.Container, *testutils.Container, error) {
	shared.InitForTest()
	network, err := testutils.CreateTestNetwork(ctx)
	if err != nil {
		t.Fatalf("failed to create test network: %v", err)
	}
	dbctr, err := testutils.StartDBContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start db container: %v", err)
	}

	if err = dbctr.CreateKeyspace(db.TestKeyspace); err != nil {
		t.Fatalf("failed to create keyspace: %v", err)
	}

	port, err := dbctr.Container.MappedPort(context.Background(), "9042")
	if err != nil {
		t.Fatalf("failed to get mapped db port: %v", err)
	}

	err = db.DB.InitSessionForTests(port.Int(), "file://../db/migrations")
	shared.Logger.Warn("session gets initiated, but if we catch the error and do t.Fatal(err), the test panics!")
	if db.DB.Session.Session().S == nil {
		t.Fatal("session is nil")
	}

	db.DB.RunMigrations()

	temporalctr, err := testutils.StartTemporalContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start temporal container: %v", err)
	}

	natsctr, err := testutils.StartNatsIOContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start natsio container: %v", err)
	}

	dbhost, _ := dbctr.Container.ContainerIP(ctx)
	temporalhost, _ := temporalctr.Container.ContainerIP(ctx)
	natshost, _ := natsctr.Container.ContainerIP(ctx)

	shared.Logger.Info("hosts ...", "db", dbhost, "temporal", temporalhost, "nats", natshost)

	return network, dbctr, temporalctr, natsctr, nil
}
