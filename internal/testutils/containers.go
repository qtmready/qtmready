package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	TestNetwork           = "testnet-io"
	DBImage               = "cassandra:4"
	DBContainerHost       = "database"
	TemporalImage         = "temporalio/auto-setup:1.20.0"
	TemporalContainerHost = "temporal"
)

type (
	Container struct {
		testcontainers.Container
		Context context.Context
		Request testcontainers.ContainerRequest
	}

	TemporalContainer struct {
		testcontainers.Container
		Context context.Context
		Request testcontainers.ContainerRequest
	}

	ContainerEnvironment map[string]string
)

// StartDBContainer starts a Cassandra container for testing purposes.
func StartDBContainer(ctx context.Context) (*Container, error) {
	env := ContainerEnvironment{
		"CASSANDRA_CLUSTER_NAME": "ctrlplane_test",
	}

	mounts := testcontainers.ContainerMounts{
		testcontainers.ContainerMount{
			Source: testcontainers.GenericVolumeMountSource{Name: "test-db"},
			Target: "/var/lib/cassandra",
		},
	}

	req := testcontainers.ContainerRequest{
		Name:         "test-db",
		Hostname:     DBContainerHost,
		Image:        DBImage,
		Mounts:       mounts,
		Env:          env,
		Networks:     []string{TestNetwork},
		ExposedPorts: []string{"9042/tcp"},
		WaitingFor:   wait.ForListeningPort("9042/tcp").WithPollInterval(time.Second * 5).WithStartupTimeout(time.Minute * 5),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, err
	}

	return &Container{Container: ctr, Context: ctx, Request: req}, nil
}

func (d *Container) RunCQL(stmt string) error {
	cmd := []string{"cqlsh", "-e", stmt}
	_, _, err := d.Exec(d.Context, cmd)
	return err
}

func (d *Container) CreateKeyspace(keyspace string) error {
	stmt := fmt.Sprintf("create keyspace if not exists %s with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};", keyspace)
	return d.RunCQL(stmt)
}

func (d *Container) DropKeyspace(keyspace string) error {
	stmt := fmt.Sprintf("drop keyspace if exists %s;", keyspace)
	return d.RunCQL(stmt)
}

func (d *Container) ShutdownCassandra() error {
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "disablegossip"})
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "disablebinary"})
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "drain"})
	return d.Terminate(d.Context)
}
