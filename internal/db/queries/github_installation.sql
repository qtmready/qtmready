-- name: CreateInstallation :one
INSERT INTO github_installations (org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at, org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login, status;

-- name: GetInstallation :one
SELECT *
FROM github_installations
WHERE id = $1;

-- name: GetInstallationByInstallationIDAndInstallationLogin :one
SELECT *
FROM github_installations
WHERE installation_id = $1 AND installation_login = $2;

-- name: UpdateInstallation :one
UPDATE github_installations
SET 
    org_id = $2,
    installation_id = $3,
    installation_login = $4,
    installation_login_id = $5,
    installation_type = $6,
    sender_id = $7,
    sender_login = $8,
    status = $9
WHERE id = $1
RETURNING id, created_at, updated_at, org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login, status;

-- name: DeleteInstallation :one
DELETE FROM github_installations
WHERE id = $1
RETURNING id;
