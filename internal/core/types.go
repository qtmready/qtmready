package core

import "go.breu.io/ctrlplane/internal/shared"

// Workflow signal types
const (
	WorkflowSignalLockAcquired shared.WorkflowSignal = "lock_acquired"
	WorkflowSignalRequestLock  shared.WorkflowSignal = "request_lock"
	WorkflowSignalReleaseLock  shared.WorkflowSignal = "release_lock"
)
