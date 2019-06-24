// Package server provides utilities to create a new HTTP server
// with basic auth.
package server

import (
	"net/http"
	"time"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/handlers"
)

// New returns a new *http.Server with basic authentication and a
// *mux.Router with all the routes set.
func New(bolter db.Bolter, col string, username string, password string, listenAddr string) *http.Server {
	routes := []*Route{
		{
			Name:        "GET Commands",
			Method:      http.MethodGet,
			Pattern:     "/api/commands",
			HandlerFunc: handlers.GetAll(bolter, col),
		},
		{
			Name:        "GET Command",
			Method:      http.MethodGet,
			Pattern:     "/api/commands/{command}",
			HandlerFunc: handlers.Get(bolter, col),
		},
		{
			Name:        "POST Command",
			Method:      http.MethodPost,
			Pattern:     "/api/commands",
			HandlerFunc: handlers.Post(bolter, col),
		},
		{
			Name:        "PUT Command",
			Method:      http.MethodPut,
			Pattern:     "/api/commands",
			HandlerFunc: handlers.Put(bolter, col),
		},
		{
			Name:        "DELETE Command",
			Method:      http.MethodDelete,
			Pattern:     "/api/commands/{command}",
			HandlerFunc: handlers.Delete(bolter, col),
		},
		{
			Name:        "Backup DB",
			Method:      http.MethodGet,
			Pattern:     "/api/backup",
			HandlerFunc: handlers.Backup(bolter, col),
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
