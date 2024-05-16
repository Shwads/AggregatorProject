package main

import (
	"AggregatorProject/internal/database"
	"context"
	"net/http"
)

func (apicfg apiConfig) authedGetUserPosts(res http.ResponseWriter, req *http.Request, user database.User) {
	posts, getPostsErr := apicfg.DB.GetPostsByUser(context.Background(), user.ID)
	if getPostsErr != nil {
		respondWithError(res, 500, "Something went wrong getting your posts.")
		return
	}

	respondWithJSON(res, 201, posts)
}
