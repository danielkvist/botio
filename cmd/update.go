package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

// Update returns a *cobra.Command.
func Update() *cobra.Command {
	var command string
	var key string
	var response string
	var url string

	update := &cobra.Command{
		Use:     "update",
		Short:   "Updates an existing command (or adds it if not exists)",
		Example: "botio update --command start --response Hi --url :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			c, err := client.Put(u, key, command, response)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(c)
		},
	}

	update.Flags().StringVarP(&command, "command", "c", "", "command to update")
	update.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	update.Flags().StringVarP(&response, "response", "r", "", "command's new response")
	update.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return update
}
