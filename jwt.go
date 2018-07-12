package main

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var JWTSecret = []byte("aoR6E4tb6TWDgP8dQdkpcg")

func (state *RuntimeState) createandSetToken(username string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	exp := time.Now().Add(time.Minute * jwtTokenExpirationMin)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserName": username,
		"exp":      int(exp.Unix()),
		"iat":      int(time.Now().Unix()),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		log.Println(err)
		return "", err
	}
	state.tokenmutex.Lock()
	state.jwtTokenmap[username] = tokenInfo{tokenString, exp}
	state.tokenmutex.Unlock()

	return tokenString, nil
}
