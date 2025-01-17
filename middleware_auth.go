package main

import (
	"fmt"
	"net/http"

	"github.com/Jkrish1011/rss-aggregator/internal/auth"
	"github.com/Jkrish1011/rss-aggregator/internal/database"
)

// Middleware logic to handle the authentication of the user request.

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	// Defining a closure to handle the request since the type we created have an extra input parameter and differs from the http.HandlerFunc signature
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth Error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get the user: %v", err))
		}
		handler(w, r, user)
	}
}
