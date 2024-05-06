package main

import (
	"AggregatorProject/internal/database"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apicfg apiConfig) createUser(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	receivedUser := user{}

	decodeErr := decoder.Decode(&receivedUser)
	if decodeErr != nil {
		respondWithError(res, 401, decodeErr.Error())
		return
	}

	newUUID := uuid.New()

	timeNow := time.Now()

	userParams := database.CreateUserParams{
		ID:        newUUID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Name:      receivedUser.Name,
	}

	user, createUserErr := apicfg.DB.CreateUser(context.Background(), userParams)
	if createUserErr != nil {
		respondWithError(res, 501, createUserErr.Error())
		log.Print("Failed to add item to database")
		return
	}

	respondWithJSON(res, 200, user)
}
