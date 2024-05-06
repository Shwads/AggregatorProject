package main

import (
	"AggregatorProject/internal/database"
	"context"
	"net/http"
)

func (apicfg apiConfig) authedGetFeedFollows(res http.ResponseWriter, req *http.Request, user database.User) {
	feedFollows, getFeedFollowsErr := apicfg.DB.GetUserFeedFollows(context.Background(), user.ID)
	if getFeedFollowsErr != nil {
		respondWithError(res, 501, "Failed to fetch user data.")
		return
	}

	respondWithJSON(res, 201, feedFollows)
}
