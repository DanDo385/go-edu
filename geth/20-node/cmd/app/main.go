package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/20-node/internal/node/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
