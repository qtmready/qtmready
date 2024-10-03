// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2023, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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

package mutex

import (
	"time"

	"go.breu.io/quantm/internal/core/defs"
	"go.temporal.io/sdk/workflow"
)

const (
	DefaultTimeout = 0 * time.Minute // DefaultTimeout is the default timeout for the mutex.
)

const (
	WorkflowSignalPrepare        defs.Signal = "mutex__prepare"
	WorkflowSignalAcquire        defs.Signal = "mutex__acquire"
	WorkflowSignalLocked         defs.Signal = "mutex__locked"
	WorkflowSignalRelease        defs.Signal = "mutex__release"
	WorkflowSignalReleased       defs.Signal = "mutex__released"
	WorkflowSignalCleanup        defs.Signal = "mutex__cleanup"
	WorkflowSignalCleanupDone    defs.Signal = "mutex__cleanup_done"
	WorkflowSignalCleanupDoneAck defs.Signal = "mutex__cleanup_done_ack"
	WorkflowSignalShutDown       defs.Signal = "mutex__shutdown"
)

type (
	// Mutex defines the signature for the workflow mutex. This workflow is meant to control the access to a resource.
	Mutex interface {
		Prepare(ctx workflow.Context) error // Prepares the mutex for use.
		Acquire(ctx workflow.Context) error // Acquire aquires the lock.
		Release(ctx workflow.Context) error // Release releases the lock.
		Cleanup(ctx workflow.Context) error // Cleanup attempts to shutdown the mutex workflow, if it is no longer needed.
	}

	Option func(*Handler)

	// Handler is the Mutex handler.
	Handler struct {
		ResourceID string              `json:"resource_id"` // ResourceID identifies the resource being locked.
		Info       *workflow.Info      `json:"info"`        // Info holds the workflow info that requests the mutex.
		Execution  *workflow.Execution `json:"execution"`   // Info holds the workflow info that holds the mutex.
		Timeout    time.Duration       `json:"timeout"`     // Timeout sets the timeout, after which the lock is automatically released.
		logger     *MutexLogger
	}
)

// Prepare prepares the mutex for use by executing the PrepareMutexActivity. It validates the mutex configuration, sets
// up the activity options, and executes the activity. If successful, it sets the Execution field of the Handler.
//
// Usage:
//
//	mutex := New(ctx, WithResourceID("resource-id"))
//	err := mutex.Prepare(ctx)
//	if err != nil {
//	  // handle error
//	}
func (h *Handler) Prepare(ctx workflow.Context) error {
	if err := h.validate(); err != nil {
		h.logger.error(h.Info.WorkflowExecution.ID, "prepare", "validate error", err)
		return err
	}

	h.logger.info(h.Info.WorkflowExecution.ID, "prepare", "preparing mutex")

	opts := workflow.ActivityOptions{StartToCloseTimeout: h.Timeout}
	ctx = workflow.WithActivityOptions(ctx, opts)

	exe := &workflow.Execution{}
	if err := workflow.ExecuteActivity(ctx, PrepareMutexActivity, h).Get(ctx, exe); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "prepare", "Unable to prepare mutex", err)
		return NewPrepareMutexError(h.ResourceID)
	}

	h.Execution = exe

	h.logger.info(h.Info.WorkflowExecution.ID, "prepare", "mutex prepared", "id", h.Execution.ID)

	return nil
}

// Acquire attempts to acquire the lock by signaling the mutex workflow. It waits for the WorkflowSignalLocked signal to
// confirm acquisition.
//
// Usage:
//
//	err := mutex.Acquire(ctx)
//	if err != nil {
//	  // handle error
//	}
//	// Critical section - mutex is acquired
func (h *Handler) Acquire(ctx workflow.Context) error {
	h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "requesting lock")

	ok := true

	if err := workflow.
		SignalExternalWorkflow(ctx, h.Execution.ID, "", WorkflowSignalAcquire.String(), h).
		Get(ctx, nil); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "acquire", "Unable to request lock", err)
		return NewAcquireLockError(h.ResourceID)
	}

	h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "waiting for lock")
	workflow.GetSignalChannel(ctx, WorkflowSignalLocked.String()).Receive(ctx, &ok)
	h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "lock acquired")

	if ok {
		return nil
	}

	return NewAcquireLockError(h.ResourceID)
}

// Release signals the mutex workflow to release the lock. It waits for the WorkflowSignalReleased signal to confirm the
// release.
//
// Usage:
//
//	err := mutex.Release(ctx)
//	if err != nil {
//	  // handle error
//	}
func (h *Handler) Release(ctx workflow.Context) error {
	h.logger.info(h.Info.WorkflowExecution.ID, "release", "requesting release")

	orphan := false

	if err := workflow.
		SignalExternalWorkflow(ctx, h.Execution.ID, "", WorkflowSignalRelease.String(), h).
		Get(ctx, nil); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "release", "Unable to request release", err)
		return NewReleaseLockError(h.ResourceID)
	}

	h.logger.info(h.Info.WorkflowExecution.ID, "release", "waiting for release")
	workflow.GetSignalChannel(ctx, WorkflowSignalReleased.String()).Receive(ctx, &orphan)

	if orphan {
		h.logger.warn(h.Info.WorkflowExecution.ID, "release", "lock released, orphaned", nil)
	} else {
		h.logger.info(h.Info.WorkflowExecution.ID, "release", "lock released")
	}

	return nil
}

// Cleanup attempts to shut down the mutex workflow if it's no longer needed. It signals the mutex workflow and waits
// for confirmation of cleanup.
//
// Usage:
//
//	err := mutex.Cleanup(ctx)
//	if err != nil {
//	  // handle error
//	}
func (h *Handler) Cleanup(ctx workflow.Context) error {
	h.logger.info(h.Info.WorkflowExecution.ID, "cleanup", "requesting cleanup")

	shutdown := false

	if err := workflow.
		SignalExternalWorkflow(ctx, h.Execution.ID, "", WorkflowSignalCleanup.String(), h).
		Get(ctx, nil); err != nil {
		h.logger.error(h.Info.WorkflowExecution.ID, "cleanup", "Unable to clean up", err)
		return NewCleanupMutexError(h.ResourceID)
	}

	h.logger.info(h.Info.WorkflowExecution.ID, "cleanup", "waiting for cleanup")
	workflow.GetSignalChannel(ctx, WorkflowSignalCleanupDone.String()).Receive(ctx, &shutdown)

	if err := workflow.
		SignalExternalWorkflow(ctx, h.Execution.ID, "", WorkflowSignalCleanupDoneAck.String(), shutdown).
		Get(ctx, nil); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "cleanup", "Unable to acknowledge cleanup", err)
		return NewCleanupMutexError(h.ResourceID)
	}

	if shutdown {
		h.logger.info(h.Info.WorkflowExecution.ID, "cleanup", "cleanup done")
	} else {
		h.logger.warn(h.Info.WorkflowExecution.ID, "cleanup", "cleanup failed, mutex in use", nil)
	}

	return nil
}

// validate checks if the mutex is properly configured with a ResourceID and workflow Info.
//
// Usage: This method is called internally by other methods and typically doesn't need to be called directly.
func (h *Handler) validate() error {
	if h.ResourceID == "" {
		return ErrNoResourceID
	}

	if h.Info == nil {
		return ErrNilContext
	}

	return nil
}

// WithResourceID sets the resource ID for the mutex workflow. We start with the assumption that a valid resource ID
// will be provided. The lock must always be held against the ids of core entities e.g. Stack, Repo or Resource. and the
// format may look like ${entity_type}.${entity_id}.mutex
//
//   - entity type e.g stack, repo, resource
//   - entity id e.g. the database id.
//
// for some cases, this may be made easy by getting the id of the parent workflow info e.g. if we are running stack
// controller, we can get the stack controller id, which would be in the format "io.quantm.stack.${stack_id}" and then
// adding the "mutex" suffix. Alernatively this can be set explicitly as well as
// "io.quantm.repo.${repo_id}.branch.${branch_name}.mutex". This is the format that should be used when holding locks
// against specific resources like repos or artifacts or cloud resources. This is a judgement call. The goal is, we
// should be able to arrive at the lock id regardless of the context.
func WithResourceID(id string) Option {
	return func(m *Handler) {
		m.ResourceID = id
	}
}

// WithTimeout sets the timeout for the mutex workflow.
func WithTimeout(timeout time.Duration) Option {
	return func(m *Handler) {
		m.Timeout = timeout
	}
}

// New returns a new Mutex. It should always be called with WithResourceID option. If WithTimeout is not called, it
// defaults to DefaultTimeout.
//
// Usage:
//
//	m := mutex.New(
//	  ctx,
//	  mutex.WithResourceID("id"),
//	  mutex.WithTimeout(30*time.Minute), // Optional
//	)
//	if err := m.Prepare(ctx); err != nil {/*handle error*/}
//	if err := m.Acquire(ctx); err != nil {/*handle error*/}
//	if err := m.Release(ctx); err != nil {/*handle error*/}
//	if err := m.Cleanup(ctx); err != nil {/*handle error*/}
func New(ctx workflow.Context, opts ...Option) Mutex {
	h := &Handler{Timeout: DefaultTimeout}
	for _, opt := range opts {
		opt(h)
	}

	h.Info = workflow.GetInfo(ctx)
	h.logger = NewMutexHandlerLogger(ctx, h.ResourceID)

	h.logger.info(h.Info.WorkflowExecution.ID, "create", "creating new mutex")

	return h
}
