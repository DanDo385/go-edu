package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/30-build-tags-conditional-compilation/internal/buildtagsconditionalcompilation/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
