package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/danielkvist/botio/proto"
	"github.com/danielkvist/botio/server"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TODO: Check flags description in both commands

// Server returns a *cobra.Command.
func Server() *cobra.Command {
	return serverCmd(serverWithBoltDB(), serverWithPostgresDB())
}

func serverCmd(commands ...*cobra.Command) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "server contains some subcommands to initialize a server with different databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	for _, cmd := range commands {
		serverCmd.AddCommand(cmd)
	}

	return serverCmd
}

func serverWithBoltDB() *cobra.Command {
	// var key string
	var collection string
	var database string
	var httpPort string
	var port string
	var sslca string
	var sslcrt string
	var sslkey string

	s := &cobra.Command{
		Use:     "bolt",
		Short:   "Starts a server with a BoltDB database to manage your commands with HTTP methods",
		Example: "botio server bolt --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey",
		RunE: func(cmd *cobra.Command, args []string) error {
			quit := make(chan error, 2)
			defer close(quit)

			serverOptions := []server.Option{
				server.WithBoltDB(database, collection),
			}

			s, err := server.New(serverOptions...)
			if err != nil {
				return fmt.Errorf("while creating a new Server: %v", err)
			}

			listener, err := net.Listen("tcp", port)
			if err != nil {
				return fmt.Errorf("while creating a new listener: %v", err)
			}
			defer listener.Close()

			cert, err := tls.LoadX509KeyPair(sslcrt, sslkey)
			if err != nil {
				return fmt.Errorf("while loading SSL key pair: %v", err)
			}

			certPool := x509.NewCertPool()
			ca, err := ioutil.ReadFile(sslca)
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

			srv := grpc.NewServer(grpc.Creds(creds))
			proto.RegisterBotioServer(srv, s)

			go func() {
				if err := srv.Serve(listener); err != nil {
					quit <- fmt.Errorf("while listeting to requests on %v: %v", listener.Addr().String(), err)
				}
			}()

			go func() {
				if err := runHTTPEndpoint(httpPort); err != nil {
					quit <- err
				}
			}()

			return <-quit
		},
	}

	// s.Flags().StringVar(&key, "key", "", "authentication key to generate a jwt token")
	s.Flags().StringVar(&collection, "collection", "commands", "collection used to store commands")
	s.Flags().StringVar(&database, "database", "./botio.db", "database path")
	s.Flags().StringVar(&httpPort, "http", ":8081", "port for HTTP server")
	s.Flags().StringVar(&port, "port", ":9091", "port for gRPC server")
	s.Flags().StringVar(&sslca, "sslca", "./ca.crt", "ssl client certification file")
	s.Flags().StringVar(&sslcrt, "sslcrt", "./server.crt", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "./server.key", "ssl certification key file")

	return s
}

func serverWithPostgresDB() *cobra.Command {
	// var key string
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
		Short:   "Starts a server with that connects to a PostgreSQL database to manage your commands with HTTP methods",
		Example: "botio server postgres --user postgres --password toor --database botio --table commands --key mysupersecretkey",
		RunE: func(cmd *cobra.Command, args []string) error {
			quit := make(chan error, 2)
			defer close(quit)

			serverOptions := []server.Option{
				server.WithPostgresDB(host, pport, database, table, user, password),
			}

			s, err := server.New(serverOptions...)
			if err != nil {
				log.Fatalf("while creating a new Server: %v", err)
			}

			listener, err := net.Listen("tcp", port)
			if err != nil {
				return fmt.Errorf("while creating a new listener: %v", err)
			}
			defer listener.Close()

			cert, err := tls.LoadX509KeyPair(sslcrt, sslkey)
			if err != nil {
				return fmt.Errorf("while loading SSL key pair: %v", err)
			}

			certPool := x509.NewCertPool()
			ca, err := ioutil.ReadFile(sslca)
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

			srv := grpc.NewServer(grpc.Creds(creds))
			proto.RegisterBotioServer(srv, s)

			go func() {
				if err := srv.Serve(listener); err != nil {
					quit <- fmt.Errorf("while listeting to requests on %v: %v", listener.Addr().String(), err)
				}
			}()

			go func() {
				if err := runHTTPEndpoint(httpPort); err != nil {
					quit <- err
				}
			}()

			return <-quit
		},
	}

	// s.Flags().StringVar(&key, "key", "", "authentication key to generate a jwt token")
	s.Flags().StringVar(&database, "database", "botio", "PostgreSQL database name")
	s.Flags().StringVar(&host, "host", "postgres", "host of the PostgreSQL database")
	s.Flags().StringVar(&httpPort, "http", ":8081", "port for HTTP server")
	s.Flags().StringVar(&password, "password", "", "password for the user of the PostgreSQL database")
	s.Flags().StringVar(&port, "port", ":9091", "port for gRPC server")
	s.Flags().StringVar(&pport, "postgresPort", "5432", "port of the PostgreSQL database host")
	s.Flags().StringVar(&sslca, "sslca", "./ca.crt", "ssl client certification file")
	s.Flags().StringVar(&sslcrt, "sslcrt", "./server.crt", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "./server.key", "ssl certification key file")
	s.Flags().StringVar(&table, "table", "commands", "table of the PostgreSQL database")
	s.Flags().StringVar(&user, "user", "", "user of the PostgreSQL database")

	return s
}

func runHTTPEndpoint(port string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	options := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if err := proto.RegisterBotioHandlerFromEndpoint(ctx, mux, port, options); err != nil {
		return fmt.Errorf("while registering gRPC HTTP endpoint: %v", err)
	}

	return http.ListenAndServe(port, mux)
}
