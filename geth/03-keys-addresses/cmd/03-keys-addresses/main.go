package main

import (
	"flag"
	"fmt"
	"log"

	"geth-edu/03-keys-addresses/exercise"
)

func main() {
	outDir := flag.String("out", "./keystore-demo", "directory to store the encrypted key")
	pass := flag.String("pass", "changeit", "keystore passphrase (demo only)")
	flag.Parse()

	result, err := exercise.Run(exercise.Config{
		OutputDir:  *outDir,
		Passphrase: *pass,
	})
	if err != nil {
		log.Fatalf("key generation failed: %v", err)
	}

	fmt.Println("🔐 Generated Ethereum keypair")
	fmt.Printf("  Address:         %s\n", result.Address.Hex())
	fmt.Printf("  Private key hex: 0x%s\n", result.PrivateKeyHex)
	fmt.Printf("  Keystore file:   %s\n", result.KeystorePath)
	fmt.Println("\nReminder: NEVER commit private keys or weak demo passphrases in real projects.")
}
