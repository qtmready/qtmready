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
SET
    org_id = $2,
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
SELECT
  repo.*, 
  CASE 
    WHEN msg.id IS NOT NULL AND msg.link_to IS NOT NULL THEN TRUE
    ELSE FALSE
  END AS has_msging,
  CASE 
    WHEN msg.id IS NOT NULL THEN msg.data->>'channel_name'
    ELSE ''
  END AS channel_name
FROM
  repos AS repo
LEFT JOIN 
  messaging AS msg
ON 
  repo.id = msg.link_to
WHERE 
  repo.org_id = $1
ORDER BY 
  repo.updated_at DESC;

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

-- name: GetRepo :one
SELECT
  sqlc.embed(repo),
  sqlc.embed(msg),
  sqlc.embed(org)
FROM
  github_repos gr
JOIN
  repos repo ON gr.id = repo.hook_id
LEFT JOIN
  messaging msg ON msg.link_to = repo.id
JOIN
  orgs org ON repo.org_id = org.id
WHERE
  gr.installation_id = $1 AND gr.github_id = $2;
