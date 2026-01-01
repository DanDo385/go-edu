package main

import (
	"fmt"

	"geth/03-keys-addresses/internal/keysaddresses"
)

/*
Debug Harness for Keys and Addresses Module
*/
func main() {
	fmt.Println("=== Keys and Addresses Debug Harness ===")
	fmt.Println()

	// BREAKPOINT: Step into Run() with no private key (generate new)
	fmt.Println("Test 1: Generate new key pair...")
	result1, err := keysaddresses.Run(keysaddresses.Config{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Printf("Generated Address: %s\n", result1.Address.Hex())
	fmt.Println()

	// BREAKPOINT: Step into Run() with known private key
	fmt.Println("Test 2: Derive from known private key...")
	// This is a well-known test private key - NEVER use in production
	testKey := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	result2, err := keysaddresses.Run(keysaddresses.Config{
		PrivateKeyHex: testKey,
	})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Printf("Derived Address: %s\n", result2.Address.Hex())
	fmt.Println()

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. Private keys are 256-bit random numbers")
	fmt.Println("2. Public keys are derived via secp256k1 elliptic curve")
	fmt.Println("3. Addresses are last 20 bytes of Keccak256(public key)")
	fmt.Println()
	fmt.Println("Next: Proceed to geth/04-accounts-balances")
}
