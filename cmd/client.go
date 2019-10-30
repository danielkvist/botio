package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

// Client returns a *cobra.Command with multiple subcommands.
func Client() *cobra.Command {
	return clientCmd(add(), print(), list(), update(), delete())
}

func clientCmd(commands ...*cobra.Command) *cobra.Command {
	clientCmd := &cobra.Command{
		Use:   "client",
		Short: "Client contains some subcommands to manage your bot's commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := checkURL(url, false, false)
			if err != nil {
				return fmt.Errorf("while parsing URL: %v", err)
			}

			c, err := client.New(u, grpc.WithInsecure())
			if err != nil {
				return fmt.Errorf("while creating a new client to add command %q: %v", command, err)
			}

			if _, err := c.AddCommand(context.TODO(), &proto.BotCommand{
				Cmd: &proto.Command{
					Command: command,
				},
				Resp: &proto.Response{
					Response: response,
				},
			}); err != nil {
				return fmt.Errorf("while adding command %q: %v", command, err)
			}

			log.Printf("command %q added!", command)
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := checkURL(url, false, false)
			if err != nil {
				return fmt.Errorf("while parsing URL: %v", err)
			}

			c, err := client.New(u, grpc.WithInsecure())
			if err != nil {
				return fmt.Errorf("while creating a new client to print command %q: %v", command, err)
			}

			botCommand, err := c.GetCommand(context.TODO(), &proto.Command{
				Command: command,
			})
			if err != nil {
				return fmt.Errorf("while getting command %q for printing: %v", command, err)
			}

			printCommand(botCommand)
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := checkURL(url, false, false)
			if err != nil {
				return fmt.Errorf("while parsing URL: %v", err)
			}

			c, err := client.New(u, grpc.WithInsecure())
			if err != nil {
				return fmt.Errorf("while creating a new client to print list of commands: %v", err)
			}

			botCommands, err := c.ListCommands(context.TODO(), &empty.Empty{})
			for _, bc := range botCommands.GetCommands() {
				printCommand(bc)
			}

			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := checkURL(url, false, false)
			if err != nil {
				return fmt.Errorf("while parsing URL: %v", err)
			}

			c, err := client.New(u, grpc.WithInsecure())
			if err != nil {
				return fmt.Errorf("while creating a new client to update command %q: %v", command, err)
			}

			if _, err := c.UpdateCommand(context.TODO(), &proto.BotCommand{
				Cmd: &proto.Command{
					Command: command,
				},
				Resp: &proto.Response{
					Response: response,
				},
			}); err != nil {
				return fmt.Errorf("while updating command %q: %v", command, err)
			}

			log.Printf("command %q updated!", command)
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := checkURL(url, false, false)
			if err != nil {
				return fmt.Errorf("while parsing URL: %v", err)
			}

			c, err := client.New(u, grpc.WithInsecure())
			if err != nil {
				return fmt.Errorf("while creating a new client to delete command %q: %v", command, err)
			}

			if _, err := c.DeleteCommand(context.TODO(), &proto.Command{
				Command: command,
			}); err != nil {
				return fmt.Errorf("while deleting command %q: %v", command, err)
			}

			log.Printf("command %q updated!", command)
			return nil
		},
	}

	delete.Flags().StringVarP(&command, "command", "c", "", "command to delete")
	delete.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	delete.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return delete
}
