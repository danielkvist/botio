package cmd

import (
	"log"

	"github.com/danielkvist/botio/bot"

	"github.com/spf13/cobra"
)

// Bot returns a *cobra.Command
func Bot() *cobra.Command {
	// var jwtToken string
	var goroutines int
	var platform string
	var serverName string
	var sslca string
	var sslcrt string
	var sslkey string
	var token string
	var url string

	b := &cobra.Command{
		Use:     "bot",
		Short:   "Initializes a bot for a supported platform (telegram and discord for the moment)",
		Example: "botio bot --platform telegram --token <telegram-token> --url :9090 --jwt <jwt-token>",
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := checkURL(url, false, false)
			if err != nil {
				return err
			}

			c, err := getClient(u, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
			}

			b, err := bot.Create(platform)
			if err != nil {
				log.Fatalf("%v", err)
			}

			b.Connect(c, u, token, goroutines)
			b.Listen()
			defer b.Stop()

			return b.Start()
		},
	}

	// b.Flags().StringVarP(&jwtToken, "jwt", "j", "", "jwt authenticaton token")
	b.Flags().IntVar(&goroutines, "goroutines", 10, "number of goroutines")
	b.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	b.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	b.Flags().StringVar(&sslcrt, "sslkey", "", "ssl certification key file")
	b.Flags().StringVar(&token, "token", "", "bot's token")
	b.Flags().StringVar(&url, "url", "", "botio's server URL")
	b.Flags().StringVar(&platform, "platform", "", "platform (discord or telegram)")

	return b
}
