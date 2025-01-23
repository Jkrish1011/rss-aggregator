package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Jkrish1011/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%s", err))
		return
	}

	feedFollows, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create Feed Follow request :%s", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollows))
}
