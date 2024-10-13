package main

import "net/http"

func handlerRediness(w http.ResponseWriter, r *http.Request) {
	// struct{}{} - will be an empty obj `{}`
	respondWithJSON(w, 200, struct{}{})
}
