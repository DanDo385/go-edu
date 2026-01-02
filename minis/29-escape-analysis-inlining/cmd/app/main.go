package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/29-escape-analysis-inlining/internal/escapeanalysisinlining/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
