package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/43-proof-of-work-demo/internal/proofofworkdemo/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
