-- name: CreateOrg :one
INSERT INTO orgs (name, domain, slug)
VALUES ($1, LOWER($2), $3)
RETURNING *;

-- name: UpdateOrg :one
UPDATE orgs
SET name = $2, domain = LOWER($3), slug = $4
WHERE id = $1
RETURNING *;

-- name: GetOrgByDomain :one
SELECT *
FROM orgs
WHERE LOWER(domain) = LOWER($1);
