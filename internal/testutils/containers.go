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

package testutils

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

const (
	DBImage                 = "cassandra:4"
	TestNetworkName         = "testnet"
	DBContainerHost         = "test-db"
	TemporalImage           = "temporalio/auto-setup:1.20.0"
	TemporalContainerHost   = "test-temporal"
	NatsIOImage             = "nats:2.9.15"
	NatsIOContainerHost     = "test-natsio"
	AirImage                = "cosmtrek/air:v1.42.0"
	APIContainerHost        = "test-api"
	MothershipContainerHost = "test-mothership"
)

type (
	// Container is a wrapper around testcontainers.Container.
	Container struct {
		testcontainers.Container
		Context context.Context
		Request testcontainers.ContainerRequest
	}

	// ContainerEnvironment is a map of environment variables for a container.
	ContainerEnvironment map[string]string
)

// CreateTestNetwork creates a test network for testing purposes.
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
		"JVM_EXTRA_OPTS":         "-Dcassandra.skip_wait_for_gossip_to_settle=0",
	}

	_, caller, _, _ := runtime.Caller(0)
	pkgroot := path.Join(path.Dir(caller), "..", "..")
	cassandrayaml := path.Join(pkgroot, "deploy", "cassandra", "cassandra.yaml")

	mounts := testcontainers.ContainerMounts{
		testcontainers.ContainerMount{
			Source: testcontainers.GenericVolumeMountSource{Name: "test-db"},
			Target: "/var/lib/cassandra",
		},
		testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{HostPath: cassandrayaml},
			Target: "/etc/cassandra/cassandra.yaml",
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
	pkgroot := path.Join(path.Dir(caller), "..", "..")
	dynamicconfigpath := path.Join(pkgroot, "deploy", "temporal", "dynamicconfig")
	mounts := testcontainers.ContainerMounts{
		testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{HostPath: dynamicconfigpath},
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

// StartNatsIOContainer starts a NatsIO container for testing purposes.
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

// StartAPIContainer starts the API container for testing purposes.
func StartAPIContainer(ctx context.Context, secret string) (*Container, error) {
	return StartAirContainer(
		ctx,
		APIContainerHost,
		secret,
		"/api",
		"api.toml",
		"8000/tcp",
		wait.ForListeningPort("8000/tcp").WithPollInterval(time.Second*5).WithStartupTimeout(time.Minute*5),
	)
}

// StartMothershipContainer starts the Mothership container for testing purposes.
func StartMothershipContainer(ctx context.Context, secret string) (*Container, error) {
	return StartAirContainer(
		ctx,
		MothershipContainerHost,
		secret,
		"/mothership",
		"mothership.toml",
		"8080/tcp",
		wait.ForLog("Started Worker").WithPollInterval(time.Second*5).WithStartupTimeout(time.Minute*5),
	)
}

// StartAirContainer sets up comsmtrek/air to quickly compile different containers for services written in go.
func StartAirContainer(ctx context.Context, name string, secret string, workdir string, config string, port string, waiting wait.Strategy) (*Container, error) {
	env := ContainerEnvironment{
		"DEBUG":              "true",
		"SECRET":             secret,
		"EVENTS_SERVERS_URL": fmt.Sprintf("nats://%s:4222", NatsIOContainerHost),
		"TEMPORAL_HOST":      TemporalContainerHost,
		"CASSANDRA_HOSTS":    DBContainerHost,
		"CASSANDRA_KEYSPACE": db.TestKeyspace,
		"GOCOVERDIR":         ".coverage",
		"GOEXPERIMENT":       "coverageredesign",
	}
	_, caller, _, _ := runtime.Caller(0)
	pkgroot := path.Join(path.Dir(caller), "..", "..")
	mounts := testcontainers.ContainerMounts{
		testcontainers.ContainerMount{
			Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("%s-test--go-pkg-mod", name)},
			Target: "/go/pkg/mod",
		},
		testcontainers.ContainerMount{
			Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("%s-test--go-build-cache", name)},
			Target: "/root/.cache/go-build",
		},
		testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{HostPath: pkgroot},
			Target: testcontainers.ContainerMountTarget(workdir),
		},
		testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{HostPath: path.Join(pkgroot, "deploy", "air", config)},
			Target: testcontainers.ContainerMountTarget(path.Join(workdir, ".air.toml")),
		},
	}

	req := testcontainers.ContainerRequest{
		Name:           name,
		Hostname:       name,
		Image:          AirImage,
		Mounts:         mounts,
		Env:            env,
		Networks:       []string{TestNetworkName},
		ExposedPorts:   []string{port},
		WaitingFor:     waiting,
		ConfigModifier: func(config *container.Config) { config.WorkingDir = workdir },
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

// ShutdownCassandra first disables gossip and binary protocols, then drains the node, and finally terminates the container.
func (d *Container) ShutdownCassandra() error {
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "disablegossip"})
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "disablebinary"})
	_, _, _ = d.Exec(d.Context, []string{"nodetool", "drain"})

	return d.Container.Terminate(d.Context)
}

// Shutdown gracefully shuts down the container.
func (d *Container) Shutdown() error {
	timeout := 5 * time.Second
	return d.Container.Stop(d.Context, &timeout)
}
