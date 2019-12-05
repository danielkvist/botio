// Package server defines a gRPC server.
package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/danielkvist/botio/cache"
	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/proto"
	"github.com/dgrijalva/jwt-go"

	"github.com/golang/protobuf/ptypes/empty"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	CloseList()
}

type server struct {
	db            db.DB
	dbPlatform    string
	cache         cache.Cache
	cachePlatform string
	srv           *grpc.Server
	ssl           bool
	listener      net.Listener
	httpPort      string
	key           string
	jwt           string
	log           *logrus.Logger
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
		s.dbPlatform = "BoltDB"

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
		s.dbPlatform = "PostgreSQL"

		return nil
	}
}

// WithTestDB returns an Option to a new server that assigns to its
// db field an in-memory database for testing.
func WithTestDB() Option {
	return func(s *server) error {
		database := db.Create("testing")
		s.db = database
		s.dbPlatform = "In Memory"

		return nil
	}
}

// WithRistrettoCache returns an Option to a new Server that assigns to its cache
// a Ristretto's based cache.
func WithRistrettoCache(cap int) Option {
	return func(s *server) error {
		c := cache.Create("ristretto")
		err := c.Init(cap)
		if err != nil {
			return err
		}

		s.cache = c
		s.cachePlatform = "Ristretto"

		return nil
	}
}

// WithHTTPPort returns an Option to a new Server that assigns to its httpPort
// field the received port.
func WithHTTPPort(port string) Option {
	return func(s *server) error {
		s.httpPort = port
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

		s.srv = grpc.NewServer(
			grpc.Creds(creds),
			grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					grpc_auth.UnaryServerInterceptor(s.jwtAuth),
					grpc_recovery.UnaryServerInterceptor(),
				),
			),
		)
		s.ssl = true
		return nil
	}
}

// WithInsecureGRPCServer returns an Option to a new Server that assigns to its
// srv field a insecure gRPC server.
func WithInsecureGRPCServer() Option {
	return func(s *server) error {
		s.srv = grpc.NewServer(
			grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					grpc_auth.UnaryServerInterceptor(s.jwtAuth),
					grpc_recovery.UnaryServerInterceptor(),
				),
			),
		)
		return nil
	}
}

// WithTextLogger returns an Option to a new Server with a text
// based logger.
func WithTextLogger(out io.Writer) Option {
	return func(s *server) error {
		s.log = logrus.New()
		s.log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  time.RFC850,
			DisableSorting:   true,
			QuoteEmptyFields: true,
		})
		s.log.Out = out

		return nil
	}
}

// WithJSONLogger returns an Option to a new Server with a JSON
// based logger.
func WithJSONLogger(out io.Writer) Option {
	return func(s *server) error {
		s.log = logrus.New()
		s.log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC850,
			PrettyPrint:     true,
		})
		s.log.Out = out

		return nil
	}
}

// WithJWTAuthToken returns an Option to a new Server that creates from
// the received key a JWT that is assigned to the Server.
func WithJWTAuthToken(key string) Option {
	return func(s *server) error {
		token := jwt.New(jwt.SigningMethodHS256)
		tokenStr, err := token.SignedString([]byte(key))
		if err != nil {
			return errors.Wrap(err, "while signing JWT token for authentication")
		}

		s.key = key
		s.jwt = tokenStr
		return nil
	}
}

// New should receive one or more Options to apply then to a new Server
// that will be return completely initialized.
func New(options ...Option) (Server, error) {
	start := time.Now()

	errMsg := "while creating a new Server"
	if len(options) == 0 {
		return nil, errors.Errorf("%s: no options provided", errMsg)
	}

	s := &server{}
	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, errors.Wrapf(err, "%s", errMsg)
		}
	}

	switch {
	case s.cache == nil:
		return nil, errors.Errorf("%s: no Cache provided", errMsg)
	case s.db == nil:
		return nil, errors.Errorf("%s: no DB provided", errMsg)
	case s.jwt == "":
		return nil, errors.Errorf("%s: no key to generate a valid JWT provided", errMsg)
	case s.listener == nil:
		return nil, errors.Errorf("%s: no net.Listener provided", errMsg)
	case s.log == nil:
		return nil, errors.Errorf("%s: no logger provided", errMsg)
	case s.srv == nil:
		return nil, errors.Errorf("%s: no gRPC server provided", errMsg)
	}

	proto.RegisterBotioServer(s.srv, s)

	s.logInfo(
		"server",
		"New",
		fmt.Sprintf("created a new Server with JWT auth token: %q", s.jwt),
		time.Since(start),
	)

	return s, nil
}

// Serve accepts incoming gRPC connections using the Server's listeners and also
// listens to HTTP requests using a JSON gateway. It returns an error if one of the
// two process fail.
func (s *server) Serve() error {
	errCh := make(chan error, 2)
	defer close(errCh)

	go func() {
		if err := s.srv.Serve(s.listener); err != nil {
			errCh <- errors.Wrapf(err, "while listening to gRPC requests on %q", s.listener.Addr().String())
		}
	}()

	if s.httpPort != "" {
		go s.serveJSONGateway(errCh)
	}

	return <-errCh
}

func (s *server) serveJSONGateway(errCh chan<- error) {
	errCh <- errors.Wrapf(s.jsonGateway(), "while listening to HTTP requests on %q", s.httpPort)
}

// Connect tries to connect the Server to its database.
func (s *server) Connect() error {
	if err := s.db.Connect(); err != nil {
		s.logFatal("db", "Connect", err.Error(), "while connecting Server to database")
		return err
	}

	return nil
}

// CloseList closes the Server's listener.
func (s *server) CloseList() {
	s.logInfo("server", "CloseList", "closing Server's listener", 0*time.Second)
	s.listener.Close()
}
