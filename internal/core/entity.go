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
	"encoding/json"
	"errors"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"

	"go.breu.io/ctrlplane/internal/db"
)

func (stack *Stack) PreCreate() error { stack.Slug = db.CreateSlug(stack.Name); return nil }
func (stack *Stack) PreUpdate() error { return nil }

func (repo *Repo) PreCreate() error { return nil }
func (repo *Repo) PreUpdate() error { return nil }

// TODO: move these entities to be generated by the code generator

var (
	workloadColumns = []string{
		"id",
		"stack_id",
		"repo_id",
		"repo_path",
		"name",
		"kind",
		"container",
		"builder",
		"created_at",
		"updated_at",
	}

	workloadMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "workloads",
			Columns: workloadColumns,
		},
	}

	workloadTable = itable.New(*workloadMeta.M)

	resourceColumns = []string{
		"id",
		"stack_id",
		"repo_id",
		"name",
		"provider",
		"driver",
		"is_immutable",
		"created_at",
		"updated_at",
	}

	resourceMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "resources",
			Columns: resourceColumns,
		},
	}

	resourceTable = itable.New(*resourceMeta.M)

	blueprintColumns = []string{
		"id",
		"stack_id",
		"repo_id",
		"name",
		"regions",
		"repo_branch",
		"rollout_budget",
		"created_at",
		"updated_at",
	}

	blueprintMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "blueprints",
			Columns: blueprintColumns,
		},
	}

	blueprintTable = itable.New(*blueprintMeta.M)

	changesetColums = []string{
		"id",
		"stack_id",
		"calver",
		"repo_markers",
		"created_by",
		"created_at",
		"updated_at",
	}

	changesetMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "changesets",
			Columns: changesetColums,
		},
	}

	changesetTable = itable.New(*changesetMeta.M)

	rolloutColumns = []string{
		"id",
		"stack_id",
		"blueprint_id",
		"trigger",
		"state",
		"created_at",
		"updated_at",
	}

	rolloutMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "rollouts",
			Columns: rolloutColumns,
		},
	}

	rolloutTable = itable.New(*rolloutMeta.M)
)

type (
	// Workload defines a workload for the app. See Workload.Kind for type.
	Workload struct {
		ID        gocql.UUID `json:"id" cql:"id"`
		StackID   gocql.UUID `json:"stack_id" cql:"stack_id"`
		RepoID    gocql.UUID `json:"repo_id" cql:"repo_id"`
		RepoPath  string     `json:"repo_path" cql:"repo_path"`
		Name      string     `json:"name" cql:"name"`
		Kind      string     `json:"kind" cql:"kind"`           // "default" | "worker" | "job" | "cronjob"
		Container string     `json:"container" cql:"container"` // json with keys: "image" | "command" | "environment" | "dependencies"
		Builder   string     `json:"builder" cql:"builder"`     // json with keys: "buildpack" | "dockerfile" | "script" | "external"
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}

	// Resource defines the cloud provider resources for the app e.g. s3, sqs, etc.
	Resource struct {
		ID          gocql.UUID `json:"id" cql:"id"`
		StackID     gocql.UUID `json:"stack_id" cql:"stack_id"`
		RepoID      gocql.UUID `json:"repo_id" cql:"repo_id"`
		Name        string     `json:"name" cql:"name"`
		Provider    string     `json:"provider" cql:"provider"` // "aws" | "gcp" | "azure"
		Driver      string     `json:"driver" cql:"driver"`     // "s3" | "sqs" | "sns" | "dynamodb" | "postgres" | "mysql" etc.
		IsImmutable bool       `json:"is_immutable" cql:"is_immutable"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
	}

	// Blueprint contains a collection of Workload & Resource to define one single release.
	Blueprint struct {
		ID            gocql.UUID       `json:"id" cql:"id"`
		StackID       gocql.UUID       `json:"stack_id" cql:"stack_id"`
		RepoID        gocql.UUID       `json:"repo_id" cql:"repo_id"`
		RepoBranch    string           `json:"repo_branch" cql:"repo_branch"`
		Name          string           `json:"name" cql:"name"`
		Regions       BluePrintRegions `json:"regions" cql:"regions"`
		RolloutBudget int              `json:"rollout_budget" cql:"rollout_budget"`
		CreatedAt     time.Time        `json:"created_at"`
		UpdatedAt     time.Time        `json:"updated_at"`
	}

	// ChangeSet records the state of the stack at a given point in time.
	// For a poly-repo BluePrint, a PR on one repo can trigger a release for the BluePrint.
	ChangeSet struct {
		ID          gocql.UUID           `json:"id" cql:"id"`
		StackID     gocql.UUID           `json:"stack_id" cql:"stack_id"`
		Calver      string               `json:"calver" cql:"calver"`
		RepoMarkers ChangeSetRepoMarkers `json:"repo_markers" cql:"repo_markers"`
		CreatedBy   string               `json:"created_by" cql:"created_by"`
		CreatedAt   time.Time            `json:"created_at"`
		UpdatedAt   time.Time            `json:"updated_at"`
	}

	Rollout struct {
		ID          gocql.UUID   `json:"id" cql:"id"`
		StackID     gocql.UUID   `json:"stack_id" cql:"stack_id"`
		BlueprintID gocql.UUID   `json:"blueprint_id" cql:"blueprint_id"`
		ChangeSetID gocql.UUID   `json:"changeset_id" cql:"changeset_id"`
		State       RolloutState `json:"state" cql:"state"` // "in_progress" | "live" | "rejected"
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
	}
)

func (workload *Workload) GetTable() itable.ITable { return workloadTable }
func (workload *Workload) PreCreate() error        { return nil }
func (workload *Workload) PreUpdate() error        { return nil }

func (resource *Resource) GetTable() itable.ITable { return resourceTable }
func (resource *Resource) PreCreate() error        { return nil }
func (resource *Resource) PreUpdate() error        { return nil }

func (blueprint *Blueprint) GetTable() itable.ITable { return blueprintTable }
func (blueprint *Blueprint) PreCreate() error        { return nil }
func (blueprint *Blueprint) PreUpdate() error        { return nil }

func (changeset *ChangeSet) GetTable() itable.ITable { return changesetTable }
func (changeset *ChangeSet) PreCreate() error        { return nil }
func (changeset *ChangeSet) PreUpdate() error        { return nil }

func (rollout *Rollout) GetTable() itable.ITable { return rolloutTable }
func (rollout *Rollout) PreCreate() error        { return nil }
func (rollout *Rollout) PreUpdate() error        { return nil }

type (

	// BluePrintRegions sets the cloud regions where a blueprint can be deployed.
	BluePrintRegions struct {
		GCP     []string `json:"gcp"`
		AWS     []string `json:"aws"`
		Azure   []string `json:"azure"`
		Default string   `json:"default"`
	}

	// RolloutState is the state of a rollout.
	RolloutState        string
	RolloutStateMapType map[string]RolloutState

	ChangeSetRepoMarker struct {
		Provider   string `json:"provider"`
		CommitID   string `json:"commit_id"`
		HasChanged bool   `json:"changed"`
	}

	ChangeSetRepoMarkers []ChangeSetRepoMarker
)

const (
	RolloutStateQueued     RolloutState = "queued"
	RolloutStateInProgress RolloutState = "in_progress"
	RolloutStateCompleted  RolloutState = "completed"
	RolloutStateRejected   RolloutState = "rejected"
)

var (
	RolloutStateMap = RolloutStateMapType{
		RolloutStateQueued.String():     RolloutStateQueued,
		RolloutStateInProgress.String(): RolloutStateInProgress,
		RolloutStateCompleted.String():  RolloutStateCompleted,
		RolloutStateRejected.String():   RolloutStateRejected,
	}
)

func (config StackConfig) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(config)
}

func (config *StackConfig) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, config)
}

func (regions BluePrintRegions) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(regions)
}

func (regions *BluePrintRegions) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, regions)
}

func (marker ChangeSetRepoMarkers) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(marker)
}

func (marker *ChangeSetRepoMarkers) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, marker)
}

func (rs RolloutState) String() string {
	return string(rs)
}

func (rs RolloutState) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.String())
}

func (rs *RolloutState) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val, ok := RolloutStateMap[s]
	if !ok {
		return errors.New("invalid rollout state")
	}

	*rs = val

	return nil
}

func (rs RolloutState) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(rs)
}

func (rs *RolloutState) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, rs)
}

func (rp RepoProvider) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(rp)
}

func (rp *RepoProvider) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, rp)
}
