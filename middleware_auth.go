package main

import (
	"fmt"
	"net/http"

	"github.com/BigPeanutFromStudio/bingo/internal/auth"
	"github.com/BigPeanutFromStudio/bingo/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)


func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		token, err := auth.GetToken(r.Header)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByToken(r.Context(), token)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error getting user: %v", err))
			return
		}

	handler(w, r, user)
	}
}