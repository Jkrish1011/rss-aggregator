package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	/*
		json.Marshal in Go is a function used to encode a Go value to JSON. It takes a Go value as input and returns a byte slice representing the JSON-encoded data.
		- Input: Can accept any Go value that can be represented as JSON, including basic data types (e.g., int, string, bool), structs, slices, and maps.
		- Output: Returns a byte slice containing the JSON-encoded data.
		- Error Handling: Returns an error if the input value cannot be encoded as JSON (e.g., due to unsupported types or circular references).
	*/
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal JSON response: %v\n", payload)
		w.WriteHeader(500)
		return
	}
	// Responding with particular Content-Type
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	// For all internal bugs/problems
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	type ErrorResponse struct {
		Error string `json:"error"` // indicating the when marshalling, the key of this string will be error
	}
	respondWithJSON(w, code, ErrorResponse{
		Error: msg,
	})
}
