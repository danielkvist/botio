package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	PrintCmd.Flags().String("command", "", "command to search for")
	PrintCmd.Flags().String("password", "toor", "password for basic auth")
	PrintCmd.Flags().String("url", "", "URL where the botio's server is listening")
	PrintCmd.Flags().String("user", "admin", "username for basic auth")
}

// PrintCmd is a cobra.Command to print a specified command from the botio's commands server.
var PrintCmd = &cobra.Command{
	Use:     "print",
	Short:   "Prints the specified botio's command with his response",
	Example: "botio print --command start --url localhost:9090 --user myuser --password mypassword",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		command, _ := cmd.Flags().GetString("command")
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// GET command
		c, err := client.Get(url+"/"+command, user, password)
		if err != nil {
			log.Fatalf("while getting command from %q: %v", url, err)
		}

		printCommands(c)
	},
}
