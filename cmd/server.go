package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/handlers"
	"github.com/danielkvist/botio/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/cobra"
)

// Server returns a *cobra.Command
func Server() *cobra.Command {
	var collection string
	var database string
	var http string
	var https string
	var key string
	var sslcert string
	var sslkey string

	s := &cobra.Command{
		Use:     "server",
		Short:   "Starts a server to manage the commands with simple HTTP methods.",
		Example: "botio server --database ./data/botio.db --collection commands --http :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			c := checkFlag("collection", collection, false)
			d := checkFlag("database", database, false)
			portHTTP := checkFlag("http", http, false)
			portHTTPS := checkFlag("https", https, false)
			k := checkFlag("key", key, false)
			sslCert := checkFlag("sslcert", sslcert, true)
			sslKey := checkFlag("sslkey", sslkey, true)

			var tls bool
			if sslCert != "" && sslKey != "" {
				tls = true
			}

			env := "production"
			bdb := db.Factory(env)
			err := bdb.Open(d, c)
			if err != nil {
				log.Fatalf("while opening a connection with database: %v", err)
			}

			done := make(chan struct{}, 1)
			quit := make(chan struct{}, 1)

			r := newRouter(bdb, collection)
			serverOptions := []server.Option{
				server.WithListenAddr(portHTTP),
				server.WithHandler(r),
				server.WithJWTAuth(k),
				server.WithGracefulShutdown(done, quit),
			}

			if tls {
				serverOptions = append(serverOptions, server.WithListenAddr(portHTTPS))
				serverOptions = append(serverOptions, server.WithTLS())
			}

			s, err := server.New(serverOptions...)
			if err != nil {
				log.Fatalf("while creating a new server: %v", err)
			}

			go listenAndServe(s, tls, sslCert, sslKey, done)
			<-quit
		},
	}

	s.Flags().StringVarP(&collection, "collection", "c", "commands", "collection used to store commands")
	s.Flags().StringVarP(&database, "database", "d", "./commands.db", "database path")
	s.Flags().StringVar(&http, "http", ":80", "port for HTTP connections")
	s.Flags().StringVar(&https, "https", ":443", "port for HTTPS connections")
	s.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	s.Flags().StringVar(&sslcert, "sslcert", "", "ssl certification file")
	s.Flags().StringVar(&sslkey, "sslkey", "", "ssl key file")

	return s
}

func newRouter(database db.DB, col string) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.URLFormat)

	// Routes
	r.Route("/api/commands", func(r chi.Router) {
		r.Get("/", handlers.GetAll(database, col))
		r.Get("/{command}", handlers.Get(database, col))
		r.Post("/", handlers.Post(database, col))
		r.Put("/{command}", handlers.Put(database, col))
		r.Delete("/{command}", handlers.Delete(database, col))
	})

	r.Route("/api/backup", func(r chi.Router) {
		r.Get("/", handlers.Backup(database, col))
	})

	return r
}

func listenAndServe(s *http.Server, tls bool, sslCert, sslKey string, done chan<- struct{}) {
	var err error
	if tls {
		err = s.ListenAndServeTLS(sslCert, sslKey)
	} else {
		err = s.ListenAndServe()
	}

	if err != nil {
		log.Printf("%v", err)
		done <- struct{}{}
	}
}
