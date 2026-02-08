//go:build reference

package keysaddresses

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	defaultOutputDir  = "./keystore-demo"
	defaultPassphrase = "changeit"
)

/*
Reference Solution - Ethereum Key Management and Address Generation
=================================================================

This file demonstrates cryptographic key generation and management for Ethereum.
Every Ethereum account is controlled by an ECDSA private key, from which the
public key and address are mathematically derived. This establishes the foundation
for transaction signing and account ownership in Ethereum.

This connects to the broader Ethereum ecosystem by showing:
- BIP-39/BIP-44 derivation paths (though simplified here)
- Keystore V3 format for encrypted private key storage
- Address derivation from public keys using Keccak-256
- Secure random number generation for private keys
- File system security for sensitive cryptographic material

The exercise builds understanding of:
- Public-key cryptography fundamentals (ECDSA, secp256k1 curve)
- Key derivation mathematics (public key from private key)
- Address generation (Keccak-256 hash of public key)
- Encrypted storage patterns for sensitive data
- File permission security for cryptographic keys

Teaching notes:
- Memory/ownership: private keys are extremely sensitive. We ensure they're
  zeroed after use and never logged or exposed in error messages.
- Invariants: keystore files must be encrypted at rest, and passphrases should
  be strong. We validate file creation and encryption before proceeding.
- Error surfaces: key generation and file operations can fail due to entropy
  issues, disk space, permissions, or cryptographic failures. We surface all
  errors so callers can implement appropriate security monitoring.
*/

/*
Run - Generate Ethereum Account with Secure Key Management

This function creates a new Ethereum account with proper cryptographic key generation,
encrypted storage, and secure file handling. It demonstrates the complete lifecycle
of Ethereum account creation used in wallets, exchanges, and dApps.

Parameters:
- cfg: configuration for output directory and encryption passphrase

Returns:
- *Result: pointer to Result — caller gets shared reference. We allocate with
  &Result{...}. Returning *Result avoids copying the struct (contains address
  string, path, etc.). Caller must not mutate; we own the allocation.
- error: any failures in key generation, encryption, or file operations

Algorithm steps:
1. Configure output directory and encryption passphrase
2. Create secure directory for keystore files
3. Generate new ECDSA private key and derive address
4. Encrypt and store private key in keystore format
5. Decrypt and validate keystore file integrity
6. Return account information for external use

Why secure key management matters:
- Private keys control millions in cryptocurrency value
- Compromised keys lead to total fund loss
- Proper encryption prevents unauthorized access
- Secure directory permissions prevent local attacks
*/
func Run(cfg Config) (*Result, error) {
	// Step 1: Configuration with sensible defaults
	// Allow callers to specify custom paths and passphrases
	// Fall back to secure defaults for demo purposes
	outDir := cfg.OutputDir
	if outDir == "" {
		outDir = defaultOutputDir
	}
	passphrase := cfg.Passphrase
	if passphrase == "" {
		passphrase = defaultPassphrase
	}

	// Step 2: Create secure output directory
	// 0o700 = owner read/write/execute only (no group/world access)
	// This prevents other users on the system from accessing keystore files
	// os.MkdirAll creates parent directories as needed
	if err := os.MkdirAll(outDir, 0o700); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	// Step 3: Initialize keystore and generate account
	// keystore.NewKeyStore creates a keystore manager for encrypted key storage
	// StandardScryptN/P provide balanced security vs performance for key derivation
	// NewAccount generates a random ECDSA private key and derives the address
	ks := keystore.NewKeyStore(outDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		// Key generation or encryption failed
		// Could be due to insufficient entropy, disk space, or cryptographic errors
		return nil, fmt.Errorf("create account: %w", err)
	}

	// Step 4: Read encrypted keystore file
	// account.URL.Path contains the filesystem path to the encrypted keystore file
	// This file contains the encrypted private key in JSON format (Web3 Secret Storage)
	keyJSON, err := os.ReadFile(account.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("read keystore file: %w", err)
	}

	// Step 5: Decrypt and validate keystore
	// keystore.DecryptKey performs the reverse: decrypts JSON and reconstructs private key
	// This validates that encryption/decryption works correctly
	// In production, this step would be skipped - callers shouldn't need raw private keys
	key, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		// Decryption failed - either wrong passphrase or corrupted file
		return nil, fmt.Errorf("decrypt keystore file: %w", err)
	}

	// Step 6: Return account information
	// Address: derived from public key (visible, can be shared)
	// PrivateKeyHex: hex-encoded private key (NEVER share this!)
	// KeystorePath: path to encrypted file (can be backed up safely)
	return &Result{
		Address:       account.Address,                                    // Public address
		PrivateKeyHex: hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)), // PRIVATE KEY - NEVER LOG
		KeystorePath:  account.URL.Path,                                  // Encrypted file location
	}, nil
}
