-- name: CreateGithubRepo :one
INSERT INTO github_repos (repo_id, installation_id, github_id, name, full_name, url, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetGithubRepoByID :one
SELECT *
FROM github_repos
WHERE id = $1;

-- name: UpdateGithubRepo :one
UPDATE github_repos
SET repo_id = $2, 
    installation_id = $3, 
    github_id = $4, 
    name = $5, 
    full_name = $6, 
    url = $7, 
    is_active = $8
WHERE id = $1
RETURNING *;

-- name: DeleteGithubRepo :one
DELETE FROM github_repos
WHERE id = $1
RETURNING id;

-- name: GetGithubRepoByFullName :one
SELECT *
FROM github_repos
WHERE full_name = $1;

-- name: GetGithubRepoByName :one
SELECT *
FROM github_repos
WHERE name = $1;

-- name: GetGithubRepoByInstallationIDAndGithubID :one
SELECT *
FROM github_repos
WHERE installation_id = $1 AND github_id = $2;

-- name: GetGithubReposWithCoreRepo :one
SELECT 
    g.*, 
    json_build_object(
        'repo_id', r.id,
        'repo_name', r.name,
        'provider', r.provider,
        'provider_id', r.provider_id,
        'default_branch', r.default_branch,
        'is_monorepo', r.is_monorepo,
        'threshold', r.threshold,
        'stale_duration', r.stale_duration
    ) AS repo
FROM 
    github_repos g
LEFT JOIN 
    repos r ON g.repo_id = r.id
WHERE 
    g.id = $1 -- TODO - based on intallation id or some other
LIMIT 1;

-- name: GetGithubRepo :one
SELECT *
FROM github_repos
WHERE name = $1 AND full_name = $2 AND github_id = $3; 