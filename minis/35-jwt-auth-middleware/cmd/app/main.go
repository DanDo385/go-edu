package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/35-jwt-auth-middleware/internal/jwtauthmiddleware/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
