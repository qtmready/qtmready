package models

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"go.breu.io/ctrlplane/internal/conf"
)

var teamMeta = table.Metadata{
	Name: "orgs",
	Columns: []string{
		"id",
		"name",
		"slug",
		"created_at",
		"updated_at",
	},
}

var teamTable = table.New(teamMeta)

type Team struct {
	ID        gocql.UUID `cql:"id"`
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Team) Create() error {
	t.ID, _ = gocql.RandomUUID()

	now := time.Now()
	t.Slug = slugify(t.Name)
	t.CreatedAt = now
	t.UpdatedAt = now

	query := conf.DB.Session.Query(teamTable.Insert()).BindStruct(t)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}
