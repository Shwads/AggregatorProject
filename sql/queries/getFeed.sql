-- name: GetFeedByID :one
SELECT * FROM feeds
WHERE id = $1;