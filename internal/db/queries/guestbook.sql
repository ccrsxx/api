-- name: CreateGuestbook :one
WITH new_guestbook AS (
    INSERT INTO guestbook (text, user_id)
    VALUES ($1, $2)
    RETURNING *
)
SELECT g.id,
    g.text,
    u.name,
    u.image,
    a.username,
    g.created_at
FROM new_guestbook AS g
    JOIN users AS u ON g.user_id = u.id
    LEFT JOIN account AS a ON u.id = a.user_id;

-- name: GetGuestbookByID :one
SELECT *
FROM guestbook
WHERE id = $1
LIMIT 1;

-- name: DeleteGuestbook :exec
DELETE FROM guestbook
WHERE id = $1;

-- name: ListGuestbook :many
SELECT g.id,
    g.text,
    u.name,
    u.image,
    a.username,
    g.created_at
FROM guestbook AS g
    JOIN users AS u ON g.user_id = u.id
    LEFT JOIN account AS a ON u.id = a.user_id
ORDER BY g.created_at DESC;