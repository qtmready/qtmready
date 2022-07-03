package models

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"go.breu.io/ctrlplane/internal/conf"
)

var orgMeta = table.Metadata{
	Name: "orgs",
	Columns: []string{
		"id",
		"name",
		"website",
		"created_at",
		"updated_at",
	},
}

var orgTable = table.New(orgMeta)

type Team struct {
	ID        gocql.UUID `sql:"id"`
	Name      string
	Website   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Team) Create() error {
	t.ID, _ = gocql.RandomUUID()

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := conf.DB.Session.Query(orgTable.Insert()).BindStruct(t)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}
