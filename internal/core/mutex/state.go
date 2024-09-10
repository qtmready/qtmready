// Copyright © 2024, Breu, Inc. <info@breu.io>
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

// Package mutex provides a distributed mutex implementation for Temporal workflows.
package mutex

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

type (
	// MutexStatus represents the current state of the mutex.
	MutexStatus string

	// MutexState encapsulates the state of the mutex workflow.
	MutexState struct {
		Status  MutexStatus   `json:"status"`
		Handler *Handler      `json:"handler"`
		Pool    *Pool         `json:"pool"`
		Orphans *Pool         `json:"orphans"`
		Timeout time.Duration `json:"timeout"`
		Persist bool          `json:"persist"`

		mutex  workflow.Mutex
		logger *MutexLogger
	}
)

const (
	MutexStatusAcquiring MutexStatus = "mutex__acquiring"
	MutexStatusLocked    MutexStatus = "mutex__locked"
	MutexStatusReleasing MutexStatus = "mutex__releasing"
	MutexStatusReleased  MutexStatus = "mutex__released"
	MutexStatusTimeout   MutexStatus = "mutex__timeout"
)

const (
	WorkflowQueryState shared.WorkflowSignal = "query__mutex__state"
)

func (s *MutexState) set_query_state(ctx workflow.Context) error {
	return workflow.SetQueryHandler(ctx, WorkflowQueryState.String(), func() (*MutexState, error) {
		return s, nil
	})
}

// on_prepare handles the preparation of lock requests.
// This signal originates from a client attempting to prepare for lock acquisition.
func (s *MutexState) on_prepare(_ workflow.Context) func(workflow.Context) {
	return func(ctx workflow.Context) {
		for {
			rx := &Handler{}
			workflow.GetSignalChannel(ctx, WorkflowSignalPrepare.String()).Receive(ctx, rx)

			s.logger.info(rx.Info.WorkflowExecution.ID, "prepare", "init")
			s.Pool.add(ctx, rx.Info.WorkflowExecution.ID, rx.Timeout)
			s.logger.info(rx.Info.WorkflowExecution.ID, "prepare", "done")
		}
	}
}

// on_acquire handles the acquisition of locks.
// This signal originates from a client attempting to acquire the lock.
func (s *MutexState) on_acquire(ctx workflow.Context) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Handler{}
		channel.Receive(ctx, rx)

		s.logger.info(rx.Info.WorkflowExecution.ID, "acquire", "init")

		timeout, _ := s.Pool.get(rx.Info.WorkflowExecution.ID)
		s.set_handler(ctx, rx)
		s.set_timeout(ctx, timeout)

		_ = workflow.SignalExternalWorkflow(ctx, rx.Info.WorkflowExecution.ID, "", WorkflowSignalLocked.String(), true).Get(ctx, nil)

		s.logger.info(rx.Info.WorkflowExecution.ID, "acquire", "done")
	}
}

// on_release handles the release of locks.
// This signal originates from a client that currently holds the lock and wants to release it.
func (s *MutexState) on_release(ctx workflow.Context) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Handler{}
		channel.Receive(ctx, rx)

		s.logger.info(rx.Info.WorkflowExecution.ID, "release", "init")

		if rx.Info.WorkflowExecution.ID == s.Handler.Info.WorkflowExecution.ID {
			s.to_releasing(ctx)
			s.Pool.remove(ctx, s.Handler.Info.WorkflowExecution.ID)
			s.to_released(ctx)

			_ = workflow.
				SignalExternalWorkflow(ctx, s.Handler.Info.WorkflowExecution.ID, "", WorkflowSignalReleased.String(), true).
				Get(ctx, nil)

			s.logger.info(rx.Info.WorkflowExecution.ID, "release", "done")
		}
	}
}

// on_abort handles the timeout and abortion of locks.
// This is triggered internally when a lock timeout occurs.
func (s *MutexState) on_abort(ctx workflow.Context) shared.FutureHandler {
	return func(future workflow.Future) {
		s.logger.info(s.Handler.Info.WorkflowExecution.ID, "abort", "init")

		if s.Status == MutexStatusLocked && s.Status != MutexStatusReleasing && s.Timeout > 0 {
			s.Pool.remove(ctx, s.Handler.Info.WorkflowExecution.ID)
			s.Orphans.add(ctx, s.Handler.Info.WorkflowExecution.ID, s.Timeout)
			s.to_timeout(ctx)
			s.logger.info(s.Handler.Info.WorkflowExecution.ID, "abort", "done")
		}
	}
}

// on_cleanup handles the cleanup process.
// This signal originates from an external system or administrator initiating a cleanup.
func (s *MutexState) on_cleanup(_ workflow.Context, fn workflow.Settable) func(workflow.Context) {
	shutdown := false

	return func(ctx workflow.Context) {
		for !shutdown {
			rx := &Handler{}
			workflow.GetSignalChannel(ctx, WorkflowSignalCleanup.String()).Receive(ctx, rx)

			s.logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "init")

			if s.Pool.size() == 0 {
				fn.Set(rx, nil)

				shutdown = true

				s.logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "shutdown", "pool_size", s.Pool.size())
			} else {
				s.logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "continue", "pool_size", s.Pool.size())
			}

			_ = workflow.
				SignalExternalWorkflow(ctx, rx.Info.WorkflowExecution.ID, "", WorkflowSignalCleanupDone.String(), shutdown).
				Get(ctx, nil)

			workflow.GetSignalChannel(ctx, WorkflowSignalCleanupDoneAck.String()).Receive(ctx, nil)
		}
	}
}

// on_terminate handles the termination process.
// This is triggered internally when the workflow is being shut down.
func (s *MutexState) on_terminate(ctx workflow.Context) shared.FutureHandler {
	return func(future workflow.Future) {
		rx := &Handler{}
		_ = future.Get(ctx, rx)

		s.logger.info(rx.Info.WorkflowExecution.ID, "terminate", "init")
		s.stop_persisting(ctx)
		s.logger.info(rx.Info.WorkflowExecution.ID, "terminate", "done")
	}
}

// to_locked transitions the state to Locked.
func (s *MutexState) to_locked(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.Status = MutexStatusLocked
	s.logger.info(s.Handler.Info.WorkflowExecution.ID, "transition", "to locked")
}

// to_releasing transitions the state to Releasing.
func (s *MutexState) to_releasing(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.Status = MutexStatusReleasing
	s.logger.info(s.Handler.Info.WorkflowExecution.ID, "transition", "to releasing")
}

// to_released transitions the state to Released.
func (s *MutexState) to_released(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.Status = MutexStatusReleased
	s.logger.info(s.Handler.Info.WorkflowExecution.ID, "transition", "to released")
}

// to_timeout transitions the state to Timeout.
func (s *MutexState) to_timeout(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.Status = MutexStatusTimeout
	s.logger.info(s.Handler.Info.WorkflowExecution.ID, "transition", "to timeout")
}

// to_acquiring transitions the state to Acquiring.
func (s *MutexState) to_acquiring(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.Status = MutexStatusAcquiring
	s.logger.info(s.Handler.Info.WorkflowExecution.ID, "transition", "to acquiring")
}

// set_handler updates the current handler.
func (s *MutexState) set_handler(ctx workflow.Context, handler *Handler) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.Handler = handler
}

// set_timeout updates the current timeout.
func (s *MutexState) set_timeout(ctx workflow.Context, timeout time.Duration) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.Timeout = timeout
}

// stop_persisting sets the persist flag to false.
func (s *MutexState) stop_persisting(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.Persist = false
	s.logger.info(s.Handler.Info.WorkflowExecution.ID, "persist", "stopped")
}

// restore reinitializes the mutex and logger, and restores the Pool and Orphans.
// It should be called after deserializing a MutexState instance.
func (s *MutexState) restore(ctx workflow.Context) {
	s.mutex = workflow.NewMutex(ctx)
	s.logger = NewMutexControllerLogger(ctx, s.Handler.ResourceID)

	if s.Pool == nil {
		s.Pool = NewPool(ctx)
	} else {
		s.Pool.restore(ctx)
	}

	if s.Orphans == nil {
		s.Orphans = NewPool(ctx)
	} else {
		s.Orphans.restore(ctx)
	}
}

// NewMutexState creates a new MutexState instance.
func NewMutexState(ctx workflow.Context, handler *Handler) *MutexState {
	state := &MutexState{
		Status:  MutexStatusAcquiring,
		Handler: handler,
		Pool:    NewPool(ctx),
		Orphans: NewPool(ctx),
		Timeout: 0,
		Persist: true,
		mutex:   workflow.NewMutex(ctx),
		logger:  NewMutexControllerLogger(ctx, handler.ResourceID),
	}

	return state
}
