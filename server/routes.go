package server

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) routes() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.URLFormat)

	// Routes
	r.Route("/api/commands", func(r chi.Router) {
		r.Get("/", s.jwtMiddleware(s.handleGetAll()))
		r.Get("/{command}", s.jwtMiddleware(s.handleGet()))
		r.Post("/", s.jwtMiddleware(s.handlePost()))
		r.Put("/{command}", s.jwtMiddleware(s.handlePut()))
		r.Delete("/{command}", s.jwtMiddleware(s.handleDelete()))
	})

	s.router = r
}
