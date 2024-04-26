package auth

import (
	"errors"
	"net/http"
	"strings"
)

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