package cmd

import (
	"log"

	"github.com/danielkvist/botio/bot"

	"github.com/spf13/cobra"
)

// Bot returns a *cobra.Command
func Bot() *cobra.Command {
	var platform string
	var jwtToken string
	var token string
	var url string

	b := &cobra.Command{
		Use:     "bot",
		Short:   "Initializes a bot for a supported platform (telegram and discord for the moment)",
		Example: "botio bot --platform telegram --token <telegram-token> --url :9090 --jwt <jwt-token>",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			b, err := bot.Create(platform)
			if err != nil {
				log.Fatalf("%v", err)
			}

			b.Connect(token, 10)
			b.Listen(u, jwtToken)
			defer b.Stop()

			if err := b.Start(); err != nil {
				log.Fatalf("%v", err)
			}
		},
	}

	b.Flags().StringVarP(&platform, "platform", "p", "", "platform (discord or telegram)")
	b.Flags().StringVarP(&jwtToken, "jwt", "j", "", "jwt authenticaton token")
	b.Flags().StringVarP(&token, "token", "t", "", "bot's token")
	b.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return b
}
