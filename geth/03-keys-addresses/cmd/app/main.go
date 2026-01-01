package main

import (
	"fmt"
	"log"
	"os"

	"github.com/example/go-10x-minis/geth/03-keys-addresses/internal/keysaddresses"
)

/*
===================================================================
Ethereum Keys and Addresses - Cryptographic Identity
===================================================================

This program demonstrates how to generate Ethereum private keys,
derive public addresses, and securely store them in encrypted
keystore files using the go-ethereum keystore format.

USAGE:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go [demo]
    go run ./cmd/app/main.go [demo] [output_dir] [passphrase]

Arguments:
    [demo]       - Demo to run: all, basic, derive, keystore, security (default: "all")
    [output_dir] - Directory for keystore files (default: "./keystore-demo")
    [passphrase] - Passphrase for encryption (default: "changeit")

Examples:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go basic
    go run ./cmd/app/main.go derive
    go run ./cmd/app/main.go keystore
    go run ./cmd/app/main.go all ./my-keys mySecurePassword123

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter keysaddresses package functions
- Watch Variables panel to see key generation and derivation
- Inspect Address, PrivateKey, and KeystorePath structures

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the keysaddresses package from internal/keysaddresses
- exercise.go, solution.reference.go are in package keysaddresses
- The keysaddresses.Run() function handles key generation and storage

KEY CONCEPTS:
- secp256k1: The elliptic curve used by Ethereum (same as Bitcoin)
- Private Key: 256-bit random number that controls an account
- Public Key: Derived from private key using elliptic curve multiplication
- Address: Last 20 bytes of Keccak-256 hash of public key
- Keystore: Encrypted JSON file containing private key (scrypt + AES)

SECURITY NOTES:
- NEVER share or commit private keys
- NEVER log private keys in production
- ALWAYS use strong passphrases for keystores
- ALWAYS keep encrypted keystores backed up
- This demo uses weak defaults for educational purposes only
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the demo type (optional)
	// DEBUG: os.Args[2] is the output directory (optional)
	// DEBUG: os.Args[3] is the passphrase (optional)

	demo := "all"
	outputDir := ""
	passphrase := ""

	if len(os.Args) > 1 {
		demo = os.Args[1]
	}
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}
	if len(os.Args) > 3 {
		passphrase = os.Args[3]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'demo', 'outputDir', and 'passphrase' variables
	fmt.Printf("=== Ethereum Keys and Addresses ===\n")
	fmt.Printf("Demo: %s\n", demo)
	if outputDir != "" {
		fmt.Printf("Output Directory: %s\n", outputDir)
	} else {
		fmt.Printf("Output Directory: ./keystore-demo (default)\n")
	}
	if passphrase != "" {
		fmt.Printf("Passphrase: [provided]\n")
	} else {
		fmt.Printf("Passphrase: changeit (default)\n")
	}
	fmt.Println()

	// ============================================================================
	// STEP 2: Example 1 - Basic Key Generation
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through key generation
	// DEBUG: This demonstrates the core cryptographic operation
	// DEBUG: Watch how a private key is generated and address is derived

	if demo == "all" || demo == "basic" {
		fmt.Println("=== Example 1: Basic Key Generation and Address Derivation ===")
		fmt.Println("Generating a new Ethereum account...")
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: We're about to generate a cryptographically secure random private key
		cfg := keysaddresses.Config{
			OutputDir:  outputDir,
			Passphrase: passphrase,
		}

		// BREAKPOINT: Set breakpoint here before calling keysaddresses.Run
		// DEBUG: Step Into (F11) to see the implementation in exercise.go or solution.reference.go
		// DEBUG: You'll see crypto.GenerateKey() create a secp256k1 private key
		// DEBUG: Then crypto.PubkeyToAddress() derive the Ethereum address
		result, err := keysaddresses.Run(cfg)
		if err != nil {
			log.Fatalf("Failed to generate key: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful generation
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at Address (20 bytes), PrivateKeyHex (64 hex chars = 32 bytes)
		// DEBUG: Notice KeystorePath shows where the encrypted key was saved

		fmt.Println("✓ Account Generated Successfully!")
		fmt.Println()
		fmt.Println("Account Details:")
		fmt.Printf("  Address:      %s\n", result.Address.Hex())
		fmt.Printf("  Private Key:  %s\n", result.PrivateKeyHex)
		fmt.Printf("  Keystore:     %s\n\n", result.KeystorePath)

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This output shows the relationship between private key and address
		// DEBUG: The address is deterministically derived from the private key
		fmt.Println("SECURITY WARNING:")
		fmt.Println("  The private key above controls this account")
		fmt.Println("  Anyone with this key can spend all funds in this account")
		fmt.Println("  NEVER share your private key with anyone")
		fmt.Println("  NEVER commit private keys to version control")
		fmt.Println()
	}

	// ============================================================================
	// STEP 3: Example 2 - Understanding Address Derivation
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand cryptographic derivation
	// DEBUG: This demonstrates the mathematical relationship between keys and addresses

	if demo == "all" || demo == "derive" {
		fmt.Println("=== Example 2: Understanding Address Derivation ===")
		fmt.Println()

		fmt.Println("The Ethereum Address Derivation Process:")
		fmt.Println()
		fmt.Println("STEP 1: Generate Private Key")
		fmt.Println("  - Generate 256-bit (32-byte) random number")
		fmt.Println("  - Use cryptographically secure random number generator")
		fmt.Println("  - Example: 0x1a2b3c4d5e6f... (64 hex characters)")
		fmt.Println()

		fmt.Println("STEP 2: Derive Public Key")
		fmt.Println("  - Apply secp256k1 elliptic curve multiplication")
		fmt.Println("  - Public Key = Private Key * Generator Point (G)")
		fmt.Println("  - Result: 512-bit (64-byte) uncompressed public key")
		fmt.Println("  - This operation is one-way (cannot reverse)")
		fmt.Println()

		fmt.Println("STEP 3: Hash Public Key")
		fmt.Println("  - Apply Keccak-256 hash function to public key")
		fmt.Println("  - Result: 256-bit (32-byte) hash")
		fmt.Println("  - Keccak is different from SHA-3 (pre-standardization)")
		fmt.Println()

		fmt.Println("STEP 4: Extract Address")
		fmt.Println("  - Take the last 20 bytes of the hash")
		fmt.Println("  - Prepend with '0x' for hex representation")
		fmt.Println("  - Result: Ethereum address (e.g., 0x1234...5678)")
		fmt.Println()

		// Demonstrate with actual key generation
		// BREAKPOINT: Set breakpoint here
		// DEBUG: We'll generate a key and show each step of the derivation
		cfg := keysaddresses.Config{
			OutputDir:  outputDir,
			Passphrase: passphrase,
		}

		result, err := keysaddresses.Run(cfg)
		if err != nil {
			log.Fatalf("Failed to generate key: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect the result to see the final derived address
		fmt.Println("Example Derivation Result:")
		fmt.Printf("  Private Key (32 bytes): %s\n", result.PrivateKeyHex)
		fmt.Printf("  ↓ secp256k1 curve multiplication\n")
		fmt.Printf("  Public Key (64 bytes)\n")
		fmt.Printf("  ↓ Keccak-256 hash\n")
		fmt.Printf("  Hash (32 bytes)\n")
		fmt.Printf("  ↓ Take last 20 bytes\n")
		fmt.Printf("  Address (20 bytes):     %s\n", result.Address.Hex())
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This demonstrates the one-way nature of the derivation
		fmt.Println("Key Property:")
		fmt.Println("  ✓ Private Key → Address (easy, deterministic)")
		fmt.Println("  ✗ Address → Private Key (computationally infeasible)")
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 3 - Keystore Format and Encryption
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand keystore encryption
	// DEBUG: This shows how private keys are securely stored on disk

	if demo == "all" || demo == "keystore" {
		fmt.Println("=== Example 3: Keystore Format and Encryption ===")
		fmt.Println()

		fmt.Println("Understanding Ethereum Keystore Files:")
		fmt.Println()
		fmt.Println("1. Purpose:")
		fmt.Println("   - Safely store private keys on disk")
		fmt.Println("   - Encrypted with user's passphrase")
		fmt.Println("   - Standard format (Web3 Secret Storage)")
		fmt.Println()

		fmt.Println("2. Encryption Process:")
		fmt.Println("   STEP A: Passphrase → Key Derivation (scrypt)")
		fmt.Println("     - scrypt: memory-hard function (resistant to GPUs/ASICs)")
		fmt.Println("     - Parameters: N=262144, r=8, p=1 (standard)")
		fmt.Println("     - Result: 32-byte encryption key")
		fmt.Println()
		fmt.Println("   STEP B: Encrypt Private Key (AES-128-CTR)")
		fmt.Println("     - AES: Advanced Encryption Standard")
		fmt.Println("     - CTR mode: Counter mode (stream cipher)")
		fmt.Println("     - Result: Encrypted private key (ciphertext)")
		fmt.Println()
		fmt.Println("   STEP C: Generate MAC (Message Authentication Code)")
		fmt.Println("     - MAC: Keccak-256(encryption_key + ciphertext)")
		fmt.Println("     - Purpose: Detect tampering or wrong passphrase")
		fmt.Println()

		fmt.Println("3. File Format (JSON):")
		fmt.Println("   {")
		fmt.Println("     \"address\": \"...\",        // Ethereum address")
		fmt.Println("     \"crypto\": {")
		fmt.Println("       \"cipher\": \"aes-128-ctr\",")
		fmt.Println("       \"ciphertext\": \"...\",   // Encrypted private key")
		fmt.Println("       \"kdf\": \"scrypt\",        // Key derivation function")
		fmt.Println("       \"kdfparams\": {...},      // scrypt parameters")
		fmt.Println("       \"mac\": \"...\"            // Authentication code")
		fmt.Println("     }")
		fmt.Println("   }")
		fmt.Println()

		// Generate a keystore and show its location
		// BREAKPOINT: Set breakpoint here
		// DEBUG: We'll create a keystore and show where it's saved
		cfg := keysaddresses.Config{
			OutputDir:  outputDir,
			Passphrase: passphrase,
		}

		result, err := keysaddresses.Run(cfg)
		if err != nil {
			log.Fatalf("Failed to generate key: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect result.KeystorePath to see where the file was created
		fmt.Println("Keystore Created:")
		fmt.Printf("  Address:  %s\n", result.Address.Hex())
		fmt.Printf("  Location: %s\n", result.KeystorePath)
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This explains the security properties of the keystore
		fmt.Println("Security Properties:")
		fmt.Println("  ✓ Private key is encrypted (AES)")
		fmt.Println("  ✓ Passphrase is never stored (only derived key)")
		fmt.Println("  ✓ scrypt makes brute-force attacks expensive")
		fmt.Println("  ✓ MAC prevents undetected tampering")
		fmt.Println("  ✓ Industry-standard format (compatible with all wallets)")
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Example 4 - Security Best Practices
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand security implications
	// DEBUG: This demonstrates critical security concepts for key management

	if demo == "all" || demo == "security" {
		fmt.Println("=== Example 4: Security Best Practices ===")
		fmt.Println()

		fmt.Println("Critical Security Considerations:")
		fmt.Println()

		fmt.Println("1. Private Key Protection:")
		fmt.Println("   ✓ DO: Store keys in encrypted keystores")
		fmt.Println("   ✓ DO: Use hardware wallets for significant funds")
		fmt.Println("   ✓ DO: Keep backups in secure, offline locations")
		fmt.Println("   ✗ DON'T: Store keys in plain text")
		fmt.Println("   ✗ DON'T: Commit keys to version control")
		fmt.Println("   ✗ DON'T: Share keys via email, chat, or screenshots")
		fmt.Println()

		fmt.Println("2. Passphrase Selection:")
		fmt.Println("   ✓ DO: Use long, random passphrases (20+ characters)")
		fmt.Println("   ✓ DO: Use a password manager")
		fmt.Println("   ✓ DO: Consider using a BIP-39 mnemonic phrase")
		fmt.Println("   ✗ DON'T: Use common words or patterns")
		fmt.Println("   ✗ DON'T: Reuse passwords from other services")
		fmt.Println("   ✗ DON'T: Use 'changeit' (this demo's default!)")
		fmt.Println()

		fmt.Println("3. Key Generation Entropy:")
		fmt.Println("   ✓ DO: Use crypto/rand (cryptographically secure)")
		fmt.Println("   ✓ DO: Ensure sufficient system entropy")
		fmt.Println("   ✓ DO: Generate keys on secure, offline machines")
		fmt.Println("   ✗ DON'T: Use math/rand (not cryptographically secure)")
		fmt.Println("   ✗ DON'T: Use predictable seeds or timestamps")
		fmt.Println()

		fmt.Println("4. Operational Security:")
		fmt.Println("   ✓ DO: Test with small amounts first")
		fmt.Println("   ✓ DO: Use testnet for development")
		fmt.Println("   ✓ DO: Verify addresses before sending funds")
		fmt.Println("   ✗ DON'T: Run untrusted code with real keys")
		fmt.Println("   ✗ DON'T: Enter keys on untrusted computers")
		fmt.Println()

		// Demonstrate the importance of secure key generation
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This shows multiple key generations to emphasize uniqueness
		fmt.Println("Demonstration: Key Uniqueness")
		fmt.Println("Generating 3 different keys to show randomness:")
		fmt.Println()

		for i := 1; i <= 3; i++ {
			cfg := keysaddresses.Config{
				OutputDir:  outputDir,
				Passphrase: passphrase,
			}

			// BREAKPOINT: Set breakpoint here inside the loop
			// DEBUG: Each iteration generates a completely different key
			result, err := keysaddresses.Run(cfg)
			if err != nil {
				log.Fatalf("Failed to generate key %d: %v\n", i, err)
			}

			// BREAKPOINT: Set breakpoint here
			// DEBUG: Notice how each address is completely different
			fmt.Printf("Key #%d:\n", i)
			fmt.Printf("  Address: %s\n", result.Address.Hex())
			fmt.Printf("  Private: %s...%s\n",
				result.PrivateKeyHex[:8],
				result.PrivateKeyHex[len(result.PrivateKeyHex)-8:])
			fmt.Println()
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This emphasizes that each key is unique and unpredictable
		fmt.Println("Each key is completely random and unique!")
		fmt.Println("The probability of generating the same key twice is ~1/2^256")
		fmt.Println("That's more than the number of atoms in the observable universe!")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Understanding the Public-Key Cryptography Foundation
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand the math behind Ethereum
	// DEBUG: This explains the fundamental cryptographic concepts

	if demo == "all" {
		fmt.Println("=== Understanding Public-Key Cryptography ===")
		fmt.Println()

		fmt.Println("secp256k1 Elliptic Curve:")
		fmt.Println("  Equation: y² = x³ + 7 (over a finite field)")
		fmt.Println("  Used by: Bitcoin, Ethereum, and many other blockchains")
		fmt.Println("  Security: ~128-bit security level (comparable to AES-128)")
		fmt.Println()

		fmt.Println("Key Properties:")
		fmt.Println("  1. One-way Function:")
		fmt.Println("     - Easy: PrivateKey × G = PublicKey")
		fmt.Println("     - Hard: PublicKey ÷ G = PrivateKey (discrete log problem)")
		fmt.Println()
		fmt.Println("  2. Digital Signatures:")
		fmt.Println("     - Sign: Use private key to create signature")
		fmt.Println("     - Verify: Use public key to verify signature")
		fmt.Println("     - Property: Proves ownership without revealing private key")
		fmt.Println()
		fmt.Println("  3. Deterministic Derivation:")
		fmt.Println("     - Same private key always produces same address")
		fmt.Println("     - Enables address verification and recovery")
		fmt.Println()

		fmt.Println("Why This Matters for Ethereum:")
		fmt.Println("  - Transactions are signed with private keys")
		fmt.Println("  - Network verifies signatures using addresses (derived from public keys)")
		fmt.Println("  - No central authority needed to verify ownership")
		fmt.Println("  - Foundation of trustless, decentralized systems")
		fmt.Println()
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter keysaddresses package functions
	// 5. Watch Variables panel to see key generation
	// 6. Inspect result to see Address, PrivateKeyHex, and KeystorePath
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - secp256k1: Why this specific elliptic curve?
	// - Keccak-256: How is it different from SHA-256?
	// - Keystore Format: Why scrypt instead of PBKDF2?
	// - Address Derivation: Why only use last 20 bytes?
	// - Public-Key Cryptography: How do digital signatures work?
	//
	// STEPPING INTO KEYSADDRESSES PACKAGE:
	// - Set breakpoint at keysaddresses.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see: crypto.GenerateKey(), crypto.PubkeyToAddress()
	// - Watch keystore creation: keystore.NewKeyStore(), ImportECDSA()
	// - Observe the unlock/lock pattern for secure key handling
	//
	// COMMON ISSUES TO DEBUG:
	// - Permission errors: Check directory write permissions
	// - Keystore conflicts: Multiple keys with same timestamp
	// - Weak passphrase warnings: Use strong passphrases in production
	// - File not found: Ensure output directory exists
	//
	// SECURITY TESTING:
	// - Try different passphrases and verify encryption
	// - Examine keystore JSON file structure
	// - Verify that same private key produces same address
	// - Test keystore unlock/lock functionality
}
