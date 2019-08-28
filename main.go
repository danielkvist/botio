package main

import (
	"log"

	"github.com/danielkvist/botio/cmd"
)

func main() {
	if err := cmd.Root(
		cmd.Add(),
		cmd.Delete(),
		cmd.Discord(),
		cmd.List(),
		cmd.Print(),
		cmd.Server(),
		cmd.Telegram(),
		cmd.Update(),
	); err != nil {
		log.Fatalf("%v", err)
	}
}
