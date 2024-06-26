// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: getUser.sql

package database

import (
	"context"
)

const getUserByKey = `-- name: GetUserByKey :one
SELECT id, created_at, updated_at, name, api_key FROM users
WHERE api_key = $1
`

func (q *Queries) GetUserByKey(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByKey, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ApiKey,
	)
	return i, err
}
