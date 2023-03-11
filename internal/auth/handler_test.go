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
	containers struct {
		network    testcontainers.Network
		db         *testutils.Container
		temporal   *testutils.Container
		nats       *testutils.Container
		api        *testutils.Container
		mothership *testutils.Container
	}
)

func (c *containers) shutdown(ctx context.Context) {
	shared.Logger.Info("graceful shutdown test environment ...")
	db.DB.Session.Close()
	_ = c.temporal.Shutdown()
	_ = c.nats.Shutdown()
	_ = c.api.Shutdown()
	_ = c.db.ShutdownCassandra()
	_ = c.network.Remove(ctx)
	shared.Logger.Info("graceful shutdown complete.")
}

func TestHandler(t *testing.T) {
	ctx := context.Background()
	ctrs := setup(ctx, t)

	t.Cleanup(func() {
		ctrs.shutdown(ctx)
	})
}

func setup(ctx context.Context, t *testing.T) *containers {
	shared.InitForTest()
	shared.Logger.Info("setting up test environment ...")
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

	apictr, err := testutils.StartAPIContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start api container: %v", err)
	}

	dbhost, _ := dbctr.Container.ContainerIP(ctx)
	temporalhost, _ := temporalctr.Container.ContainerIP(ctx)
	natshost, _ := natsctr.Container.ContainerIP(ctx)
	apihost, _ := apictr.Container.ContainerIP(ctx)

	shared.Logger.Info("hosts ...", "db", dbhost, "temporal", temporalhost, "nats", natshost, "api", apihost)

	return &containers{
		network:  network,
		db:       dbctr,
		temporal: temporalctr,
		nats:     natsctr,
		api:      apictr,
	}
}
