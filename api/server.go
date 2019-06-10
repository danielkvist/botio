package api

import (
	"net/http"
	"time"

	"github.com/danielkvist/botio/db"
)

func NewServer(bolter db.Bolter, listenAddr string, username string, password string) *http.Server {
	routes := []*Route{
		&Route{
			Name:        "GET Commands",
			Method:      http.MethodGet,
			Pattern:     "/api/commands",
			HandlerFunc: GetAll(bolter, "commands"),
		},
		&Route{
			Name:        "GET Command",
			Method:      http.MethodGet,
			Pattern:     "/api/commands/{command}",
			HandlerFunc: Get(bolter, "commands"),
		},
		&Route{
			Name:        "POST Command",
			Method:      http.MethodPost,
			Pattern:     "/api/commands",
			HandlerFunc: Post(bolter, "commands"),
		},
		&Route{
			Name:        "PUT Command",
			Method:      http.MethodPut,
			Pattern:     "/api/commands",
			HandlerFunc: Put(bolter, "commands"),
		},
		&Route{
			Name:        "DELETE Command",
			Method:      http.MethodDelete,
			Pattern:     "/api/commands/{command}",
			HandlerFunc: Delete(bolter, "commands"),
		},
	}

	r := NewRouter(routes)
	return &http.Server{
		Addr:         listenAddr,
		Handler:      basicAuth(username, password, r),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func basicAuth(username string, password string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if username != user || password != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
