package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/41-signed-transactions-ed25519/internal/signedtransactionsed25519/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
