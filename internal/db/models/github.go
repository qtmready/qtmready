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
		"installation_id",
		"installation_login",
		"installation_type",
		"sender_id",
		"sender_login",
		"ctrlplane_team_id",
	},
}

var githubInstallationTable = table.New(githubInstallMeta)

type GithubInstallation struct {
	ID                      gocql.UUID `cql:"id"`
	TeamID                  gocql.UUID `cql:"team_id"`
	GithubInstallationID    int64      `cql:"installation_id"`
	GithubSenderID          int64      `cql:"sender_id"`
	GithubInstallationLogin string
	GithubInstallationType  string
	GithubSenderLogin       string
}

func (gi *GithubInstallation) Create() error {
	gi.ID, _ = gocql.RandomUUID()

	query := conf.DB.Session.Query(githubInstallationTable.Insert()).BindStruct(gi)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}
