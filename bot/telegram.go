package bot

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/logger"

	"github.com/yanzay/tbot/v2"
)

// Telegram is a wrapper for a yanzay/tbot client
// that satifies the Bot interface.
type Telegram struct {
	s  *tbot.Server
	c  *tbot.Client
	r  chan *Response
	wg sync.WaitGroup
}

// Connect receives a token with which tries to indentify,
// setups everything necessary and initializes a goroutine
// to send the responses from the responses channel to the respective
// clients.
func (t *Telegram) Connect(token string, cap int) error {
	s := tbot.New(token)
	c := s.Client()
	responses := make(chan *Response, cap)

	t.s = s
	t.c = c
	t.r = responses

	t.wg.Add(1)
	go func() {
		for r := range t.r {
			t.c.SendMessage(r.id, r.text)
		}
		t.wg.Done()
	}()

	return nil
}

// Listen handles all the messages sent to the Telegram bot
// and tries to get the response for the asked command from the botio's server
// and submit it to the responses channel, which eventually should send
// the response back to the client.
func (t *Telegram) Listen(url, key string) error {
	l := logger.New()
	t.s.HandleMessage(".", func(m *tbot.Message) {
		start := time.Now()
		msg := strings.TrimPrefix(m.Text, "/")
		resp := &Response{
			id: m.Chat.ID,
		}

		cmd, err := client.Get(url+"/"+msg, key)
		if err != nil {
			resp.text = "I'm sorry. I didn't understand you. Bzz"
			t.r <- resp
			return
		}

		resp.text = cmd.Response
		t.r <- resp

		l.Info(fmt.Sprintf("platform=%s id=%v msg=%q response=%q in=%v", "telegram", m.Chat.ID, msg, resp.text, time.Since(start)))
		return
	})

	return nil
}

// Start opens a connection to Telegram.
func (t *Telegram) Start() error {
	t.s.Start()
	return nil
}

// Stop waits until the responses channel is closed
// and then stops the Telegram session.
func (t *Telegram) Stop() error {
	close(t.r)
	t.wg.Wait()
	t.s.Stop()
	return nil
}
