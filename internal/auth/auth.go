package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
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
	MaxAge = 86400 * 30 // 30 days
	isProd = false
)

func NewAuth(){

	godotenv.Load()

	key := os.Getenv("SECRET")

	if key == ""{
		log.Fatal("SECRET variable not found in environment")
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


func GetID(headers http.Header) (string, error) {
	id := headers.Get("Authorization")
	if id == "" {
		return "", errors.New("no authentication token provided") 
	}
	
	vals := strings.Split(id, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid token format")
	}

	if vals[0] != "Token" {
		return "", errors.New("invalid token type")
	}

	// godotenv.Load()

	// key := os.Getenv("SECRET")

	// if key == ""{
	// 	log.Fatal("SECRET variable not found in environment")
	// }

	// userID, err := Decrypt([]byte(key), []byte(vals[1]))


	// if err != nil{
	// 	return "", errors.New("decryption error")
	// }

	// userIDstring := string(userID)

	return vals[1], nil



}

func Encrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    b := base64.StdEncoding.EncodeToString(text)
    ciphertext := make([]byte, aes.BlockSize+len(b))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
    return ciphertext, nil
}

func Decrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    if len(text) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }
    iv := text[:aes.BlockSize]
    text = text[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(text, text)
    data, err := base64.StdEncoding.DecodeString(string(text))
    if err != nil {
        return nil, err
    }
    return data, nil
}