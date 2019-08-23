// Package tbot exports a wrapper for a yanzay/tbot's Client and Server.
package tbot

import (
	"log"
	"strings"
	"sync"

	"github.com/danielkvist/botio/client"

	"github.com/yanzay/tbot/v2"
)

// Bot wraps a tbot's Client and Server along with a channel
// to send responses concurrently.
type Bot struct {
	s  *tbot.Server
	c  *tbot.Client
	r  chan *Response
	wg sync.WaitGroup
}

// Response represents a Bot's response.
type Response struct {
	id   string
	text string
}

// New returns a *Bot with the tbot's Server, Client and the response channel set up.
// It receives a Telegram's token to initialize the tbot's Server and a capacity
// for the response channel.
// It also initializes a goroutine to send all the responses from the
// responses channel to the respective clients.
func New(token string, cap int) *Bot {
	server := tbot.New(token)
	c := server.Client()
	responses := make(chan *Response, cap)

	bot := &Bot{
		s: server,
		c: c,
		r: responses,
	}

	bot.wg.Add(1)
	go func() {
		for r := range bot.r {
			bot.c.SendMessage(r.id, r.text)
		}
		bot.wg.Done()
	}()

	return bot
}

// Listen handles all the messages received from the Bot tbot's Server
// and tries to get the asked command from the received server's URL using
// an user and a password to authenticate.
func (b *Bot) Listen(url, key string) {
	b.s.HandleMessage(".", func(m *tbot.Message) {
		log.Printf("%s\t%s\t%s", m.Chat.ID, m.Chat.Username, m.Text)
		msg := strings.TrimPrefix(m.Text, "/")
		resp := &Response{
			id: m.Chat.ID,
		}

		cmd, err := client.Get(url+"/"+msg, key)
		if err != nil {
			resp.text = "I'm sorry. I didn't understand you. Bzz"
			b.r <- resp
			return
		}

		resp.text = cmd.Response
		b.r <- resp
	})
}

// Start initializes the Bot tbot's Server.
func (b *Bot) Start() error {
	return b.s.Start()
}

// Stop waits until the channel for the Bot's responses is closed
// and then closes the Bot tbot's Server.
func (b *Bot) Stop() {
	close(b.r)
	b.wg.Wait()
	b.s.Stop()
}
