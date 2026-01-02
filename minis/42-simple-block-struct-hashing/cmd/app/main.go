package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/42-simple-block-struct-hashing/internal/simpleblockstructhashing/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
