package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/03-csv-stats/internal/csvstats/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
