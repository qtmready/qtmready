// Copyright © 2023, 2024, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package testutils

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.breu.io/quantm/internal/db"
)

const (
	DBImage                 = "scylladb/scylla:5.4"
	DBContainerHost         = "test-db"
	TemporalImage           = "temporalio/auto-setup:1.22.5"
	TemporalContainerHost   = "test-temporal"
	NatsIOImage             = "nats:2.9.15"
	NatsIOContainerHost     = "test-natsio"
	AirImage                = "cosmtrek/air:v1.50.0"
	APIContainerHost        = "test-api"
	MothershipContainerHost = "test-mothership"
)

var (
	TestNetworkName string
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
func CreateTestNetwork(ctx context.Context) (*testcontainers.DockerNetwork, error) {
	net, err := network.New(ctx)
	if err != nil {
		return nil, err
	}

	TestNetworkName = net.Name

	return net, nil
}

// StartDBContainer starts a Cassandra container for testing purposes.
func StartDBContainer(ctx context.Context) (*Container, error) {
	env := ContainerEnvironment{
		// "CASSANDRA_CLUSTER_NAME": "ctrlplane_test",
		// "JVM_EXTRA_OPTS":         "-Dcassandra.skip_wait_for_gossip_to_settle=0",
	}

	mounts := testcontainers.Mounts(
		testcontainers.VolumeMount("test-db", "/var/lib/scylla"),
	)

	req := testcontainers.ContainerRequest{
		Name:         DBContainerHost,
		Hostname:     DBContainerHost,
		Image:        DBImage,
		Mounts:       mounts,
		Env:          env,
		Cmd:          []string{"--smp", "1", "--skip-wait-for-gossip-to-settle", "1"},
		Networks:     []string{TestNetworkName},
		ExposedPorts: []string{"9042/tcp"},
		WaitingFor: wait.ForListeningPort("9042/tcp").
			WithPollInterval(time.Second * 5).
			WithStartupTimeout(time.Minute * 5),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		// Logger:           shared.Logger(),
		Started: true,
		Reuse:   true,
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

	req := testcontainers.ContainerRequest{
		Name:         TemporalContainerHost,
		Hostname:     TemporalContainerHost,
		Image:        TemporalImage,
		Env:          env,
		Networks:     []string{TestNetworkName},
		ExposedPorts: []string{"7233/tcp", "7234/tcp", "7239/tcp"},
		WaitingFor:   wait.ForListeningPort("7233/tcp").WithPollInterval(time.Second * 5).WithStartupTimeout(time.Minute * 5),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: dynamicconfigpath,
				Target: "/etc/temporal/config/dynamicconfig",
			})
		},
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		// Logger:           shared.Logger(),
		Started: true,
		Reuse:   true,
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
		// Logger:           shared.Logger(),
		Started: true,
		Reuse:   true,
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
func StartAirContainer(ctx context.Context, name, secret, workdir, config, port string, waiting wait.Strategy) (*Container, error) {
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

	req := testcontainers.ContainerRequest{
		Name:         name,
		Hostname:     name,
		Image:        AirImage,
		Env:          env,
		Networks:     []string{TestNetworkName},
		ExposedPorts: []string{port},
		WaitingFor:   waiting,
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      path.Join(pkgroot, "deploy", "air", config),
				ContainerFilePath: path.Join(workdir, ".air.toml"),
				FileMode:          0644,
			},
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts,
				mount.Mount{
					Type:   mount.TypeVolume,
					Source: fmt.Sprintf("%s-test--go-pkg-mod", name),
					Target: "/go/pkg/mod",
				},
				mount.Mount{
					Type:   mount.TypeVolume,
					Source: fmt.Sprintf("%s-test--go-build-cache", name),
					Target: "/root/.cache/go-build",
				},
				mount.Mount{
					Type:   mount.TypeBind,
					Source: pkgroot,
					Target: workdir,
				},
			)
		},
		ConfigModifier: func(config *container.Config) { config.WorkingDir = workdir },
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		// Logger:           shared.Logger(),
		Started: true,
		Reuse:   true,
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
	stmt := fmt.Sprintf(
		"create keyspace if not exists %s with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};",
		keyspace,
	)

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
