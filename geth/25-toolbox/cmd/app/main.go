package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/25-toolbox/internal/toolbox/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
