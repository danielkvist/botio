// Package cmd exports a function to create easily a CLI based on cobra.
package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Root creates a root *cobra.Command and then adds to it
// the received *cobra.Commands, then it executes the root command
// returning an error if any.
func Root(commands ...*cobra.Command) error {
	examples := []string{
		"botio server bolt --database ./data/commands.db --collection commands",
		"botio bot --platform telegram --token <telegram-token> --addr :9091",
		"botio client add --command start --response Hi",
	}

	root := &cobra.Command{
		Use:   "botio",
		Short: "Botio is a CLI to create and manage easily chatbots for different platforms with the possibility of using differents databases.",
		Long: `Botio is a CLI to create and manage easily chatbots for different platforms such as Telegram or Discord.
It also let's you use different databases to manage their available commands wuch as BoltDB or PostgreSQL.
		
Botio is a project in development so use it with caution!`,
		Example: strings.Join(examples, "\n"),
	}

	for _, cmd := range commands {
		root.AddCommand(cmd)
	}

	return root.Execute()
}

func checkURL(url string, prefix bool, suffix bool) (string, error) {
	if url == "" {
		return "", errors.New("server URL cannot be an empty string")
	}

	if !strings.HasPrefix(url, "https://") && prefix {
		url = "https://" + url
	}

	if !strings.HasSuffix(url, "/api/commands") && suffix {
		url = url + "/api/commands"
	}

	return url, nil
}
