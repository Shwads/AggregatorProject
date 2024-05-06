package main

import (
	"AggregatorProject/internal/database"
	"context"
	"net/http"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apicfg apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		apiKey := req.Header.Get("Authorization")
		authParams := strings.Split(apiKey, " ")
		if len(authParams) < 2 {
			respondWithError(res, 501, "Please provide a valid auth header of the format 'Authorization: ApiKey XXX' where XXX is your api key.")
			return
		}
		apiKey = authParams[1]

		user, getUserErr := apicfg.DB.GetUserByKey(context.Background(), apiKey)
		if getUserErr != nil {
			respondWithError(res, 501, getUserErr.Error())
			return
		}

		handler(res, req, user)
	})
}
