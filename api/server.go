package api

import (
	"net/http"
	"time"

	"github.com/danielkvist/botio/db"
)

func NewServer(bolter db.Bolter, listenAddr string) *http.Server {
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
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
