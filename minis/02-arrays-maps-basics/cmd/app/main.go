package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/02-arrays-maps-basics/internal/arraysmapsbasics/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
