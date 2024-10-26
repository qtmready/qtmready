-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = LOWER($1);

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
  u.email = LOWER($1)
GROUP BY
  u.id;

-- name: GetFullUserByID :one
SELECT
  u.id,
  u.created_at,
  u.updated_at,
  u.first_name,
  u.last_name,
  u.email,
  u.org_id,
  ARRAY_AGG(
    DISTINCT ROW(
      team.id,
      team.created_at,
      team.updated_at,
      team.org_id,
      team.name,
      team.slug
    )
  ) AS teams,
  ARRAY_AGG(
    DISTINCT ROW(
      account.id,
      account.created_at,
      account.updated_at,
      account.user_id,
      account.provider,
      account.provider_account_id,
      account.expires_at,
      account.type
      )
  ) as accounts,
  ROW(org.id, org.created_at, org.updated_at, org.name, org.domain, org.slug) AS org
FROM users AS u
LEFT JOIN team_users AS team_user
  ON u.id = team_user.user_id
LEFT JOIN teams AS team
  ON team_user.team_id = team.id
LEFT JOIN oauth_accounts AS account
  ON u.id = account.user_id
LEFT JOIN orgs AS org
  ON u.org_id = org.id
WHERE
  u.id = $1
GROUP BY u.id, org.id, org.created_at, org.updated_at, org.name, org.domain, org.slug;

-- name: GetFullUserByEmail :one
SELECT
  u.id,
  u.created_at,
  u.updated_at,
  u.first_name,
  u.last_name,
  u.email,
  u.org_id,
  ARRAY_AGG(
    DISTINCT ROW(
      team.id,
      team.created_at,
      team.updated_at,
      team.org_id,
      team.name,
      team.slug
    )
  ) AS teams,
  ARRAY_AGG(
    DISTINCT ROW(
      account.id,
      account.created_at,
      account.updated_at,
      account.user_id,
      account.provider,
      account.provider_account_id,
      account.expires_at,
      account.type
      )
  ) as accounts,
  ROW(org.id, org.created_at, org.updated_at, org.name, org.domain, org.slug) AS org
FROM users AS u
LEFT JOIN team_users AS team_user
  ON u.id = team_user.user_id
LEFT JOIN teams AS team
  ON team_user.team_id = team.id
LEFT JOIN oauth_accounts AS account
  ON u.id = account.user_id
LEFT JOIN orgs AS org
  ON u.org_id = org.id
WHERE
  u.email = LOWER($1)
GROUP BY u.id, org.id, org.created_at, org.updated_at, org.name, org.domain, org.slug;

-- name: GetUserByProviderAccount :one
SELECT
  u.*
FROM users as u
WHERE u.id IN (
  SELECT user_id
  FROM oauth_accounts
  WHERE provider = $1 AND provider_account_id = $2
);

-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3, email = LOWER($4), org_id = $5
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1;
