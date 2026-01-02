package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/15-receipts/internal/receipts/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
