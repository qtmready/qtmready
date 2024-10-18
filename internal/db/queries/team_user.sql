-- name: CreateTeamUser :one
INSERT INTO team_users (team_id, user_id, role, is_active, is_admin)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, team_id, user_id, role, is_active, is_admin;

-- name: GetTeamUser :one
SELECT *
FROM team_users
WHERE user_id = $1;
