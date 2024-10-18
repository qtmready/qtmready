-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;


-- name: GetUserByEmailFull :one
SELECT
  u.*,
  array_agg(t.*) AS teams,
  array_agg(oa.*) AS oauth_accounts,
  array_agg(o.*) AS orgs
FROM users AS u
LEFT JOIN team_users AS tu
  ON u.id = tu.user_id
LEFT JOIN teams AS t
  ON tu.team_id = t.id
LEFT JOIN oauth_accounts AS oa
  ON u.id = oa.user_id
LEFT JOIN orgs AS o
  ON u.org_id = o.id
WHERE
  u.email = $1
GROUP BY
  u.id;

-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password) 
VALUES ($1, $2, $3, $4) 
RETURNING id, first_name, last_name, email;

-- name: GetUsersByEmail :many
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByProviderAccount :one
SELECT 
  u.*, 
  array_agg(oa.*) AS oauth_accounts
FROM users u
INNER JOIN oauth_accounts a ON u.id = a.user_id
WHERE a.provider = $1 AND a.provider_account_id = $2;
