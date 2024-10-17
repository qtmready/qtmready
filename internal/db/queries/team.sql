-- name: GetTeam :one

SELECT *
FROM teams
WHERE id = $1
LIMIT 1;
