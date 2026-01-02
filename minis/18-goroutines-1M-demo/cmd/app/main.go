package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/18-goroutines-1M-demo/internal/goroutines1mdemo/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
