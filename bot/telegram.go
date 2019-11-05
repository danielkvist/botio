package bot

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/proto"

	"github.com/sirupsen/logrus"
	"github.com/yanzay/tbot/v2"
)

// Telegram is a wrapper for a yanzay/tbot client
// that satifies the Bot interface.
type Telegram struct {
	session   *tbot.Server
	client    client.Client
	tclient   *tbot.Client
	responses chan *Response
	log       *logrus.Logger
	wg        sync.WaitGroup
}

// Connect receives a token with which tries to indentify,
// setups everything necessary and initializes a goroutine
// to send the responses from the responses channel to the respective
// clients.
func (t *Telegram) Connect(c client.Client, addr string, token string, cap int) error {
	session := tbot.New(token)
	tclient := session.Client()
	responses := make(chan *Response, cap)

	t.session = session
	t.client = c
	t.tclient = tclient
	t.responses = responses
	t.log = logrus.New()
	t.log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  "02-01-2006 15:04:05",
	})
	t.log.Out = os.Stdout

	t.wg.Add(1)
	go func() {
		for r := range t.responses {
			t.tclient.SendMessage(r.id, r.text)
		}
		t.wg.Done()
	}()

	return nil
}

// Listen handles all the messages sent to the Telegram bot
// and tries to get the response for the asked command from the botio's server
// and submit it to the responses channel, which eventually should send
// the response back to the client.
func (t *Telegram) Listen() error {
	t.session.HandleMessage(".", func(m *tbot.Message) {
		msg := strings.TrimPrefix(m.Text, "/")
		resp := &Response{
			id: m.Chat.ID,
		}

		cmd, err := t.client.GetCommand(context.TODO(), &proto.Command{Command: msg})
		if err != nil {
			resp.text = "I'm sorry. I didn't understand you. Bzz"
			t.responses <- resp
			return
		}

		resp.text = cmd.GetResp().GetResponse()
		t.responses <- resp

		log(t.log, "telegram", m.Chat.ID, msg, resp.text)
		return
	})

	return nil
}

// Start opens a connection to Telegram.
func (t *Telegram) Start() error {
	t.session.Start()
	return nil
}

// Stop waits until the responses channel is closed
// and then stops the Telegram session.
func (t *Telegram) Stop() error {
	close(t.responses)
	t.wg.Wait()
	t.session.Stop()
	return nil
}
