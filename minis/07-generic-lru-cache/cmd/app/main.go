package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/07-generic-lru-cache/internal/genericlrucache/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
