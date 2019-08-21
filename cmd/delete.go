package cmd

import (
	"fmt"
	"log"

	"github.com/danielkvist/botio/client"
	"github.com/spf13/cobra"
)

func init() {
	DeleteCmd.Flags().String("command", "", "command to delete")
	DeleteCmd.Flags().String("password", "toor", "password for basic auth")
	DeleteCmd.Flags().String("url", "", "URL where the botio's server is listening")
	DeleteCmd.Flags().String("user", "admin", "username for basic auth")
}

// DeleteCmd is a cobra.Command to delete commands from the botio's commands server.
var DeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Deletes the specified botio's command from the botio's server",
	Example: "botio delete --command start --url localhost:9090 --user myuser --password mypassword",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		command, _ := cmd.Flags().GetString("command")
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")

		// Check command
		if command == "" {
			log.Fatal("either command or response cannot be an empty string")
		}

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// DELETE command
		if err := client.Delete(url+"/"+command, user, password); err != nil {
			log.Fatalf("while deleting command %q from %q: %v", command, url, err)
		}

		fmt.Printf("command %q deleted successfully\n", command)
	},
}
