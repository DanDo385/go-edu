package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/12-proofs/internal/proofs/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
