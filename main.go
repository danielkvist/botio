package main

import (
	"log"

	"github.com/danielkvist/botio/cmd"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:          "botio",
		Short:        "",
		Long:         "",
		SilenceUsage: true,
	}

	root.AddCommand(cmd.ServerCmd)

	if err := root.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
