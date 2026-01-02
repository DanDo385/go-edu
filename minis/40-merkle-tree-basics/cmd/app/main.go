package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/40-merkle-tree-basics/internal/merkletreebasics/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
