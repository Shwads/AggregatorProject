// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: getFeed.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getFeedByID = `-- name: GetFeedByID :one
SELECT id, created_at, updated_at, name, url, user_id FROM feeds
WHERE id = $1
`

func (q *Queries) GetFeedByID(ctx context.Context, id uuid.UUID) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByID, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}
