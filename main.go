package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *database.Queries
}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT variable not found in environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatal("DB_URL variable not found in environment")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil{
		log.Fatal("error connecting to database: ", err)
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

	v1Router := chi.NewRouter()

	// Goes first
	v1Router.Use(middleware.Logger)

	//Handlers
	v1Router.Post("/generate", generateBoard)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

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
	err = srv.ListenAndServe()
	if err != nil{
		log.Fatal(err)
	}
}