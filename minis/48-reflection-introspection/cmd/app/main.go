package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/48-reflection-introspection/internal/reflectionintrospection/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
