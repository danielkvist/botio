// Package jwt exports basic functions to work with JWT easily.
package jwt

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// Generate receives a string and returns a signed JWT token
// or an error if something went wrong.
func Generate(key string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tkStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("while generating the authentication JWT token: %v", err)
	}

	return tkStr, nil
}

// Middleware adds JWT authentication to the received handler checking
// if the received key is a valid JWT token.
func Middleware(key string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("while parsing a JWT token there was an error")
			}

			return []byte(key), nil
		})

		if err != nil {
			http.Error(w, "error while parsing JWT token", http.StatusInternalServerError)
		}

		if token.Valid {
			h.ServeHTTP(w, r)
		}
	})
}
