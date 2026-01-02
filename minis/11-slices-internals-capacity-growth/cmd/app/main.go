package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/11-slices-internals-capacity-growth/internal/slicesinternalscapacitygrowth/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
