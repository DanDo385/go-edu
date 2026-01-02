package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/07-eth-call/internal/ethcall/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
