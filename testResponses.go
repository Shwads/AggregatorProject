package main

import "net/http"

func readinessHandler(res http.ResponseWriter, req *http.Request) {
	responseObj := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	respondWithJSON(res, 200, responseObj)
}

func errorHandler(res http.ResponseWriter, req *http.Request) {
	respondWithError(res, 500, "Internal Server Error")
}
