package entities

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
)

var (
	teamColumns = []string{
		"id",
		"name",
		"slug",
		"created_at",
		"updated_at",
	}

	teamMeta = table.Metadata{
		Name:    "teams",
		Columns: teamColumns,
		PartKey: []string{},
		SortKey: []string{},
	}

	teamTable = table.New(teamMeta)
)

type Team struct {
	ID        gocql.UUID `json:"id" cql:"id"`
	Name      string     `json:"name" validate:"required"`
	Slug      string     `json:"slug"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (t *Team) GetTable() *table.Table { return teamTable }
func (t *Team) PreCreate() error       { return nil }
func (t *Team) PreUpdate() error       { return nil }
