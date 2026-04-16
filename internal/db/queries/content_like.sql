-- name: GetContentLikeStatus :one
SELECT COUNT(cl.id)::int AS likes,
       COUNT(cl.id) FILTER (WHERE cl.ip_address_id = @ip_address_id)::int AS user_likes
FROM content_like cl
JOIN content c ON c.id = cl.content_id
WHERE c.slug = @slug;
-- name: CreateContentLike :exec
INSERT INTO content_like (content_id, ip_address_id)
VALUES ($1, $2);
-- name: GetUserLikeCount :one
SELECT COUNT(*)::int AS user_likes
FROM content_like cl
JOIN content c ON c.id = cl.content_id
WHERE c.slug = $1
  AND cl.ip_address_id = $2;
