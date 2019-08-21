package cmd

import (
	"log"

	"github.com/danielkvist/botio/tbot"

	"github.com/spf13/cobra"
)

func init() {
	TelegramBotCmd.Flags().String("password", "toor", "password for basic auth")
	TelegramBotCmd.Flags().String("token", "", "Telegram's token")
	TelegramBotCmd.Flags().String("url", "", "URL where the botio's server is listening for requests")
	TelegramBotCmd.Flags().String("user", "admin", "username for basic auth")
}

var TelegramBotCmd = &cobra.Command{
	Use:   "tbot",
	Short: "Initializes a Telegram's bot that extracts the commands from the botio's server.",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		password, _ := cmd.Flags().GetString("password")
		token, _ := cmd.Flags().GetString("token")
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")

		// Check Telegram's token
		if token == "" {
			log.Fatalf("token for Telegram bot cannot be an empty string")
		}

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Bot initialization
		b := tbot.New(token, 10)
		b.Listen(url, user, password)
		defer b.Stop()

		if err := b.Start(); err != nil {
			log.Fatalf("%v", err)
		}
	},
}
