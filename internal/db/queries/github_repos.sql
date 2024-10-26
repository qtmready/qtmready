-- name: CreateGitHubRepo :one
INSERT INTO github_repos (repo_id, installation_id, github_id, name, full_name, url, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active;

-- name: GetGitHubRepoByID :one
SELECT *
FROM github_repos
WHERE id = $1;

-- name: UpdateGitHubRepo :one
UPDATE github_repos
SET repo_id = $2, 
    installation_id = $3, 
    github_id = $4, 
    name = $5, 
    full_name = $6, 
    url = $7, 
    is_active = $8
WHERE id = $1
RETURNING id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active;

-- name: DeleteGitHubRepo :one
DELETE FROM github_repos
WHERE id = $1
RETURNING id;

-- name: GetGitHubRepoByRepoID :many
SELECT *
FROM github_repos
WHERE repo_id = $1;

-- name: GetGitHubRepoByFullName :one
SELECT *
FROM github_repos
WHERE full_name = $1;

-- name: GetGitHubRepoByName :one
SELECT *
FROM github_repos
WHERE name = $1;

-- name: GetGitHubRepoByInstallationIDAndGitHubID :one
SELECT *
FROM github_repos
WHERE installation_id = $1 AND github_id = $2;

