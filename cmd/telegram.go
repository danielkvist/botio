package cmd

import (
	"log"

	"github.com/danielkvist/botio/bot"
	"github.com/spf13/cobra"
)

func init() {
	TelegramBotCmd.Flags().String("key", "", "authentication key for JWT")
	TelegramBotCmd.Flags().String("token", "", "telegram's token")
	TelegramBotCmd.Flags().String("url", "", "botio's server URL")
}

// TelegramBotCmd is a cobra.Command to manage a Telegram bot.
var TelegramBotCmd = &cobra.Command{
	Use:     "telegram",
	Short:   "Initializes a Telegram bot that extracts the commands from the botio's server.",
	Example: "botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		key := checkFlag(cmd, "key", false)
		token := checkFlag(cmd, "token", false)
		url := checkFlag(cmd, "url", false)

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Bot initialization
		b := bot.Factory("telegram")
		b.Connect(token, 10)
		b.Listen(url, key)
		defer b.Stop()

		if err := b.Start(); err != nil {
			log.Fatalf("%v", err)
		}
	},
}
