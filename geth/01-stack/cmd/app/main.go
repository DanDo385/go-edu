package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/01-stack/internal/stack/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
