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

	commands, err := db.GetAll(collection)
	if err != nil {
		log.Fatalf("while getting commands: %v", err)
	}

	actions := make(map[string]string, len(commands))
	for _, cmd := range commands {
		actions[cmd.Cmd] = cmd.Response
	}

	bot.HandleMessage("/.", func(m *tbot.Message) {
		log.Print(m.Chat.ID, m.Text)

		req := strings.TrimPrefix(m.Text, "/")
		c.SendChatAction(m.Chat.ID, tbot.ActionTyping)

		resp, ok := actions[req]
		if !ok {
			c.SendMessage(m.Chat.ID, " I'm sorry. I didn't understand you. Bzz")
		}

		time.Sleep(1 * time.Second)
		c.SendMessage(m.Chat.ID, resp)
	})

	log.Fatal(bot.Start())
}
