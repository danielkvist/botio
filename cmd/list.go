package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	ListCmd.Flags().String("password", "toor", "password for basic auth")
	ListCmd.Flags().String("url", "", "URL where the botio's server is listening")
	ListCmd.Flags().String("user", "admin", "username for basic auth")
}

// ListCmd is a cobra.Command to print all the commands available on the
// botio's commands server.
var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Prints a list with all the botio's commands",
	Example: "botio list --url localhost:9090 --user myuser --password mypassword",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// GET commands
		commands, err := client.GetAll(url, user, password)
		if err != nil {
			log.Fatalf("%v", err)
		}

		printCommands(commands...)
	},
}
