package main

import (
	"time"

	"github.com/google/uuid"
)

type responseFeed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt time.Time
}
