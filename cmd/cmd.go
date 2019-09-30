// Package cmd exports a function to create easily a CLI based on cobra.
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/danielkvist/botio/models"

	"github.com/spf13/cobra"
)

// Root creates a root *cobra.Command and then adds to it
// the received *cobra.Commands, then it executes the root command
// returning an error if any.
func Root(commands ...*cobra.Command) error {
	examples := []string{
		"botio server bolt --database ./data/commands.db --collection commands --http :9090 --key mysupersecretkey",
		"botio bot --platform telegram --token <telegram-token> --url :9090 --key mysupersecretkey",
		"botio client print --command start --url :9090 --key mysupersecretkey",
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

func checkURL(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("server URL cannot be an empty string")
	}

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	if !strings.HasSuffix(url, "/api/commands") {
		url = url + "/api/commands"
	}

	return url, nil
}

func printCommands(commands ...*models.Command) {
	const format = "%s\t\t%s\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "command", "response")
	fmt.Fprintf(tw, format, "-------", "--------")

	for _, c := range commands {
		fmt.Fprintf(tw, format, c.Cmd, c.Response)
	}

	tw.Flush()
}
