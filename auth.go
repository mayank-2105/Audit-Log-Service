/* This file contains the authentication mechanisms for our API. 
I am using JWT(token based authentication) where user first needs to send a login request and obtain the token. 
The user can use this token for authentication in subsequent requests. 
The token is valid upto 24 hours and I am using jwt-go library for implementing this*/ 

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
)



/* For simplicity I am hard coding an easy secret key here, 
In real time systems we can store these credentials in an 
env file*/
const (
    secretKey = "your_secret_key" // Replace with a strong secret key
    tokenDuration = 24 * time.Hour
)

// GenerateToken generates a new JWT token
func GenerateToken() (string, error) {
    claims := jwt.StandardClaims{
        ExpiresAt: time.Now().Add(tokenDuration).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// ValidateToken validates a JWT token and returns the claims if valid
func ValidateToken(tokenStr string) (*jwt.StandardClaims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

/*TokenAuthMiddleware handles token-based authentication, we extract the
auth token from header and then validate it, the request is only allowed to
proceed only if the validation succeeds*/
func TokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenStr := r.Header.Get("Authorization") // This will extract  
        if tokenStr == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        claims, err := ValidateToken(tokenStr)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

		fmt.Println(claims.Subject)

        // If we want, you can use the claims to store user information
        // For example: userID := claims.Subject

        next(w, r)
    }
}

/* This func returns auth token which needs to be used in the header of all our upcoming requests
It calls the GenerateToken method and returns the token as a json*/
func Login(w http.ResponseWriter, r *http.Request) {
    token, err := GenerateToken()
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Return the token in the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}