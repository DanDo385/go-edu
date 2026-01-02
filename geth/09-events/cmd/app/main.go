package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/09-events/internal/events/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
