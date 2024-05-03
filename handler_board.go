package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/google/uuid"
)

type event struct{
	Name string
}

func (apiCfg *apiConfig)handlerCreateBoard(w http.ResponseWriter, r *http.Request, user database.User){
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

	board, err := apiCfg.DB.CreateBoard(r.Context(), database.CreateBoardParams{
		ID: uuid.New(),
		Name: params.Name,
		Events: Events,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		OwnerID: user.ID,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error creating board: %v", err))
		return
	}

	respondWithJSON(w, 201, board)
}

func (apiCfg *apiConfig)handlerGetBoards(w http.ResponseWriter, r *http.Request, user database.User){
	boards, err := apiCfg.DB.GetUserBoards(r.Context(), user.ID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error getting boards: %v", err))
		return
	}

	respondWithJSON(w, 200, boards)
}