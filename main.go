package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/danielkvist/botio/db"

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

	bot := tbot.New(ttoken)
	c := bot.Client()

	db, err := db.Open(database)
	if err != nil {
		log.Fatalf("while connecting to the database: %v", err)
	}

	bot.HandleMessage(".", func(m *tbot.Message) {
		log.Printf("%v: %s", m.Chat.ID, m.Text)

		resp := make(chan string)
		req := strings.TrimPrefix(m.Text, "/")
		c.SendChatAction(m.Chat.ID, tbot.ActionTyping)

		go func() {
			cmd, err := db.Get(collection, req)
			if err != nil || cmd.Response == "" {
				resp <- " I'm sorry. I didn't understand you. Bzz"
			}

			resp <- cmd.Response
		}()

		time.Sleep(1 * time.Second)
		c.SendMessage(m.Chat.ID, <-resp)
	})

	log.Fatal(bot.Start())
}
