// Package server defines a gRPC server implementation
// with options for its creation and logging.
package server

import (
	"context"
	"fmt"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
)

// Server represents a gRPC BotioServer with a method to connect
// to its database.
type Server interface {
	AddCommand(context.Context, *proto.BotCommand) (*empty.Empty, error)
	GetCommand(context.Context, *proto.Command) (*proto.BotCommand, error)
	ListCommands(context.Context, *empty.Empty) (*proto.BotCommands, error)
	UpdateCommand(context.Context, *proto.BotCommand) (*empty.Empty, error)
	DeleteCommand(context.Context, *proto.Command) (*empty.Empty, error)
	Connect() error
}

type server struct {
	db db.DB
}

// Option represents an option for a new *server.
type Option func(s *server) error

// WithBoltDB receives a path and a colletion to create a BoltDB client
// and assign it to the new server. If something goes wrong while
// configuring the client it panics.
func WithBoltDB(path, col string) Option {
	return func(s *server) error {
		database := db.Create("local")
		bdb, ok := database.(*db.Bolt)
		if !ok {
			return fmt.Errorf("while connecting BoltDB database a fatal error happened")
		}

		bdb.Path = path
		bdb.Col = col

		s.db = bdb

		return nil
	}
}

// WithPostgresDB receives a set of parameters to create a PostgreSQL client
// and assign it to the new server. If something goes wrong while
// configuring the client it panics.
func WithPostgresDB(host, port, dbName, table, user, password string) Option {
	return func(s *server) error {
		database := db.Create("postgres")
		ps, ok := database.(*db.Postgres)
		if !ok {
			return fmt.Errorf("while creating a PostgreSQL database a fatal error happened")
		}

		ps.Host = host
		ps.Port = port
		ps.User = user
		ps.Password = password
		ps.DB = dbName
		ps.Table = table

		s.db = ps

		return nil
	}
}

// WithTestDB returns an Option to a new server that assigns to its
// db field an in-memory database for testing.
func WithTestDB() Option {
	return func(s *server) error {
		database := db.Create("testing")
		s.db = database

		return nil
	}
}

// New should receive one or more Options to apply then to a new Server
// that will be return completely initialized.
func New(options ...Option) (Server, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("no options received for creating a new Server")
	}

	s := &server{}
	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("while creating a new Server: %v", err)
		}
	}

	return s, nil
}

// Connect tries to connect the Server to its database.
func (s *server) Connect() error {
	if err := s.db.Connect(); err != nil {
		return fmt.Errorf("while connecting Server to database: %v", err)
	}

	return nil
}
