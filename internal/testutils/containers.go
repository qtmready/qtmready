package testutils

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	TestNetwork = "testnet-io"
)

type (
	DatabaseContainer struct {
		testcontainers.Container
		Context context.Context
		Request testcontainers.ContainerRequest
	}

	TemporalContainer struct{}
)

// StartDBContainer starts a Cassandra container for testing purposes.
func StartDBContainer(ctx context.Context) (*DatabaseContainer, error) {
	req := testcontainers.ContainerRequest{
		Hostname:     "database",
		Image:        "cassandra:4",
		Networks:     []string{TestNetwork},
		ExposedPorts: []string{"9042/tcp"},
		WaitingFor:   wait.ForListeningPort("9042/tcp").WithStartupTimeout(time.Minute * 5),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            false,
	})
	if err != nil {
		return nil, err
	}

	return &DatabaseContainer{Container: ctr, Context: ctx, Request: req}, nil
}

func (d *DatabaseContainer) GetHost() (string, error) {
	return d.Host(d.Context)
}

func (d *DatabaseContainer) Stop() error {
	return d.Terminate(d.Context)
}

func (d *DatabaseContainer) Setup() error {
	return nil
}
