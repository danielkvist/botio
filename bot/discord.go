package bot

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/proto"
	"github.com/pkg/errors"

	dg "github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Discord is a wrapper for a bwmawwin/discordgo session
// that satisfies the Bot interface.
type Discord struct {
	id              string
	session         *dg.Session
	responses       chan *Response
	defaultResponse string
	cancel          chan struct{}
	log             *logrus.Logger
	wg              sync.WaitGroup
	client          client.Client
}

// Connect receives a token with which tries to identify, setups
// everything necessary and initializes a goroutine to send
// the responses from the responses channel to the respective clients.
func (d *Discord) Connect(c client.Client, addr string, token string, cap int, defaultResponse string) error {
	session, err := dg.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("while creating a new Discord session: %v", err)
	}

	d.client = c
	id, err := session.User("@me")
	if err != nil {
		return fmt.Errorf("while extracting the current user ID for the bot: %v", err)
	}

	responses := make(chan *Response, cap)
	cancel := make(chan struct{})

	d.id = id.ID
	d.session = session
	d.responses = responses
	d.defaultResponse = defaultResponse
	d.cancel = cancel

	d.log = logrus.New()
	d.log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  time.RFC850,
		DisableSorting:   true,
		QuoteEmptyFields: true,
	})
	d.log.Out = os.Stdout

	d.wg.Add(1)
	go func() {
		for r := range d.responses {
			d.session.ChannelMessageSend(r.id, r.text)
		}

		d.wg.Done()
	}()

	return nil
}

// Listen handles all the messages sent to the Discord bot
// and tries to get the response for the asked command from the botio's server
// and submit it to the responses channel, which eventually
// should send the response back to the client.
func (d *Discord) Listen() error {
	d.session.AddHandler(func(s *dg.Session, m *dg.MessageCreate) {
		start := time.Now()

		if m.Author.Bot {
			return
		}

		resp := &Response{id: m.ChannelID}
		msg := strings.Fields(m.Content)
		if len(msg) != 2 {
			return
		}

		botID := "<@" + d.id + ">"
		if msg[0] != botID {
			return
		}

		cmd, err := d.client.GetCommand(context.TODO(), &proto.Command{Command: msg[1]})
		if err != nil {
			resp.text = d.defaultResponse
			d.responses <- resp

			logError(
				d.log,
				"Discord",
				"client",
				"GetCommand",
				m.ChannelID,
				m.Content,
				err.Error(),
				"error while responding to command",
			)
			return
		}

		resp.text = cmd.GetResp().GetResponse()
		d.responses <- resp

		logInfo(
			d.log,
			"Discord",
			m.ChannelID,
			m.Content,
			resp.text,
			"command responded successfully",
			time.Since(start),
		)
		return
	})

	return nil
}

// Start opens the connection to Discord.
func (d *Discord) Start() error {
	if err := d.session.Open(); err != nil {
		return errors.Wrap(err, "while opening a new Discord session")
	}

	<-d.cancel
	return nil
}

// Stop waits until the responses channel is closed
// and then closes the Discord session.
func (d *Discord) Stop() error {
	close(d.responses)
	close(d.cancel)
	d.wg.Wait()
	if err := d.session.Close(); err != nil {
		return errors.Wrap(err, "while closing a Discord session")
	}

	return nil
}
