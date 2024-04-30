package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(res http.ResponseWriter, statusCode int, errorMessage string) {
	err := struct {
		Error string `json:"error"`
	}{
		Error: errorMessage,
	}

	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		log.Printf("Encountered Error: %s. in func respondWithError(). Failed to marshal json data:", marshalErr)
		res.WriteHeader(500)
		return
	}

	res.WriteHeader(statusCode)
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func respondWithJSON(res http.ResponseWriter, statusCode int, payload interface{}) {
	jsonData, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		log.Printf("Encountered Error: %s. in func respondWithJSON(). Failed to marshal json data:", marshalErr)
		log.Print(payload)
		res.WriteHeader(500)
		return
	}

	res.WriteHeader(statusCode)
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}
