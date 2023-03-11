package testutils

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.breu.io/ctrlplane/internal/shared"
)

const (
	DBImage               = "cassandra:4"
	TestNetworkName       = "testnet"
	DBContainerHost       = "test-db"
	TemporalImage         = "temporalio/auto-setup:1.20.0"
	TemporalContainerHost = "test-temporal"
	NatsIOImage           = "nats:2.9.15"
	NatsIOContainerHost   = "test-natsio"
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

func CreateTestNetwork(ctx context.Context) (testcontainers.Network, error) {
	req := testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Name: TestNetworkName, CheckDuplicate: true},
	}
	network, err := testcontainers.GenericNetwork(ctx, req)
	if err != nil {
		return nil, err
	}
	return network, nil
}

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
		Name:         DBContainerHost,
		Hostname:     DBContainerHost,
		Image:        DBImage,
		Mounts:       mounts,
		Env:          env,
		Networks:     []string{TestNetworkName},
		ExposedPorts: []string{"9042/tcp"},
		WaitingFor:   wait.ForListeningPort("9042/tcp").WithPollInterval(time.Second * 5).WithStartupTimeout(time.Minute * 5),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Logger:           shared.Logger,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, err
	}

	return &Container{Container: ctr, Context: ctx, Request: req}, nil
}

// StartTemporalContainer starts a Temporal container for testing purposes.
func StartTemporalContainer(ctx context.Context) (*Container, error) {
	env := ContainerEnvironment{
		"CASSANDRA_SEEDS":          DBContainerHost,
		"DYNAMIC_CONFIG_FILE_PATH": "config/dynamicconfig/development-cass.yaml",
	}

	_, caller, _, _ := runtime.Caller(0)
	hostpath := path.Join(path.Dir(caller), "..", "..", "deploy", "temporal", "dynamicconfig")
	shared.Logger.Info("Mounting volume", "path", hostpath)

	mounts := testcontainers.ContainerMounts{
		testcontainers.ContainerMount{
			// NOTE: We are assuming that we are running the tests from the root of the project.
			Source: testcontainers.GenericBindMountSource{HostPath: hostpath},
			Target: "/etc/temporal/config/dynamicconfig",
		},
	}

	req := testcontainers.ContainerRequest{
		Name:         TemporalContainerHost,
		Hostname:     TemporalContainerHost,
		Image:        TemporalImage,
		Mounts:       mounts,
		Env:          env,
		Networks:     []string{TestNetworkName},
		ExposedPorts: []string{"7233/tcp", "7234/tcp", "7239/tcp"},
		WaitingFor:   wait.ForListeningPort("7233/tcp").WithPollInterval(time.Second * 5).WithStartupTimeout(time.Minute * 5),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Logger:           shared.Logger,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, err
	}

	return &Container{Container: ctr, Context: ctx, Request: req}, nil
}

func StartNatsIOContainer(ctx context.Context) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Name:         NatsIOContainerHost,
		Hostname:     NatsIOContainerHost,
		Image:        NatsIOImage,
		Networks:     []string{TestNetworkName},
		ExposedPorts: []string{"4222/tcp"},
		WaitingFor:   wait.ForLog("Server is ready"),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Logger:           shared.Logger,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, err
	}

	return &Container{Container: ctr, Context: ctx, Request: req}, nil
}

// RunCQL runs a CQL statement against the Cassandra container.
func (d *Container) RunCQL(stmt string) error {
	cmd := []string{"cqlsh", "-e", stmt}
	_, _, err := d.Exec(d.Context, cmd)
	return err
}

// CreateKeyspace creates a keyspace in the Cassandra container.
func (d *Container) CreateKeyspace(keyspace string) error {
	stmt := fmt.Sprintf("create keyspace if not exists %s with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};", keyspace)
	return d.RunCQL(stmt)
}

// DropKeyspace drops a keyspace in the Cassandra container.
func (d *Container) DropKeyspace(keyspace string) error {
	stmt := fmt.Sprintf("drop keyspace if exists %s;", keyspace)
	return d.RunCQL(stmt)
}

// ShutdownCassandra shuts down the Cassandra container.
func (d *Container) ShutdownCassandra() error {
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "disablegossip"})
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "disablebinary"})
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "drain"})
	return d.Terminate(d.Context)
}

func (d *Container) Shutdown() error {
	return d.Terminate(d.Context)
}
