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

func init() {
	ServerCmd.Flags().String("col", "commands", "collection used to store the commands")
	ServerCmd.Flags().String("db", "./botio/botio.db", "path to the database")
	ServerCmd.Flags().String("http", ":80", "port for HTTP connections")
	ServerCmd.Flags().String("https", ":443", "port for HTTPS connections")
	ServerCmd.Flags().String("key", "", "authentication key for JWT")
	ServerCmd.Flags().String("sslcert", "", "ssl certification")
	ServerCmd.Flags().String("sslkey", "", "ssl key")
}

// ServerCmd is a cobra.Command to manage the botio's commands server.
var ServerCmd = &cobra.Command{
	Use:     "server",
	Short:   "Starts a botio's server to manage the botio's commands with simple HTTP methods.",
	Example: "botio server --db ./data/botio.db --col commands --http :9090 --key mysupersecretkey",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		collection := checkFlag(cmd, "col", false)
		database := checkFlag(cmd, "db", false)
		httpPort := checkFlag(cmd, "http", false)
		httpsPort := checkFlag(cmd, "https", false)
		key := checkFlag(cmd, "key", false)
		sslCert := checkFlag(cmd, "sslcert", true)
		sslKey := checkFlag(cmd, "sslkey", true)

		// TLS
		var tls bool
		if sslCert != "" && sslKey != "" {
			tls = true
		}

		// Database initialization
		env := "production"
		bdb := db.Factory(env)
		err := bdb.Open(database, collection)
		if err != nil {
			log.Fatalf("while opening a connection with database: %v", err)
		}

		// Server initialization
		done := make(chan struct{}, 1)
		quit := make(chan struct{}, 1)

		r := newRouter(bdb, collection)
		serverOptions := []server.Option{
			server.WithListenAddr(httpPort),
			server.WithHandler(r),
			server.WithJWTAuth(key),
			server.WithGracefulShutdown(done, quit),
		}

		if tls {
			serverOptions = append(serverOptions, server.WithListenAddr(httpsPort))
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
