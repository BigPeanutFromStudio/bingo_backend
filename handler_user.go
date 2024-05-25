package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BigPeanutFromStudio/bingo/internal/database"
	"github.com/markbates/goth"
)

//DO NOT RETURN SENSITIVE DATA LMAO

func randomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }
  return hex.EncodeToString(bytes), nil
}

func (apiCfg *apiConfig)handlerSetGoogleUserNickname(w http.ResponseWriter, r *http.Request, user database.User){

	type parameters struct{
		Nickname string `json:"nickname"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	nickID, err := randomHex(5)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error generating id: %v", err))
		return
	}

	updatedUser, err := apiCfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: user.ID,
		Nickname: params.Nickname + "#" + nickID,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error updating user: %v", err))
		return
	}

	respondWithJSON(w, 200, updatedUser)
}

// How does it work when redirecting
func (apiCfg *apiConfig)handlerCreateGoogleUser(w http.ResponseWriter, r *http.Request, userData goth.User){

	// GET USER MAYBE?

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: userData.UserID,
		Nickname: "Temporarily not working LMAO",
		Email: userData.Email,
		PictureUrl: userData.AvatarURL,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	// CREATE USER
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	// type userToReturn struct{
	// 	Nickname string `json:"nickname"`
	// 	Email string `json:"email"`
	// 	PictureUrl string `json:"picture_url"`
	// }

	// godotenv.Load()

	// key := os.Getenv("SECRET")

	// if key == ""{
	// 	log.Fatal("SECRET variable not found in environment")
	// }

	// cipherID, err := auth.Encrypt([]byte(key), []byte(user.ID))

	// if err != nil{
	// 	respondWithError(w, 400, fmt.Sprintf("Error encrypting ID: %v", err))
	// 	return
	// }

	token := "Token " + user.ID

	w.Header().Add("Authorization", token)

	fmt.Printf("%v", token)

	user.ID = "Lol you tried"

	//userInfo, _ := json.Marshal(userToReturn{Nickname: user.Nickname, Email: user.Email, PictureUrl: user.PictureUrl})

	respondWithJSON(w, 201, user)
} 

// func (apiCfg *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request) {

// 	type parameters struct{
// 		Nickname string `json:"nickname"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil{
// 		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
// 		return
// 	}

// 	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
// 		ID: "TEMPORARY VALUE",
// 		Nickname: params.Nickname,
// 		CreatedAt: time.Now().UTC(),
// 		UpdatedAt: time.Now().UTC(),
// 	})

// 	if err != nil{
// 		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
// 		return
// 	}

// 	respondWithJSON(w, 201, user)
// }

func (apiCfg *apiConfig)handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	

	respondWithJSON(w, 200, user)
}