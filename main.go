package main

import (
	"log"
	"strings"

	"github.com/danielkvist/botio/cmd"

	"github.com/spf13/cobra"
)

func main() {
	examples := []string{
		"botio server --db ./data/commands.db --col commands --addr localhost:9090 --key mysupersecretkey",
		"botio tbot --token <telegram-token> --url localhost:9090 --key mysupersecretkey",
		"botio print --command start --url localhost:9090 --key mysupersecretkey",
	}

	root := &cobra.Command{
		Use:          "botio",
		Short:        "Simple CLI tool to create and manage easily bots for different platforms.",
		Example:      strings.Join(examples, "\n"),
		SilenceUsage: true,
	}

	root.AddCommand(cmd.ServerCmd)
	root.AddCommand(cmd.TelegramBotCmd)
	root.AddCommand(cmd.AddCmd)
	root.AddCommand(cmd.PrintCmd)
	root.AddCommand(cmd.ListCmd)
	root.AddCommand(cmd.UpdateCmd)
	root.AddCommand(cmd.DeleteCmd)

	if err := root.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
