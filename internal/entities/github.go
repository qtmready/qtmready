package entities

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	githubInstallationColumns = []string{
		"id",
		"team_id",
		"installation_id",
		"installation_login",
		"installation_type",
		"sender_id",
		"sender_login",
		"status",
		"created_at",
		"updated_at",
	}

	githubInstallationMeta = table.Metadata{
		Name:    "github_installations",
		Columns: githubInstallationColumns,
		PartKey: []string{"id"},
	}

	githubInstallationTable = table.New(githubInstallationMeta)
)

type GithubInstallation struct {
	ID                gocql.UUID `json:"id" cql:"id"`
	TeamID            gocql.UUID `json:"team_id" cql:"team_id"`
	InstallationID    int64      `json:"installation_id" cql:"installation_id" validate:"required,db_unique"`
	InstallationLogin string     `json:"installation_login" cql:"installation_login" validate:"required,db_unique"`
	InstallationType  string     `json:"installation_type" cql:"installation_type"`
	SenderID          int64      `json:"sender_id" cql:"sender_id"`
	SenderLogin       string     `json:"sender_login" cql:"sender_login"`
	Status            string     `json:"status" cql:"status"`
	CreatedAt         time.Time  `json:"created_at" cql:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" cql:"updated_at"`
}

func (g GithubInstallation) GetTable() *table.Table { return githubInstallationTable }
func (g GithubInstallation) PreCreate() error       { return nil }
func (g GithubInstallation) PreUpdate() error       { return nil }

var (
	githubRepoColumns = []string{
		"id",
		"github_id",
		"team_id",
		"name",
		"full_name",
		"created_at",
		"updated_at",
	}

	githubRepoMeta = table.Metadata{
		Name:    "github_repos",
		Columns: githubRepoColumns,
		PartKey: []string{"id"},
	}

	githubRepoTable = table.New(githubRepoMeta)
)

type GithubRepo struct {
	ID        gocql.UUID `json:"id" cql:"id"`
	GithubID  int64      `json:"github_id" cql:"github_id" validate:"required"`
	TeamID    gocql.UUID `json:"team_id" cql:"team_id"`
	Name      string     `json:"name" cql:"name" validate:"required,db_unique"`
	FullName  string     `json:"full_name" cql:"full_name" validate:"required,db_unique"`
	CreatedAt time.Time  `json:"created_at" cql:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" cql:"updated_at"`
}

func (g GithubRepo) GetTable() *table.Table { return githubRepoTable }
func (g GithubRepo) PreCreate() error       { return nil }
func (g GithubRepo) PreUpdate() error       { return nil }
