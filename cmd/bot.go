package cmd

import (
	"log"

	"github.com/danielkvist/botio/bot"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Bot returns a *cobra.Command
func Bot() *cobra.Command {
	var jwtToken string
	var addr string
	var goroutines int
	var platform string
	var serverName string
	var defaultResp string
	var sslca string
	var sslcrt string
	var sslkey string
	var token string

	b := &cobra.Command{
		Use:     "bot",
		Short:   "Starts a chatbot for the specified platform.",
		Example: "botio bot --platform telegram --token <telegram-token>",
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := checkURL(addr, false, false)
			if err != nil {
				return err
			}

			c, err := getClient(u, token, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
			}

			b, err := bot.Create(platform)
			if err != nil {
				return errors.Wrapf(err, "while creating a new chatbot for platform %q: %v", platform, err)
			}

			b.Connect(c, u, token, goroutines, defaultResp)
			b.Listen()
			defer b.Stop()

			log.Printf("chatbot for platform %q initialized!\n", platform)
			if err := b.Start(); err != nil {
				return errors.Wrapf(b.Start(), "while starting chatbot for platform %q", platform)
			}

			return nil
		},
	}

	b.Flags().StringVar(&jwtToken, "jwt", "", "authenticaton token")
	b.Flags().IntVar(&goroutines, "goroutines", 10, "number of goroutines")
	b.Flags().StringVar(&addr, "addr", ":9091", "botio's gRPC server address")
	b.Flags().StringVar(&platform, "platform", "", "platform (discord or telegram)")
	b.Flags().StringVar(&defaultResp, "resp", "I'm sorry but something's happened and I can't answer that command rigth now", "default response for when the bot fails to respond to a command")
	b.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	b.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	b.Flags().StringVar(&sslcrt, "sslkey", "", "ssl certification key file")
	b.Flags().StringVar(&token, "token", "", "bot's token")

	return b
}
