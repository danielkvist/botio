package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/danielkvist/botio/bot"
	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/handlers"
	"github.com/danielkvist/botio/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// Flags
	ttoken := flag.String("token", "", "telegram's bot token")
	database := flag.String("db", "./data/commands.db", "where the database is supposed to be or should be")
	listenAddr := flag.String("address", ":9090", "TCP address to listen on for requests")
	username := flag.String("username", "admin", "username for basic authentication")
	password := flag.String("password", "toor", "password for basic authentication")
	flag.Parse()

	commands := "commands"
	env := "production"

	if *ttoken == "" {
		log.Fatal("it's needed a valid token for a telegram's bot")
	}

	// Database initialization
	bdb := db.DBFactory(env)
	err := bdb.Open(*database, commands)
	if err != nil {
		log.Fatalf("while connecting to the database: %v", err)
	}

	done := make(chan struct{}, 2)
	quit := make(chan struct{}, 1)

	// Telegram and Server initialization
	b := bot.New(*ttoken, 10)
	r := newRouter(bdb, commands)
	s, err := server.New(
		server.WithListenAddr(*listenAddr),
		server.WithHandler(r),
		server.WithGracefulShutdown(done, quit),
		server.WithBasicAuth(*username, *password, r),
	)
	if err != nil {
		log.Fatalf("while creating new server: %v", err)
	}

	// Telegram Bot
	go func() {
		b.HandlerMessage(".", bdb, commands)

		if err := b.Start(); err != nil {
			log.Printf("%v", err)
			b.Stop()
			done <- struct{}{}
		}
	}()

	// HTTP server
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("%v", err)
			done <- struct{}{}
		}
	}()

	<-quit
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
