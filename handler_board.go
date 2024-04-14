package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type event struct{
	Name string
}

func generateBoard(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		BoardSize int `json:"size"`
		Events []event `json:"events"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	events := params.Events
	size := params.BoardSize
	rand.Shuffle(len(events), func (i, j int) {events[i], events[j] = events[j], events[i]})
	if(len(events) > size){
		events = events[:size]
	}

	respondWithJSON(w, 201, parameters{Events: events})
}