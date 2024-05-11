package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key = "bda413dfbbdf78d06947bf6a9e60e83c40cb975f74ff47159c63f14376b97611"
	MaxAge = 86400 * 30 // 30 days
	isProd = false
)

func NewAuth(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:8000/v1/auth/google/callback"),
	)

}

func GetToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("no authentication token provided") 
	}
	
	vals := strings.Split(token, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid token format")
	}

	if vals[0] != "Token" {
		return "", errors.New("invalid token type")
	}

	return vals[1], nil



}