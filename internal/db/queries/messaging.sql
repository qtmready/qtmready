-- name: CreateMessaging :one
INSERT INTO messaging (provider, kind, link_to, data)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at, provider, kind, link_to, data;

