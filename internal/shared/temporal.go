// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
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

package shared

import (
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/client"
)

var (
	Temporal = &temporal{
		Queues: Queues{
			CoreQueue:      &queue{CoreQueue, "ai.ctrlplane.core"},
			ProvidersQueue: &queue{ProvidersQueue, "ai.ctrlplane.providers"},
		},
	}
)

type (
	QueueName string

	Queue interface {
		CreateWorkflowID(sender string, args ...string) string
		GetWorkflowOptions(sender string, args ...string) client.StartWorkflowOptions
		GetName() string
	}

	Queues map[QueueName]Queue
)

// TODO: The greater plan is to move each tenant in its own namespace.
const (
	CoreQueue      QueueName = "core"      // core queue.
	ProvidersQueue QueueName = "providers" // messaging related to providers
)

func (q QueueName) ToString() string {
	return string(q)
}

type (
	// queue holds the queue name and prefix for workflow id.
	queue struct {
		Name   QueueName // The name of the queue.
		Prefix string    // The prefix to create the workflow ID.
	}

	// temporal holds the temporal client and client.
	//
	// FIXME: The current design is not ideal for a central multi-tenannt solution. Temporal provides strong isolation via
	// namespaces. Ideally, each tenant should have its own namespace. That would require a change in the struct to have a
	// map. The map would be keyed by tenant ID and the value would be the temporal client. A Client(id string)
	// should either get the client from the map or create a new one (singleton) if it does not exist.
	temporal struct {
		ServerHost string `env:"TEMPORAL_HOST" env-default:"temporal"`
		ServerPort string `env:"TEMPORAL_PORT" env-default:"7233"`
		Client     client.Client
		Queues     Queues
	}
)

// GetName gets the name as string from QueueName.
func (q *queue) GetName() string {
	return q.Name.ToString()
}

// CreateWorkflowID creates the unique workflow ID from the workflow sender and appropriate arguments.
//
// TODO: document the grand scheme of things.
func (q *queue) CreateWorkflowID(sender string, args ...string) string {
	return q.Prefix + "." + sender + "." + strings.Join(args, ".")
}

func (q *queue) GetWorkflowOptions(sender string, args ...string) client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:        q.CreateWorkflowID(sender, args...),
		TaskQueue: q.GetName(),
		// WorkflowIDReusePolicy: 3, // client.WorkflowIDReusePolicyRejectDuplicate
	}
}

func (t *temporal) ReadEnv() {
	if err := cleanenv.ReadEnv(t); err != nil {
		Logger.Error("temporal: error", "error", err)
	}
}

func (t *temporal) GetConnectionString() string {
	return t.ServerHost + ":" + t.ServerPort
}

func (t *temporal) InitClient() {
	Logger.Info("temporal: connecting ...", "host", t.ServerHost, "port", t.ServerPort)

	options := client.Options{HostPort: t.GetConnectionString(), Logger: Logger}
	retryTemporal := func() error {
		clt, err := client.Dial(options)
		if err != nil {
			return err
		}

		t.Client = clt

		Logger.Info("temporal: connected")

		return nil
	}

	if err := retry.Do(retryTemporal, retry.Attempts(10), retry.Delay(1*time.Second)); err != nil {
		Logger.Error("temporal: connection error", "error", err)
	}
}
