-- name: GetContentBySlug :one
SELECT *
FROM content
WHERE slug = $1
LIMIT 1;

-- name: ListContentByType :many
SELECT c.slug,
    c.type,
    COALESCE(SUM(cm.views), 0)::bigint AS views,
    COALESCE(SUM(cm.likes), 0)::bigint AS likes
FROM content c
    LEFT JOIN content_meta AS cm ON c.id = cm.content_id
WHERE (
        @type::text = ''
        OR c.type = @type
    )
GROUP BY c.id,
    c.slug
ORDER BY c.created_at;

-- name: GetContentStatsByType :one
SELECT COUNT(DISTINCT c.id)::bigint AS total_posts,
    COALESCE(SUM(cm.views), 0)::bigint AS total_views,
    COALESCE(SUM(cm.likes), 0)::bigint AS total_likes
FROM content AS c
    LEFT JOIN content_meta AS cm ON c.id = cm.content_id
WHERE (
        @type::text = ''
        OR c.type = @type
    );

-- name: UpsertContent :one
INSERT INTO content (slug, type)
VALUES ($1, $2) ON CONFLICT (slug) DO
UPDATE
SET slug = EXCLUDED.slug
RETURNING *;