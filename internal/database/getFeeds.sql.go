// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: getFeeds.sql

package database

import (
	"context"
)

const getAllFeeds = `-- name: GetAllFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds
`

func (q *Queries) GetAllFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
