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


// https://www.googleapis.com/oauth2/v3/userinfo?access_token={access_token} 
// use this to get id and get user in database

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		id, err := auth.GetID(r.Header)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByID(r.Context(), id)

		
		if err != nil {
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
	_, err = apiCfg.DB.GetUserByID(r.Context(), user.UserID)

	//defer http.Redirect(w, r, "http://localhost:5173/", 200)

	if err != nil{
		apiCfg.handlerCreateGoogleUser(w, r, user)
		return
	}

	// godotenv.Load()

	// key := os.Getenv("SECRET")

	// if key == ""{
	// 	log.Fatal("SECRET variable not found in environment")
	// }

	// cipherID, err := auth.Encrypt([]byte(key), []byte(user.UserID))

	// if err != nil{
	// 	respondWithError(w, 400, fmt.Sprintf("Error encrypting ID: %v", err))
	// 	return
	// }

	// fmt.Printf("ID: %v", string(cipherID[:]))


	token := "Token " + user.UserID

	w.Header().Add("Authorization", token)

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