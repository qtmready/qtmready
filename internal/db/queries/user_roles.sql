-- name: CreateUserRole :one
INSERT INTO user_roles (name, user_id, org_id)
VALUES ($1, $2, $3) RETURNING *;
