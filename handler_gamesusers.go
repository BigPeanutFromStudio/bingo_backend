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
	

	gamesUsersIDstr := chi.URLParam(r, "gamesuesersID")

	gamesUsersID, err := uuid.Parse(gamesUsersIDstr)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing UUID: %v", err))
		return
	}

	_, err = apiCfg.DB.GetGame(r.Context(), gamesUsersID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting game: %v", err))
		return
	}

	gameuser, err := apiCfg.DB.CreateGameUser(r.Context(), database.CreateGameUserParams{
		ID: uuid.New(),
		UserID: user.ID,
		GameID: gamesUsersID,
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

func (apiCfg *apiConfig)handlerAddUsersToGame(w http.ResponseWriter, r *http.Request, user database.User) {
	
	gameIDstr := chi.URLParam(r, "gameID")

	gameID, err := uuid.Parse(gameIDstr)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing UUID: %v", err))
		return
	}

	type parameters struct{
		UserID []string `json:"user_nicknames"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	game, err := apiCfg.DB.GetGame(r.Context(), gameID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting game: %v", err))
		return
	}

	if game.AdminID != user.ID{
		respondWithError(w, 400, "User is not an admin of the game")
		return
	}

	var gameusers []database.GamesUser
	for _, userNickname := range(params.UserID){

		userToAdd, err := apiCfg.DB.GetUserByNickname(r.Context(), userNickname)

		if err != nil{
			respondWithError(w, 400, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		gameuser, err := apiCfg.DB.CreateGameUser(r.Context(), database.CreateGameUserParams{
			ID: uuid.New(),
			UserID: userToAdd.ID,
			GameID: game.ID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})

		if err != nil{
			respondWithError(w, 400, fmt.Sprintf("Error creating gameuser: %v", err))
			return
		}
		
		gameusers = append(gameusers, gameuser)
	}


	respondWithJSON(w, 200, gameusers)
}


