package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type apiConfig struct{
	//DB *database.Queries
	//but first sqlc generate and goose up to migrate
	//but first first create schema and queries
}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT variable not found in environment")
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

	v1Router := chi.NewRouter()

	// Goes first
	v1Router.Use(middleware.Logger)

	//Handlers
	v1Router.Post("/generate", generateBoard)

	router.Mount("/v1", v1Router)

	// 404 & 405 handling
	router.NotFound(func(w http.ResponseWriter, r *http.Request){
		respondWithError(w, 404, "route does not exist")
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request){
		respondWithError(w, 405, "method is not valid")
	})

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe()
	if err != nil{
		log.Fatal(err)
	}
}