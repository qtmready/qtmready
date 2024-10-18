-- name: GetTeam :one
SELECT *
FROM teams
WHERE id = $1
LIMIT 1;

-- name: CreateTeam :one
INSERT INTO teams (name) 
VALUES ($1) 
RETURNING id, name;

-- name: GetTeamByName :one
SELECT id, name 
FROM teams 
WHERE name = $1;

-- name: UpdateTeam :one
UPDATE teams
SET org_id = $2, name = $3, slug = $4
WHERE id = $1
RETURNING id, created_at, name, slug;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;