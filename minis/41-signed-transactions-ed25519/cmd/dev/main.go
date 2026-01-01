package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Transaction represents a simple transaction structure
type Transaction struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Nonce     int64   `json:"nonce"`
	Timestamp int64   `json:"timestamp"`
}

// SignedTransaction contains a transaction and its signature
type SignedTransaction struct {
	Transaction
	Signature string `json:"signature"`
}

// Serialize converts a transaction to bytes for signing
func (tx *Transaction) Serialize() []byte {
	data, err := json.Marshal(tx)
	if err != nil {
		panic(err)
	}
	return data
}

// GenerateKeypair creates a new Ed25519 keypair
func GenerateKeypair() (publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	return pub, priv
}

// SignTransaction signs a transaction with a private key
func SignTransaction(tx *Transaction, privateKey ed25519.PrivateKey) []byte {
	txBytes := tx.Serialize()
	signature := ed25519.Sign(privateKey, txBytes)
	return signature
}

// VerifyTransaction verifies a transaction signature
func VerifyTransaction(tx *Transaction, signature []byte, publicKey ed25519.PublicKey) bool {
	txBytes := tx.Serialize()
	return ed25519.Verify(publicKey, txBytes, signature)
}

// PrintSeparator prints a section separator
func PrintSeparator() {
	fmt.Println()
	for i := 0; i < 80; i++ {
		fmt.Print("=")
	}
	fmt.Println()
}

func main() {
	fmt.Println("Ed25519 Digital Signature Demo")
	fmt.Println("================================")

	// Demo 1: Keypair Generation
	PrintSeparator()
	fmt.Println("DEMO 1: Keypair Generation")
	fmt.Println("==========================")

	pub, priv := GenerateKeypair()

	fmt.Printf("\nPublic Key (32 bytes):\n  %s\n", hex.EncodeToString(pub))
	fmt.Printf("\nPrivate Key (64 bytes):\n  %s\n", hex.EncodeToString(priv))

	fmt.Printf("\nKey Sizes:\n")
	fmt.Printf("  Public:  %d bytes (%d hex characters)\n", len(pub), len(hex.EncodeToString(pub)))
	fmt.Printf("  Private: %d bytes (%d hex characters)\n", len(priv), len(hex.EncodeToString(priv)))

	// Demo 2: Creating and Signing a Transaction
	PrintSeparator()
	fmt.Println("DEMO 2: Creating and Signing a Transaction")
	fmt.Println("===========================================")

	tx := Transaction{
		From:      hex.EncodeToString(pub)[:40] + "...", // Shortened for display
		To:        "recipient_address_xyz123",
		Amount:    10.5,
		Nonce:     1,
		Timestamp: time.Now().Unix(),
	}

	fmt.Printf("\nTransaction Details:\n")
	fmt.Printf("  From:      %s\n", tx.From)
	fmt.Printf("  To:        %s\n", tx.To)
	fmt.Printf("  Amount:    %.2f BTC\n", tx.Amount)
	fmt.Printf("  Nonce:     %d\n", tx.Nonce)
	fmt.Printf("  Timestamp: %d (%s)\n", tx.Timestamp, time.Unix(tx.Timestamp, 0).Format(time.RFC3339))

	txBytes := tx.Serialize()
	fmt.Printf("\nSerialized Transaction (%d bytes):\n  %s\n", len(txBytes), string(txBytes))

	signature := SignTransaction(&tx, priv)
	fmt.Printf("\nSignature (64 bytes):\n  %s\n", hex.EncodeToString(signature))
	fmt.Printf("\nSignature Size: %d bytes\n", len(signature))

	// Demo 3: Verifying a Signature
	PrintSeparator()
	fmt.Println("DEMO 3: Signature Verification")
	fmt.Println("===============================")

	valid := VerifyTransaction(&tx, signature, pub)
	fmt.Printf("\nVerifying with correct public key...\n")
	fmt.Printf("Result: %v ✓\n", valid)

	if valid {
		fmt.Println("\nSignature is VALID!")
		fmt.Println("This proves:")
		fmt.Println("  1. Transaction was signed by the holder of the private key")
		fmt.Println("  2. Transaction has not been modified since signing")
		fmt.Println("  3. Sender authenticity is guaranteed")
	}

	// Demo 4: Tampering Detection
	PrintSeparator()
	fmt.Println("DEMO 4: Tampering Detection")
	fmt.Println("===========================")

	// Try to modify the transaction
	tamperedTx := tx
	tamperedTx.Amount = 1000.0

	fmt.Printf("\nOriginal transaction amount: %.2f BTC\n", tx.Amount)
	fmt.Printf("Tampered transaction amount: %.2f BTC\n", tamperedTx.Amount)

	validTampered := VerifyTransaction(&tamperedTx, signature, pub)
	fmt.Printf("\nVerifying tampered transaction...\n")
	fmt.Printf("Result: %v ✗\n", validTampered)

	if !validTampered {
		fmt.Println("\nSignature is INVALID!")
		fmt.Println("The signature verification detected the tampering.")
		fmt.Println("This prevents attackers from modifying signed transactions.")
	}

	// Demo 5: Wrong Public Key
	PrintSeparator()
	fmt.Println("DEMO 5: Wrong Public Key Detection")
	fmt.Println("===================================")

	// Generate another keypair
	wrongPub, _ := GenerateKeypair()

	fmt.Printf("\nOriginal public key:  %s...\n", hex.EncodeToString(pub)[:40])
	fmt.Printf("Different public key: %s...\n", hex.EncodeToString(wrongPub)[:40])

	validWrongKey := VerifyTransaction(&tx, signature, wrongPub)
	fmt.Printf("\nVerifying with wrong public key...\n")
	fmt.Printf("Result: %v ✗\n", validWrongKey)

	if !validWrongKey {
		fmt.Println("\nSignature is INVALID!")
		fmt.Println("Only the correct public key can verify the signature.")
		fmt.Println("This prevents impersonation attacks.")
	}

	// Demo 6: Complete Signed Transaction
	PrintSeparator()
	fmt.Println("DEMO 6: Complete Signed Transaction")
	fmt.Println("====================================")

	signedTx := SignedTransaction{
		Transaction: tx,
		Signature:   hex.EncodeToString(signature),
	}

	signedJSON, _ := json.MarshalIndent(signedTx, "", "  ")
	fmt.Printf("\nSigned Transaction (JSON):\n%s\n", string(signedJSON))

	fmt.Println("\nThis signed transaction can now be:")
	fmt.Println("  1. Broadcast to a blockchain network")
	fmt.Println("  2. Verified by any node with the public key")
	fmt.Println("  3. Included in a block by miners")
	fmt.Println("  4. Permanently recorded in the blockchain")

	// Demo 7: Signing Multiple Transactions
	PrintSeparator()
	fmt.Println("DEMO 7: Multiple Transactions with Nonces")
	fmt.Println("==========================================")

	fmt.Println("\nWhy nonces? Prevents replay attacks!")
	fmt.Println("Without nonces, an attacker could resubmit the same transaction multiple times.")
	fmt.Println("\nGenerating 3 transactions with different nonces:")

	for i := int64(1); i <= 3; i++ {
		tx := Transaction{
			From:      hex.EncodeToString(pub)[:20] + "...",
			To:        "recipient_xyz",
			Amount:    5.0,
			Nonce:     i,
			Timestamp: time.Now().Unix(),
		}

		sig := SignTransaction(&tx, priv)

		fmt.Printf("\nTransaction #%d:\n", i)
		fmt.Printf("  Amount: %.2f BTC\n", tx.Amount)
		fmt.Printf("  Nonce:  %d\n", tx.Nonce)
		fmt.Printf("  Signature: %s...\n", hex.EncodeToString(sig)[:40])

		// Each signature is different due to different nonce
		// Even though amount is the same!
	}

	// Demo 8: Signing Hashes (Performance Optimization)
	PrintSeparator()
	fmt.Println("DEMO 8: Signing Hashes for Large Data")
	fmt.Println("======================================")

	// Simulate large data
	largeData := make([]byte, 1024*1024) // 1 MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	fmt.Printf("\nLarge data size: %d bytes (1 MB)\n", len(largeData))

	// Method 1: Sign the data directly
	start := time.Now()
	sigDirect := ed25519.Sign(priv, largeData)
	directTime := time.Since(start)

	fmt.Printf("\nMethod 1 - Direct signing:\n")
	fmt.Printf("  Time: %v\n", directTime)
	fmt.Printf("  Signature: %s...\n", hex.EncodeToString(sigDirect)[:40])

	// Method 2: Sign the hash (standard practice)
	start = time.Now()
	hash := sha256.Sum256(largeData)
	sigHash := ed25519.Sign(priv, hash[:])
	hashTime := time.Since(start)

	fmt.Printf("\nMethod 2 - Hash then sign:\n")
	fmt.Printf("  Hash: %s...\n", hex.EncodeToString(hash[:])[:40])
	fmt.Printf("  Time: %v\n", hashTime)
	fmt.Printf("  Signature: %s...\n", hex.EncodeToString(sigHash)[:40])

	fmt.Printf("\nPerformance comparison:\n")
	fmt.Printf("  Direct:     %v\n", directTime)
	fmt.Printf("  Hash-first: %v\n", hashTime)
	fmt.Printf("  Speedup:    %.2fx faster\n", float64(directTime)/float64(hashTime))

	fmt.Println("\nNote: For large data, signing the hash is standard practice.")
	fmt.Println("Used in: Bitcoin, Ethereum, Git commits, file signatures, etc.")

	// Demo 9: Batch Verification
	PrintSeparator()
	fmt.Println("DEMO 9: Batch Transaction Processing")
	fmt.Println("=====================================")

	fmt.Println("\nSimulating a block of transactions:")

	transactions := make([]Transaction, 5)
	signatures := make([][]byte, 5)

	for i := 0; i < 5; i++ {
		transactions[i] = Transaction{
			From:      hex.EncodeToString(pub)[:20] + "...",
			To:        fmt.Sprintf("recipient_%d", i),
			Amount:    float64(i+1) * 1.5,
			Nonce:     int64(i + 1),
			Timestamp: time.Now().Unix(),
		}
		signatures[i] = SignTransaction(&transactions[i], priv)
	}

	fmt.Printf("\nVerifying %d transactions:\n", len(transactions))
	allValid := true
	for i := range transactions {
		valid := VerifyTransaction(&transactions[i], signatures[i], pub)
		status := "✓"
		if !valid {
			status = "✗"
			allValid = false
		}
		fmt.Printf("  Transaction %d: %s (%.2f BTC)\n", i+1, status, transactions[i].Amount)
	}

	if allValid {
		fmt.Println("\nAll transactions verified successfully!")
		fmt.Println("Block can be added to the blockchain.")
	}

	// Summary
	PrintSeparator()
	fmt.Println("SUMMARY")
	fmt.Println("=======")
	fmt.Println("\nEd25519 Digital Signatures provide:")
	fmt.Println("  ✓ Authentication  - Prove who created the transaction")
	fmt.Println("  ✓ Integrity       - Detect any modifications")
	fmt.Println("  ✓ Non-repudiation - Signer cannot deny signing")
	fmt.Println("  ✓ Performance     - 50,000+ signatures/sec")
	fmt.Println("  ✓ Security        - 128-bit security level")
	fmt.Println("  ✓ Compact         - 32-byte keys, 64-byte signatures")

	fmt.Println("\nUsed by:")
	fmt.Println("  • Cryptocurrencies (Bitcoin, Ethereum, Solana)")
	fmt.Println("  • Secure messaging (Signal, WhatsApp)")
	fmt.Println("  • SSH authentication")
	fmt.Println("  • Code signing (Git commits, software packages)")
	fmt.Println("  • TLS certificates")

	fmt.Println("\nKey takeaways:")
	fmt.Println("  1. Private keys must NEVER be shared")
	fmt.Println("  2. Public keys can be freely distributed")
	fmt.Println("  3. Signatures prove ownership without revealing secrets")
	fmt.Println("  4. Nonces prevent replay attacks")
	fmt.Println("  5. Tampering is always detected")

	PrintSeparator()
	fmt.Println("\nDemo completed successfully!")
}
