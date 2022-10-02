// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package entities

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	teamUserColumns = []string{
		"id",
		"user_id",
		"team_id",
		"created_at",
		"updated_at",
	}

	teamUserMeta = table.Metadata{
		Name:    "team_users",
		Columns: teamUserColumns,
	}

	teamUserTable = table.New(teamUserMeta)
)

type (
	TeamUser struct {
		ID        gocql.UUID `json:"id" cql:"id"`
		UserID    gocql.UUID `json:"user_id" cql:"user_id"`
		TeamID    gocql.UUID `json:"team_id" cql:"team_id"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}
)

func (tu *TeamUser) GetTable() *table.Table { return teamUserTable }
func (tu *TeamUser) PreCreate() error       { return nil }
func (tu *TeamUser) PreUpdate() error       { return nil }
