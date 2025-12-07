//go:build solution
// +build solution

package exercise

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	defaultOutDir   = "./keystore-demo"
	defaultPassword = "changeit"
)

// Run contains the reference solution for module 03-keys-addresses.
func Run(cfg Config) (*Result, error) {
	if cfg.OutputDir == "" {
		cfg.OutputDir = defaultOutDir
	}
	if cfg.Passphrase == "" {
		cfg.Passphrase = defaultPassword
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("generate secp256k1 key: %w", err)
	}

	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	if err := os.MkdirAll(cfg.OutputDir, 0o700); err != nil {
		return nil, fmt.Errorf("create keystore directory: %w", err)
	}

	ks := keystore.NewKeyStore(cfg.OutputDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.ImportECDSA(privKey, cfg.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("import key into keystore: %w", err)
	}

	if err := ks.Unlock(account, cfg.Passphrase); err != nil {
		return nil, fmt.Errorf("unlock keystore: %w", err)
	}
	defer ks.Lock(account.Address)

	keystorePath := account.URL.Path
	if keystorePath == "" {
		keystorePath = filepath.Join(cfg.OutputDir, filepath.Base(account.URL.String()))
	} else {
		keystorePath = filepath.Clean(keystorePath)
	}

	return &Result{
		Address:       addr,
		PrivateKeyHex: hex.EncodeToString(crypto.FromECDSA(privKey)),
		KeystorePath:  keystorePath,
	}, nil
}
