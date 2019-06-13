package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danielkvist/botio/db"
	"github.com/danielkvist/botio/server"

	"github.com/joho/godotenv"
	"github.com/yanzay/tbot/v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("while loading env variables: %v", err)
	}

	ttoken := os.Getenv("TELEGRAM_TOKEN")
	database := os.Getenv("DATABASE")
	collection := os.Getenv("COLLECTION")
	listenAddr := os.Getenv("LISTEN_ADDRESS")
	username := os.Getenv("API_USERNAME")
	password := os.Getenv("API_PASSWORD")

	bot := tbot.New(ttoken)
	c := bot.Client()

	bdb, err := db.Connect(database, collection)
	if err != nil {
		log.Fatalf("while connecting to the database: %v", err)
	}

	done := make(chan struct{}, 2)
	quit := make(chan struct{}, 1)

	s := server.New(bdb, collection, username, password, listenAddr)

	go func() {
		bot.HandleMessage(".", func(m *tbot.Message) {
			log.Printf("%v: %s", m.Chat.ID, m.Text)

			resp := make(chan string)
			req := strings.TrimPrefix(m.Text, "/")
			c.SendChatAction(m.Chat.ID, tbot.ActionTyping)

			go func() {
				cmd, err := bdb.Get(collection, req)
				if err != nil {
					resp <- " I'm sorry. I didn't understand you. Bzz"
					log.Printf("%v: %v", m.Chat.ID, err)
					return
				}

				resp <- cmd.Response
			}()

			time.Sleep(1 * time.Second)
			c.SendMessage(m.Chat.ID, <-resp)
		})

		if err := bot.Start(); err != nil {
			log.Printf("%v", err)
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
