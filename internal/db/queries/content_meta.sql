-- name: GetContentLikeStatus :one
SELECT COALESCE(SUM(likes), 0)::bigint AS likes,
    COALESCE(
        SUM(likes) FILTER (
            WHERE ip_address_id = $1
        ),
        0
    )::bigint AS user_likes
FROM content_meta
WHERE content_id = $2;

-- name: GetTotalContentMeta :one
SELECT COALESCE(SUM(views), 0)::bigint AS total_views,
    COALESCE(SUM(likes), 0)::bigint AS total_likes
FROM content_meta
WHERE content_id = $1;

-- name: IncrementContentView :one
INSERT INTO content_meta (content_id, ip_address_id, views)
VALUES ($1, $2, 1) ON CONFLICT (content_id, ip_address_id) DO
UPDATE
SET views = content_meta.views + 1,
    updated_at = NOW()
RETURNING views,
    likes;

-- name: IncrementContentLike :one
INSERT INTO content_meta (content_id, ip_address_id, likes)
VALUES ($1, $2, 1) ON CONFLICT (content_id, ip_address_id) DO
UPDATE
SET likes = LEAST(content_meta.likes + 1, 5),
    updated_at = NOW()
RETURNING views,
    likes;