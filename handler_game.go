package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/google/uuid"
)

// handlerCreateGame handles the HTTP request for creating a new game.
// It expects a JSON payload in the request body with the following parameters:
// - name: The name of the game (string)
// - end_time: The end time of the game (time.Time)
// - preset_id: The ID of the preset for the game (uuid.UUID)
//
// Example HTTP request:
// POST /games
// Body:
// {
//   "name": "My Game",
//   "end_time": "2022-12-31T23:59:59Z",
//   "preset_id": "123e4567-e89b-12d3-a456-426614174000"
// }
//
// The function decodes the JSON payload, validates the parameters, and creates a new game in the database.
// If successful, it returns the created game as a JSON response with status code 200.
// If there is an error parsing the JSON or creating the game, it returns an error response with status code 400.
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