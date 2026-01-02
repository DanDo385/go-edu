package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/19-devnets/internal/devnets/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
