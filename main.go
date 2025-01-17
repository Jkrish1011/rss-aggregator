package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Jkrish1011/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	log.Println("::::RSS AGGREGATOR::::")
	godotenv.Load(".env")
	// godotenv.Load() - This will also work
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment variable. Set your environment variable::PORT")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Database URL is not found in the environment variable. Set your environment variable::DB_URL")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to Database:", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Now to connect this router to a http server.
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	v1Router := chi.NewRouter()
	// to check if the server is alive and is running.
	v1Router.Get("/healthz", handlerRediness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	// Hooking up a router under v1 so that it is easy to debug, make new release etc.
	router.Mount("/v1", v1Router)

	log.Printf("::::Server starting on port : %s::::\n", port)
	// ListenAndServe Function will block, i.e., it works like that.
	err1 := server.ListenAndServe()
	if err1 != nil {
		log.Fatal(err)
	}
}
