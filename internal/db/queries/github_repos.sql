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

-- name: SetRepoIDonGihubRepo :exec
UPDATE github_repos
SET repo_id = $2
WHERE id = $1;

-- name: SuspendedGithubRepo :exec
UPDATE github_repos
SET is_active = false
WHERE id = $1;

-- name: ActivateGithubRepo :exec
UPDATE github_repos
SET is_active = true
WHERE id = $1;

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
