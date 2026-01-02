package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/46-generics-map-reduce/internal/genericsmapreduce/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
