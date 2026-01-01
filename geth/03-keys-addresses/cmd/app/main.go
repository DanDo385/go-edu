package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/go-10x-minis/geth/03-keys-addresses/internal/keysaddresses"
)

func main() {
	// geth/03-keys-addresses
	//
	// This module does NOT need an RPC endpoint. It demonstrates:
	// - generating a keypair
	// - deriving an Ethereum address
	// - writing an encrypted keystore file
	//
	// Usage:
	//   go run ./geth/03-keys-addresses/cmd/app --out ./keystore-demo --pass changeit
	//
	// BREAKPOINT: parse flags
	outDir := flag.String("out", "./keystore-demo", "Directory to write keystore file")
	pass := flag.String("pass", "changeit", "Keystore passphrase (demo default is intentionally weak)")
	flag.Parse()

	res, err := keysaddresses.Run(keysaddresses.Config{OutputDir: *outDir, Passphrase: *pass})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("Address:", res.Address.Hex())
	fmt.Println("PrivateKey (hex):", res.PrivateKeyHex)
	fmt.Println("KeystorePath:", res.KeystorePath)
}
