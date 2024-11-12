-- name: CreateGithubRepo :one
INSERT INTO github_repos (installation_id, github_id, name, full_name, url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetGithubRepoByID :one
SELECT *
FROM github_repos
WHERE id = $1;

-- name: GetGithubRepoByInstallationIDAndGithubID :one
SELECT *
FROM github_repos
WHERE installation_id = $1 AND github_id = $2;

-- name: SuspendedGithubRepo :exec
UPDATE github_repos
SET is_active = false
WHERE id = $1;

-- name: ActivateGithubRepo :exec
UPDATE github_repos
SET is_active = true
WHERE id = $1;

