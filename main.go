package main

import (
	"log"

	"github.com/danielkvist/botio/cmd"
)

func main() {
	if err := cmd.Root(
		cmd.Bot(),
		cmd.Server(),
		cmd.Client(),
	); err != nil {
		log.Fatalf("%v", err)
	}
}
