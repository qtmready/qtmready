package core

import (
	"github.com/gocql/gocql"
	"go.breu.io/ctrlplane/internal/shared"
)

// Workflow signal types.
const (
	// mutex workflow signals.
	WorkflowSignalLockAcquired shared.WorkflowSignal = "lock_acquired"
	WorkflowSignalRequestLock  shared.WorkflowSignal = "request_lock"
	WorkflowSignalReleaseLock  shared.WorkflowSignal = "release_lock"

	// PR workflow signals.
	WorkflowSignalAssetsRetrieved     shared.WorkflowSignal = "assets_retreived"
	WorkflowSignalInfraProvisioned    shared.WorkflowSignal = "infra_created"
	WorkflowSignalDeploymentCompleted shared.WorkflowSignal = "deployment_completed"
	WorkflowSignalManaulOverride      shared.WorkflowSignal = "manual_override"
)

const (
	GettingAssets State = iota
	GotAssets
	ProvisioningInfra
	InfraProvisioned
	CreatingDeployment
)

type (
	SlicedResult[T any] struct {
		Data []T `json:"data"`
	}

	ResourceData struct{}

	ChildWorkflowIDs struct {
		GetAssets      string
		ProvisionInfra string
		Deployment     string
	}

	State int64

	DeploymentData struct {
		State       State
		WorkflowIDs ChildWorkflowIDs
	}

	DeploymentsData map[gocql.UUID]*DeploymentData // changesetID and deploymentData map
	AssetsMap       map[int64]*Assets

	// Assets contains all the assets fetched from DB against a stack.
	Assets struct {
		Repos           []Repo
		Resources       []Resource
		Workloads       []Workload
		Blueprint       Blueprint
		ResourcesConfig []ResourceData
		ChangesetID     gocql.UUID
	}
)

func (a *Assets) Create() {
	a.Repos = make([]Repo, 0)
	a.Resources = make([]Resource, 0)
	a.Workloads = make([]Workload, 0)
	a.ResourcesConfig = make([]ResourceData, 0)

}
