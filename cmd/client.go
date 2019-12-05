package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/proto"
	"github.com/pkg/errors"

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
		Short: "Client provides subcommands to manage your commands.",
	}

	for _, cmd := range commands {
		clientCmd.AddCommand(cmd)
	}

	return clientCmd
}

func add() *cobra.Command {
	var addr string
	var command string
	var response string
	var serverName string
	var sslca string
	var sslcrt string
	var sslkey string
	var token string

	add := &cobra.Command{
		Use:     "add",
		Short:   "Adds a new command.",
		Example: "botio client add --command start --response Hello",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(addr, token, serverName, sslcrt, sslkey, sslca)
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
				return errors.Wrapf(err, "while adding command %q with response %q", command, response)
			}

			log.Printf("command %q with response %q added successfully!\n", command, response)
			return nil
		},
		SilenceUsage: true,
	}

	add.Flags().StringVar(&addr, "addr", ":9091", "botio's gRPC server address")
	add.Flags().StringVar(&command, "command", "", "command to add")
	add.Flags().StringVar(&response, "response", "", "command's response")
	add.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	add.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	add.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	add.Flags().StringVar(&token, "token", "", "authentication token")

	return add
}

func print() *cobra.Command {
	var addr string
	var command string
	var serverName string
	var sslca string
	var sslcrt string
	var sslkey string
	var token string

	print := &cobra.Command{
		Use:     "print",
		Short:   "Prints the requested command.",
		Example: "botio client print --command start",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(addr, token, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
			}

			botCommand, err := c.GetCommand(context.TODO(), &proto.Command{
				Command: command,
			})
			if err != nil {
				return errors.Wrapf(err, "while getting command %q", command)
			}

			printCommand(botCommand)
			return nil
		},
		SilenceUsage: true,
	}

	print.Flags().StringVar(&addr, "addr", ":9091", "botio's gRPC server address")
	print.Flags().StringVar(&command, "command", "", "command to print")
	print.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	print.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	print.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	print.Flags().StringVar(&token, "token", "", "uthentication token")

	return print
}

func list() *cobra.Command {
	var addr string
	var serverName string
	var sslca string
	var sslcrt string
	var sslkey string
	var token string

	list := &cobra.Command{
		Use:     "list",
		Short:   "List all the commands.",
		Example: "botio client list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(addr, token, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
			}

			botCommands, err := c.ListCommands(context.TODO(), &empty.Empty{})
			if err != nil {
				return errors.Wrap(err, "while listing commands")
			}

			for _, bc := range botCommands.GetCommands() {
				printCommand(bc)
			}

			return nil
		},
		SilenceUsage: true,
	}

	list.Flags().StringVar(&addr, "addr", ":9091", "botio's gRPC server address")
	list.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	list.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	list.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	list.Flags().StringVar(&token, "token", "", "authentication token")

	return list
}

func update() *cobra.Command {
	var addr string
	var command string
	var response string
	var serverName string
	var sslca string
	var sslcrt string
	var sslkey string
	var token string

	update := &cobra.Command{
		Use:     "update",
		Short:   "Updates the requested command or adds it if don't exists.",
		Example: "botio client update --command start --response Hi",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(addr, token, serverName, sslcrt, sslkey, sslca)
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
				return errors.Wrapf(err, "while updating command %q with response %q", command, response)
			}

			log.Printf("command %q updated with response %q successfully!", command, response)
			return nil
		},
		SilenceUsage: true,
	}

	update.Flags().StringVar(&addr, "addr", ":9091", "botio's gRPC server address")
	update.Flags().StringVar(&command, "command", "", "command to update")
	update.Flags().StringVar(&response, "response", "", "command's new response")
	update.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	update.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	update.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	update.Flags().StringVar(&token, "token", "", "authentication token")

	return update
}

func delete() *cobra.Command {
	var addr string
	var command string
	var serverName string
	var sslca string
	var sslcrt string
	var sslkey string
	var token string

	delete := &cobra.Command{
		Use:     "delete",
		Short:   "Deletes the requested command",
		Example: "botio client delete --command start",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(addr, token, serverName, sslcrt, sslkey, sslca)
			if err != nil {
				return err
			}

			if _, err := c.DeleteCommand(context.TODO(), &proto.Command{
				Command: command,
			}); err != nil {
				return errors.Wrapf(err, "while deleting command %q", command)
			}

			log.Printf("command %q deleted successfully!", command)
			return nil
		},
		SilenceUsage: true,
	}

	delete.Flags().StringVar(&addr, "addr", ":9091", "botio's gRPC server address")
	delete.Flags().StringVar(&command, "command", "", "command to delete")
	delete.Flags().StringVar(&sslca, "sslca", "", "ssl client certification file")
	delete.Flags().StringVar(&sslcrt, "sslcrt", "", "ssl certification file")
	delete.Flags().StringVar(&sslkey, "sslkey", "", "ssl certification key file")
	delete.Flags().StringVar(&token, "token", "t", "authentication token")

	return delete
}

func getClient(url, token, server, crt, key, ca string) (client.Client, error) {
	var c client.Client
	var u string
	var err error

	u, err = checkURL(url, false, false)
	if err != nil {
		return nil, errors.Wrapf(err, "while parsing URL %q", url)
	}

	if crt == "" || key == "" || ca == "" {
		c, err = insecureClient(u, token)
	} else {
		c, err = securedClient(u, token, server, crt, key, ca)
	}

	if err != nil {
		return nil, errors.Wrap(err, "while creating gRPC client")
	}

	return c, nil
}

func insecureClient(url, token string) (client.Client, error) {
	return client.New(url, token, client.WithInsecureConn(url))
}

func securedClient(url, token, server, crt, key, ca string) (client.Client, error) {
	return client.New(url, token, client.WithTLSSecureConn(url, server, crt, key, ca))
}

// FIXME:
func printCommand(cmd *proto.BotCommand) {
	fmt.Printf("%q: %q\n", cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse())
}
