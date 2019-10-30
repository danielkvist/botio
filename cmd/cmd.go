// Package cmd exports a function to create easily a CLI based on cobra.
package cmd

import (
	"fmt"
	"strings"

	"github.com/danielkvist/botio/proto"
	"github.com/spf13/cobra"
)

// Root creates a root *cobra.Command and then adds to it
// the received *cobra.Commands, then it executes the root command
// returning an error if any.
func Root(commands ...*cobra.Command) error {
	examples := []string{
		"botio server bolt --database ./data/commands.db --collection commands --http :9090 --key mysupersecretkey",
		"botio bot --platform telegram --token <telegram-token> --url :9090 --token <jwt-token>",
		"botio client print --command start --url :9090 --jwt <jwt-token>",
	}

	root := &cobra.Command{
		Use:          "botio",
		Short:        "Botio is a simple and opinionated CLI to create and manage easily bots for differents platforms.",
		Example:      strings.Join(examples, "\n"),
		SilenceUsage: true,
	}

	for _, cmd := range commands {
		root.AddCommand(cmd)
	}

	return root.Execute()
}

func checkURL(url string, prefix bool, suffix bool) (string, error) {
	if url == "" {
		return "", fmt.Errorf("server URL cannot be an empty string")
	}

	if !strings.HasPrefix(url, "https://") && prefix {
		url = "https://" + url
	}

	if !strings.HasSuffix(url, "/api/commands") && suffix {
		url = url + "/api/commands"
	}

	return url, nil
}

func printCommand(cmd *proto.BotCommand) {
	fmt.Printf("%q: %q", cmd.GetCmd().GetCommand(), cmd.GetResp().GetResponse())
}
