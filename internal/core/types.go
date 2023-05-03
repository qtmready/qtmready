package core

import "go.breu.io/ctrlplane/internal/shared"

// Workflow signal types
const (
	// mutex workflow signals
	WorkflowSignalLockAcquired shared.WorkflowSignal = "lock_acquired"
	WorkflowSignalRequestLock  shared.WorkflowSignal = "request_lock"
	WorkflowSignalReleaseLock  shared.WorkflowSignal = "release_lock"

	// PR workflow signals
	WorkflowSignalAssetsRetrieved     shared.WorkflowSignal = "assets_retreived"
	WorkflowSignalInfraProvisioned    shared.WorkflowSignal = "infra_created"
	WorkflowSignalDeploymentCompleted shared.WorkflowSignal = "deployment_completed"
	WorkflowSignalManaulOverride      shared.WorkflowSignal = "manual_override"
)
