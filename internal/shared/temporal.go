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

package shared

import (
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"

	tmp "go.temporal.io/sdk/temporal"
)

var (
	Temporal = &temporal{
		Queues: Queues{
			CoreQueue:      &queue{CoreQueue, "ai.ctrlplane.core"},
			ProvidersQueue: &queue{ProvidersQueue, "ai.ctrlplane.providers"},
		},

		WorkflowTools: workflowTools{},
	}
)

type (
	QueueName string

	Queue interface {
		// CreateWorkflowID creates an idempotency key. Sometimes we need to signal the workflow from a completely
		// disconnected part of the application. For us, it is important to arrive at the same workflow ID regardless
		// of the conditions.
		CreateWorkflowID(sender string, args ...string) string

		// GetWorkflowOptions returns the workflow options for the queue.
		GetWorkflowOptions(sender string, args ...string) client.StartWorkflowOptions

		// GetName gets the name of the queue as string.
		GetName() string

		// GetChildWorkflowOptions gets the child workflow options.
		GetChildWorkflowOptions(sender string, args ...string) workflow.ChildWorkflowOptions
	}

	Queues map[QueueName]Queue
)

const (
	CoreQueue      QueueName = "core"      // core queue
	ProvidersQueue QueueName = "providers" // messaging related to providers
	MutexQueue     QueueName = "mutex"     // mutex workflow queue
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

	// workflowTools holds the helper methods for ctrlplane workflows. TODO: See how it evolves and comeup with a better solution.
	workflowTools struct{}

	// temporal holds the temporal server host and port, the client and all the available queues.
	//
	// TODO: The current design is not be ideal for a central multi-tenant solution due to the lack of strong isolation
	// for each tenant. For complaince, e.g. GDPR, SOC2, ISO27001, HIPAA, etc. we require strong tennant isolation. As
	// temporal.io provides strong namespace isolation, we can leverage this feature to implement a new design where
	// the client.Client field is replaced with a map of client.Client organized by tenant ID. A thread-safe method should
	// be added to the temporal struct to instantiate and return the appropriate client for a specific tenant. For
	// subsequent requests, the already instantiated client should be returned. This would allow for separate clients to
	// be created for each tenant, providing strong isolation and meeting compliance requirements.
	temporal struct {
		ServerHost    string `env:"TEMPORAL_HOST" env-default:"temporal"`
		ServerPort    string `env:"TEMPORAL_PORT" env-default:"7233"`
		Client        client.Client
		Queues        Queues
		WorkflowTools workflowTools
	}
)

func (w *workflowTools) GetStackWorkflowName(stackID string) string {
	return Temporal.Queues[CoreQueue].CreateWorkflowID("stack", stackID)
}

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
	opts := client.StartWorkflowOptions{
		ID:        q.CreateWorkflowID(sender, args...),
		TaskQueue: q.GetName(),
		// WorkflowIDReusePolicy: 3, // client.WorkflowIDReusePolicyRejectDuplicate
	}
	retryPolicy := &tmp.RetryPolicy{MaximumAttempts: WorkflowMaxAttempts}
	opts.RetryPolicy = retryPolicy
	return opts
}

func (q *queue) GetChildWorkflowOptions(sender string, args ...string) workflow.ChildWorkflowOptions {
	opts := workflow.ChildWorkflowOptions{
		WorkflowID: q.CreateWorkflowID(sender, args...),
	}

	retryPolicy := &tmp.RetryPolicy{MaximumAttempts: WorkflowMaxAttempts}
	opts.RetryPolicy = retryPolicy
	return opts
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
