// Package routes manages the router for the server.
package routes

import (
	"net/http"
	"time"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Routes receives a db.DB and a collection and initializes
// all the routes. Finally it returns an http.Handler.
func Routes(database db.DB, col string) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.URLFormat)

	// Routes
	r.Route("/api/commands", func(r chi.Router) {
		r.Get("/", handlers.GetAll(database, col))
		r.Get("/{command}", handlers.Get(database, col))
		r.Post("/", handlers.Post(database, col))
		r.Put("/{command}", handlers.Put(database, col))
		r.Delete("/{command}", handlers.Delete(database, col))
	})

	r.Route("/api/backup", func(r chi.Router) {
		r.Get("/", handlers.Backup(database, col))
	})

	return r
}
