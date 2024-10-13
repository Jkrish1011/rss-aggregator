package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("::::RSS AGGREGATOR::::")
	godotenv.Load(".env")
	// godotenv.Load() - This will also work
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment variable. Set your environment variable")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		// AllowedCredentials: false,
		MaxAge: 300,
	}))

	// Now to connect this router to a http server.
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	// We need to hook up our handler function
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerRediness)
	v1Router.Get("/err", handlerErr)
	// Hooking up a router under v1 so that it is easy to debug, make new release etc.
	router.Mount("/v1", v1Router)

	log.Printf("::::Server starting on port : %s::::\n", port)
	// ListenAndServe Function will block, i.e., it works like that.
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
