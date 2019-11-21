package bot

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/proto"

	"github.com/sirupsen/logrus"
	"github.com/yanzay/tbot/v2"
)

// Telegram is a wrapper for a yanzay/tbot client
// that satifies the Bot interface.
type Telegram struct {
	tclient         *tbot.Client
	session         *tbot.Server
	responses       chan *Response
	defaultResponse string
	log             *logrus.Logger
	wg              sync.WaitGroup
	client          client.Client
}

// Connect receives a token with which tries to indentify,
// setups everything necessary and initializes a goroutine
// to send the responses from the responses channel to the respective
// clients.
func (t *Telegram) Connect(c client.Client, addr string, token string, cap int, defaultResponse string) error {
	session := tbot.New(token)
	tclient := session.Client()
	responses := make(chan *Response, cap)

	t.session = session
	t.tclient = tclient
	t.responses = responses
	t.defaultResponse = defaultResponse
	t.client = c

	t.log = logrus.New()
	t.log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  time.RFC850,
		QuoteEmptyFields: true,
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
		start := time.Now()

		msg := strings.TrimPrefix(m.Text, "/")
		resp := &Response{
			id: m.Chat.ID,
		}

		cmd, err := t.client.GetCommand(context.TODO(), &proto.Command{Command: msg})
		if err != nil {
			resp.text = t.defaultResponse
			t.responses <- resp

			logError(
				t.log,
				"Telegram",
				"client",
				"GetCommand",
				m.Chat.ID,
				msg,
				err.Error(),
				"error while responding to command",
			)
			return
		}

		resp.text = cmd.GetResp().GetResponse()
		t.responses <- resp

		logInfo(
			t.log,
			"Telegram",
			m.Chat.ID,
			msg,
			resp.text,
			"command responded successfully",
			time.Since(start),
		)
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
