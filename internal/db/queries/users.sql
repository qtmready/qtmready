-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password, picture, org_id)
VALUES ($1, $2, LOWER($3), $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = lower($1);

-- name: GetAuthUserByID :one
SELECT
  json_build_object(
    'id', usr.id,
    'created_at', usr.created_at,
    'updated_at', usr.updated_at,
    'org_id', usr.org_id,
    'email', usr.email,
    'first_name', usr.first_name,
    'last_name', usr.last_name,
    'picture', usr.picture,
    'is_active', usr.is_active,
    'is_verified', usr.is_verified
  ) AS user,
  json_build_object(
    'id', org.id,
    'created_at', org.created_at,
    'updated_at', org.updated_at,
    'name', org.name,
    'domain', org.domain,
    'slug', org.slug
  ) AS org,
  json_agg(role.name) AS roles,
  json_agg(team.*) AS teams,
  json_agg(account.*) AS oauth_accounts
FROM users AS usr
LEFT JOIN team_users AS tu
  ON usr.id = tu.user_id
LEFT JOIN teams AS team
  ON tu.team_id = team.id
LEFT JOIN oauth_accounts AS account
  ON usr.id = account.user_id
JOIN orgs AS org
  ON usr.org_id = org.id
LEFT JOIN user_roles AS role
  ON usr.id = role.user_id
WHERE
  usr.id = $1
GROUP BY
  usr.id, org.id;

-- name: GetAuthUserByEmail :one
SELECT
  json_build_object(
    'id', usr.id,
    'created_at', usr.created_at,
    'updated_at', usr.updated_at,
    'org_id', usr.org_id,
    'email', usr.email,
    'first_name', usr.first_name,
    'last_name', usr.last_name,
    'picture', usr.picture,
    'is_active', usr.is_active,
    'is_verified', usr.is_verified
  ) AS user,
  json_build_object(
    'id', org.id,
    'created_at', org.created_at,
    'updated_at', org.updated_at,
    'name', org.name,
    'domain', org.domain,
    'slug', org.slug
  ) AS org,
  json_agg(role.name) AS roles,
  json_agg(team.*) AS teams,
  json_agg(account.*) AS oauth_accounts
FROM users AS usr
LEFT JOIN team_users AS tu
  ON usr.id = tu.user_id
LEFT JOIN teams AS team
  ON tu.team_id = team.id
LEFT JOIN oauth_accounts AS account
  ON usr.id = account.user_id
JOIN orgs AS org
  ON usr.org_id = org.id
LEFT JOIN user_roles AS role
  ON usr.id = role.user_id
WHERE
  usr.email = lower($1)
GROUP BY
  usr.id, org.id;

-- name: GetUserByProviderAccount :one
SELECT
  usr.*
FROM
  users usr
JOIN
  oauth_accounts act ON usr.id = act.user_id
WHERE
  act.provider = $1 AND act.provider_account_id = $2;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3, email = lower($4), org_id = $5
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1;
