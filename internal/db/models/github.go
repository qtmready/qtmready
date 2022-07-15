package models

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"go.breu.io/ctrlplane/internal/db"
)

var githubInstallMeta = table.Metadata{
	Name: "github_installations",
	Columns: []string{
		"id",
		"team_id",
		"installation_id",
		"installation_login",
		"installation_type",
		"sender_id",
		"sender_login",
		"created_at",
		"updated_at",
	},
}

var githubInstallationTable = table.New(githubInstallMeta)

type GithubInstallation struct {
	ID                gocql.UUID `cql:"id"`
	TeamID            gocql.UUID `cql:"team_id"`
	InstallationID    int64      `cql:"installation_id"`
	SenderID          int64      `cql:"sender_id"`
	InstallationLogin string
	InstallationType  string
	SenderLogin       string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (g *GithubInstallation) Create() error {
	g.ID, _ = gocql.RandomUUID()

	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now

	query := db.DB.Session.Query(githubInstallationTable.Insert()).BindStruct(g)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func (g *GithubInstallation) Get(params interface{}) error {
	query := db.DB.Session.Query(githubInstallationTable.Select()).BindStruct(params)

	if err := query.GetRelease(&g); err != nil {
		return err
	}

	return nil
}
