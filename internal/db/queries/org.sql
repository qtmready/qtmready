-- name: CreateOrg :one
INSERT INTO orgs (name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateOrg :one
UPDATE orgs 
SET name = $2, slug = $3
WHERE id = $1 
RETURNING id, created_at, updated_at, name, slug;

-- name: DeleteOrg :exec
DELETE FROM orgs 
WHERE id = $1;
