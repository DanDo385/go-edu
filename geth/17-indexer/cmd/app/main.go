package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/17-indexer/internal/indexer/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
