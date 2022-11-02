// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package entities

import (
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"

	"go.breu.io/ctrlplane/internal/db"
)

var (
	teamColumns = []string{
		"id",
		"name",
		"slug",
		"created_at",
		"updated_at",
	}

	teamMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "teams",
			Columns: teamColumns,
		},
	}

	teamTable = itable.New(*teamMeta.M)
)

type (
	Team struct {
		ID        gocql.UUID `json:"id" cql:"id"`
		Name      string     `json:"name" validate:"required"`
		Slug      string     `json:"slug"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}
)

func (t *Team) GetTable() itable.ITable { return teamTable }
func (t *Team) PreCreate() error        { t.Slug = db.CreateSlug(t.Name); return nil }
func (t *Team) PreUpdate() error        { return nil }
