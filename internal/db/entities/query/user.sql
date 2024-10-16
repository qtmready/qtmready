-- name: GetAuthor :one

SELECT *
FROM users
WHERE id = $1
LIMIT 1;
