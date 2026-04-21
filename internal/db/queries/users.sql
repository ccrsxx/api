-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserWithAccountByID :one
SELECT *
FROM users AS u
    JOIN account AS a ON u.id = a.user_id
WHERE u.id = $1
LIMIT 1;

-- name: GetAccountByProvider :one
SELECT *
FROM account
WHERE provider = $1
    AND provider_account_id = $2
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (name, image, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $1,
    image = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: CreateAccount :one
INSERT INTO account (user_id, provider, provider_account_id, username)
VALUES ($1, $2, $3, $4)
RETURNING *;