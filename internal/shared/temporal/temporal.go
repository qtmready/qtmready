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

package temporal

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"go.breu.io/quantm/internal/shared/queue"
	"go.breu.io/slog-utils/calldepth"
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
			logger := calldepth.New(
				calldepth.WithLogger(t.logger),
				calldepth.WithCallDepth(6), // exactly pin points from where the task was called.
			).WithGroup("temporal")

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
		options := worker.Options{OnFatalError: func(err error) { t.logger.Error("Fatal error during worker execution %v", err) }}
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

		options := client.Options{HostPort: t.GetConnectionString(), Logger: t.logger}
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
