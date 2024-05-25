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

type event struct{
	Name string
	Description string
}

func (apiCfg *apiConfig)handlerCreatePreset(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct{
		Name string `json:"name"`
		Events []event `json:"events"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	Events, err := json.Marshal(params.Events)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	preset, err := apiCfg.DB.CreatePreset(r.Context(), database.CreatePresetParams{
		ID: uuid.New(),
		Name: params.Name,
		Events: Events,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		OwnerID: user.ID,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error creating preset: %v", err))
		return
	}

	respondWithJSON(w, 201, preset)
}

func (apiCfg *apiConfig)handlerGetPresets(w http.ResponseWriter, r *http.Request, user database.User){
	presets, err := apiCfg.DB.GetUserPresets(r.Context(), user.ID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting presets: %v", err))
		return
	}

	respondWithJSON(w, 200, presets)
}

func (apiCfg *apiConfig)handlerGetPresetByID(w http.ResponseWriter, r *http.Request){

	presetIDstr := chi.URLParam(r, "presetid")

	presetID, err := uuid.Parse(presetIDstr)

	presets, err := apiCfg.DB.GetUserPresetByID(r.Context(), presetID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting presets: %v", err))
		return
	}

	respondWithJSON(w, 200, presets)
}