package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	ListCmd.Flags().String("key", "", "authentication key for JWT")
	ListCmd.Flags().String("url", "", "botio's server URL")
}

// ListCmd is a cobra.Command to print all the commands available on the
// botio's commands server.
var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Prints a list with all the botio's commands",
	Example: "botio list --url :9090 --key mysupersecretkey",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		key := checkFlag(cmd, "key", false)
		url := checkFlag(cmd, "url", false)

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// GET commands
		commands, err := client.GetAll(url, key)
		if err != nil {
			log.Fatalf("%v", err)
		}

		printCommands(commands...)
	},
}
