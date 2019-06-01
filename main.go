package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/yanzay/tbot/v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("while loading env variables: %v", err)
	}

	ttoken := os.Getenv("TELEGRAM_TOKEN")
	bot := tbot.New(ttoken)
	c := bot.Client()

	bot.HandleMessage("/start", func(m *tbot.Message) {
		c.SendChatAction(m.Chat.ID, tbot.ActionTyping)
		time.Sleep(1 * time.Second)
		c.SendMessage(m.Chat.ID, "Hello, World!")
	})

	log.Fatal(bot.Start())
}
