package main

import (
	"fmt"

	"github.com/example/go-10x-minis/geth/03-keys-addresses/internal/keysaddresses"
)

func main() {
	// Deterministic-ish harness: fixed output dir and passphrase.
	// (Key generation is random by design.)
	//
	// BREAKPOINT: change inputs
	res, err := keysaddresses.Run(keysaddresses.Config{OutputDir: "./keystore-demo", Passphrase: "changeit"})
	if err != nil {
		panic(err)
	}

	fmt.Println("Address:", res.Address.Hex())
	fmt.Println("KeystorePath:", res.KeystorePath)
}
