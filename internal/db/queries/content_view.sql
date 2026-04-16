-- name: GetContentViewCount :one
SELECT COUNT(*)::int AS views
FROM content_view cv
JOIN content c ON c.id = cv.content_id
WHERE c.slug = $1;
-- name: CreateContentView :exec
INSERT INTO content_view (content_id, ip_address_id)
VALUES ($1, $2);
