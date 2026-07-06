-- name: CreateContact :one
INSERT INTO contact (name, email, message)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateContactDeliveredAtByID :one
UPDATE contact
SET delivered_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;