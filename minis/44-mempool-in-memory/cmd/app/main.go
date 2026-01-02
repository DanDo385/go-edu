package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/44-mempool-in-memory/internal/mempoolinmemory/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
