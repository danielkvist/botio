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
			k := checkFlag("key", key, false)
			t := checkFlag("token", token, false)
			u := checkFlag("url", url, false)

			u, err := checkURL(u)
			if err != nil {
				log.Fatalf("%v", err)
			}

			b := bot.Factory("telegram")
			b.Connect(t, 10)
			b.Listen(u, k)
			defer b.Stop()

			if err := b.Start(); err != nil {
				log.Fatalf("%v", err)
			}
		},
		Args: cobra.ExactArgs(3),
	}

	t.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	t.Flags().StringVarP(&token, "token", "t", "", "telegram's token")
	t.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return t
}
