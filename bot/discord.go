package bot

import (
	"fmt"
	"strings"
	"sync"

	dg "github.com/bwmarrin/discordgo"
	"github.com/danielkvist/botio/client"
)

// Discord is a wrapper for a bwmawwin/discordgo session
// that satisfies the Bot interface.
type Discord struct {
	id string
	s  *dg.Session
	r  chan *Response
	c  chan struct{}
	wg sync.WaitGroup
}

// Connect receives a token with which tries to identify, setups
// everything necessary and initializes a goroutine to send
// the responses from the responses channel to the respective clients.
func (d *Discord) Connect(token string, cap int) error {
	s, err := dg.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("while creating a new Discord session: %v", err)
	}

	id, err := s.User("@me")
	if err != nil {
		return fmt.Errorf("while extracting the current user ID for the bot: %v", err)
	}

	responses := make(chan *Response, cap)
	c := make(chan struct{})

	d.id = id.ID
	d.s = s
	d.c = c
	d.r = responses

	d.wg.Add(1)
	go func() {
		for r := range d.r {
			d.s.ChannelMessageSend(r.id, r.text)
		}
		d.wg.Done()
	}()

	return nil
}

// Listen handles all the messages sent to the Discord bot
// and tries to get the response for the asked command from the botio's server
// and submit it to the responses channel, which eventually
// should send the response back to the client.
func (d *Discord) Listen(url, key string) error {
	d.s.AddHandler(func(s *dg.Session, m *dg.MessageCreate) {
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

		cmd, err := client.Get(url+"/"+msg[1], key)
		if err != nil {
			return
		}

		resp.text = cmd.Response
		d.r <- resp
		return
	})
	return nil
}

// Start opens the connection to Discord.
func (d *Discord) Start() error {
	if err := d.s.Open(); err != nil {
		return fmt.Errorf("while opening a connection: %v", err)
	}

	<-d.c
	return nil
}

// Stop waits until the responses channel is closed
// and then closes the Discord session.
func (d *Discord) Stop() error {
	close(d.r)
	close(d.c)
	d.wg.Wait()
	if err := d.s.Close(); err != nil {
		return fmt.Errorf("while closing a connection: %v", err)
	}

	return nil
}
