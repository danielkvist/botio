package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

// List returns a *cobra.Command.
func List() *cobra.Command {
	var key string
	var url string

	list := &cobra.Command{
		Use:     "list",
		Short:   "Prints a list with all the commands",
		Example: "botio list --url :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			key := checkFlag("key", key, false)
			url := checkFlag("url", url, false)

			url, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			commands, err := client.GetAll(url, key)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(commands...)
		},
	}

	list.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	list.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return list
}
