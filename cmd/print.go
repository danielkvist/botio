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
			c := checkFlag("command", command, false)
			k := checkFlag("key", key, false)
			u := checkFlag("url", url, false)

			u, err := checkURL(u)
			if err != nil {
				log.Fatalf("%v", err)
			}

			command, err := client.Get(u+"/"+c, k)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(command)
		},
		Args: cobra.ExactArgs(3),
	}

	print.Flags().StringVarP(&command, "command", "c", "", "command to print")
	print.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	print.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return print
}
