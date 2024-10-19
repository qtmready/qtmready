-- name: CreateOrg :one
INSERT INTO orgs (name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateOrg :one
UPDATE orgs
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteOrg :exec
DELETE FROM orgs
WHERE id = $1;
