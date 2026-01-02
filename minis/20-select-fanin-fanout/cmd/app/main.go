package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/20-select-fanin-fanout/internal/selectfaninfanout/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
