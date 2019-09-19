package main

import (
	"log"

	"github.com/danielkvist/botio/cmd"
)

func main() {
	if err := cmd.Root(
		cmd.Bot(),
		cmd.Add(),
		cmd.Delete(),
		cmd.List(),
		cmd.Print(),
		cmd.Server(),
		cmd.Update(),
	); err != nil {
		log.Fatalf("%v", err)
	}
}
