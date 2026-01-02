package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/14-methods-value-vs-pointer-receivers/internal/methodsvaluevspointerreceivers/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
