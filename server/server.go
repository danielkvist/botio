// Package server defines a gRPC server implementation
// with options for its creation and logging.
package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"

	"github.com/danielkvist/botio/cache"
	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
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
	cache    cache.Cache
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
			return errors.Errorf("while connecting BoltDB database on path %q a fatal error happened", path)
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
			return errors.Errorf("while connecting to PostgreSQL server on %s:%s a fatal error happened", host, port)
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

// WithCache returns an Option to a new Server that assigns to its cache
// field an in-memory concurrently-safe cache.
func WithCache(counters, cost, bufferItems int64) Option {
	return func(s *server) error {
		c, err := cache.New(counters, cost, bufferItems)
		if err != nil {
			return err
		}

		s.cache = c

		return nil
	}
}

// WithListener returns an Option to a new Server that assigns to its
// listener field a TCP listener with the received address.
func WithListener(addr string) Option {
	return func(s *server) error {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return errors.Wrapf(err, "while creating a new tcp listener for addr %q", addr)
		}

		s.listener = listener
		return nil
	}
}

// WithSecuredGRPCServer returns an Option to a new Server that assigns to its
// srv field a secured gRPC server with TLS.
func WithSecuredGRPCServer(crt, key, ca string) Option {
	return func(s *server) error {
		cert, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return errors.Wrapf(err, "while loading SSL key pair")
		}

		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(ca)
		if err != nil {
			return errors.Wrapf(err, "while reading CA certificate")
		}

		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return errors.New("error while appending client certificates to cert pool")
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

// WithInsecureGRPCServer returns an Option to a new Server that assigns to its
// srv field a insecure gRPC server.
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
		return nil, errors.New("no options received for creating a new Server")
	}

	s := &server{}
	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, errors.Wrapf(err, "while creating a new Server")
		}
	}

	if s.srv == nil {
		insecureOpt := WithInsecureGRPCServer()
		if err := insecureOpt(s); err != nil {
			return nil, errors.Wrapf(err, "while creating a new Server")
		}
	}

	if s.cache == nil {
		// 262,144,000 its the number of bytes for the cache capacity.
		// 262144000 Bytes => 250 Megabytes
		cacheOpt := WithCache(1e7, 262144000, 64)
		if err := cacheOpt(s); err != nil {
			return nil, errors.Wrapf(err, "while creating a new Server")
		}
	}

	proto.RegisterBotioServer(s.srv, s)
	return s, nil
}

// Serve accepts incoming connections on the Server's listener.
func (s *server) Serve() error {
	return s.srv.Serve(s.listener)
}

// Connect tries to connect the Server to its database.
func (s *server) Connect() error {
	if err := s.db.Connect(); err != nil {
		return errors.Wrapf(err, "while connecting Server to database")
	}

	return nil
}

// CloseList closes the Server's listener.
func (s *server) CloseList() {
	s.listener.Close()
}
