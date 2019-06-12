package main

import (
	"log"
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

	go func() {
		bot.HandleMessage(".", func(m *tbot.Message) {
			log.Printf("%v: %s", m.Chat.ID, m.Text)

			resp := make(chan string)
			req := strings.TrimPrefix(m.Text, "/")
			c.SendChatAction(m.Chat.ID, tbot.ActionTyping)

			go func() {
				cmd, err := bdb.Get(collection, req)
				if err != nil || cmd.Response == "" {
					resp <- " I'm sorry. I didn't understand you. Bzz"
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
		s := server.New(bdb, collection, username, password, listenAddr)
		if err := s.ListenAndServe(); err != nil {
			log.Printf("%v", err)
			done <- struct{}{}
		}
	}()

	<-done
}
