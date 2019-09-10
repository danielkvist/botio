package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

// Print returns a *cobra.Command.
func Print() *cobra.Command {
	var command string
	var key string
	var url string

	print := &cobra.Command{
		Use:     "print",
		Short:   "Prints the specified command and his response",
		Example: "botio print --command start --url :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			c, err := client.Get(u+"/"+command, key)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(c)
		},
	}

	print.Flags().StringVarP(&command, "command", "c", "", "command to print")
	print.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	print.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return print
}
