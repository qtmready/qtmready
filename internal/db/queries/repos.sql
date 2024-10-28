-- name: CreateRepo :one
INSERT INTO repos (org_id, name, provider, provider_id, default_branch, is_monorepo, threshold, stale_duration)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetRepoByID :one
SELECT *
FROM repos
WHERE id = $1;

-- name: UpdateRepo :one
UPDATE repos
SET org_id = $2,
    name = $3,
    provider = $4,
    provider_id = $5,
    default_branch = $6,
    is_monorepo = $7,
    threshold = $8,
    stale_duration = $9
WHERE id = $1
RETURNING *;

-- name: DeleteRepo :exec
DELETE FROM repos
WHERE id = $1;

-- name: ListRepos :many
SELECT *
FROM repos
ORDER BY created_at DESC;

-- name: GetOrgRepos :many
SELECT *
FROM repos 
WHERE org_id = $1; 
