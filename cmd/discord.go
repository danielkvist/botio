package cmd

import (
	"log"

	"github.com/danielkvist/botio/bot"

	"github.com/spf13/cobra"
)

func init() {
	DiscordBotCmd.Flags().String("key", "", "authentication key for JWT")
	DiscordBotCmd.Flags().String("token", "", "discord's token")
	DiscordBotCmd.Flags().String("url", "", "botio's server URL")
}

// DiscordBotCmd is a cobra.Command to manage a Discord bot.
var DiscordBotCmd = &cobra.Command{
	Use:     "discord",
	Short:   "Initializes a Discord bot that extracts the commands from the botio's server",
	Example: "botio discord --token <discord-token> --url :9090 --key mysupersecretkey",
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
		b := bot.Factory("discord")
		b.Connect(token, 10)
		b.Listen(url, key)
		defer b.Stop()

		if err := b.Start(); err != nil {
			log.Fatalf("%v", err)
		}
	},
}
