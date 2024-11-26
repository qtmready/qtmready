-- name: CreateMessaging :one
INSERT INTO messaging (hook, kind, link_to, data)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMessagesByLinkTo :one
SELECT *
FROM messaging
WHERE link_to = $1; 
