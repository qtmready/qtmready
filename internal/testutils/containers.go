package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	TestNetwork = "testnet-io"
	DBImage     = "cassandra:4"
)

type (
	DatabaseContainer struct {
		testcontainers.Container
		Context context.Context
		Request testcontainers.ContainerRequest
	}

	TemporalContainer    struct{}
	ContainerEnvironment map[string]string
)

// StartDBContainer starts a Cassandra container for testing purposes.
func StartDBContainer(ctx context.Context) (*DatabaseContainer, error) {
	env := ContainerEnvironment{
		"CASSANDRA_CLUSTER_NAME": "ctrlplane_test",
	}

	mnts := testcontainers.ContainerMounts{
		testcontainers.ContainerMount{
			Source: testcontainers.GenericVolumeMountSource{Name: "test-db"},
			Target: "/var/lib/cassandra",
		},
	}

	req := testcontainers.ContainerRequest{
		Name:         "test-db",
		Hostname:     "database",
		Image:        DBImage,
		Mounts:       mnts,
		Env:          env,
		Networks:     []string{TestNetwork},
		ExposedPorts: []string{"9042/tcp"},
		// WaitingFor:   wait.ForListeningPort("9042/tcp").WithStartupTimeout(time.Minute * 5),
		WaitingFor: wait.ForListeningPort("9042/tcp").WithPollInterval(time.Second * 5).WithStartupTimeout(time.Minute * 5),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, err
	}

	return &DatabaseContainer{Container: ctr, Context: ctx, Request: req}, nil
}

func (d *DatabaseContainer) RunCQL(stmt string) error {
	fmt.Printf("Running CQL: %s\n", stmt)
	cmd := []string{"cqlsh", "-e", stmt}
	_, _, err := d.Exec(d.Context, cmd)
	return err
}

func (d *DatabaseContainer) CreateKeyspace(keyspace string) error {
	stmt := fmt.Sprintf("create keyspace if not exists %s with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};", keyspace)
	stmt = fmt.Sprintf("%q", stmt)
	return d.RunCQL(stmt)
}

func (d *DatabaseContainer) DropKeyspace(keyspace string) error {
	stmt := fmt.Sprintf("drop keyspace if exists %s;", keyspace)
	return d.RunCQL(stmt)
}

func (d *DatabaseContainer) Stop() error {
	return d.Terminate(d.Context)
}
