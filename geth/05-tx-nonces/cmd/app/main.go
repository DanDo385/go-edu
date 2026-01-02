package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/05-tx-nonces/internal/txnonces/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
