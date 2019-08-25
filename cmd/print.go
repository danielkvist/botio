package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

func init() {
	PrintCmd.Flags().String("command", "", "command to search for")
	PrintCmd.Flags().String("key", "", "authentication key for JWT")
	PrintCmd.Flags().String("url", "", "botio's server URL")
}

// PrintCmd is a cobra.Command to print a specified command from the botio's commands server.
var PrintCmd = &cobra.Command{
	Use:     "print",
	Short:   "Prints the specified botio's command with his response",
	Example: "botio print --command start --url :9090 --key mysupersecretkey",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		command := checkFlag(cmd, "command", false)
		key := checkFlag(cmd, "key", false)
		url := checkFlag(cmd, "url", false)

		// Check URL
		url, err := checkURL(url)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// GET command
		c, err := client.Get(url+"/"+command, key)
		if err != nil {
			log.Fatalf("%v", err)
		}

		printCommands(c)
	},
}
