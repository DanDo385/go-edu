package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/18-reorgs/internal/reorgs/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
