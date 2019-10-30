package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/danielkvist/botio/proto"
	"github.com/danielkvist/botio/server"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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
	var collection string
	var database string
	var port string
	// var key string
	var sslcrt string
	var sslkey string

	s := &cobra.Command{
		Use:     "bolt",
		Short:   "Starts a server with a BoltDB database to manage your commands with HTTP methods",
		Example: "botio server bolt --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			creds, err := credentials.NewServerTLSFromFile(sslcrt, sslkey)
			if err != nil {
				return fmt.Errorf("while creating TLS credentials: %v", err)
			}

			srv := grpc.NewServer(grpc.Creds(creds))
			proto.RegisterBotioServer(srv, s)

			if err := srv.Serve(listener); err != nil {
				return fmt.Errorf("while listeting to requests on %v: %v", listener.Addr().String(), err)
			}

			return nil
		},
	}

	s.Flags().StringVar(&collection, "collection", "commands", "collection used to store commands")
	s.Flags().StringVar(&database, "database", "./botio.db", "database path")
	s.Flags().StringVar(&port, "port", ":443", "port for HTTPS connections")
	// s.Flags().StringVar(&key, "key", "", "authentication key to generate a jwt token")
	s.Flags().StringVar(&sslcrt, "sslcrt", "./server.crt", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "./server.key", "ssl certification key file")

	return s
}

func serverWithPostgresDB() *cobra.Command {
	var host string
	var pport string
	var user string
	var password string
	var table string
	var database string
	var port string
	// var key string
	var sslcrt string
	var sslkey string

	s := &cobra.Command{
		Use:     "postgres",
		Short:   "Starts a server with that connects to a PostgreSQL database to manage your commands with HTTP methods",
		Example: "botio server postgres --user postgres --password toor --database botio --table commands --key mysupersecretkey",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			creds, err := credentials.NewServerTLSFromFile(sslcrt, sslkey)
			if err != nil {
				return fmt.Errorf("while creating TLS credentials: %v", err)
			}

			srv := grpc.NewServer(grpc.Creds(creds))
			proto.RegisterBotioServer(srv, s)

			if err := srv.Serve(listener); err != nil {
				return fmt.Errorf("while listeting to requests on %v: %v", listener.Addr().String(), err)
			}

			return nil

		},
	}

	s.Flags().StringVar(&host, "host", "postgres", "host of the PostgreSQL database")
	s.Flags().StringVar(&pport, "postgresPort", "5432", "port of the PostgreSQL database host")
	s.Flags().StringVar(&database, "database", "botio", "PostgreSQL database name")
	s.Flags().StringVar(&table, "table", "commands", "table of the PostgreSQL database")
	s.Flags().StringVar(&user, "user", "", "user of the PostgreSQL database")
	s.Flags().StringVar(&password, "password", "", "password for the user of the PostgreSQL database")
	s.Flags().StringVar(&port, "port", ":443", "port for HTTPS connections")
	// s.Flags().StringVar(&key, "key", "", "authentication key to generate a jwt token")
	s.Flags().StringVar(&sslcrt, "sslcrt", "./server.crt", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "./server.key", "ssl certification key file")

	return s
}
