package main

import (
	"AggregatorProject/internal/database"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (apicfg apiConfig) authedDeleteFeedFollow(res http.ResponseWriter, req *http.Request, user database.User) {
	feedFollowId := req.PathValue("feedFollowID")

	uuidFeedFollowId, feedFollowIdParseErr := uuid.Parse(feedFollowId)
	if feedFollowIdParseErr != nil {
		respondWithError(res, 403, "Please provide a valid feedFollowId.")
		return
	}

	deleteErr := apicfg.DB.DeleteFeedFollow(context.Background(), uuidFeedFollowId)
	if deleteErr != nil {
		respondWithError(res, 501, fmt.Sprintf("Error deleting feed follow: %s", feedFollowId))
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: "OK. Delete successful.",
	}

	respondWithJSON(res, 201, response)
}
