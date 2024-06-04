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

// Stack Signals.
const (
	StackSignalAssetsRerieved     shared.WorkflowSignal = "stack__assets_retrieved"
	StackSignalInfraProvisioned   shared.WorkflowSignal = "stack__infra_provisioned"
	StackSignalDeploymentComplete shared.WorkflowSignal = "stack__deployment_completed"
	StackSignalManualOverride     shared.WorkflowSignal = "stack__manual_override"
)

const (
	GettingAssets State = iota
	GotAssets
	ProvisioningInfra
	InfraProvisioned
	CreatingDeployment
)

// RepoIO payloads.
type (
	RepoIOGetLatestCommitPayload struct {
		RepoID     string `json:"repo_id"`
		BranchName string `json:"branch_name"`
	}

	RepoIODeployChangesetPayload struct {
		RepoID      string      `json:"repo_id"`
		ChangesetID *gocql.UUID `json:"changeset_id"`
	}

	RepoIOTagCommitPayload struct {
		RepoID     string `json:"repo_id"`
		CommitSHA  string `json:"commit_sha"`
		TagName    string `json:"tag_name"`
		TagMessage string `json:"tag_message"`
	}

	RepoIOCreateBranchPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoID         string       `json:"repo_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		Commit         string       `json:"target_commit"`
		BranchName     string       `json:"branch_name"`
	}

	RepoIODeleteBranchPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		BranchName     string       `json:"branch_name"`
	}

	RepoIOMergeBranchPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		BaseBranch     string       `json:"base_branch"`
		TargetBranch   string       `json:"target_branch"`
	}

	RepoIORebaseAndMergePayload struct {
		RepoOwner        string       `json:"repo_owner"`
		RepoName         string       `json:"repo_name"`
		TargetBranchName string       `json:"target_branch_name"`
		InstallationID   shared.Int64 `json:"installation_id"`
	}

	RepoIODetectChangePayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		DefaultBranch  string       `json:"default_branch"`
		TargetBranch   string       `json:"target_branch"`
	}

	RepoIOTriggerCIActionPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoOwner      string       `json:"repo_owner"`
		RepoName       string       `json:"repo_name"`
		TargetBranch   string       `json:"target_branch"`
	}

	RepoIOGetRepoTeamIDPayload struct {
		RepoID string `json:"repo_id"`
	}

	RepoIOGetAllRelevantActionsPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
	}

	RepoIOGetRepoByProviderIDPayload struct {
		ProviderID string `json:"provider_id"`
	}

	RepoIOUpdateRepoHasRarlyWarningPayload struct {
		ProviderID string `json:"provider_id"`
	}
)

// MessageIO payloads.
type (
	MessageIOSendStaleBranchMessagePayload struct {
		TeamID string        `json:"team_id"`
		Stale  *LatestCommit `json:"slate"`
	}

	MessageIOSendNumberOfLinesExceedMessagePayload struct {
		TeamID        string         `json:"team_id"`
		RepoName      string         `json:"repo_name"`
		BranchName    string         `json:"branch_name"`
		Threshold     int            `json:"threshold"`
		BranchChnages *BranchChanges `json:"branch_chnages"`
	}

	MessageIOSendMergeConflictsMessagePayload struct {
		TeamID string        `json:"team_id"`
		Merge  *LatestCommit `json:"merge"`
	}

	MessageIOCompleteOauthResponsePayload struct {
		Code string `json:"code"`
	}
)

type (
	Infra               map[gocql.UUID]CloudResource // Map of resource Name and provider
	JsonInfra           map[gocql.UUID][]byte
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
		OldInfra  JsonInfra
		NewInfra  JsonInfra
	}

	Deployments     map[gocql.UUID]*Deployment // deployments against a changesetID.
	ChangesetAssets map[gocql.UUID]*Assets     // assets against a changesetID.

	// Assets contains all the assets fetched from DB against a stack.
	Assets struct {
		Repos       []Repo     // stack repos
		Resources   []Resource // stack cloud resources
		Workloads   []Workload // stack workloads
		Blueprint   Blueprint  // stack blueprint
		ChangesetID gocql.UUID
		Infra       JsonInfra
	}

	GetAssetsPayload struct {
		StackID       string
		RepoID        gocql.UUID
		ChangeSetID   gocql.UUID
		Image         string
		ImageRegistry string
		Digest        string
	}
)

func NewAssets() *Assets {
	return &Assets{
		Repos:     make([]Repo, 0),
		Resources: make([]Resource, 0),
		Workloads: make([]Workload, 0),
		Infra:     make(JsonInfra),
	}
}

func NewDeployment() *Deployment {
	d := new(Deployment)
	d.NewInfra = make(JsonInfra)
	d.OldInfra = make(JsonInfra)
	// d.rwpair = make(map[string][]Workload)
	return d
}
