package cmd

import (
	"fmt"
	"log"

	"github.com/danielkvist/botio/client"
	"github.com/spf13/cobra"
)

// Delete returns a *cobra.Command.
func Delete() *cobra.Command {
	var command string
	var key string
	var url string

	delete := &cobra.Command{
		Use:     "delete",
		Short:   "Deletes the specified command",
		Example: "botio delete --command start --url :9090 --key mysupersecretkey",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			if err := client.Delete(u+"/"+command, key); err != nil {
				log.Fatalf("%v", err)
			}

			fmt.Printf("command %q deleted successfully\n", command)
		},
	}

	delete.Flags().StringVarP(&command, "command", "c", "", "command to delete")
	delete.Flags().StringVarP(&key, "key", "k", "", "authentication key")
	delete.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return delete
}
