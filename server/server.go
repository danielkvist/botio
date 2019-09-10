// Package server exports a struct called Server that satisfies the
// http.Handler interface.
package server

import (
	"log"
	"net/http"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/logger"

	"github.com/go-chi/chi"
)

// Server is an abstraction to manage differents aspects of a mux router.
type Server struct {
	key    string
	db     db.DB
	router *chi.Mux
	logger *logger.Logger
}

// Option represents an option to a *Server.
type Option func(s *Server)

// WithDB receives an environment, a path and a collection
// and returns an Option that creates, connects and
// assigns a db.DB to the Server's db.
func WithDB(env, path, col string) Option {
	return func(s *Server) {
		s.db = db.Create(env)
		s.db.Open(path, col)
	}
}

// WithBoltDB receives a path and a collection and returns
// an Option that creates, connects and assigns
// a BoltDB db.DB to the Server's db.
func WithBoltDB(path, col string) Option {
	return func(s *Server) {
		s.db = db.Create("local")
		s.db.Open(path, col)
	}
}

// WithJWTMiddleware receives a key and returns an Option
// that assigns it to the Server's key.
func WithJWTMiddleware(key string) Option {
	return func(s *Server) {
		s.key = key
	}
}

// New should receive one or more Options to apply then to a new *Server
// that will be return completely initialized.
func New(options ...Option) *Server {
	if len(options) == 0 {
		log.Fatalf("no options received for creating a new Server")
	}
	s := &Server{}
	s.routes()
	s.logger = logger.New()

	for _, opt := range options {
		opt(s)
	}

	return s
}

// ServeHTTP makes the Server type to satisfy the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
