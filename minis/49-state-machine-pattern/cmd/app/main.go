package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/49-state-machine-pattern/internal/statemachinepattern/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
