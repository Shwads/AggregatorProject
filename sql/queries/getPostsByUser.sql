-- name: GetPostsByUser :many
SELECT p.* FROM posts p
JOIN feed_follows ff ON p.feed_id = ff.feed_id
JOIN users u ON ff.user_id = u.id
WHERE u.id = $1
ORDER BY p.published_at DESC
LIMIT $2;