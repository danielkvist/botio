package main

import (
	"log"

	"github.com/danielkvist/botio/cmd"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:          "botio",
		Short:        "Botio is a simple CLI tool to create and manage easily Telegram Bots.",
		SilenceUsage: true,
	}

	root.AddCommand(cmd.ServerCmd)
	root.AddCommand(cmd.TelegramBotCmd)
	root.AddCommand(cmd.ListCmd)

	if err := root.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
