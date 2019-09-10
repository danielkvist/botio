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
	var collection string
	var database string
	var porthttp string
	var porthttps string
	var key string
	var sslcert string
	var sslkey string

	s := &cobra.Command{
		Use:     "server",
		Short:   "Starts a server to manage the commands with simple HTTP methods.",
		Example: "botio server --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey",
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
			done := make(chan struct{}, 1)

			go func() {
				err := listenAndServe(porthttp, porthttps, s, tls, sslcert, sslkey)
				if err != nil {
					log.Printf("%v", err)
					done <- struct{}{}
				}
			}()
			<-done
		},
	}

	s.Flags().StringVarP(&collection, "collection", "c", "commands", "collection used to store commands")
	s.Flags().StringVarP(&database, "database", "d", "./commands.db", "database path")
	s.Flags().StringVar(&porthttp, "http", ":80", "port for HTTP connections")
	s.Flags().StringVar(&porthttps, "https", ":443", "port for HTTPS connections")
	s.Flags().StringVarP(&key, "key", "k", "", "authentication key")
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
