package entities

import (
	"time"

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

	appMeta = table.Metadata{
		Name:    "apps",
		Columns: appColumns,
	}

	appTable = table.New(appMeta)
)

// App defines an application.
// An application is a collection of repos, workloads & resources defined by an environment.
type App struct {
	ID        gocql.UUID `json:"id" cql:"id"`
	TeamID    gocql.UUID `json:"team_id" cql:"team_id"`
	Name      string     `json:"name" validate:"required"`
	Slug      string     `json:"slug"`
	Config    AppConfig  `json:"config" cql:"config"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (app *App) GetTable() *table.Table { return appTable }
func (app *App) PreCreate() error       { app.Slug = db.CreateSlug(app.Name); return nil }
func (app *App) PreUpdate() error       { return nil }

var (
	repoColumns = []string{
		"id",
		"app_id",
		"repo_id",
		"default_branch",
		"is_monorepo",
		"driver",
		"created_at",
		"updated_at",
	}

	repoMeta = table.Metadata{
		Name:    "repos",
		Columns: repoColumns,
	}

	repoTable = table.New(repoMeta)
)

// Repo is the git repository for an app
type Repo struct {
	ID            gocql.UUID `json:"id" cql:"id"`
	AppID         gocql.UUID `json:"app_id" cql:"app_id"`
	RepoID        gocql.UUID `json:"repo_id" cql:"repo_id"`
	DefaultBranch string     `json:"default_branch" cql:"default_branch"` // The default branch to keep track of major releases.
	IsMonorepo    bool       `json:"is_monorepo" cql:"is_monorepo"`       // app can have multiple repos
	Driver        string     `json:"driver" cql:"driver"`                 // can be github, gitlab, bitbucket, etc
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (repo *Repo) GetTable() *table.Table { return repoTable }
func (repo *Repo) PreCreate() error       { return nil }
func (repo *Repo) PreUpdate() error       { return nil }

var (
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

	workloadMeta = table.Metadata{
		Name:    "workloads",
		Columns: workloadColumns,
	}

	workloadTable = table.New(workloadMeta)
)

// Workload defines a workload for the app. See Workload.Kind for type
type Workload struct {
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

func (workload *Workload) GetTable() *table.Table { return workloadTable }
func (workload *Workload) PreCreate() error       { return nil }
func (workload *Workload) PreUpdate() error       { return nil }

var (
	resourceColumns = []string{
		"id",
		"app_id",
		"repo_id",
		"name",
		"driver",
		"kind",
		"is_immutable",
		"created_at",
		"updated_at",
	}

	resourceMeta = table.Metadata{
		Name:    "resources",
		Columns: resourceColumns,
	}

	resourceTable = table.New(resourceMeta)
)

// Resource defines the cloud provider resources for the app e.g. s3, sqs, etc
type Resource struct {
	ID          gocql.UUID `json:"id" cql:"id"`
	AppID       gocql.UUID `json:"app_id" cql:"app_id"`
	RepoID      gocql.UUID `json:"repo_id" cql:"repo_id"`
	Name        string     `json:"name" cql:"name"`
	Driver      string     `json:"driver" cql:"driver"` // "aws" | "gcp" | "azure"
	Kind        string     `json:"kind" cql:"kind"`     // "s3" | "sqs" | "sns" | "dynamodb" | "postgres" | "mysql" etc.
	IsImmutable bool       `json:"is_immutable" cql:"is_immutable"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (resource *Resource) GetTable() *table.Table { return resourceTable }
func (resource *Resource) PreCreate() error       { return nil }
func (resource *Resource) PreUpdate() error       { return nil }

var (
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

	blueprintMeta = table.Metadata{
		Name:    "blueprints",
		Columns: blueprintColumns,
	}

	blueprintTable = table.New(blueprintMeta)
)

type Blueprint struct {
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

func (blueprint *Blueprint) GetTable() *table.Table { return blueprintTable }
func (blueprint *Blueprint) PreCreate() error       { return nil }
func (blueprint *Blueprint) PreUpdate() error       { return nil }

var (
	rolloutColumns = []string{
		"id",
		"app_id",
		"blueprint_id",
		"trigger",
		"state",
		"created_at",
		"updated_at",
	}

	rolloutMeta = table.Metadata{
		Name:    "rollouts",
		Columns: rolloutColumns,
	}

	rolloutTable = table.New(rolloutMeta)
)

type Rollout struct {
	ID          gocql.UUID `json:"id" cql:"id"`
	AppID       gocql.UUID `json:"app_id" cql:"app_id"`
	BlueprintID string     `json:"blueprint_id" cql:"blueprint_id"`
	Trigger     string     `json:"trigger" cql:"trigger"`
	State       string     `json:"state" cql:"state"` // "in_progress" | "live" | "rejected"
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (rollout *Rollout) GetTable() *table.Table { return rolloutTable }
func (rollout *Rollout) PreCreate() error       { return nil }
func (rollout *Rollout) PreUpdate() error       { return nil }
