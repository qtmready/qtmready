-- name: CreateOrg :one
INSERT INTO orgs (name, domain, slug)
VALUES ($1, LOWER($2), $3)
RETURNING *;

-- name: GetOrgByDomain :one
SELECT *
FROM orgs
WHERE LOWER(domain) = LOWER($1);

-- name: GetOrgSlugByID :one
SELECT slug
FROM orgs
WHERE id = $1;

-- name: SetOrgHooks :exec
UPDATE orgs
SET hooks = $2
WHERE id = $1;

-- name: UpdateOrg :one
UPDATE orgs
SET name = $2, domain = LOWER($3), slug = $4
WHERE id = $1
RETURNING *;
