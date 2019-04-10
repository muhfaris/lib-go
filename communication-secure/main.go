package main

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

const secretKey = "secret123 me"

type myClaim struct {
	jwt.StandardClaims
	Identity identityData `json:"identity"`
}

type identityData struct {
	AccessToken string `json:"access_token"`
	ID          string `json:"id"`
}

// dont forget for this law
// structer jwt must same
// if not same, token can decode but not valid
func main() {
	claim := myClaim{
		/*
			StandardClaims: jwt.StandardClaims{
				Issuer:    "cron-microservice",
				ExpiresAt: time.Now().Add(90).Unix(),
			},
		*/
		Identity: identityData{
			AccessToken: "EAAPP8YZCCyhsBAFbA48gaDgvd5NUMN9dwAtveZCsetTxMm8ZAZCnadRUG064dbE0KvQ3zyjv2mvU6ke6Ai6isM8uSZCSyvqZAJL7T8ZCNAfCPw8GFhZCpyrZCfXJhwWYgac2WxbZC5CeRPkgcGbuCwrvbDCxyYfZBIcoXBhCEr2PziGLqJXwoiCz7HP",
			ID:          "5a73efdc435c8a0adccc9ade",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenSigned, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(tokenSigned)
}
