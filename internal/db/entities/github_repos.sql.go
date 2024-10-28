// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: github_repos.sql

package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createGithubRepo = `-- name: CreateGithubRepo :one
INSERT INTO github_repos (repo_id, installation_id, github_id, name, full_name, url, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
`

type CreateGithubRepoParams struct {
	RepoID         pgtype.UUID `json:"repo_id"`
	InstallationID uuid.UUID   `json:"installation_id"`
	GithubID       int64       `json:"github_id"`
	Name           string      `json:"name"`
	FullName       string      `json:"full_name"`
	Url            string      `json:"url"`
	IsActive       pgtype.Bool `json:"is_active"`
}

func (q *Queries) CreateGithubRepo(ctx context.Context, arg CreateGithubRepoParams) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, createGithubRepo,
		arg.RepoID,
		arg.InstallationID,
		arg.GithubID,
		arg.Name,
		arg.FullName,
		arg.Url,
		arg.IsActive,
	)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const deleteGithubRepo = `-- name: DeleteGithubRepo :one
DELETE FROM github_repos
WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteGithubRepo(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, deleteGithubRepo, id)
	err := row.Scan(&id)
	return id, err
}

const getGithubRepoByFullName = `-- name: GetGithubRepoByFullName :one
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE full_name = $1
`

func (q *Queries) GetGithubRepoByFullName(ctx context.Context, fullName string) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, getGithubRepoByFullName, fullName)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getGithubRepoByID = `-- name: GetGithubRepoByID :one
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE id = $1
`

func (q *Queries) GetGithubRepoByID(ctx context.Context, id uuid.UUID) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, getGithubRepoByID, id)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getGithubRepoByInstallationIDAndGithubID = `-- name: GetGithubRepoByInstallationIDAndGithubID :one
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE installation_id = $1 AND github_id = $2
`

type GetGithubRepoByInstallationIDAndGithubIDParams struct {
	InstallationID uuid.UUID `json:"installation_id"`
	GithubID       int64     `json:"github_id"`
}

func (q *Queries) GetGithubRepoByInstallationIDAndGithubID(ctx context.Context, arg GetGithubRepoByInstallationIDAndGithubIDParams) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, getGithubRepoByInstallationIDAndGithubID, arg.InstallationID, arg.GithubID)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getGithubRepoByName = `-- name: GetGithubRepoByName :one
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE name = $1
`

func (q *Queries) GetGithubRepoByName(ctx context.Context, name string) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, getGithubRepoByName, name)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getGithubRepoByRepoID = `-- name: GetGithubRepoByRepoID :many
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE repo_id = $1
`

func (q *Queries) GetGithubRepoByRepoID(ctx context.Context, repoID pgtype.UUID) ([]GithubRepo, error) {
	rows, err := q.db.Query(ctx, getGithubRepoByRepoID, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GithubRepo
	for rows.Next() {
		var i GithubRepo
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.RepoID,
			&i.InstallationID,
			&i.GithubID,
			&i.Name,
			&i.FullName,
			&i.Url,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGithubRepo = `-- name: UpdateGithubRepo :one
UPDATE github_repos
SET repo_id = $2, 
    installation_id = $3, 
    github_id = $4, 
    name = $5, 
    full_name = $6, 
    url = $7, 
    is_active = $8
WHERE id = $1
RETURNING id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
`

type UpdateGithubRepoParams struct {
	ID             uuid.UUID   `json:"id"`
	RepoID         pgtype.UUID `json:"repo_id"`
	InstallationID uuid.UUID   `json:"installation_id"`
	GithubID       int64       `json:"github_id"`
	Name           string      `json:"name"`
	FullName       string      `json:"full_name"`
	Url            string      `json:"url"`
	IsActive       pgtype.Bool `json:"is_active"`
}

func (q *Queries) UpdateGithubRepo(ctx context.Context, arg UpdateGithubRepoParams) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, updateGithubRepo,
		arg.ID,
		arg.RepoID,
		arg.InstallationID,
		arg.GithubID,
		arg.Name,
		arg.FullName,
		arg.Url,
		arg.IsActive,
	)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}
