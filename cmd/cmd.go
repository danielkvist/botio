package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/danielkvist/botio/models"
)

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
