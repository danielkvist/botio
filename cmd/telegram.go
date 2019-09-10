package cmd

import (
	"log"

	"github.com/danielkvist/botio/bot"
	"github.com/spf13/cobra"
)

// Telegram returns a *cobra.Command.
func Telegram() *cobra.Command {
	var key string
	var token string
	var url string

	t := &cobra.Command{
		Use:     "telegram",
		Short:   "Initializes a Telegram bot",
		Example: "botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			b := bot.Create("telegram")
			b.Connect(token, 10)
			b.Listen(u, key)
			defer b.Stop()

			if err := b.Start(); err != nil {
				log.Fatalf("%v", err)
			}
		},
	}

	t.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	t.Flags().StringVarP(&token, "token", "t", "", "telegram's token")
	t.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return t
}
