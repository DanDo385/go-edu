package exercise

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"
)

// TestGenerateWallet tests wallet generation
func TestGenerateWallet(t *testing.T) {
	wallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("GenerateWallet() error = %v", err)
	}

	if wallet == nil {
		t.Fatal("GenerateWallet() returned nil wallet")
	}

	// Check public key
	if len(wallet.PublicKey) != ed25519.PublicKeySize {
		t.Errorf("PublicKey size = %d, want %d", len(wallet.PublicKey), ed25519.PublicKeySize)
	}

	// Check private key
	if len(wallet.PrivateKey) != ed25519.PrivateKeySize {
		t.Errorf("PrivateKey size = %d, want %d", len(wallet.PrivateKey), ed25519.PrivateKeySize)
	}

	// Check address is derived
	if wallet.Address == "" {
		t.Error("Wallet.Address is empty")
	}

	// Check address format (should be 40 hex characters = 20 bytes)
	if len(wallet.Address) != 40 {
		t.Errorf("Address length = %d, want 40 (20 bytes hex)", len(wallet.Address))
	}

	// Check address is valid hex
	_, err = hex.DecodeString(wallet.Address)
	if err != nil {
		t.Errorf("Address is not valid hex: %v", err)
	}
}

// TestDeriveAddress tests address derivation
func TestDeriveAddress(t *testing.T) {
	wallet, _ := GenerateWallet()
	address := DeriveAddress(wallet.PublicKey)

	// Should be 40 hex characters
	if len(address) != 40 {
		t.Errorf("Address length = %d, want 40", len(address))
	}

	// Should be valid hex
	decoded, err := hex.DecodeString(address)
	if err != nil {
		t.Errorf("Address is not valid hex: %v", err)
	}

	// Should be 20 bytes
	if len(decoded) != 20 {
		t.Errorf("Decoded address length = %d, want 20", len(decoded))
	}

	// Should match first 20 bytes of public key
	expectedHex := hex.EncodeToString(wallet.PublicKey[:20])
	if address != expectedHex {
		t.Errorf("Address = %s, want %s", address, expectedHex)
	}

	// Deterministic: same public key -> same address
	address2 := DeriveAddress(wallet.PublicKey)
	if address != address2 {
		t.Error("DeriveAddress is not deterministic")
	}
}

// TestNewTransaction tests transaction creation
func TestNewTransaction(t *testing.T) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient_address", 10.5, 1)

	if tx == nil {
		t.Fatal("NewTransaction() returned nil")
	}

	// Check fields
	expectedFrom := hex.EncodeToString(wallet.PublicKey)
	if tx.From != expectedFrom {
		t.Errorf("Transaction.From = %s, want %s", tx.From, expectedFrom)
	}

	if tx.To != "recipient_address" {
		t.Errorf("Transaction.To = %s, want %s", tx.To, "recipient_address")
	}

	if tx.Amount != 10.5 {
		t.Errorf("Transaction.Amount = %f, want %f", tx.Amount, 10.5)
	}

	if tx.Nonce != 1 {
		t.Errorf("Transaction.Nonce = %d, want %d", tx.Nonce, 1)
	}

	// Timestamp should be recent
	now := time.Now().Unix()
	if tx.Timestamp < now-10 || tx.Timestamp > now+10 {
		t.Errorf("Transaction.Timestamp = %d, should be close to %d", tx.Timestamp, now)
	}
}

// TestTransactionSerialize tests transaction serialization
func TestTransactionSerialize(t *testing.T) {
	tx := &Transaction{
		From:      "sender_address",
		To:        "recipient_address",
		Amount:    10.5,
		Nonce:     1,
		Timestamp: 1234567890,
	}

	data, err := tx.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	if len(data) == 0 {
		t.Error("Serialize() returned empty data")
	}

	// Should be valid JSON
	var decoded Transaction
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Serialized data is not valid JSON: %v", err)
	}

	// Should deserialize to same values
	if decoded.From != tx.From {
		t.Errorf("Deserialized From = %s, want %s", decoded.From, tx.From)
	}
	if decoded.To != tx.To {
		t.Errorf("Deserialized To = %s, want %s", decoded.To, tx.To)
	}
	if decoded.Amount != tx.Amount {
		t.Errorf("Deserialized Amount = %f, want %f", decoded.Amount, tx.Amount)
	}

	// Deterministic: same transaction -> same serialization
	data2, _ := tx.Serialize()
	if string(data) != string(data2) {
		t.Error("Serialize() is not deterministic")
	}
}

// TestWalletSign tests transaction signing
func TestWalletSign(t *testing.T) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient", 10.5, 1)

	signedTx, err := wallet.Sign(tx)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	if signedTx == nil {
		t.Fatal("Sign() returned nil")
	}

	// Check signature is present
	if signedTx.Signature == "" {
		t.Error("Signature is empty")
	}

	// Check signature is valid hex
	sigBytes, err := hex.DecodeString(signedTx.Signature)
	if err != nil {
		t.Errorf("Signature is not valid hex: %v", err)
	}

	// Check signature length (64 bytes = 128 hex characters)
	if len(sigBytes) != ed25519.SignatureSize {
		t.Errorf("Signature size = %d, want %d", len(sigBytes), ed25519.SignatureSize)
	}

	// Check public key is present
	if signedTx.PublicKey == "" {
		t.Error("PublicKey is empty")
	}

	// Check public key matches wallet
	expectedPubKey := hex.EncodeToString(wallet.PublicKey)
	if signedTx.PublicKey != expectedPubKey {
		t.Errorf("PublicKey = %s, want %s", signedTx.PublicKey, expectedPubKey)
	}

	// Check transaction is preserved
	if signedTx.Transaction.Amount != tx.Amount {
		t.Error("Transaction not preserved in SignedTransaction")
	}
}

// TestSignedTransactionVerify tests signature verification
func TestSignedTransactionVerify(t *testing.T) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient", 10.5, 1)
	signedTx, _ := wallet.Sign(tx)

	valid, err := signedTx.Verify()
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if !valid {
		t.Error("Verify() = false, want true for valid signature")
	}

	// Test tampering detection
	t.Run("tampered_amount", func(t *testing.T) {
		tampered := *signedTx
		tampered.Transaction.Amount = 1000.0

		valid, err := tampered.Verify()
		if err != nil {
			t.Fatalf("Verify() error = %v", err)
		}
		if valid {
			t.Error("Verify() = true, want false for tampered transaction")
		}
	})

	t.Run("tampered_recipient", func(t *testing.T) {
		tampered := *signedTx
		tampered.Transaction.To = "attacker_address"

		valid, err := tampered.Verify()
		if err != nil {
			t.Fatalf("Verify() error = %v", err)
		}
		if valid {
			t.Error("Verify() = true, want false for tampered recipient")
		}
	})

	t.Run("invalid_signature", func(t *testing.T) {
		tampered := *signedTx
		tampered.Signature = "0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000"

		valid, err := tampered.Verify()
		if err != nil {
			t.Fatalf("Verify() error = %v", err)
		}
		if valid {
			t.Error("Verify() = true, want false for invalid signature")
		}
	})

	t.Run("wrong_public_key", func(t *testing.T) {
		wrongWallet, _ := GenerateWallet()
		tampered := *signedTx
		tampered.PublicKey = hex.EncodeToString(wrongWallet.PublicKey)

		valid, err := tampered.Verify()
		if err != nil {
			t.Fatalf("Verify() error = %v", err)
		}
		if valid {
			t.Error("Verify() = true, want false for wrong public key")
		}
	})
}

// TestVerifyOwnership tests ownership verification
func TestVerifyOwnership(t *testing.T) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient", 10.5, 1)
	signedTx, _ := wallet.Sign(tx)

	if !signedTx.VerifyOwnership() {
		t.Error("VerifyOwnership() = false, want true for matching sender")
	}

	// Test mismatched sender
	t.Run("mismatched_sender", func(t *testing.T) {
		wrongWallet, _ := GenerateWallet()
		tampered := *signedTx
		tampered.Transaction.From = hex.EncodeToString(wrongWallet.PublicKey)

		if tampered.VerifyOwnership() {
			t.Error("VerifyOwnership() = true, want false for mismatched sender")
		}
	})
}

// TestGetTransactionID tests transaction ID generation
func TestGetTransactionID(t *testing.T) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient", 10.5, 1)
	signedTx, _ := wallet.Sign(tx)

	txID := signedTx.GetTransactionID()

	// Should be 32 hex characters (16 bytes)
	if len(txID) != 32 {
		t.Errorf("Transaction ID length = %d, want 32", len(txID))
	}

	// Should be valid hex
	_, err := hex.DecodeString(txID)
	if err != nil {
		t.Errorf("Transaction ID is not valid hex: %v", err)
	}

	// Deterministic: same signature -> same ID
	txID2 := signedTx.GetTransactionID()
	if txID != txID2 {
		t.Error("GetTransactionID() is not deterministic")
	}

	// Different signatures -> different IDs
	tx2 := NewTransaction(wallet, "recipient", 10.5, 2) // Different nonce
	signedTx2, _ := wallet.Sign(tx2)
	txID3 := signedTx2.GetTransactionID()
	if txID == txID3 {
		t.Error("Different transactions produced same transaction ID")
	}
}

// TestNoncePreventsDuplicates tests that nonce prevents duplicate transactions
func TestNoncePreventsDuplicates(t *testing.T) {
	wallet, _ := GenerateWallet()

	tx1 := NewTransaction(wallet, "recipient", 10.5, 1)
	tx2 := NewTransaction(wallet, "recipient", 10.5, 2)

	signedTx1, _ := wallet.Sign(tx1)
	signedTx2, _ := wallet.Sign(tx2)

	// Signatures should be different
	if signedTx1.Signature == signedTx2.Signature {
		t.Error("Same signature for different nonces (should be different)")
	}

	// Transaction IDs should be different
	if signedTx1.GetTransactionID() == signedTx2.GetTransactionID() {
		t.Error("Same transaction ID for different nonces")
	}
}

// TestMultiSigTransaction tests multi-signature transactions
func TestMultiSigTransaction(t *testing.T) {
	// Create 3 wallets
	wallet1, _ := GenerateWallet()
	wallet2, _ := GenerateWallet()
	wallet3, _ := GenerateWallet()

	// Create transaction
	tx := NewTransaction(wallet1, "recipient", 100.0, 1)

	// Create 2-of-3 multi-sig
	multiSig := NewMultiSigTransaction(tx, 2)

	if multiSig == nil {
		t.Fatal("NewMultiSigTransaction() returned nil")
	}

	if multiSig.Required != 2 {
		t.Errorf("Required = %d, want 2", multiSig.Required)
	}

	// Add first signature
	if err := multiSig.AddSignature(wallet1); err != nil {
		t.Fatalf("AddSignature(wallet1) error = %v", err)
	}

	// Only 1 signature, should not verify
	valid, err := multiSig.Verify()
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if valid {
		t.Error("Verify() = true with 1 signature, want false (need 2)")
	}

	// Add second signature
	if err := multiSig.AddSignature(wallet2); err != nil {
		t.Fatalf("AddSignature(wallet2) error = %v", err)
	}

	// Now should verify
	valid, err = multiSig.Verify()
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if !valid {
		t.Error("Verify() = false with 2 signatures, want true")
	}

	// Add third signature (more than required is OK)
	if err := multiSig.AddSignature(wallet3); err != nil {
		t.Fatalf("AddSignature(wallet3) error = %v", err)
	}

	// Should still verify
	valid, err = multiSig.Verify()
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if !valid {
		t.Error("Verify() = false with 3 signatures, want true")
	}

	// Check signature count
	if len(multiSig.Signatures) != 3 {
		t.Errorf("Signature count = %d, want 3", len(multiSig.Signatures))
	}
}

// TestDeterministicSignatures tests that Ed25519 signatures are deterministic
func TestDeterministicSignatures(t *testing.T) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient", 10.5, 1)

	// Sign multiple times
	signedTx1, _ := wallet.Sign(tx)
	signedTx2, _ := wallet.Sign(tx)

	// Signatures should be identical
	if signedTx1.Signature != signedTx2.Signature {
		t.Error("Ed25519 signatures should be deterministic (same message -> same signature)")
	}
}

// TestJSONRoundTrip tests JSON encoding and decoding
func TestJSONRoundTrip(t *testing.T) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient", 10.5, 1)
	signedTx, _ := wallet.Sign(tx)

	// Encode to JSON
	jsonData, err := json.Marshal(signedTx)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Decode from JSON
	var decoded SignedTransaction
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Should match original
	if decoded.Signature != signedTx.Signature {
		t.Error("Signature mismatch after JSON round-trip")
	}
	if decoded.PublicKey != signedTx.PublicKey {
		t.Error("PublicKey mismatch after JSON round-trip")
	}
	if decoded.Transaction.Amount != signedTx.Transaction.Amount {
		t.Error("Transaction.Amount mismatch after JSON round-trip")
	}

	// Should still verify
	valid, err := decoded.Verify()
	if err != nil {
		t.Fatalf("Verify() after round-trip error = %v", err)
	}
	if !valid {
		t.Error("Signature invalid after JSON round-trip")
	}
}

// TestConcurrentSigning tests concurrent transaction signing
func TestConcurrentSigning(t *testing.T) {
	wallet, _ := GenerateWallet()

	const numTx = 100
	results := make(chan *SignedTransaction, numTx)

	// Sign transactions concurrently
	for i := 0; i < numTx; i++ {
		go func(nonce int64) {
			tx := NewTransaction(wallet, "recipient", 10.5, nonce)
			signedTx, _ := wallet.Sign(tx)
			results <- signedTx
		}(int64(i))
	}

	// Collect results
	signedTxs := make([]*SignedTransaction, 0, numTx)
	for i := 0; i < numTx; i++ {
		signedTxs = append(signedTxs, <-results)
	}

	// Verify all signatures
	for i, signedTx := range signedTxs {
		valid, err := signedTx.Verify()
		if err != nil {
			t.Errorf("Transaction %d: Verify() error = %v", i, err)
		}
		if !valid {
			t.Errorf("Transaction %d: invalid signature", i)
		}
	}
}

// BenchmarkGenerateWallet benchmarks wallet generation
func BenchmarkGenerateWallet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateWallet()
	}
}

// BenchmarkSign benchmarks transaction signing
func BenchmarkSign(b *testing.B) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient", 10.5, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wallet.Sign(tx)
	}
}

// BenchmarkVerify benchmarks signature verification
func BenchmarkVerify(b *testing.B) {
	wallet, _ := GenerateWallet()
	tx := NewTransaction(wallet, "recipient", 10.5, 1)
	signedTx, _ := wallet.Sign(tx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signedTx.Verify()
	}
}

// BenchmarkSignAndVerify benchmarks the full sign-verify cycle
func BenchmarkSignAndVerify(b *testing.B) {
	wallet, _ := GenerateWallet()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx := NewTransaction(wallet, "recipient", 10.5, int64(i))
		signedTx, _ := wallet.Sign(tx)
		signedTx.Verify()
	}
}
