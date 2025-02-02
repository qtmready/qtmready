-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (user_id, provider, provider_account_id, expires_at, type)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetOAuthAccountByID :one
SELECT *
FROM oauth_accounts
WHERE id = $1;

-- name: GetOAuthAccountsByUserID :many
SELECT *
FROM oauth_accounts
WHERE user_id = $1;

-- name: GetOAuthAccountByProviderAccountID :one
SELECT *
FROM oauth_accounts
WHERE provider_account_id = $1 and provider = $2;
