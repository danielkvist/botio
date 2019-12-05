package cmd

import (
	"os"

	"github.com/danielkvist/botio/server"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	var key string
	var cacheCap int
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
			serverOptions := []server.Option{
				server.WithTextLogger(os.Stdout),
				server.WithHTTPPort(httpPort),
				server.WithListener(port),
				server.WithBoltDB(database, collection),
				server.WithRistrettoCache(cacheCap),
			}

			if sslcrt == "" || sslkey == "" || sslca == "" {
				serverOptions = append(serverOptions, server.WithInsecureGRPCServer())
			} else {
				serverOptions = append(serverOptions, server.WithSecuredGRPCServer(sslcrt, sslkey, sslca))
			}

			if key != "" {
				serverOptions = append(serverOptions, server.WithJWTAuthToken(key))
			}

			s, err := server.New(serverOptions...)
			if err != nil {
				return errors.Wrap(err, "while creating a new Botio server with BoltDB")
			}

			if err := s.Connect(); err != nil {
				return errors.Wrapf(err, "while connectign server to BoltDB")
			}

			if err := s.Serve(); err != nil {
				return errors.Wrap(err, "while listening to requests")
			}

			return nil
		},
	}

	s.Flags().StringVar(&key, "key", "", "key to generate a JWT token for authentication")
	s.Flags().IntVar(&cacheCap, "cache", 262144000, "capacity of the in-memory cache in bytes")
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
	var cacheCap int
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
			serverOptions := []server.Option{
				server.WithTextLogger(os.Stdout),
				server.WithHTTPPort(httpPort),
				server.WithListener(port),
				server.WithRistrettoCache(cacheCap),
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

			if err := s.Connect(); err != nil {
				return errors.Wrapf(err, "while connectign server to PostgreSQL")
			}

			if err := s.Serve(); err != nil {
				return errors.Wrap(err, "while listening to requests")
			}

			return nil
		},
	}

	// s.Flags().StringVar(&key, "key", "", "authentication key to generate a jwt token")
	s.Flags().IntVar(&cacheCap, "cache", 262144000, "capacity of the in-memory cache in bytes")
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
