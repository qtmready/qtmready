// Copyright © 2023, Breu, Inc. <info@breu.io>
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

package temporal

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"

	"go.breu.io/quantm/internal/shared/queue"
)

var (
	once sync.Once
)

type (
	Temporal interface {
		GetConnectionString() string
		Queue(queue.Name) queue.Queue
		Client() client.Client
		Worker(queue.Name) worker.Worker
	}

	// Config holds the temporal server host and port, the client and all the available queues.
	//
	// TODO: The current design is not be ideal for a central multi-tenant solution due to the lack of strong isolation
	// for each tenant. For complaince, e.g. GDPR, SOC2, ISO27001, HIPAA, etc. we require strong tennant isolation. As
	// temporal.io provides strong namespace isolation, we can leverage this feature to implement a new design where
	// the client.Client field is replaced with a map of client.Client organized by tenant ID. A thread-safe method should
	// be added to the temporal struct to instantiate and return the appropriate client for a specific tenant. For
	// subsequent requests, the already instantiated client should be returned. This would allow for separate clients to
	// be created for each tenant, providing strong isolation and meeting compliance requirements.
	Config struct {
		ServerHost string `env:"TEMPORAL_HOST" env-default:"temporal"`
		ServerPort string `env:"TEMPORAL_PORT" env-default:"7233"`
		client     client.Client
		queues     queue.Queues
		logger     *slog.Logger
		workers    queue.Workers
	}

	ConfigOption func(*Config)
)

func (t *Config) GetConnectionString() string {
	return fmt.Sprintf("%s:%s", t.ServerHost, t.ServerPort)
}

func (t *Config) Queue(name queue.Name) queue.Queue {
	return t.queues[name]
}

func (t *Config) Worker(name queue.Name) worker.Worker {
	return t.workers[name]
}

func (t *Config) Client() client.Client {
	if t.client == nil {
		once.Do(func() {
			t.logger.Info("temporal instantiating ....")

			logger := log.NewStructuredLogger(t.logger)
			options := client.Options{HostPort: t.GetConnectionString(), Logger: logger}

			retryTemporal := func() error {
				clt, err := client.Dial(options)
				if err != nil {
					return err
				}

				t.client = clt

				t.logger.Info("temporal: connected")

				return nil
			}

			if err := retry.Do(
				retryTemporal,
				retry.Attempts(10),
				retry.Delay(1*time.Second),
				retry.OnRetry(func(n uint, err error) {
					t.logger.Info("temporal: failed to connect. retrying connection ...", "attempt", n, "error", err)
				}),
			); err != nil {
				panic(fmt.Errorf("failed to connect to temporal: %w", err))
			}
		})
	}

	return t.client
}

// WithQueue adds a new queue and worker to the Config.
func WithQueue(name queue.Name) ConfigOption {
	return func(t *Config) {
		options := worker.Options{
			OnFatalError:        func(err error) { t.logger.Error("Fatal error during worker execution", "error", err.Error()) },
			EnableSessionWorker: true,
		}

		t.queues[name] = queue.NewQueue(queue.WithName(name))
		t.workers[name] = worker.New(t.client, name.String(), options)
	}
}

// WithLogger sets the logger for the Config.
func WithLogger(logger *slog.Logger) ConfigOption {
	return func(t *Config) {
		t.logger = logger
	}
}

// FromEnvironment reads the environment variables.
func FromEnvironment() ConfigOption {
	return func(t *Config) {
		if err := cleanenv.ReadEnv(t); err != nil {
			panic(fmt.Errorf("failed to read environment variables: %w", err))
		}
	}
}

// WithClientCreation initializes the Temporal client.
func WithClientCreation() ConfigOption {
	return func(t *Config) {
		t.logger.Info("temporal: connecting ...", "host", t.ServerHost, "port", t.ServerPort)

		logger := log.NewStructuredLogger(t.logger)
		options := client.Options{HostPort: t.GetConnectionString(), Logger: logger}

		retryTemporal := func() error {
			clt, err := client.Dial(options)
			if err != nil {
				return err
			}

			t.client = clt

			t.logger.Info("temporal: connected")

			return nil
		}

		if err := retry.Do(
			retryTemporal,
			retry.Attempts(10),
			retry.Delay(1*time.Second),
			retry.OnRetry(func(n uint, err error) {
				t.logger.Info("temporal: failed to connect. retrying connection ...", "attempt", n, "error", err)
			}),
		); err != nil {
			panic(fmt.Errorf("failed to connect to temporal: %w", err))
		}
	}
}

// New creates a new Temporal instance.
func New(opts ...ConfigOption) Temporal {
	t := &Config{queues: make(queue.Queues), workers: make(queue.Workers)}
	for _, opt := range opts {
		opt(t)
	}

	return t
}
