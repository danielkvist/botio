// Package bot exports a wrapper to work with Telegram bots.
package bot

import (
	"log"
	"strings"
	"sync"

	"github.com/danielkvist/botio/db"

	"github.com/yanzay/tbot/v2"
)

// Bot wraps a bot client and server along with a channel
// to dispatch concurrently the responses.
type Bot struct {
	s  *tbot.Server
	c  *tbot.Client
	r  chan *Response
	wg sync.WaitGroup
}

// Response represents a bot response with an ID
// and a text.
type Response struct {
	id   string
	text string
}

// New returns an initialized *Bot ready to respond with
// the server, the client and the channel for the responses
// already set up.
func New(token string, cap int) *Bot {
	server := tbot.New(token)
	client := server.Client()
	responses := make(chan *Response, cap)

	bot := &Bot{
		s: server,
		c: client,
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

// HandlerMessage creates a handler to manage the specified message
// and send a response to the channel for responses. Which eventually
// will send the response to the user.
//
// Actually it only works with commands.
func (b *Bot) HandlerMessage(msg string, bolter db.Bolter, col string) {
	b.s.HandleMessage(msg, func(m *tbot.Message) {
		log.Printf("%s\t%s\t%s", m.Chat.ID, m.Chat.Username, m.Text)
		req := strings.TrimPrefix(m.Text, "/")
		response := &Response{
			id: m.Chat.ID,
		}

		cmd, err := bolter.Get(col, req)
		if err != nil {
			response.text = "I'm sorry. I didn't understand you. Bzz"
			b.r <- response
			return
		}

		response.text = cmd.Response
		b.r <- response
	})
}

// Stop waits until the channel for the responses is closed
// and then closes the bot server.
func (b *Bot) Stop() {
	close(b.r)
	b.wg.Wait()
	b.s.Stop()
}

// Start starts the bot server.
func (b *Bot) Start() error {
	return b.s.Start()
}
