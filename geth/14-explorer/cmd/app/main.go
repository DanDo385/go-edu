package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/14-explorer/internal/explorer/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
