package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/06-eip1559/internal/eip1559/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
