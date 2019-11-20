package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/danielkvist/botio/proto"
	"github.com/danielkvist/botio/server"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// Server returns a *cobra.Command.
func Server() *cobra.Command {
	return serverCmd(serverWithBoltDB(), serverWithPostgresDB())
}

func serverCmd(commands ...*cobra.Command) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Server provides subcommands to initialize a server with differents databases.",
	}

	for _, cmd := range commands {
		serverCmd.AddCommand(cmd)
	}

	return serverCmd
}

func serverWithBoltDB() *cobra.Command {
	// var key string
	var cacheCap int64
	var collection string
	var database string
	var httpPort string
	var port string
	var sslca string
	var sslcrt string
	var sslkey string

	s := &cobra.Command{
		Use:     "bolt",
		Short:   "Starts a Botio server with BoltDB.",
		Example: "botio server bolt --database ./data/botio.db --collection commands --http :8081 --port :9091",
		RunE: func(cmd *cobra.Command, args []string) error {
			quit := make(chan error, 2)
			defer close(quit)

			serverOptions := []server.Option{
				server.WithListener(port),
				server.WithBoltDB(database, collection),
				server.WithCache(cacheCap),
			}

			if sslcrt == "" || sslkey == "" || sslca == "" {
				serverOptions = append(serverOptions, server.WithInsecureGRPCServer())
			} else {
				serverOptions = append(serverOptions, server.WithSecuredGRPCServer(sslcrt, sslkey, sslca))
			}

			s, err := server.New(serverOptions...)
			if err != nil {
				return errors.Wrap(err, "while creating a new Botio server with BoltDB")
			}

			go func() {
				if err := s.Connect(); err != nil {
					quit <- errors.Wrapf(err, "while connectign server to BoltDB")
					return
				}

				if err := s.Serve(); err != nil {
					quit <- errors.Wrap(err, "while listening to requests")
				}
			}()

			go func() {
				if err := runHTTPEndpoint(httpPort); err != nil {
					quit <- err
				}
			}()

			log.Printf("server with BoltDB listening to HTTP requests on %q and to gRPC requests on %q!", httpPort, port)
			return <-quit
		},
	}

	// s.Flags().StringVar(&key, "key", "", "authentication key to generate a jwt token")
	s.Flags().Int64Var(&cacheCap, "cache", 262144000, "capacity of the in-memory cache in bytes")
	s.Flags().StringVar(&collection, "collection", "commands", "collection used to store commands")
	s.Flags().StringVar(&database, "database", "./botio.db", "database path")
	s.Flags().StringVar(&httpPort, "http", ":8081", "port for HTTP server")
	s.Flags().StringVar(&port, "port", ":9091", "port for gRPC server")
	s.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	s.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")

	return s
}

func serverWithPostgresDB() *cobra.Command {
	// var key string
	var cacheCap int64
	var database string
	var host string
	var httpPort string
	var password string
	var port string
	var pport string
	var sslca string
	var sslcrt string
	var sslkey string
	var table string
	var user string

	s := &cobra.Command{
		Use:     "postgres",
		Short:   "Starts a Botio server with PostgreSQL.",
		Example: "botio server postgres --user postgres --password toor --database botio --table commands --http :8081 --port :9091",
		RunE: func(cmd *cobra.Command, args []string) error {
			quit := make(chan error, 2)
			defer close(quit)

			serverOptions := []server.Option{
				server.WithListener(port),
				server.WithCache(cacheCap),
				// TODO: Clean pport
				server.WithPostgresDB(host, pport, database, table, user, password),
			}

			if sslcrt == "" || sslkey == "" || sslca == "" {
				serverOptions = append(serverOptions, server.WithInsecureGRPCServer())
			} else {
				serverOptions = append(serverOptions, server.WithSecuredGRPCServer(sslcrt, sslkey, sslca))
			}

			s, err := server.New(serverOptions...)
			if err != nil {
				return errors.Wrap(err, "while creating a new Botio server with PostgreSQL")
			}

			go func() {
				if err := s.Connect(); err != nil {
					quit <- errors.Wrapf(err, "while connectign server to PostgreSQL")
					return
				}

				if err := s.Serve(); err != nil {
					quit <- errors.Wrap(err, "while listening to requests")
				}
			}()

			go func() {
				if err := runHTTPEndpoint(httpPort); err != nil {
					quit <- err
				}
			}()

			log.Printf("server with PostgreSQL listening to HTTP requests on %q and to gRPC requests on %q!", httpPort, port)
			return <-quit
		},
	}

	// s.Flags().StringVar(&key, "key", "", "authentication key to generate a jwt token")
	s.Flags().Int64Var(&cacheCap, "cache", 262144000, "capacity of the in-memory cache in bytes")
	s.Flags().StringVar(&database, "database", "botio", "PostgreSQL database name")
	s.Flags().StringVar(&host, "host", "postgres", "host of the PostgreSQL database")
	s.Flags().StringVar(&httpPort, "http", ":8081", "port for HTTP server")
	s.Flags().StringVar(&password, "password", "", "password for the user of the PostgreSQL database")
	s.Flags().StringVar(&port, "port", ":9091", "port for gRPC server")
	s.Flags().StringVar(&pport, "postgresPort", "5432", "port of the PostgreSQL database host")
	s.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	s.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	s.Flags().StringVar(&table, "table", "commands", "table of the PostgreSQL database")
	s.Flags().StringVar(&user, "user", "", "user of the PostgreSQL database")

	return s
}

// FIXME:
func runHTTPEndpoint(port string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	options := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if err := proto.RegisterBotioHandlerFromEndpoint(ctx, mux, port, options); err != nil {
		return errors.Wrapf(err, "while registering gRPC HTTP endpoint")
	}

	return http.ListenAndServe(port, mux)
}
