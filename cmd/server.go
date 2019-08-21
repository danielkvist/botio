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
	ServerCmd.Flags().String("addr", "localhost:9090", "address where the server should listen for requests")
	ServerCmd.Flags().String("password", "toor", "password for basic auth")
	ServerCmd.Flags().String("user", "admin", "username for basic auth")
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "server initializes a botio's server to manage the Botio's commands with simple HTTP methods.",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		collection, _ := cmd.Flags().GetString("col")
		database, _ := cmd.Flags().GetString("db")
		listenAddr, _ := cmd.Flags().GetString("addr")
		password, _ := cmd.Flags().GetString("password")
		username, _ := cmd.Flags().GetString("user")

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

		s, err := server.New(
			server.WithListenAddr(listenAddr),
			server.WithHandler(r),
			server.WithBasicAuth(username, password),
			server.WithGracefulShutdown(done, quit),
		)
		if err != nil {
			log.Fatalf("while creating a new server: %v", err)
		}

		go func() {
			if err := s.ListenAndServe(); err != nil {
				log.Printf("%v", err)
				done <- struct{}{}
			}
		}()

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
		r.Get("/backup", handlers.Backup(database, col))
		r.Post("/", handlers.Post(database, col))
		r.Put("/{command}", handlers.Put(database, col))
		r.Delete("/{command}", handlers.Delete(database, col))
	})

	return r
}
