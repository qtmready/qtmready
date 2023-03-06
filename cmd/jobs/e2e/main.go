package main

import (
	"context"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
	"go.breu.io/ctrlplane/internal/testutils"
)

func main() {
	ctx := context.Background()
	shared.InitForTest()
	dbcon, err := testutils.StartDBContainer(ctx)
	if err != nil {
		shared.Logger.Error("db: failed to connect", "error", err)
	}

	if err = dbcon.CreateKeyspace(db.TestKeyspace); err != nil {
		shared.Logger.Error("db: unable to create keyspace", "error", err)
	}

	port, err := dbcon.Container.MappedPort(context.Background(), "9042")
	if err != nil {
		shared.Logger.Error("db: unable to get mapped port", "error", err)
	}

	if err = db.DB.InitSessionForTests(port.Int(), "file://internal/db/migrations"); err != nil {
		shared.Logger.Error("db: unable to initialize session", "error", err)
	}
}
