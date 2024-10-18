-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (user_id, provider, provider_account_id, expires_at, type)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, provider, provider_account_id, expires_at, type;

-- name: GetOAuthAccountByID :one
SELECT *
FROM oauth_accounts
WHERE id = $1;

-- name: GetOAuthAccountsByUserID :many
SELECT *
FROM oauth_accounts
WHERE user_id = $1;

-- name: GetOAuthAccountsByProviderAccountID :many
SELECT *
FROM oauth_accounts
WHERE provider_account_id = $1;
