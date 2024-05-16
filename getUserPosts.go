package main

import (
	"AggregatorProject/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type postNum struct {
	Limit int `json:"limit"`
}

func (apicfg apiConfig) authedGetUserPosts(res http.ResponseWriter, req *http.Request, user database.User) {

	lim := postNum{}

	decoder := json.NewDecoder(req.Body)

	decoderErr := decoder.Decode(&lim)
	if decoderErr != nil {
		if decoderErr.Error() == "EOF" {
			log.Printf("We caught it, it was just an EOF!")
		} else {
			respondWithError(res, 402, fmt.Sprintf("%+v", decoderErr.Error()))
			return
		}
	}

	if lim.Limit == 0 {
		lim.Limit = 10
	}

	params := database.GetPostsByUserParams{
		ID:    user.ID,
		Limit: int32(lim.Limit),
	}

	posts, getPostsErr := apicfg.DB.GetPostsByUser(context.Background(), params)
	if getPostsErr != nil {
		respondWithError(res, 500, "Something went wrong getting your posts.")
		return
	}

	respondWithJSON(res, 201, posts)
}
