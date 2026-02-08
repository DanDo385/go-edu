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
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

func Run(cfg Config) (*Result, error) {
	outDir := cfg.OutputDir
	if outDir == "" {
		outDir = defaultOutputDir
	}
	passphrase := cfg.Passphrase
	if passphrase == "" {
		passphrase = defaultPassphrase
	}

	if err := os.MkdirAll(outDir, 0o700); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	ks := keystore.NewKeyStore(outDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	keyJSON, err := os.ReadFile(account.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("read keystore file: %w", err)
	}

	key, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt keystore file: %w", err)
	}

	return &Result{
		Address:       account.Address,
		PrivateKeyHex: hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)),
		KeystorePath:  account.URL.Path,
	}, nil
}
