package cmd

import (
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

// Add returns a *cobra.Command.
func Add() *cobra.Command {
	var command string
	var key string
	var response string
	var url string

	add := &cobra.Command{
		Use:     "add",
		Short:   "Adds a new command",
		Example: "botio add --command start --response Hello --url :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			c := checkFlag("command", command, false)
			k := checkFlag("key", key, false)
			r := checkFlag("response", response, false)
			u := checkFlag("url", url, false)

			u, err := checkURL(u)
			if err != nil {
				log.Fatalf("%v", err)
			}

			command, err := client.Post(u, k, c, r)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(command)
		},
		Args: cobra.ExactArgs(4),
	}

	add.Flags().StringVarP(&command, "command", "c", "", "command to add")
	add.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	add.Flags().StringVarP(&response, "response", "r", "", "command's response")
	add.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return add
}
