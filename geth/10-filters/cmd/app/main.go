package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/10-filters/internal/filters/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
