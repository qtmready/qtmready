-- name: CreateGithubOrg :one
INSERT INTO github_orgs (installation_id, github_org_id, name)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at, installation_id, github_org_id, name;

-- name: GetGithubOrgByID :one
SELECT *
FROM github_orgs
WHERE id = $1;

-- name: GetGithubOrgByGithubOrgID :one
SELECT *
FROM github_orgs
WHERE github_org_id = $1;

-- name: GetGithubOrgByInstallationID :many
SELECT *
FROM github_orgs
WHERE installation_id = $1;

-- name: UpdateGithubOrg :one
UPDATE github_orgs
SET installation_id = $2, github_org_id = $3, name = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, installation_id, github_org_id, name;

-- name: DeleteGithubOrg :exec
DELETE FROM github_orgs
WHERE id = $1;
