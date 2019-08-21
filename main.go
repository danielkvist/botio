package main

import (
	"log"

	"github.com/danielkvist/botio/cmd"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "botio",
		Short: "Simple CLI tool to create and manage easily Telegram Bots.",
		// TODO: Long
		SilenceUsage: true, // TODO:
	}

	root.AddCommand(cmd.ServerCmd)
	root.AddCommand(cmd.TelegramBotCmd)
	root.AddCommand(cmd.AddCmd)
	root.AddCommand(cmd.PrintCmd)
	root.AddCommand(cmd.ListCmd)
	root.AddCommand(cmd.DeleteCmd)

	if err := root.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
