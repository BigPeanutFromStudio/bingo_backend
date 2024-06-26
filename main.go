package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BigPeanutFromStudio/bingo/internal/auth"
	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/olahol/melody"

	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *database.Queries
}

//IMPORTANT: FIX IRREGULAR NAMING CONVENTION IN DATABASE
//BOARDS TABLE, EVENTS IN JSON, ONE TO MANY
//MAYBE DON'T STORE THE GOOGLE ID??

func main() {

	auth.NewAuth()
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT variable not found in environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatal("DB_URL variable not found in environment")
	}

	m := melody.New()

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
	v1Router.Post("/presets", apiCfg.middlewareAuth(apiCfg.handlerCreatePreset))
	v1Router.Get("/presets", apiCfg.middlewareAuth(apiCfg.handlerGetPresets))
	v1Router.Get("/presets/{presetid}", apiCfg.handlerGetPresetByID)
	v1Router.Put("/presets/{presetid}", apiCfg.middlewareAuth(apiCfg.handlerEditPreset))

	//v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Get("/users/list", apiCfg.handlerGetAllUsers) //temporary
	v1Router.Put("/users", apiCfg.middlewareAuth(apiCfg.handlerSetGoogleUserNickname))

	v1Router.Post("/games", apiCfg.middlewareAuth(apiCfg.handlerCreateGame)) //REWORK
	
	//v1Router.Post("/games/join/{gamesuesersID}", apiCfg.middlewareAuth(apiCfg.handlerCreateGamesUsers))
	v1Router.Put("/games/join/{gameID}", apiCfg.middlewareAuth(apiCfg.handlerAddUsersToGame)) //REWORK
	v1Router.Get("/games", apiCfg.middlewareAuth(apiCfg.handlerGetGamesUsers))
	v1Router.Get("/games/admin", apiCfg.middlewareAuth(apiCfg.handlerGetAdminedGames))
	v1Router.Delete("/games/{gamesuesersID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteGamesUsers))

	v1Router.Get("/auth/{provider}/callback", apiCfg.getAuthCallbackFunction)
	v1Router.Get("/logout", logoutHandler)
	v1Router.Get("/auth/{provider}", beginAuthHandler)

	v1Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})


	//figure out how to respond only to the messenger
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if(string(msg) == "lol"){
			m.Broadcast([]byte("Not funny"))
		}else{
			m.Broadcast(msg)
		}
		fmt.Printf("Handled message %v", string(msg))
	})


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