package exercise

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestRunCreatesKeystore(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pass := "horse battery staple"

	res, err := Run(Config{
		OutputDir:  dir,
		Passphrase: pass,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if res == nil {
		t.Fatalf("Run returned nil result")
	}
	if res.Address == (common.Address{}) {
		t.Fatalf("expected non-zero address")
	}
	if res.PrivateKeyHex == "" {
		t.Fatalf("expected private key hex to be populated")
	}
	if res.KeystorePath == "" {
		t.Fatalf("expected keystore path to be populated")
	}

	keyJSON, err := os.ReadFile(res.KeystorePath)
	if err != nil {
		t.Fatalf("read keystore file: %v", err)
	}

	key, err := keystore.DecryptKey(keyJSON, pass)
	if err != nil {
		t.Fatalf("decrypt keystore: %v", err)
	}

	gotPriv := crypto.FromECDSA(key.PrivateKey)
	if hexKey := strings.ToLower(res.PrivateKeyHex); hexKey != strings.ToLower(hex.EncodeToString(gotPriv)) {
		t.Fatalf("private key mismatch: want %s got %s", hex.EncodeToString(gotPriv), res.PrivateKeyHex)
	}

	if key.Address != res.Address {
		t.Fatalf("address mismatch: keystore %s vs result %s", key.Address.Hex(), res.Address.Hex())
	}

	if gotDir := filepath.Dir(res.KeystorePath); gotDir != dir {
		t.Fatalf("expected keystore in %s, got %s", dir, gotDir)
	}
}

func TestRunUsesDefaults(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("relative path permissions behave differently on Windows")
	}

	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldDir)
	})

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}

	res, err := Run(Config{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if res == nil {
		t.Fatalf("Run returned nil result")
	}
	if !strings.Contains(res.KeystorePath, "keystore-demo") {
		t.Fatalf("expected default keystore dir to include keystore-demo, got %s", res.KeystorePath)
	}

	if _, err := os.Stat(res.KeystorePath); err != nil {
		t.Fatalf("keystore file not found at %s: %v", res.KeystorePath, err)
	}
}
