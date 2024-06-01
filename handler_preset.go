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
	ID string
	Name string
	Description string
}


//This repeats, maybe create utils

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

	for i, _ := range(params.Events){
		eventID, err := randomHex(5)

		if err != nil{
			respondWithError(w, 400, fmt.Sprintf("Error generating id: %v", err))
			return
		}
		params.Events[i].ID = eventID
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

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing UUID: %v", err))
		return
	}

	preset, err := apiCfg.DB.GetUserPresetByID(r.Context(), presetID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting preset: %v", err))
		return
	}

	respondWithJSON(w, 200, preset)
}

func (apiCfg *apiConfig)handlerEditPreset(w http.ResponseWriter, r *http.Request, user database.User){

	presetIDstr := chi.URLParam(r, "presetid")

	presetID, err := uuid.Parse(presetIDstr)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing UUID: %v", err))
		return
	}

	type parameters struct{
		Events []event `json:"events"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}


	preset, err := apiCfg.DB.GetUserPresetByID(r.Context(), presetID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting presets: %v", err))
		return
	}

	if preset.OwnerID != user.ID{
		respondWithError(w, 400, "User not authorized")
		return
	}

	presetEventsJson := preset.Events
	var presetEvents []event

	err = json.Unmarshal(presetEventsJson, &presetEvents)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	for i, event := range(presetEvents){
		for _, eventToEdit := range(params.Events){
			if eventToEdit.ID == event.ID{
				presetEvents[i] = eventToEdit
			}
		}
	}

	eventsJson, err := json.Marshal(presetEvents)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	newPreset, err := apiCfg.DB.UpdatePresetEvents(r.Context(), database.UpdatePresetEventsParams{
		Events: eventsJson,
		ID: presetID,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error updating preset: %v", err))
		return
	}

	respondWithJSON(w, 200, newPreset)
}