package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/25-atomic-counters-vs-mutex/internal/atomiccountersvsmutex/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
