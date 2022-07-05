package models

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"go.breu.io/ctrlplane/internal/conf"
)

var githubInstallMeta = table.Metadata{
	Name: "github_installations",
	Columns: []string{
		"id",
		"team_id",
		"github_installation_id",
		"github_installation_login",
		"github_installation_type",
		"github_sender_id",
		"github_sender_login",
	},
}

var githubInstallationTable = table.New(githubInstallMeta)

type GithubInstallation struct {
	ID                      gocql.UUID `cql:"id"`
	TeamID                  gocql.UUID `cql:"team_id"`
	GithubInstallationID    int64      `cql:"github_installation_id"`
	GithubSenderID          int64      `cql:"github_sender_id"`
	GithubInstallationLogin string
	GithubInstallationType  string
	GithubSenderLogin       string
}

func (g *GithubInstallation) Create() error {
	g.ID, _ = gocql.RandomUUID()

	query := conf.DB.Session.Query(githubInstallationTable.Insert()).BindStruct(g)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func (g *GithubInstallation) Get(params struct{}) error {
	query := conf.DB.Session.Query(githubInstallationTable.Select()).BindStruct(params)

	if err := query.GetRelease(&g); err != nil {
		return err
	}

	return nil
}
