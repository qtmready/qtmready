-- name: CreateRepo :one
INSERT INTO repos (org_id, name, hook, hook_id, url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRepoByID :one
SELECT *
FROM repos
WHERE id = $1;

-- name: UpdateRepo :one
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
RETURNING *;

-- name: DeleteRepo :exec
DELETE FROM repos
WHERE id = $1;

-- name: ListRepos :many
SELECT *
FROM repos
ORDER BY created_at DESC;

-- name: GetOrgReposByOrgID :many
SELECT *
FROM repos 
WHERE org_id = $1; 

-- name: GetReposByHookAndHookID :one
SELECT *
FROM repos 
WHERE hook = $1 AND hook_id = $2; 

-- name: SuspendedRepoByHookID :exec
UPDATE repos
SET is_active = false
WHERE hook_id = $1;

-- name: ActivateRepoByHookID :exec
UPDATE repos
SET is_active = true
WHERE hook_id = $1;

-- name: GetRepoByInstallationIDAndGithubID :one
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
  r.is_active
FROM 
  github_repos gr
JOIN 
  repos r ON gr.id = r.hook_id
WHERE 
  gr.installation_id = $1 AND gr.github_id = $2;
