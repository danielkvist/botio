package main

import (
	"log"
	"strings"

	"github.com/danielkvist/botio/cmd"

	"github.com/spf13/cobra"
)

func main() {
	examples := []string{
		"botio server --db ./data/commands.db --col commands --http :9090 --key mysupersecretkey",
		"botio telegram --token <telegram-token> --url :9090 --key mysupersecretkey",
		"botio print --command start --url :9090 --key mysupersecretkey",
	}

	root := &cobra.Command{
		Use:          "botio",
		Short:        "Simple CLI tool to create and manage easily bots for different platforms.",
		Example:      strings.Join(examples, "\n"),
		SilenceUsage: true,
	}

	// Server
	root.AddCommand(cmd.ServerCmd)

	// Bots
	root.AddCommand(cmd.TelegramBotCmd)
	root.AddCommand(cmd.DiscordBotCmd)

	// Client
	root.AddCommand(cmd.AddCmd)
	root.AddCommand(cmd.PrintCmd)
	root.AddCommand(cmd.ListCmd)
	root.AddCommand(cmd.UpdateCmd)
	root.AddCommand(cmd.DeleteCmd)

	if err := root.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
