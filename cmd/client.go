package cmd

import (
	"fmt"
	"log"

	"github.com/danielkvist/botio/client"

	"github.com/spf13/cobra"
)

// Client returns a *cobra.Command with multiple subcommands.
func Client() *cobra.Command {
	return clientCmd(add(), print(), list(), update(), delete())
}

func clientCmd(commands ...*cobra.Command) *cobra.Command {
	clientCmd := &cobra.Command{
		Use:                   "client",
		Short:                 "Client contains some subcommands to manage your bot's commands",
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
	}

	for _, cmd := range commands {
		clientCmd.AddCommand(cmd)
	}

	return clientCmd
}

func add() *cobra.Command {
	var command string
	var token string
	var response string
	var url string

	add := &cobra.Command{
		Use:     "add",
		Short:   "Adds a new command",
		Example: "botio client add --command start --response Hello --url :9090 --token <jwt-token>",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			c, err := client.Post(u, token, command, response)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(c)
		},
	}

	add.Flags().StringVarP(&command, "command", "c", "", "command to add")
	add.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	add.Flags().StringVarP(&response, "response", "r", "", "command's response")
	add.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return add
}

func print() *cobra.Command {
	var command string
	var token string
	var url string

	print := &cobra.Command{
		Use:     "print",
		Short:   "Prints the specified command and his response",
		Example: "botio client print --command start --url :9090 --token <jwt-token>",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			c, err := client.Get(u+"/"+command, token)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(c)
		},
	}

	print.Flags().StringVarP(&command, "command", "c", "", "command to print")
	print.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	print.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return print
}

func list() *cobra.Command {
	var token string
	var url string

	list := &cobra.Command{
		Use:     "list",
		Short:   "Prints a list with all the commands",
		Example: "botio client list --url :9090 --token <jwt-token>",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			commands, err := client.GetAll(u, token)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(commands...)
		},
	}

	list.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	list.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return list
}

func update() *cobra.Command {
	var command string
	var token string
	var response string
	var url string

	update := &cobra.Command{
		Use:     "update",
		Short:   "Updates an existing command (or adds it if not exists)",
		Example: "botio client update --command start --response Hi --url :9090 --token <jwt-token>",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			c, err := client.Put(u, token, command, response)
			if err != nil {
				log.Fatalf("%v", err)
			}

			printCommands(c)
		},
	}

	update.Flags().StringVarP(&command, "command", "c", "", "command to update")
	update.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	update.Flags().StringVarP(&response, "response", "r", "", "command's new response")
	update.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return update
}

func delete() *cobra.Command {
	var command string
	var token string
	var url string

	delete := &cobra.Command{
		Use:     "delete",
		Short:   "Deletes the specified command",
		Example: "botio client delete --command start --url :9090 --token <jwt-authentication>",
		Run: func(cmd *cobra.Command, args []string) {
			u, err := checkURL(url)
			if err != nil {
				log.Fatalf("%v", err)
			}

			if err := client.Delete(u+"/"+command, token); err != nil {
				log.Fatalf("%v", err)
			}

			fmt.Printf("command %q deleted successfully\n", command)
		},
	}

	delete.Flags().StringVarP(&command, "command", "c", "", "command to delete")
	delete.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	delete.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return delete
}
