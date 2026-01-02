package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/13-interfaces-duck-typing/internal/interfacesducktyping/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
