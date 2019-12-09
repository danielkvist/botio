package cmd

import (
	"os"
	"time"

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
	var cacheCap int
	var collection string
	var database string
	var httpPort string
	var jsonOutput bool
	var key string
	var port string
	var sslca string
	var sslcrt string
	var sslkey string

	s := &cobra.Command{
		Use:     "bolt",
		Short:   "Starts a Botio server with BoltDB.",
		Example: "botio server bolt --database ./data/botio.db --collection commands --key mysupersecretkey",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverOptions := []server.Option{
				server.WithBoltDB(database, collection),
				server.WithHTTPPort(httpPort),
				server.WithListener(port),
				server.WithRistrettoCache(cacheCap),
				server.WithTextLogger(os.Stdout),
				server.WithJWTAuthToken(key),
			}

			if sslcrt == "" || sslkey == "" || sslca == "" {
				serverOptions = append(serverOptions, server.WithInsecureGRPCServer())
			} else {
				serverOptions = append(serverOptions, server.WithSecuredGRPCServer(sslcrt, sslkey, sslca))
			}

			if jsonOutput {
				serverOptions = append(serverOptions, server.WithJSONLogger(os.Stdout))
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
		SilenceUsage: true,
	}

	s.Flags().BoolVar(&jsonOutput, "json", false, "enables JSON formatted logs")
	s.Flags().IntVar(&cacheCap, "cache", 262144000, "capacity of the in-memory cache in bytes")
	s.Flags().StringVar(&collection, "collection", "commands", "collection used to store commands")
	s.Flags().StringVar(&database, "database", "./data/botio.db", "database path")
	s.Flags().StringVar(&httpPort, "http", ":8081", "port for HTTP server")
	s.Flags().StringVar(&key, "key", "", "key to generate a JWT token for authentication")
	s.Flags().StringVar(&port, "port", ":9091", "port for gRPC server")
	s.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	s.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")

	return s
}

func serverWithPostgresDB() *cobra.Command {
	var cacheCap int
	var database string
	var host string
	var httpPort string
	var jsonOutput bool
	var key string
	var maxConnLifetime time.Duration
	var maxConns int
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
		Example: "botio server postgres --user postgres --password toor --database botio --table commands --key mysupersecretkey",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverOptions := []server.Option{
				server.WithHTTPPort(httpPort),
				server.WithListener(port),
				// TODO: Clean pport
				server.WithPostgresDB(host, pport, database, table, user, password, maxConns, maxConnLifetime),
				server.WithRistrettoCache(cacheCap),
				server.WithTextLogger(os.Stdout),
				server.WithJWTAuthToken(key),
			}

			if sslcrt == "" || sslkey == "" || sslca == "" {
				serverOptions = append(serverOptions, server.WithInsecureGRPCServer())
			} else {
				serverOptions = append(serverOptions, server.WithSecuredGRPCServer(sslcrt, sslkey, sslca))
			}

			if jsonOutput {
				serverOptions = append(serverOptions, server.WithJSONLogger(os.Stdout))
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
		SilenceUsage: true,
	}

	s.Flags().BoolVar(&jsonOutput, "json", false, "enables JSON formatted logs")
	s.Flags().DurationVar(&maxConnLifetime, "maxConnLifetime", 2*time.Minute, "sets the lifetime of idle connections")
	s.Flags().IntVar(&cacheCap, "cache", 262144000, "capacity of the in-memory cache in bytes")
	s.Flags().IntVar(&maxConns, "maxConns", 5, "maximum number of open connections")
	s.Flags().StringVar(&database, "database", "botio", "PostgreSQL database name")
	s.Flags().StringVar(&host, "host", "postgres", "host of the PostgreSQL database")
	s.Flags().StringVar(&httpPort, "http", ":8081", "port for HTTP server")
	s.Flags().StringVar(&key, "key", "", "authentication key to generate a jwt token")
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
