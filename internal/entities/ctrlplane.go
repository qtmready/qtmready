package entities

import (
	"time"

	"github.com/scylladb/gocqlx/v2/table"
)

var (
	appColumns = []string{
		"id",
		"team_id",
		"name",
		"created_at",
		"updated_at",
	}

	appMeta = table.Metadata{
		Name:    "apps",
		Columns: appColumns,
	}

	appTable = table.New(appMeta)
)

type App struct {
	ID        string    `json:"id" cql:"id"`
	TeamID    string    `json:"team_id" cql:"team_id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *App) GetTable() *table.Table { return appTable }
func (c *App) PreCreate() error       { return nil }
func (c *App) PreUpdate() error       { return nil }

var (
	appRepoColumns = []string{
		"id",
		"app_id",
		"repo_id",
		"default_branch",
		"is_monorepo",
		"git_provider",
		"created_at",
		"updated_at",
	}

	appRepoMeta = table.Metadata{
		Name:    "app_repos",
		Columns: appRepoColumns,
	}

	appRepoTable = table.New(appRepoMeta)
)

type AppRepo struct {
	ID            string    `json:"id" cql:"id"`
	AppID         string    `json:"app_id" cql:"app_id"`
	RepoID        string    `json:"repo_id" cql:"repo_id"`
	DefaultBranch string    `json:"default_branch" cql:"default_branch"`
	IsMonorepo    bool      `json:"is_monorepo" cql:"is_monorepo"`   // an app can have multiple repos, our of which one can be a monorepo.
	GitProvider   string    `json:"git_provider" cql:"git_provider"` // can be github, gitlab, bitbucket, etc
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (c *AppRepo) GetTable() *table.Table { return appRepoTable }
func (c *AppRepo) PreCreate() error       { return nil }
func (c *AppRepo) PreUpdate() error       { return nil }
