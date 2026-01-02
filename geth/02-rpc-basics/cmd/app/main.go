package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/02-rpc-basics/internal/rpcbasics/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
