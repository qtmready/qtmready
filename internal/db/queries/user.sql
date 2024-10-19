-- name: GetUserByID :one
SELECT id, created_at, updated_at, first_name, last_name, email, org_id
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, first_name, last_name, email, org_id
FROM users
WHERE email = $1;


-- name: GetUserByEmailFull :one
SELECT
  u.id, u.created_at, u.updated_at, u.first_name, u.last_name, u.email, u.org_id,
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
RETURNING id, created_at, updated_at, first_name, last_name, email, org_id;

-- name: GetUserByProviderAccount :one
SELECT
  u.id, u.created_at, u.updated_at, u.first_name, u.last_name, u.email, u.org_id
FROM users u
WHERE u.id IN (
  SELECT user_id
  FROM oauth_accounts
  WHERE provider = $1 AND provider_account_id = $2
);
