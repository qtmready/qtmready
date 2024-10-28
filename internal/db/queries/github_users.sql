-- name: CreateGithubUser :one
INSERT INTO github_users (user_id, github_id, github_org_id, login)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGithubUserByID :one
SELECT *
FROM github_users
WHERE id = $1;

-- name: GetGithubUserByGithubID :one
SELECT *
FROM github_users
WHERE github_id = $1;

-- name: GetGithubUserByGithubOrgID :one
SELECT *
FROM github_users
WHERE github_org_id = $1;

-- name: GetGithubUserByUserID :one
SELECT *
FROM github_users
WHERE user_id = $1;

-- name: GetGithubUserByLogin :one
SELECT *
FROM github_users
WHERE login = $1;

-- name: UpdateGithubUser :one
UPDATE github_users
SET user_id = $2, github_id = $3, github_org_id = $4, login = $5
WHERE id = $1
RETURNING *;

-- name: DeleteGithubUser :one
DELETE FROM github_users
WHERE id = $1
RETURNING id;
