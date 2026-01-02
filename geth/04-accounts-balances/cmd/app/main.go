package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/04-accounts-balances/internal/accountsbalances/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
