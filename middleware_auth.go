package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BigPeanutFromStudio/bingo/internal/auth"
	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)


// https://www.googleapis.com/oauth2/v3/userinfo?access_token={access_token} 
// use this to get id and get user in database

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		token, err := auth.GetToken(r.Header)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}

		url := "https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token 

		resp, err := http.Get(url)

		if err != nil{
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil{
			respondWithError(w, 400, fmt.Sprintf("Error reading response body: %v", err))
			return
		}

		var userData map[string]interface{}
		err = json.Unmarshal(body, &userData)

		if err != nil{
			respondWithError(w, 400, fmt.Sprintf("Error unmarshaling response body: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByID(r.Context(), userData["sub"].(string))


		if err != nil {
			//apiCfg.handlerCreateGoogleUser(w, r, userData)
			respondWithError(w, 400, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		handler(w, r, user)
	}
}

//idk if this should go here
func (apiCfg *apiConfig)getAuthCallbackFunction(w http.ResponseWriter, r *http.Request){

	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, r)
		return
	}

	// fmt.Printf("AccessToken: %s\nAccessTokenSecret: %s\nRefreshToken: %s\nExpiresAt: %s\nRawData: %s\n",
	// 	user.AccessToken, user.AccessTokenSecret, user.RefreshToken, user.ExpiresAt, user.RawData)
	apiCfg.handlerCreateGoogleUser(w, r, user)

	//@BigPeanutFromStudio REMEMBER TO PASS THE ACCESS TOKEN TO THE FRONT-END

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