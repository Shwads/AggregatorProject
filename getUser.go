package main

import (
	"AggregatorProject/internal/database"
	"net/http"
)

func authedGetUser(res http.ResponseWriter, req *http.Request, user database.User) {
	respondWithJSON(res, 201, user)
}
