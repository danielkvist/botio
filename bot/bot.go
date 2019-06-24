package bot

import (
	"log"
	"strings"
	"sync"

	"github.com/danielkvist/botio/db"

	"github.com/yanzay/tbot/v2"
)

type Bot struct {
	s  *tbot.Server
	c  *tbot.Client
	r  chan *Response
	wg sync.WaitGroup
}

type Response struct {
	id   string
	text string
}

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

func (b *Bot) Stop() {
	close(b.r)
	b.wg.Wait()
	b.s.Stop()
}

func (b *Bot) Start() error {
	return b.s.Start()
}
