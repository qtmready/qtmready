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

package core

import (
	"github.com/gocql/gocql"

	"go.breu.io/quantm/internal/shared"
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

	ResourceConfig struct{}

	ChildWorkflows struct {
		GetAssets      string
		ProvisionInfra string
		Deploy         string
	}

	State int64

	Deployment struct {
		state     State
		workflows ChildWorkflows
	}

	Deployments     map[gocql.UUID]*Deployment // deployments against a changesetID.
	ChangesetAssets map[gocql.UUID]*Assets     // assets against a changesetID.

	// Assets contains all the assets fetched from DB against a stack.
	Assets struct {
		Repos           []Repo
		Resources       []Resource
		Workloads       []Workload
		Blueprint       Blueprint
		ResourcesConfig []ResourceConfig
		ChangesetID     gocql.UUID
	}
)

func NewAssets() *Assets {
	return &Assets{
		Repos:           make([]Repo, 0),
		Resources:       make([]Resource, 0),
		Workloads:       make([]Workload, 0),
		ResourcesConfig: make([]ResourceConfig, 0),
	}
}
