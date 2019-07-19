package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/danielkvist/botio/bot"
	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/server"
)

func main() {
	ttoken := flag.String("token", "", "telegram's bot token")
	database := flag.String("db", "./data/commands.db", "where the database is supposed to be or should be")
	listenAddr := flag.String("address", "localhost:9090", "TCP address to listen on for requests")
	username := flag.String("username", "admin", "username for basic authentication")
	password := flag.String("password", "toor", "password for basic authentication")
	flag.Parse()

	commands := "commands"

	if *ttoken == "" {
		log.Fatal("it's needed a valid token for a telegram's bot")
	}

	bdb, err := db.Connect(*database, commands)
	if err != nil {
		log.Fatalf("while connecting to the database: %v", err)
	}

	done := make(chan struct{}, 2)
	quit := make(chan struct{}, 1)

	b := bot.New(*ttoken, 10)
	s := server.New(bdb, commands, *username, *password, *listenAddr)

	go func() {
		b.HandlerMessage(".", bdb, commands)

		if err := b.Start(); err != nil {
			log.Printf("%v", err)
			b.Stop()
			done <- struct{}{}
		}
	}()

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("%v", err)
			done <- struct{}{}
		}
	}()

	go graceShutdown(s, done, quit)

	<-quit
}

func graceShutdown(s *http.Server, done <-chan struct{}, quit chan<- struct{}) {
	<-done
	log.Printf("shutting down server listening on address %q", s.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("while trying to shutdown the server listening on address %q: %v", s.Addr, err)
	}

	quit <- struct{}{}
}
