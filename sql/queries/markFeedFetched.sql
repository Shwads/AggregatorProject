-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $2, updated_at = $3
WHERE id = $1
RETURNING *;