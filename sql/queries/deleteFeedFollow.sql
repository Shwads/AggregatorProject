-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE id = $1;