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

package entity

import (
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"

	"go.breu.io/ctrlplane/internal/db"
)

var (
	stackColumns = []string{
		"id",
		"team_id",
		"name",
		"slug",
		"config",
		"created_at",
		"updated_at",
	}

	stackMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "stacks",
			Columns: stackColumns,
		}}

	stackTable = itable.New(*stackMeta.M)

	repoColumns = []string{
		"id",
		"stack_id",
		"provider_id",
		"default_branch",
		"is_monorepo",
		"provider",
		"created_at",
		"updated_at",
	}

	repoMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "repos",
			Columns: repoColumns,
		},
	}

	repoTable = itable.New(*repoMeta.M)

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
	// Stack defines a group of services that are deployed together.
	// The deployment blueprint for each stack is defined under Blueprint.
	Stack struct {
		ID        gocql.UUID  `json:"id" cql:"id"`
		TeamID    gocql.UUID  `json:"team_id" cql:"team_id"`
		Name      string      `json:"name" validate:"required"`
		Slug      string      `json:"slug"`
		Config    StackConfig `json:"config" cql:"config"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt time.Time   `json:"updated_at"`
	}

	// Repo represents the git repository of an app.
	Repo struct {
		ID            gocql.UUID `json:"id" cql:"id"`
		StackID       gocql.UUID `json:"stack_id" cql:"stack_id"`
		ProviderID    gocql.UUID `json:"repo_id" cql:"repo_id"`               // The ID as provided by the provider.
		DefaultBranch string     `json:"default_branch" cql:"default_branch"` // The default branch to keep track of major releases.
		IsMonorepo    bool       `json:"is_monorepo" cql:"is_monorepo"`       // app can have multiple repos
		Provider      string     `json:"provider" cql:"provider"`             // can be github, gitlab, bitbucket, etc
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
	}

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

func (stack *Stack) GetTable() itable.ITable { return stackTable }
func (stack *Stack) PreCreate() error        { stack.Slug = db.CreateSlug(stack.Name); return nil }
func (stack *Stack) PreUpdate() error        { return nil }

func (repo *Repo) GetTable() itable.ITable { return repoTable }
func (repo *Repo) PreCreate() error        { return nil }
func (repo *Repo) PreUpdate() error        { return nil }

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
