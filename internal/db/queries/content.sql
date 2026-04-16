-- name: GetContentBySlug :one
SELECT *
FROM content
WHERE slug = $1
LIMIT 1;
-- name: UpsertContent :one
INSERT INTO content (slug, type)
VALUES ($1, $2)
ON CONFLICT (slug) DO UPDATE
SET slug = EXCLUDED.slug
RETURNING *;
-- name: ListContentByType :many
SELECT c.slug,
       COUNT(DISTINCT cv.id)::int AS views,
       COUNT(cl.id)::int AS likes
FROM content c
LEFT JOIN content_view cv ON cv.content_id = c.id
LEFT JOIN content_like cl ON cl.content_id = c.id
WHERE c.type = $1
GROUP BY c.id
ORDER BY c.created_at;
-- name: GetContentStatsByType :one
SELECT COUNT(DISTINCT c.id)::int AS total_posts,
       COUNT(DISTINCT cv.id)::int AS total_views,
       COUNT(cl.id)::int AS total_likes
FROM content c
LEFT JOIN content_view cv ON cv.content_id = c.id
LEFT JOIN content_like cl ON cl.content_id = c.id
WHERE c.type = $1;
