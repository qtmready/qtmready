// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.  

package shared

import (
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/client"
)

var Temporal = &temporal{
	Queues: Queues{
		MothershipQueue:   &queue{MothershipQueue, "ai.ctrlplane.mothership"},
		IntegrationsQueue: &queue{IntegrationsQueue, "ai.ctrlplane.drivers"},
	},
}

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
	MothershipQueue   QueueName = "mothership"
	IntegrationsQueue QueueName = "drivers"
	BuilderQueue      QueueName = "builder"
	ProvisionerQueue  QueueName = "provisioner"
	DeployerQueue     QueueName = "deployer"
)

func (q QueueName) ToString() string {
	return string(q)
}

type (
	queue struct {
		Name   QueueName // The name of the queue.
		Prefix string    // The prefix to create the workflow ID.
	}

	temporal struct {
		ServerHost string `env:"TEMPORAL_HOST" env-default:"temporal"`
		ServerPort string `env:"TEMPORAL_PORT" env-default:"7233"`
		Client     client.Client
		Queues     Queues
	}
)

// GetName gets the name as string from QueueName
func (q *queue) GetName() string {
	return q.Name.ToString()
}

// CreateWorkflowID creates the unique workflow ID from the workflow sender and appropriate arguments.
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
		Logger.Error("Failed to read environment variables", "error", err)
	}
}

func (t *temporal) GetConnectionString() string {
	return t.ServerHost + ":" + t.ServerPort
}

func (t *temporal) InitClient() {
	Logger.Info("Initializing Temporal Client ...", "host", t.ServerHost, "port", t.ServerPort)
	options := client.Options{
		HostPort: t.GetConnectionString(),
		Logger:   Logger,
	}

	retryTemporal := func() error {
		clt, err := client.Dial(options)
		if err != nil {
			return err
		}

		t.Client = clt
		Logger.Info("Initializing Temporal Client ... Done")
		return nil
	}

	if err := retry.Do(
		retryTemporal,
		retry.Attempts(10),
		retry.Delay(1*time.Second),
	); err != nil {
		Logger.Error("Failed to initialize Temporal Client", "error", err)
	}
}
