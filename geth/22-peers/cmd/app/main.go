package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/22-peers/internal/peers/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
