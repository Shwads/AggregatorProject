package main

import "AggregatorProject/internal/database"

func parseFeedStruct(feed database.Feed) responseFeed {
	response := responseFeed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}

	if feed.LastFetchedAt.Valid {
		response.LastFetchedAt = feed.LastFetchedAt.Time
	}

	return response
}
