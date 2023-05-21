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
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"

	"go.breu.io/ctrlplane/internal/shared/queue"
)

type (
	// Temporal holds the temporal server host and port, the client and all the available queues.
	//
	// TODO: The current design is not be ideal for a central multi-tenant solution due to the lack of strong isolation
	// for each tenant. For complaince, e.g. GDPR, SOC2, ISO27001, HIPAA, etc. we require strong tennant isolation. As
	// temporal.io provides strong namespace isolation, we can leverage this feature to implement a new design where
	// the client.Client field is replaced with a map of client.Client organized by tenant ID. A thread-safe method should
	// be added to the temporal struct to instantiate and return the appropriate client for a specific tenant. For
	// subsequent requests, the already instantiated client should be returned. This would allow for separate clients to
	// be created for each tenant, providing strong isolation and meeting compliance requirements.
	Temporal struct {
		ServerHost string `env:"TEMPORAL_HOST" env-default:"temporal"`
		ServerPort string `env:"TEMPORAL_PORT" env-default:"7233"`
		Client     client.Client
		Queues     queue.Queues
		Logger     log.Logger
	}

	TemporalOption func(*Temporal)
)

func (t *Temporal) GetConnectionString() string {
	return fmt.Sprintf("%s:%s", t.ServerHost, t.ServerPort)
}

// WithQueue adds a new queue to the Temporal.
func WithQueue(name queue.Name) TemporalOption {
	return func(t *Temporal) {
		t.Queues[name] = queue.NewQueue(queue.WithName(name))
	}
}

// WithLogger sets the logger for the Temporal.
func WithLogger(logger log.Logger) TemporalOption {
	return func(t *Temporal) {
		t.Logger = logger
	}
}

// WithConfigFromEnv reads the environment variables.
func WithConfigFromEnv() TemporalOption {
	return func(t *Temporal) {
		if err := cleanenv.ReadEnv(t); err != nil {
			panic(fmt.Errorf("failed to read environment variables: %w", err))
		}
	}
}

// WithClientConnection initializes the Temporal client.
func WithClientConnection() TemporalOption {
	return func(t *Temporal) {
		t.Logger.Info("temporal: connecting ...", "host", t.ServerHost, "port", t.ServerPort)

		options := client.Options{HostPort: t.GetConnectionString(), Logger: t.Logger}
		retryTemporal := func() error {
			clt, err := client.Dial(options)
			if err != nil {
				return err
			}

			t.Client = clt

			t.Logger.Info("temporal: connected")

			return nil
		}

		if err := retry.Do(
			retryTemporal,
			retry.Attempts(10),
			retry.Delay(1*time.Second),
			retry.OnRetry(func(n uint, err error) {
				t.Logger.Info("temporal: failed to connect. retrying connection ...", "attempt", n, "error", err)
			}),
		); err != nil {
			panic(fmt.Errorf("failed to connect to temporal: %w", err))
		}
	}
}

// NewTemporal creates a new Temporal instance.
func NewTemporal(opts ...TemporalOption) *Temporal {
	t := &Temporal{Queues: make(queue.Queues)}
	for _, opt := range opts {
		opt(t)
	}

	return t
}
