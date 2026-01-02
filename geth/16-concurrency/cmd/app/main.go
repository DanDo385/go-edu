package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/16-concurrency/internal/concurrency/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
