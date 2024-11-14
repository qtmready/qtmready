// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: repos.sql

package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const activateRepoByHookID = `-- name: ActivateRepoByHookID :exec
UPDATE repos
SET is_active = true
WHERE hook_id = $1
`

func (q *Queries) ActivateRepoByHookID(ctx context.Context, hookID uuid.UUID) error {
	_, err := q.db.Exec(ctx, activateRepoByHookID, hookID)
	return err
}

const createRepo = `-- name: CreateRepo :one
INSERT INTO repos (org_id, name, hook, hook_id, url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, org_id, name, hook, hook_id, default_branch, is_monorepo, threshold, stale_duration, url, is_active
`

type CreateRepoParams struct {
	OrgID  uuid.UUID `json:"org_id"`
	Name   string    `json:"name"`
	Hook   string    `json:"hook"`
	HookID uuid.UUID `json:"hook_id"`
	Url    string    `json:"url"`
}

func (q *Queries) CreateRepo(ctx context.Context, arg CreateRepoParams) (Repo, error) {
	row := q.db.QueryRow(ctx, createRepo,
		arg.OrgID,
		arg.Name,
		arg.Hook,
		arg.HookID,
		arg.Url,
	)
	var i Repo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Name,
		&i.Hook,
		&i.HookID,
		&i.DefaultBranch,
		&i.IsMonorepo,
		&i.Threshold,
		&i.StaleDuration,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const deleteRepo = `-- name: DeleteRepo :exec
DELETE FROM repos
WHERE id = $1
`

func (q *Queries) DeleteRepo(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRepo, id)
	return err
}

const getOrgReposByOrgID = `-- name: GetOrgReposByOrgID :many
SELECT id, created_at, updated_at, org_id, name, hook, hook_id, default_branch, is_monorepo, threshold, stale_duration, url, is_active
FROM repos
WHERE org_id = $1
`

func (q *Queries) GetOrgReposByOrgID(ctx context.Context, orgID uuid.UUID) ([]Repo, error) {
	rows, err := q.db.Query(ctx, getOrgReposByOrgID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Repo
	for rows.Next() {
		var i Repo
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OrgID,
			&i.Name,
			&i.Hook,
			&i.HookID,
			&i.DefaultBranch,
			&i.IsMonorepo,
			&i.Threshold,
			&i.StaleDuration,
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

const getRepo = `-- name: GetRepo :one
SELECT 
  r.id,
  r.org_id,
  r.name,
  r.hook,
  r.hook_id,
  r.default_branch,
  r.is_monorepo,
  r.threshold,
  r.stale_duration,
  r.url,
  r.is_active,
  json_build_object(
    'id', m.id,
    'hook', m.hook,
    'kind', m.kind,
    'link_to', m.link_to,
    'data', m.data
  ) AS messaging,
  json_build_object(
    'id', o.id,
    'name', o.name,
    'domain', o.domain,
    'slug', o.slug,
    'hooks', o.hooks
  ) AS org
FROM 
  github_repos gr
JOIN 
  repos r ON gr.id = r.hook_id
LEFT JOIN 
  messaging m ON m.link_to = r.id
JOIN 
  orgs o ON r.org_id = o.id
WHERE 
  gr.installation_id = $1 AND gr.github_id = $2
`

type GetRepoParams struct {
	InstallationID uuid.UUID `json:"installation_id"`
	GithubID       int64     `json:"github_id"`
}

type GetRepoRow struct {
	ID            uuid.UUID       `json:"id"`
	OrgID         uuid.UUID       `json:"org_id"`
	Name          string          `json:"name"`
	Hook          string          `json:"hook"`
	HookID        uuid.UUID       `json:"hook_id"`
	DefaultBranch string          `json:"default_branch"`
	IsMonorepo    bool            `json:"is_monorepo"`
	Threshold     int32           `json:"threshold"`
	StaleDuration pgtype.Interval `json:"stale_duration"`
	Url           string          `json:"url"`
	IsActive      bool            `json:"is_active"`
	Messaging     []byte          `json:"messaging"`
	Org           []byte          `json:"org"`
}

func (q *Queries) GetRepo(ctx context.Context, arg GetRepoParams) (GetRepoRow, error) {
	row := q.db.QueryRow(ctx, getRepo, arg.InstallationID, arg.GithubID)
	var i GetRepoRow
	err := row.Scan(
		&i.ID,
		&i.OrgID,
		&i.Name,
		&i.Hook,
		&i.HookID,
		&i.DefaultBranch,
		&i.IsMonorepo,
		&i.Threshold,
		&i.StaleDuration,
		&i.Url,
		&i.IsActive,
		&i.Messaging,
		&i.Org,
	)
	return i, err
}

const getRepoByID = `-- name: GetRepoByID :one
SELECT id, created_at, updated_at, org_id, name, hook, hook_id, default_branch, is_monorepo, threshold, stale_duration, url, is_active
FROM repos
WHERE id = $1
`

func (q *Queries) GetRepoByID(ctx context.Context, id uuid.UUID) (Repo, error) {
	row := q.db.QueryRow(ctx, getRepoByID, id)
	var i Repo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Name,
		&i.Hook,
		&i.HookID,
		&i.DefaultBranch,
		&i.IsMonorepo,
		&i.Threshold,
		&i.StaleDuration,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getReposByHookAndHookID = `-- name: GetReposByHookAndHookID :one
SELECT id, created_at, updated_at, org_id, name, hook, hook_id, default_branch, is_monorepo, threshold, stale_duration, url, is_active
FROM repos
WHERE hook = $1 AND hook_id = $2
`

type GetReposByHookAndHookIDParams struct {
	Hook   string    `json:"hook"`
	HookID uuid.UUID `json:"hook_id"`
}

func (q *Queries) GetReposByHookAndHookID(ctx context.Context, arg GetReposByHookAndHookIDParams) (Repo, error) {
	row := q.db.QueryRow(ctx, getReposByHookAndHookID, arg.Hook, arg.HookID)
	var i Repo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Name,
		&i.Hook,
		&i.HookID,
		&i.DefaultBranch,
		&i.IsMonorepo,
		&i.Threshold,
		&i.StaleDuration,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const listRepos = `-- name: ListRepos :many
SELECT id, created_at, updated_at, org_id, name, hook, hook_id, default_branch, is_monorepo, threshold, stale_duration, url, is_active
FROM repos
WHERE org_id = $1
ORDER BY updated_at DESC
`

func (q *Queries) ListRepos(ctx context.Context, orgID uuid.UUID) ([]Repo, error) {
	rows, err := q.db.Query(ctx, listRepos, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Repo
	for rows.Next() {
		var i Repo
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OrgID,
			&i.Name,
			&i.Hook,
			&i.HookID,
			&i.DefaultBranch,
			&i.IsMonorepo,
			&i.Threshold,
			&i.StaleDuration,
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

const suspendedRepoByHookID = `-- name: SuspendedRepoByHookID :exec
UPDATE repos
SET is_active = false
WHERE hook_id = $1
`

func (q *Queries) SuspendedRepoByHookID(ctx context.Context, hookID uuid.UUID) error {
	_, err := q.db.Exec(ctx, suspendedRepoByHookID, hookID)
	return err
}

const updateRepo = `-- name: UpdateRepo :one
UPDATE repos
SET org_id = $2,
    name = $3,
    hook = $4,
    hook_id = $5,
    default_branch = $6,
    is_monorepo = $7,
    threshold = $8,
    stale_duration = $9
WHERE id = $1
RETURNING id, created_at, updated_at, org_id, name, hook, hook_id, default_branch, is_monorepo, threshold, stale_duration, url, is_active
`

type UpdateRepoParams struct {
	ID            uuid.UUID       `json:"id"`
	OrgID         uuid.UUID       `json:"org_id"`
	Name          string          `json:"name"`
	Hook          string          `json:"hook"`
	HookID        uuid.UUID       `json:"hook_id"`
	DefaultBranch string          `json:"default_branch"`
	IsMonorepo    bool            `json:"is_monorepo"`
	Threshold     int32           `json:"threshold"`
	StaleDuration pgtype.Interval `json:"stale_duration"`
}

func (q *Queries) UpdateRepo(ctx context.Context, arg UpdateRepoParams) (Repo, error) {
	row := q.db.QueryRow(ctx, updateRepo,
		arg.ID,
		arg.OrgID,
		arg.Name,
		arg.Hook,
		arg.HookID,
		arg.DefaultBranch,
		arg.IsMonorepo,
		arg.Threshold,
		arg.StaleDuration,
	)
	var i Repo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Name,
		&i.Hook,
		&i.HookID,
		&i.DefaultBranch,
		&i.IsMonorepo,
		&i.Threshold,
		&i.StaleDuration,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}
