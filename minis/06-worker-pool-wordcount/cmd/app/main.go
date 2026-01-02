package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/06-worker-pool-wordcount/internal/workerpoolwordcount/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
