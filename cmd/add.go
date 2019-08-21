package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	AddCmd.Flags().String("command", "", "command to add")
	AddCmd.Flags().String("password", "toor", "password for basic auth")
	AddCmd.Flags().String("response", "", "response of the command to add")
	AddCmd.Flags().String("url", "", "URL where the botio's server is listening")
	AddCmd.Flags().String("user", "admin", "username for basic auth")

}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new command with a response to the botio's server",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		command, _ := cmd.Flags().GetString("command")
		password, _ := cmd.Flags().GetString("password")
		response, _ := cmd.Flags().GetString("response")
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")

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
		c, err := client.Post(url, user, password, command, response)
		if err != nil {
			log.Fatalf("while posting command %q to %q: %v", command, url, err)
		}

		printCommands(c)
	},
}
