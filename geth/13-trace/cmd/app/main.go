package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/13-trace/internal/trace/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
