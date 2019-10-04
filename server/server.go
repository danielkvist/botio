// Package server exports a struct called Server that satisfies the
// http.Handler interface.
package server

import (
	"log"
	"net/http"
	"os"

	"github.com/danielkvist/botio/db"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
)

// Server is an abstraction to manage differents aspects of a mux router.
type Server struct {
	key    string
	db     db.DB
	router *chi.Mux
	logger *logrus.Logger
}

// Option represents an option to a *Server.
type Option func(s *Server)

// WithBoltDB receives a path and a collection and returns
// an Option that creates, connects and assigns
// a BoltDB db.DB to the Server's db.
func WithBoltDB(path, col string) Option {
	return func(s *Server) {
		database := db.Create("local")
		bdb, ok := database.(*db.Bolt)
		if !ok {
			log.Fatalf("while creating BoltDB database a fatal error happened")
		}

		bdb.Path = path
		bdb.Col = col

		s.db = bdb
	}
}

// WithPostgresDB receives a set of parameters neccessaries to initialize
// a PostgreSQL database and return an Option to assign it to the
// Server's db.
func WithPostgresDB(host, port, dbName, table, user, password string) Option {
	return func(s *Server) {
		database := db.Create("postgres")
		ps, ok := database.(*db.Postgres)
		if !ok {
			log.Fatalf("while creating a PostgreSQL database a fatal error happened")
		}

		ps.Host = host
		ps.Port = port
		ps.User = user
		ps.Password = password
		ps.DB = dbName
		ps.Table = table

		s.db = ps
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
	s.logger = logrus.New()
	s.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  "02-01-2006 15:04:05",
	})
	s.logger.Out = os.Stdout

	for _, opt := range options {
		opt(s)
	}

	if err := s.db.Connect(); err != nil {
		log.Fatalf("while connecting Server to database: %v", err)
	}

	return s
}

// ServeHTTP makes the Server type to satisfy the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
