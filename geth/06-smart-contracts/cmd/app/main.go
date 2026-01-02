package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/06-smart-contracts/internal/smartcontracts/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
