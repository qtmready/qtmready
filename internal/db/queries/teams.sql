-- name: GetTeam :one
SELECT *
FROM teams
WHERE id = $1
LIMIT 1;

-- name: CreateTeam :one
INSERT INTO teams (name, org_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetTeamBySlug :one
SELECT id, name
FROM teams
WHERE slug = $1;

-- name: GetTeamByID :one
SELECT *
FROM teams
WHERE id = $1;

-- name: UpdateTeam :one
UPDATE teams
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;
