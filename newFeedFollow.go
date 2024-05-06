package main

import (
	"AggregatorProject/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apicfg apiConfig) authedNewFeedFollow(res http.ResponseWriter, req *http.Request, user database.User) {

	decoder := json.NewDecoder(req.Body)

	feed := struct {
		FeedId string `json:"feed_id"`
	}{}

	bodyDecodeErr := decoder.Decode(&feed)
	if bodyDecodeErr != nil {
		respondWithError(res, 401, "Error decoding request body.")
		return
	}

	uuidFeedId, feedIdParseErr := uuid.Parse(feed.FeedId)
	if feedIdParseErr != nil {
		respondWithError(res, 401, "Please provide a valid feed ID.")
		return
	}

	feedFollow, followErr := apicfg.createFollow(res, uuidFeedId, user.ID)
	if followErr != nil {
		return
	}

	respondWithJSON(res, 201, feedFollow)
}

func (apicfg apiConfig) createFollow(res http.ResponseWriter, feedId uuid.UUID, userId uuid.UUID) (database.FeedFollow, error) {
	_, feedFindErr := apicfg.DB.GetFeedByID(context.Background(), feedId)
	if feedFindErr == sql.ErrNoRows {
		respondWithError(res, 401, "No such feed available.")
		return database.FeedFollow{}, feedFindErr
	} else if feedFindErr != nil {
		respondWithError(res, 501, "Oops something went wrong.")
		return database.FeedFollow{}, feedFindErr
	}

	newUUID := uuid.New()
	timeNow := time.Now()

	feedFollowParams := database.NewFeedFollowParams{
		ID:        newUUID,
		FeedID:    feedId,
		UserID:    userId,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	feedFollow, feedFollowErr := apicfg.DB.NewFeedFollow(context.Background(), feedFollowParams)
	if feedFollowErr != nil {
		respondWithError(res, 501, fmt.Sprintf("Failed to follow feed %s", feedId))
		log.Printf("Encountered error in feed follow creation")
		return database.FeedFollow{}, feedFollowErr
	}

	return feedFollow, nil
}
