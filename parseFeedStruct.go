package main

import "AggregatorProject/internal/database"

// sqlc gives us the LastFetchedAt timme as a sql.NullTime value which results in a nested object.
// To clean up the response, we convert our Feed object into a new object when sending it in a response.
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
