-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password, org_id)
VALUES ($1, $2, $3, $4, $5)
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

-- name: GetAuthUserByID :one
SELECT
	usr.*,
	json_agg(team.*) AS teams,
  json_agg(account.*) AS oauth_accounts,
  json_build_object(
    'id', org.id,
    'created_at', org.created_at,
    'updated_at', org.updated_at,
    'name', org.name,
    'domain', org.domain,
    'slug', org.slug
  ) AS org
FROM users AS usr
LEFT JOIN team_users AS tu
  ON usr.id = tu.user_id
LEFT JOIN teams AS team
  ON tu.team_id = team.id
LEFT JOIN oauth_accounts AS account
  ON usr.id = account.user_id
JOIN orgs AS org
  ON usr.org_id = org.id
WHERE
  usr.id = $1
GROUP BY
  usr.id, org.id;

-- name: GetAuthUserByEmail :one
SELECT
	usr.*,
	json_agg(team.*) AS teams,
  json_agg(account.*) AS oauth_accounts,
  json_build_object(
    'id', org.id,
    'created_at', org.created_at,
    'updated_at', org.updated_at,
    'name', org.name,
    'domain', org.domain,
    'slug', org.slug
  ) AS org
FROM users AS usr
LEFT JOIN team_users AS tu
  ON usr.id = tu.user_id
LEFT JOIN teams AS team
  ON tu.team_id = team.id
LEFT JOIN oauth_accounts AS account
  ON usr.id = account.user_id
JOIN orgs AS org
  ON usr.org_id = org.id
WHERE
  usr.email = LOWER($1)
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
SET first_name = $2, last_name = $3, email = LOWER($4), org_id = $5
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1;
