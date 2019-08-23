package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	AddCmd.Flags().String("command", "", "command to add")
	AddCmd.Flags().String("key", "", "authentication key")
	AddCmd.Flags().String("response", "", "response of the command to add")
	AddCmd.Flags().String("url", "", "url where the botio's server is listening")

}

// AddCmd is a cobra.Command to add commands to the botio's commands server.
var AddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Adds a new command with a response to the botio's server",
	Example: "botio add --command start --response Hello --url localhost:9090 --key mysupersecretkey",
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

		// POST command
		c, err := client.Post(url, key, command, response)
		if err != nil {
			log.Fatalf("%v", err)
		}

		printCommands(c)
	},
}
