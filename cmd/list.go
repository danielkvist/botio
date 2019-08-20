package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/danielkvist/botio/client"
	"github.com/danielkvist/botio/models"

	"github.com/spf13/cobra"
)

func init() {
	ListCmd.Flags().String("password", "toor", "password for basic auth")
	ListCmd.Flags().String("url", "", "URL where the Botio's server is listening")
	ListCmd.Flags().String("user", "admin", "username for basic auth")
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List prints a list with all the Botio's commands",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		password, _ := cmd.Flags().GetString("password")
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")

		// Check URL
		if url == "" {
			log.Fatal("server URL cannot be an empty string")
		}

		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}

		url = url + "/api/commands"

		// GET commands
		commands, err := client.GetAll(url, user, password)
		if err != nil {
			log.Fatalf("while getting all the commands from %q: %v", url, err)
		}

		printCommands(commands)
	},
}

func printCommands(commands []*models.Command) {
	const format = "%q\t\t%s\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Command", "Response")
	fmt.Fprintf(tw, format, "-------", "--------")

	for _, c := range commands {
		fmt.Fprintf(tw, format, c.Cmd, c.Response)
	}

	tw.Flush()
}
