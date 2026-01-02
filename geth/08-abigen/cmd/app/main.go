package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/08-abigen/internal/abigen/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
