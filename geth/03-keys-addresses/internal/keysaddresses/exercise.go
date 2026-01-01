//go:build !solution && !reference

package keysaddresses

import (
	"errors"
)

/*
Problem: Generate an Ethereum key, derive its address, and store it in an
encrypted keystore file.

This module is your first step into cryptographic identity. You'll generate a
secp256k1 private key, derive the corresponding public key and Ethereum address,
and then learn how to securely store the private key using Geth's keystore
format. This is the foundation for creating and managing user accounts.

Computer science principles highlighted:
  - Public-key cryptography (secp256k1)
  - Hashing functions for address derivation (Keccak-256)
  - Key-stretching functions for secure password storage (scrypt)
  - Symmetric encryption for data confidentiality (AES)
*/
func Run(cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

