package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/google/uuid"
)


func (apiCfg *apiConfig)handlerCreateGame(w http.ResponseWriter, r *http.Request, user database.User) {
	
	type parameters struct{
		Name string `json:"name"`
		EndTime time.Time `json:"end_time"`
		PresetID uuid.UUID `json:"preset_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	game, err := apiCfg.DB.CreateGame(r.Context(), database.CreateGameParams{
		ID: uuid.New(),
		Name: params.Name,
		EndTime: params.EndTime.UTC(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Preset: params.PresetID,
		AdminID: user.ID,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error creating game: %v", err))
		return
	}

	respondWithJSON(w, 200, game)
}