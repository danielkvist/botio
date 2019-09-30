package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/danielkvist/botio/server"

	"github.com/spf13/cobra"
)

// Server returns a *cobra.Command.
func Server() *cobra.Command {
	return serverCmd(serverWithBoltDB(), serverWithPostgresDB())
}

func serverCmd(commands ...*cobra.Command) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "server contains some subcommands to initialize a server with different databases",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	for _, cmd := range commands {
		serverCmd.AddCommand(cmd)
	}

	return serverCmd
}

func serverWithBoltDB() *cobra.Command {
	var collection string
	var database string
	var porthttp string
	var porthttps string
	var key string
	var sslcert string
	var sslkey string

	s := &cobra.Command{
		Use:     "bolt",
		Short:   "Starts a server with a BoltDB database to manage your commands with HTTP methods",
		Example: "botio server bolt --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			var tls bool
			if sslcert != "" && sslkey != "" {
				tls = true
			}

			serverOptions := []server.Option{
				server.WithBoltDB(database, collection),
				server.WithJWTMiddleware(key),
			}

			s := server.New(serverOptions...)
			if err := listenAndServe(porthttp, porthttps, s, tls, sslcert, sslkey); err != nil {
				log.Printf("%v", err)
			}
		},
	}

	s.Flags().StringVar(&collection, "collection", "commands", "collection used to store commands")
	s.Flags().StringVar(&database, "database", "./commands.db", "database path")
	s.Flags().StringVar(&porthttp, "http", ":80", "port for HTTP connections")
	s.Flags().StringVar(&porthttps, "https", ":443", "port for HTTPS connections")
	s.Flags().StringVar(&key, "key", "", "authentication key")
	s.Flags().StringVar(&sslcert, "sslcert", "", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")

	return s
}

func serverWithPostgresDB() *cobra.Command {
	var host string
	var port string
	var user string
	var password string
	var table string
	var database string
	var porthttp string
	var porthttps string
	var key string
	var sslcert string
	var sslkey string

	s := &cobra.Command{
		Use:     "postgres",
		Short:   "Starts a server with that connects to a PostgreSQL database to manage your commands with HTTP methods",
		Example: "botio server postgres --user postgres --password toor --database botio --table commands --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			var tls bool
			if sslcert != "" && sslkey != "" {
				tls = true
			}

			serverOptions := []server.Option{
				server.WithPostgresDB(host, port, database, table, user, password),
				server.WithJWTMiddleware(key),
			}

			s := server.New(serverOptions...)
			if err := listenAndServe(porthttp, porthttps, s, tls, sslcert, sslkey); err != nil {
				log.Printf("%v", err)
			}
		},
	}

	s.Flags().StringVar(&host, "host", "postgres", "host of the PostgreSQL database")
	s.Flags().StringVar(&port, "port", "5432", "port of the PostgreSQL database host")
	s.Flags().StringVar(&database, "database", "botio", "PostgreSQL database name")
	s.Flags().StringVar(&table, "table", "", "table of the PostgreSQL database")
	s.Flags().StringVar(&user, "user", "", "user of the PostgreSQL database")
	s.Flags().StringVar(&password, "password", "", "password for the user of the PostgreSQL database")
	s.Flags().StringVar(&porthttp, "http", ":80", "port for HTTP connections")
	s.Flags().StringVar(&porthttps, "https", ":443", "port for HTTPS connections")
	s.Flags().StringVar(&key, "key", "", "authentication key")
	s.Flags().StringVar(&sslcert, "sslcert", "", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")

	return s
}

func listenAndServe(httpaddr string, httpsaddr string, h http.Handler, tls bool, sslcert string, sslkey string) error {
	s := &http.Server{
		Addr:         httpaddr,
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if tls {
		return listenAndServeHTTPS(httpsaddr, s, sslcert, sslkey)
	}

	return s.ListenAndServe()
}

func listenAndServeHTTPS(addr string, s *http.Server, sslcert string, sslkey string) error {
	tlsConf := &tls.Config{
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	s.TLSConfig = tlsConf
	s.Addr = addr

	return s.ListenAndServeTLS(sslcert, sslkey)
}
