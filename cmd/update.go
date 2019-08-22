package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	UpdateCmd.Flags().String("command", "", "command to add")
	UpdateCmd.Flags().String("password", "toor", "password for basic auth")
	UpdateCmd.Flags().String("response", "", "response of the command to add")
	UpdateCmd.Flags().String("url", "", "URL where the botio's server is listening")
	UpdateCmd.Flags().String("user", "admin", "username for basic auth")

}

// UpdateCmd is a cobra.Command to update commands on the botio's commands server.
var UpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Updates an existing command (or adds it if not exists) with a response on the botio's server",
	Example: "botio update --command start --response Hi --url localhost:9090 --user myuser --password mypassword",
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

		// PUT command
		c, err := client.Put(url, user, password, command, response)
		if err != nil {
			log.Fatalf("%v", err)
		}

		printCommands(c)
	},
}
