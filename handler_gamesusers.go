package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)


func (apiCfg *apiConfig)handlerCreateGamesUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	
	type parameters struct{
		GameID uuid.UUID `json:"game_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	_, err = apiCfg.DB.GetGame(r.Context(), params.GameID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting game: %v", err))
		return
	}

	gameuser, err := apiCfg.DB.CreateGameUser(r.Context(), database.CreateGameUserParams{
		ID: uuid.New(),
		UserID: user.ID,
		GameID: params.GameID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error creating game: %v", err))
		return
	}

	respondWithJSON(w, 200, gameuser)
}

func (apiCfg *apiConfig)handlerGetGamesUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	
	gameusers, err := apiCfg.DB.GetGameUsers(r.Context(), user.ID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting user games: %v", err))
		return
	}

	respondWithJSON(w, 200, gameusers)
}

func (apiCfg *apiConfig)handlerDeleteGamesUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	
	gamesUsersIDstr := chi.URLParam(r, "gamesuesersID")

	gamesUsersID, err := uuid.Parse(gamesUsersIDstr)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing UUID: %v", err))
		return
	}

	err = apiCfg.DB.DeleteGameUsers(r.Context(), database.DeleteGameUsersParams{
		ID: gamesUsersID,
		UserID: user.ID,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error deleting game user: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
