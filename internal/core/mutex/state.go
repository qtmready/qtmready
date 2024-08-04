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
		status  MutexStatus
		handler *Handler
		pool    *Pool
		orphans *Pool
		timeout time.Duration
		logger  *MutexLogger
		persist bool
		mutex   workflow.Mutex
	}
)

const (
	MutexStatusAcquiring MutexStatus = "acquiring"
	MutexStatusLocked    MutexStatus = "locked"
	MutexStatusReleasing MutexStatus = "releasing"
	MutexStatusReleased  MutexStatus = "released"
	MutexStatusTimeout   MutexStatus = "timeout"
)

// on_prepare handles the preparation of lock requests.
// This signal originates from a client attempting to prepare for lock acquisition.
func (s *MutexState) on_prepare(ctx workflow.Context) func(workflow.Context) {
	return func(ctx workflow.Context) {
		for {
			rx := &Handler{}
			workflow.GetSignalChannel(ctx, WorkflowSignalPrepare.String()).Receive(ctx, rx)

			s.logger.info(rx.Info.WorkflowExecution.ID, "prepare", "init")
			s.pool.add(ctx, rx.Info.WorkflowExecution.ID, rx.Timeout)
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

		timeout, _ := s.pool.get(rx.Info.WorkflowExecution.ID)
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

		if rx.Info.WorkflowExecution.ID == s.handler.Info.WorkflowExecution.ID {
			s.to_releasing(ctx)
			s.pool.remove(ctx, s.handler.Info.WorkflowExecution.ID)
			s.to_released(ctx)

			_ = workflow.SignalExternalWorkflow(ctx, s.handler.Info.WorkflowExecution.ID, "", WorkflowSignalReleased.String(), true).Get(ctx, nil)

			s.logger.info(rx.Info.WorkflowExecution.ID, "release", "done")
		}
	}
}

// on_abort handles the timeout and abortion of locks.
// This is triggered internally when a lock timeout occurs.
func (s *MutexState) on_abort(ctx workflow.Context) shared.FutureHandler {
	return func(future workflow.Future) {
		s.logger.info(s.handler.Info.WorkflowExecution.ID, "abort", "init")

		if s.status == MutexStatusLocked && s.status != MutexStatusReleasing && s.timeout > 0 {
			s.pool.remove(ctx, s.handler.Info.WorkflowExecution.ID)
			s.orphans.add(ctx, s.handler.Info.WorkflowExecution.ID, s.timeout)
			s.to_timeout(ctx)
			s.logger.info(s.handler.Info.WorkflowExecution.ID, "abort", "done")
		}
	}
}

// on_cleanup handles the cleanup process.
// This signal originates from an external system or administrator initiating a cleanup.
func (s *MutexState) on_cleanup(ctx workflow.Context, fn workflow.Settable) func(workflow.Context) {
	shutdown := false
	return func(ctx workflow.Context) {
		for !shutdown {
			rx := &Handler{}
			workflow.GetSignalChannel(ctx, WorkflowSignalCleanup.String()).Receive(ctx, rx)

			s.logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "init")

			if s.pool.size() == 0 {
				fn.Set(rx, nil)

				shutdown = true

				s.logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "shutdown", "pool_size", s.pool.size())
			} else {
				s.logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "abort", "pool_size", s.pool.size())
			}

			_ = workflow.SignalExternalWorkflow(ctx, rx.Info.WorkflowExecution.ID, "", WorkflowSignalCleanupDone.String(), false).Get(ctx, nil)
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

	s.status = MutexStatusLocked
	s.logger.info(s.handler.Info.WorkflowExecution.ID, "transition", "to Locked")
}

// to_releasing transitions the state to Releasing.
func (s *MutexState) to_releasing(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.status = MutexStatusReleasing
	s.logger.info(s.handler.Info.WorkflowExecution.ID, "transition", "to Releasing")
}

// to_released transitions the state to Released.
func (s *MutexState) to_released(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.status = MutexStatusReleased
	s.logger.info(s.handler.Info.WorkflowExecution.ID, "transition", "to Released")
}

// to_timeout transitions the state to Timeout.
func (s *MutexState) to_timeout(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.status = MutexStatusTimeout
	s.logger.info(s.handler.Info.WorkflowExecution.ID, "transition", "to Timeout")
}

// to_acquiring transitions the state to Acquiring.
func (s *MutexState) to_acquiring(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.status = MutexStatusAcquiring
	s.logger.info(s.handler.Info.WorkflowExecution.ID, "transition", "to Acquiring")
}

// set_handler updates the current handler.
func (s *MutexState) set_handler(ctx workflow.Context, handler *Handler) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.handler = handler
}

// set_timeout updates the current timeout.
func (s *MutexState) set_timeout(ctx workflow.Context, timeout time.Duration) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.timeout = timeout
}

// stop_persisting sets the persist flag to false.
func (s *MutexState) stop_persisting(ctx workflow.Context) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	s.persist = false
	s.logger.info(s.handler.Info.WorkflowExecution.ID, "persist", "stopped")
}

// NewMutexState creates a new MutexState instance.
func NewMutexState(ctx workflow.Context, starter *Handler) *MutexState {
	info := workflow.GetInfo(ctx)

	return &MutexState{
		status:  MutexStatusAcquiring,
		handler: starter,
		pool:    NewPool(ctx),
		orphans: NewPool(ctx),
		timeout: 0,
		logger:  NewMutexControllerLogger(ctx, info.WorkflowExecution.ID),
		persist: true,
		mutex:   workflow.NewMutex(ctx),
	}
}
