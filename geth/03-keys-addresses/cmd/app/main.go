package main

import (
	"fmt"
	"os"

	"github.com/example/go-10x-minis/geth/03-keys-addresses/internal/keysaddresses"
)

/*
Keys and Addresses Demo

This application demonstrates Ethereum key generation, address derivation,
and the relationship between private keys, public keys, and addresses.

Usage:

	go run ./cmd/app/main.go [PRIVATE_KEY_HEX]

Examples:

	# Generate new key pair
	go run ./cmd/app/main.go

	# Derive address from existing private key
	go run ./cmd/app/main.go 0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
*/
func main() {
	fmt.Println("=== Keys and Addresses ===")
	fmt.Println()

	var privateKeyHex string
	if len(os.Args) > 1 {
		privateKeyHex = os.Args[1]
		fmt.Printf("Using provided private key: %s...\n", privateKeyHex[:10])
	} else {
		fmt.Println("Generating new key pair...")
	}

	// BREAKPOINT: Step into Run() to see key generation and derivation
	result, err := keysaddresses.Run(keysaddresses.Config{
		PrivateKeyHex: privateKeyHex,
	})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("=== Key Information ===")
	fmt.Printf("Address:     %s\n", result.Address.Hex())
	fmt.Printf("Public Key:  %x\n", result.PublicKey)
	fmt.Println()
	fmt.Println("Key derivation: Private Key -> Public Key -> Keccak256 -> Last 20 bytes = Address")
}
