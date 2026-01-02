package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/11-storage/internal/storage/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
