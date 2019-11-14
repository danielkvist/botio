package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/proto"

	"github.com/golang/protobuf/ptypes/empty"
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
	var sslcrt string
	var sslkey string
	var sslca string
	var serverName string
	var command string
	var response string
	var token string
	var url string

	add := &cobra.Command{
		Use:     "add",
		Short:   "Adds a new command",
		Example: "botio client add --command start --response Hello --url :9090 --token <jwt-token>",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(url, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
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

	add.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	add.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	add.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	add.Flags().StringVar(&command, "command", "", "command to add")
	add.Flags().StringVar(&response, "response", "", "command's response")
	add.Flags().StringVar(&token, "token", "", "jwt authentication token")
	add.Flags().StringVar(&url, "url", "", "botio's server url")

	return add
}

func print() *cobra.Command {
	var sslcrt string
	var sslkey string
	var sslca string
	var serverName string
	var command string
	var token string
	var url string

	print := &cobra.Command{
		Use:     "print",
		Short:   "Prints the specified command and his response",
		Example: "botio client print --command start --url :9090 --token <jwt-token>",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(url, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
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

	print.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	print.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	print.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	print.Flags().StringVarP(&command, "command", "c", "", "command to print")
	print.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	print.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return print
}

func list() *cobra.Command {
	var sslcrt string
	var sslkey string
	var sslca string
	var serverName string
	var token string
	var url string

	list := &cobra.Command{
		Use:     "list",
		Short:   "Prints a list with all the commands",
		Example: "botio client list --url :9090 --token <jwt-token>",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(url, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
			}

			botCommands, err := c.ListCommands(context.TODO(), &empty.Empty{})
			for _, bc := range botCommands.GetCommands() {
				printCommand(bc)
			}

			return nil
		},
	}

	list.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	list.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	list.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	list.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	list.Flags().StringVarP(&url, "url", "u", "", "botio's server URL")

	return list
}

func update() *cobra.Command {
	var sslcrt string
	var sslkey string
	var sslca string
	var serverName string
	var command string
	var response string
	var token string
	var url string

	update := &cobra.Command{
		Use:     "update",
		Short:   "Updates an existing command (or adds it if not exists)",
		Example: "botio client update --command start --response Hi --url :9090 --token <jwt-token>",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(url, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
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

	update.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	update.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	update.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	update.Flags().StringVarP(&command, "command", "c", "", "command to update")
	update.Flags().StringVarP(&response, "response", "r", "", "command's new response")
	update.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	update.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return update
}

func delete() *cobra.Command {
	var sslcrt string
	var sslkey string
	var sslca string
	var serverName string
	var command string
	var token string
	var url string

	delete := &cobra.Command{
		Use:     "delete",
		Short:   "Deletes the specified command",
		Example: "botio client delete --command start --url :9090 --token <jwt-authentication>",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(url, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
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

	delete.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	delete.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	delete.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	delete.Flags().StringVarP(&command, "command", "c", "", "command to delete")
	delete.Flags().StringVarP(&token, "token", "t", "", "jwt authentication token")
	delete.Flags().StringVarP(&url, "url", "u", "", "botio's server url")

	return delete
}

func getClient(url, server, crt, key, ca string) (client.Client, error) {
	var c client.Client
	var u string
	var err error

	u, err = checkURL(url, false, false)
	if err != nil {
		return nil, fmt.Errorf("while parsing URL: %v", err)
	}

	if crt == "" || key == "" || ca == "" {
		c, err = insecureClient(u)
	} else {
		c, err = securedClient(u, server, crt, key, ca)
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}

func insecureClient(url string) (client.Client, error) {
	return client.New(url, client.WithInsecureConn(url))
}

func securedClient(url, server, crt, key, ca string) (client.Client, error) {
	return client.New(url, client.WithTLSSecureConn(url, server, crt, key, ca))
}

func printCommand(cmd *proto.BotCommand) {
	fmt.Printf("%q: %q\n", cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse())
}
