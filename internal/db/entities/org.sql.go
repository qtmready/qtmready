// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: org.sql

package entities

import (
	"context"
)

const createOrg = `-- name: CreateOrg :one
INSERT INTO orgs (name, slug)
VALUES ($1, $2)
RETURNING id, created_at, updated_at, name, slug
`

type CreateOrgParams struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (q *Queries) CreateOrg(ctx context.Context, arg CreateOrgParams) (Org, error) {
	row := q.db.QueryRow(ctx, createOrg, arg.Name, arg.Slug)
	var i Org
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Slug,
	)
	return i, err
}
