package main

import (
	"context"
	"net/http"
)

func (apicfg apiConfig) getFeeds(res http.ResponseWriter, req *http.Request) {
	feeds, getFeedsErr := apicfg.DB.GetAllFeeds(context.Background())
	if getFeedsErr != nil {
		respondWithError(res, 501, getFeedsErr.Error())
		return
	}

	responseFeeds := Map(feeds, parseFeedStruct)

	respondWithJSON(res, 201, responseFeeds)
}
