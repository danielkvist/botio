package cmd

import (
	"log"

	"github.com/danielkvist/botio/bot"

	"github.com/spf13/cobra"
)

// Discord returns a *cobra.Command.
func Discord() *cobra.Command {
	var key string
	var token string
	var url string

	d := &cobra.Command{
		Use:     "discord",
		Short:   "Initializes a Discord bot",
		Example: "botio discord --token <discord-token> --url :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			k := checkFlag("key", key, false)
			t := checkFlag("token", token, false)
			u := checkFlag("url", url, false)

			u, err := checkURL(u)
			if err != nil {
				log.Fatalf("%v", err)
			}

			b := bot.Factory("discord")
			b.Connect(t, 10)
			b.Listen(u, k)
			defer b.Stop()

			if err := b.Start(); err != nil {
				log.Fatalf("%v", err)
			}
		},
		Args: cobra.ExactArgs(3),
	}

	d.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	d.Flags().StringVarP(&token, "token", "t", "", "discord's token")
	d.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return d
}