package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/BigPeanutFromStudio/bingo/internal/auth"
	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
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

//idk if this should go here
func getAuthCallbackFunction(w http.ResponseWriter, r *http.Request){

	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, r)
		return
	}

	fmt.Println(user.UserID)

	http.Redirect(w, r, "http://localhost:5173/", 301)

}

func logoutHandler(w http.ResponseWriter, r *http.Request){
	gothic.Logout(w, r)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func beginAuthHandler(w http.ResponseWriter, r *http.Request){
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(w, r)	
}