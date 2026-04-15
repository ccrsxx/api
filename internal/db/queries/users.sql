-- name: ListUsers :many
SELECT *
FROM users
ORDER BY name;
-- name: CreateUser :one
INSERT INTO users (id, name, email, image, role)
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
WHERE email = $1
LIMIT 1;
-- name: UpdateUser :one
UPDATE users
SET name = coalesce(sqlc.narg('name'), name),
    image = coalesce(sqlc.narg('image'), image),
    updated_at = now()
WHERE id = $1
RETURNING *;
-- name: UpdateUserRole :one
UPDATE users
SET role = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;