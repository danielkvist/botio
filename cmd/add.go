package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	AddCmd.Flags().String("command", "", "command to add")
	AddCmd.Flags().String("key", "", "authentication key for JWT")
	AddCmd.Flags().String("response", "", "response of the command to add")
	AddCmd.Flags().String("url", "", "botio's server URL")
}

// AddCmd is a cobra.Command to add commands to the botio's commands server.
var AddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Adds a new command with a response to the botio's server",
	Example: "botio add --command start --response Hello --url :9090 --key mysupersecretkey",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		command := checkFlag(cmd, "command", false)
		key := checkFlag(cmd, "key", false)
		response := checkFlag(cmd, "response", false)
		url := checkFlag(cmd, "url", false)

		// Check command and response
		if command == "" || response == "" {
			log.Fatal("either command or response cannot be an empty string")
		}

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// POST command
		c, err := client.Post(url, key, command, response)
		if err != nil {
			log.Fatalf("%v", err)
		}

		printCommands(c)
	},
}
