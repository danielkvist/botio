// Package server defines a gRPC server implementation
// with options for its creation and logging.
package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	Serve() error
}

type server struct {
	db       db.DB
	srv      *grpc.Server
	listener net.Listener
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

func WithListener(port string) Option {
	return func(s *server) error {
		listener, err := net.Listen("tcp", port)
		if err != nil {
			return fmt.Errorf("while creating a new tcp listener on port %q: %v", port, err)
		}

		s.listener = listener
		return nil
	}
}

func WithSecuredGRPCServer(crt, key, ca string) Option {
	return func(s *server) error {
		cert, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return fmt.Errorf("while loading SSL key pair: %v", err)
		}

		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(ca)
		if err != nil {
			return fmt.Errorf("while reading CA certificate: %v", err)
		}

		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return fmt.Errorf("fail while appending client certificates")
		}

		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		})

		s.srv = grpc.NewServer(grpc.Creds(creds))
		return nil
	}
}

func WithInsecureGRPCServer() Option {
	return func(s *server) error {
		s.srv = grpc.NewServer()
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

	if s.srv == nil {
		insecureOpt := WithInsecureGRPCServer()
		if err := insecureOpt(s); err != nil {
			return nil, fmt.Errorf("while creating a new Server: %v", err)
		}
	}

	proto.RegisterBotioServer(s.srv, s)
	return s, nil
}

func (s *server) Serve() error {
	return s.srv.Serve(s.listener)
}

// Connect tries to connect the Server to its database.
func (s *server) Connect() error {
	if err := s.db.Connect(); err != nil {
		return fmt.Errorf("while connecting Server to database: %v", err)
	}

	return nil
}

func (s *server) CloseList() {
	s.listener.Close()
}
