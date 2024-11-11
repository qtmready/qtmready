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

-- name: GetGithubRepoWithRepo :one
SELECT 
  gr.id AS github_repo_id,
  gr.installation_id,
  gr.github_id,
  gr.name AS github_repo_name,
  gr.full_name,
  gr.url AS github_repo_url,
  gr.is_active AS github_repo_is_active,
  r.id AS repo_id,
  r.org_id,
  r.name,
  r.hook,
  r.hook_id,
  r.default_branch,
  r.is_monorepo,
  r.threshold,
  r.stale_duration,
  r.url AS repo_url,
  r.is_active AS repo_is_active
FROM 
  github_repos gr
JOIN 
  repos r ON gr.id = r.hook_id
WHERE 
  gr.id = $1;
