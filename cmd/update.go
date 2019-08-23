package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	UpdateCmd.Flags().String("command", "", "command to add")
	UpdateCmd.Flags().String("key", "", "authentication key")
	UpdateCmd.Flags().String("response", "", "response of the command to add")
	UpdateCmd.Flags().String("url", "", "url where the botio's server is listening")
}

// UpdateCmd is a cobra.Command to update commands on the botio's commands server.
var UpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Updates an existing command (or adds it if not exists) with a response on the botio's server",
	Example: "botio update --command start --response Hi --url localhost:9090 --key mysupersecretkey",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		command, _ := cmd.Flags().GetString("command")
		key, _ := cmd.Flags().GetString("key")
		response, _ := cmd.Flags().GetString("response")
		url, _ := cmd.Flags().GetString("url")

		// Check command and response
		if command == "" || response == "" {
			log.Fatal("either command or response cannot be an empty string")
		}

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// PUT command
		c, err := client.Put(url, key, command, response)
		if err != nil {
			log.Fatalf("%v", err)
		}

		printCommands(c)
	},
}
