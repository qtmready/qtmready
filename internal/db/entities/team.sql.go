// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: team.sql

package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createTeam = `-- name: CreateTeam :one
INSERT INTO teams (name) 
VALUES ($1) 
RETURNING id, name
`

type CreateTeamRow struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (q *Queries) CreateTeam(ctx context.Context, name string) (CreateTeamRow, error) {
	row := q.db.QueryRow(ctx, createTeam, name)
	var i CreateTeamRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteTeam = `-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1
`

func (q *Queries) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteTeam, id)
	return err
}

const getTeam = `-- name: GetTeam :one
SELECT id, created_at, updated_at, org_id, name, slug
FROM teams
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetTeam(ctx context.Context, id uuid.UUID) (Team, error) {
	row := q.db.QueryRow(ctx, getTeam, id)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Name,
		&i.Slug,
	)
	return i, err
}

const getTeamByName = `-- name: GetTeamByName :one
SELECT id, name 
FROM teams 
WHERE name = $1
`

type GetTeamByNameRow struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (q *Queries) GetTeamByName(ctx context.Context, name string) (GetTeamByNameRow, error) {
	row := q.db.QueryRow(ctx, getTeamByName, name)
	var i GetTeamByNameRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const updateTeam = `-- name: UpdateTeam :one
UPDATE teams
SET org_id = $2, name = $3, slug = $4
WHERE id = $1
RETURNING id, created_at, name, slug
`

type UpdateTeamParams struct {
	ID    uuid.UUID `json:"id"`
	OrgID uuid.UUID `json:"org_id"`
	Name  string    `json:"name"`
	Slug  string    `json:"slug"`
}

type UpdateTeamRow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
}

func (q *Queries) UpdateTeam(ctx context.Context, arg UpdateTeamParams) (UpdateTeamRow, error) {
	row := q.db.QueryRow(ctx, updateTeam,
		arg.ID,
		arg.OrgID,
		arg.Name,
		arg.Slug,
	)
	var i UpdateTeamRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Slug,
	)
	return i, err
}
