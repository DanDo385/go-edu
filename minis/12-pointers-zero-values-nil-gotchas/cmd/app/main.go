package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/12-pointers-zero-values-nil-gotchas/internal/pointerszerovaluesnilgotchas/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
