package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/16-context-cancellation-timeouts/internal/contextcancellationtimeouts/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
