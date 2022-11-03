// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
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
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  APPLICABLE LAW.

package entities

import (
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"

	"go.breu.io/ctrlplane/internal/db"
)

var (
	appColumns = []string{
		"id",
		"team_id",
		"name",
		"slug",
		"config",
		"created_at",
		"updated_at",
	}

	appMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "apps",
			Columns: appColumns,
		}}

	appTable = itable.New(*appMeta.M)

	repoColumns = []string{
		"id",
		"app_id",
		"repo_id",
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
		"app_id",
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
		"app_id",
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
		"app_id",
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

	rolloutColumns = []string{
		"id",
		"app_id",
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
	App struct {
		ID        gocql.UUID `json:"id" cql:"id"`
		TeamID    gocql.UUID `json:"team_id" cql:"team_id"`
		Name      string     `json:"name" validate:"required"`
		Slug      string     `json:"slug"`
		Config    AppConfig  `json:"config" cql:"config"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}

	// Repo represents the git repository of an app.
	Repo struct {
		ID            gocql.UUID `json:"id" cql:"id"`
		AppID         gocql.UUID `json:"app_id" cql:"app_id"`
		RepoID        gocql.UUID `json:"repo_id" cql:"repo_id"`
		DefaultBranch string     `json:"default_branch" cql:"default_branch"` // The default branch to keep track of major releases.
		IsMonorepo    bool       `json:"is_monorepo" cql:"is_monorepo"`       // app can have multiple repos
		Provider      string     `json:"provider" cql:"provider"`             // can be github, gitlab, bitbucket, etc
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
	}

	// Workload defines a workload for the app. See Workload.Kind for type.
	Workload struct {
		ID        gocql.UUID `json:"id" cql:"id"`
		AppID     gocql.UUID `json:"app_id" cql:"app_id"`
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
		AppID       gocql.UUID `json:"app_id" cql:"app_id"`
		RepoID      gocql.UUID `json:"repo_id" cql:"repo_id"`
		Name        string     `json:"name" cql:"name"`
		Provider    string     `json:"provider" cql:"provider"` // "aws" | "gcp" | "azure"
		Driver      string     `json:"driver" cql:"driver"`     // "s3" | "sqs" | "sns" | "dynamodb" | "postgres" | "mysql" etc.
		IsImmutable bool       `json:"is_immutable" cql:"is_immutable"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
	}

	Blueprint struct {
		ID            gocql.UUID       `json:"id" cql:"id"`
		AppID         gocql.UUID       `json:"app_id" cql:"app_id"`
		RepoID        gocql.UUID       `json:"repo_id" cql:"repo_id"`
		RepoBranch    string           `json:"repo_branch" cql:"repo_branch"`
		Name          string           `json:"name" cql:"name"`
		Regions       BluePrintRegions `json:"regions" cql:"regions"`
		RolloutBudget int              `json:"rollout_budget" cql:"rollout_budget"`
		CreatedAt     time.Time        `json:"created_at"`
		UpdatedAt     time.Time        `json:"updated_at"`
	}

	Rollout struct {
		ID          gocql.UUID `json:"id" cql:"id"`
		AppID       gocql.UUID `json:"app_id" cql:"app_id"`
		BlueprintID string     `json:"blueprint_id" cql:"blueprint_id"`
		Trigger     string     `json:"trigger" cql:"trigger"`
		State       string     `json:"state" cql:"state"` // "in_progress" | "live" | "rejected"
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
	}
)

func (app *App) GetTable() itable.ITable { return appTable }
func (app *App) PreCreate() error        { app.Slug = db.CreateSlug(app.Name); return nil }
func (app *App) PreUpdate() error        { return nil }

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

func (rollout *Rollout) GetTable() itable.ITable { return rolloutTable }
func (rollout *Rollout) PreCreate() error        { return nil }
func (rollout *Rollout) PreUpdate() error        { return nil }
