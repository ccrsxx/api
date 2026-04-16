-- name: UpsertIPAddress :one
INSERT INTO ip_address (ip_address)
VALUES ($1)
ON CONFLICT (ip_address) DO UPDATE
SET ip_address = EXCLUDED.ip_address
RETURNING *;
