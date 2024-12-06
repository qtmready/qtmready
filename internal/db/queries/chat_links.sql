-- name: CreateChatLink :one
INSERT INTO chat_links (hook, kind, link_to, data)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetChatLink :one
SELECT *
FROM chat_links
WHERE link_to = $1;
