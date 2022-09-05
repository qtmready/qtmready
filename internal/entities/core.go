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
	ID        gocql.UUID `json:"id" cql:"id"`
	TeamID    gocql.UUID `json:"team_id" cql:"team_id"`
	Name      string     `json:"name" validate:"required"`
	Slug      string     `json:"slug"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (app *App) GetTable() *table.Table { return appTable }
func (app *App) PreCreate() error       { app.Slug = db.CreateSlug(app.Name); return nil }
func (app *App) PreUpdate() error       { return nil }

var (
	appRepoColumns = []string{
		"id",
		"app_id",
		"repo_id",
		"default_branch",
		"is_monorepo",
		"provider",
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
	ID            gocql.UUID `json:"id" cql:"id"`
	AppID         gocql.UUID `json:"app_id" cql:"app_id"`
	RepoID        gocql.UUID `json:"repo_id" cql:"repo_id"`
	DefaultBranch string     `json:"default_branch" cql:"default_branch"` // The default branch to keep track of major releases.
	IsMonorepo    bool       `json:"is_monorepo" cql:"is_monorepo"`       // an core can have multiple repos, our of which one can be a monorepo.
	Provider      string     `json:"provider" cql:"provider"`             // can be github, gitlab, bitbucket, etc
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (c *AppRepo) GetTable() *table.Table { return appRepoTable }
func (c *AppRepo) PreCreate() error       { return nil }
func (c *AppRepo) PreUpdate() error       { return nil }
