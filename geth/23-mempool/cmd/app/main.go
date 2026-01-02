package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/23-mempool/internal/mempool/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
