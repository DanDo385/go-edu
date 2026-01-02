package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/04-jsonl-log-filter/internal/jsonllogfilter/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
