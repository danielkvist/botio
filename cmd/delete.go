package cmd

import (
	"fmt"
	"log"

	"github.com/danielkvist/botio/client"
	"github.com/spf13/cobra"
)

func init() {
	DeleteCmd.Flags().String("command", "", "command to delete")
	DeleteCmd.Flags().String("key", "", "authentication key for JWT")
	DeleteCmd.Flags().String("url", "", "botio's server URL")
}

// DeleteCmd is a cobra.Command to delete commands from the botio's commands server.
var DeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Deletes the specified botio's command from the botio's server",
	Example: "botio delete --command start --url :9090 --key mysupersecretkey",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		command := checkFlag(cmd, "command", false)
		key := checkFlag(cmd, "key", false)
		url := checkFlag(cmd, "url", false)

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
		if err := client.Delete(url+"/"+command, key); err != nil {
			log.Fatalf("%v", err)
		}

		fmt.Printf("command %q deleted successfully\n", command)
	},
}
