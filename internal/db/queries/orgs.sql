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

-- name: CreateDefaultOrg :one
INSERT INTO orgs (id, name, domain, slug)
VALUES ('00000000-0000-0000-0000-000000000001', 'No Org', 'example.com', 'no-org')
ON CONFLICT (id) DO NOTHING
RETURNING id, created_at, updated_at, name, domain, slug;

