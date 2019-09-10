package server

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func (s *Server) jwtMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			s.errResp(w, r, "unauthorized", fmt.Errorf("user not authorized"), http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error while parsing JWT token")
			}

			return []byte(s.key), nil
		})

		if err != nil {
			s.errResp(w, r, err.Error(), err, http.StatusInternalServerError)
			return
		}

		if token.Valid {
			h.ServeHTTP(w, r)
		}
	})
}
