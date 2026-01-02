package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/05-cli-todo-files/internal/clitodofiles/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
