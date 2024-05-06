package main

import (
	"AggregatorProject/internal/database"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apicfg apiConfig) authedCreateFeed(res http.ResponseWriter, req *http.Request, user database.User) {
	newUUID := uuid.New()

	currTime := time.Now()

	decoder := json.NewDecoder(req.Body)

	receivedFeed := feed{}

	decodeErr := decoder.Decode(&receivedFeed)
	if decodeErr != nil {
		respondWithError(res, 401, decodeErr.Error())
		return
	}

	feedParams := database.CreateFeedParams{
		ID:        newUUID,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      receivedFeed.Name,
		Url:       receivedFeed.Url,
		UserID:    user.ID,
	}

	feed, createFeedErr := apicfg.DB.CreateFeed(context.Background(), feedParams)
	if createFeedErr != nil {
		respondWithError(res, 501, "Failed to create feed.")
		return
	}

	feedFollow, followErr := apicfg.createFollow(res, newUUID, user.ID)
	if followErr != nil {
		return
	}

	responseBody := struct {
		Feed   database.Feed       `json:"feed"`
		Follow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:   feed,
		Follow: feedFollow,
	}

	respondWithJSON(res, 201, responseBody)
}
