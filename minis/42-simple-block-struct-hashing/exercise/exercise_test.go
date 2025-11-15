package exercise

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"time"
)

func TestNewGenesisBlock(t *testing.T) {
	genesis := NewGenesisBlock()

	if genesis == nil {
		t.Fatal("NewGenesisBlock returned nil")
	}

	if genesis.Header.Index != 0 {
		t.Errorf("Expected genesis index 0, got %d", genesis.Header.Index)
	}

	expectedPrevHash := strings.Repeat("0", 64)
	if genesis.Header.PrevHash != expectedPrevHash {
		t.Errorf("Expected genesis PrevHash to be all zeros, got %s", genesis.Header.PrevHash)
	}

	if len(genesis.Transactions) == 0 {
		t.Error("Expected genesis block to have at least one transaction")
	}

	if genesis.Hash == "" {
		t.Error("Expected genesis block to have a hash")
	}

	if genesis.Header.MerkleRoot == "" {
		t.Error("Expected genesis block to have a merkle root")
	}

	// Verify hash is correct
	computedHash := genesis.ComputeHash()
	if genesis.Hash != computedHash {
		t.Errorf("Genesis hash mismatch: stored=%s, computed=%s", genesis.Hash, computedHash)
	}
}

func TestComputeMerkleRoot(t *testing.T) {
	tests := []struct {
		name         string
		transactions []string
		wantEmpty    bool
	}{
		{
			name:         "empty transactions",
			transactions: []string{},
			wantEmpty:    true,
		},
		{
			name:         "single transaction",
			transactions: []string{"Alice sends 10 BTC to Bob"},
			wantEmpty:    false,
		},
		{
			name:         "multiple transactions",
			transactions: []string{"tx1", "tx2", "tx3"},
			wantEmpty:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merkleRoot := ComputeMerkleRoot(tt.transactions)

			if tt.wantEmpty && merkleRoot != "" {
				t.Errorf("Expected empty merkle root for empty transactions, got %s", merkleRoot)
			}

			if !tt.wantEmpty && merkleRoot == "" {
				t.Error("Expected non-empty merkle root")
			}

			// Verify it's a valid hex string
			if merkleRoot != "" {
				if _, err := hex.DecodeString(merkleRoot); err != nil {
					t.Errorf("Merkle root is not valid hex: %v", err)
				}

				// SHA-256 produces 32 bytes = 64 hex characters
				if len(merkleRoot) != 64 {
					t.Errorf("Expected merkle root length 64, got %d", len(merkleRoot))
				}
			}
		})
	}
}

func TestComputeMerkleRootDeterministic(t *testing.T) {
	transactions := []string{"tx1", "tx2", "tx3"}

	root1 := ComputeMerkleRoot(transactions)
	root2 := ComputeMerkleRoot(transactions)

	if root1 != root2 {
		t.Error("Merkle root should be deterministic (same input → same output)")
	}
}

func TestBlockSerialization(t *testing.T) {
	genesis := NewGenesisBlock()
	if genesis == nil {
		t.Fatal("NewGenesisBlock returned nil")
	}

	// Test that serialization is deterministic
	bytes1 := genesis.Serialize()
	bytes2 := genesis.Serialize()

	if len(bytes1) == 0 {
		t.Error("Serialization returned empty bytes")
	}

	if string(bytes1) != string(bytes2) {
		t.Error("Serialization should be deterministic")
	}
}

func TestComputeHash(t *testing.T) {
	genesis := NewGenesisBlock()
	if genesis == nil {
		t.Fatal("NewGenesisBlock returned nil")
	}

	hash := genesis.ComputeHash()

	// Check it's a valid hex string
	if _, err := hex.DecodeString(hash); err != nil {
		t.Errorf("Hash is not valid hex: %v", err)
	}

	// SHA-256 produces 32 bytes = 64 hex characters
	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}

	// Hash should be deterministic
	hash2 := genesis.ComputeHash()
	if hash != hash2 {
		t.Error("Hash should be deterministic")
	}
}

func TestNewBlock(t *testing.T) {
	genesis := NewGenesisBlock()
	if genesis == nil {
		t.Fatal("NewGenesisBlock returned nil")
	}

	transactions := []string{"Alice sends 10 BTC to Bob", "Bob sends 5 BTC to Carol"}
	block1 := NewBlock(genesis, transactions)

	if block1 == nil {
		t.Fatal("NewBlock returned nil")
	}

	// Check index incremented
	if block1.Header.Index != genesis.Header.Index+1 {
		t.Errorf("Expected block index %d, got %d", genesis.Header.Index+1, block1.Header.Index)
	}

	// Check previous hash link
	if block1.Header.PrevHash != genesis.Hash {
		t.Errorf("Expected PrevHash to match genesis hash")
	}

	// Check transactions
	if len(block1.Transactions) != len(transactions) {
		t.Errorf("Expected %d transactions, got %d", len(transactions), len(block1.Transactions))
	}

	// Check hash is set
	if block1.Hash == "" {
		t.Error("Expected block to have a hash")
	}

	// Check merkle root is set
	if block1.Header.MerkleRoot == "" {
		t.Error("Expected block to have a merkle root")
	}

	// Verify hash is correct
	computedHash := block1.ComputeHash()
	if block1.Hash != computedHash {
		t.Errorf("Block hash mismatch: stored=%s, computed=%s", block1.Hash, computedHash)
	}
}

func TestHashLinking(t *testing.T) {
	genesis := NewGenesisBlock()
	block1 := NewBlock(genesis, []string{"tx1"})
	block2 := NewBlock(block1, []string{"tx2"})

	// Verify hash chain
	if block1.Header.PrevHash != genesis.Hash {
		t.Error("Block 1 should link to genesis")
	}

	if block2.Header.PrevHash != block1.Hash {
		t.Error("Block 2 should link to block 1")
	}

	// Verify hashes are different
	if genesis.Hash == block1.Hash {
		t.Error("Genesis and block 1 should have different hashes")
	}

	if block1.Hash == block2.Hash {
		t.Error("Block 1 and block 2 should have different hashes")
	}
}

func TestValidateChain_ValidChain(t *testing.T) {
	// Build a valid chain
	chain := []*Block{
		NewGenesisBlock(),
	}

	chain = append(chain, NewBlock(chain[0], []string{"tx1"}))
	chain = append(chain, NewBlock(chain[1], []string{"tx2", "tx3"}))
	chain = append(chain, NewBlock(chain[2], []string{"tx4"}))

	// Validate should pass
	if err := ValidateChain(chain); err != nil {
		t.Errorf("Expected valid chain, got error: %v", err)
	}
}

func TestValidateChain_EmptyChain(t *testing.T) {
	chain := []*Block{}

	err := ValidateChain(chain)
	if err == nil {
		t.Error("Expected error for empty chain")
	}
}

func TestValidateChain_InvalidGenesisIndex(t *testing.T) {
	genesis := NewGenesisBlock()
	genesis.Header.Index = 1 // Wrong index

	chain := []*Block{genesis}

	err := ValidateChain(chain)
	if err == nil {
		t.Error("Expected error for invalid genesis index")
	}
}

func TestValidateChain_TamperedData(t *testing.T) {
	// Build a valid chain
	chain := []*Block{
		NewGenesisBlock(),
	}
	chain = append(chain, NewBlock(chain[0], []string{"tx1"}))
	chain = append(chain, NewBlock(chain[1], []string{"tx2"}))

	// Tamper with block 1's transaction
	chain[1].Transactions[0] = "tampered transaction"

	// Validation should fail (merkle root or hash mismatch)
	err := ValidateChain(chain)
	if err == nil {
		t.Error("Expected error for tampered data")
	}
}

func TestValidateChain_TamperedHash(t *testing.T) {
	// Build a valid chain
	chain := []*Block{
		NewGenesisBlock(),
	}
	block1 := NewBlock(chain[0], []string{"tx1"})
	chain = append(chain, block1)

	// Tamper with block 1's hash
	chain[1].Hash = "0000000000000000000000000000000000000000000000000000000000000000"

	// Validation should fail (hash mismatch)
	err := ValidateChain(chain)
	if err == nil {
		t.Error("Expected error for tampered hash")
	}
}

func TestValidateChain_BrokenLink(t *testing.T) {
	// Build a valid chain
	chain := []*Block{
		NewGenesisBlock(),
	}
	chain = append(chain, NewBlock(chain[0], []string{"tx1"}))
	chain = append(chain, NewBlock(chain[1], []string{"tx2"}))

	// Break the link by changing block 2's PrevHash
	chain[2].Header.PrevHash = "0000000000000000000000000000000000000000000000000000000000000000"
	chain[2].Hash = chain[2].ComputeHash() // Recompute hash with broken link

	// Validation should fail (prev hash mismatch)
	err := ValidateChain(chain)
	if err == nil {
		t.Error("Expected error for broken hash link")
	}
}

func TestValidateChain_NonSequentialIndex(t *testing.T) {
	// Build a chain with non-sequential indexes
	genesis := NewGenesisBlock()
	block1 := NewBlock(genesis, []string{"tx1"})
	block1.Header.Index = 5 // Skip to index 5
	block1.Hash = block1.ComputeHash()

	chain := []*Block{genesis, block1}

	err := ValidateChain(chain)
	if err == nil {
		t.Error("Expected error for non-sequential indexes")
	}
}

func TestValidateChain_TimestampOutOfOrder(t *testing.T) {
	genesis := NewGenesisBlock()

	time.Sleep(10 * time.Millisecond)

	block1 := NewBlock(genesis, []string{"tx1"})

	// Backdate block 1
	block1.Header.Timestamp = genesis.Header.Timestamp - 100
	block1.Hash = block1.ComputeHash()

	chain := []*Block{genesis, block1}

	err := ValidateChain(chain)
	if err == nil {
		t.Error("Expected error for timestamp out of order")
	}
}

func TestHashAvalancheEffect(t *testing.T) {
	// Small change in input should produce completely different hash
	genesis := NewGenesisBlock()

	block1 := NewBlock(genesis, []string{"Alice sends 10 BTC to Bob"})
	block2 := NewBlock(genesis, []string{"Alice sends 11 BTC to Bob"}) // Only 1 digit different

	// Hashes should be completely different
	if block1.Hash == block2.Hash {
		t.Error("Expected different hashes for different data")
	}

	// Count different characters (should be ~50% different)
	differentChars := 0
	for i := 0; i < len(block1.Hash); i++ {
		if block1.Hash[i] != block2.Hash[i] {
			differentChars++
		}
	}

	// At least 30% of characters should be different (avalanche effect)
	minDifferent := len(block1.Hash) * 30 / 100
	if differentChars < minDifferent {
		t.Errorf("Expected at least %d different characters (avalanche effect), got %d",
			minDifferent, differentChars)
	}
}

func TestMerkleRootChangesWithTransactions(t *testing.T) {
	genesis := NewGenesisBlock()

	block1 := NewBlock(genesis, []string{"tx1"})
	block2 := NewBlock(genesis, []string{"tx2"})

	// Different transactions should produce different merkle roots
	if block1.Header.MerkleRoot == block2.Header.MerkleRoot {
		t.Error("Expected different merkle roots for different transactions")
	}
}

func TestBlockHashIncludesMerkleRoot(t *testing.T) {
	genesis := NewGenesisBlock()

	block := NewBlock(genesis, []string{"tx1", "tx2"})
	originalHash := block.Hash

	// Change merkle root (simulate tampering)
	block.Header.MerkleRoot = "0000000000000000000000000000000000000000000000000000000000000000"

	// Recompute hash
	newHash := block.ComputeHash()

	// Hash should be different (proves hash includes merkle root)
	if originalHash == newHash {
		t.Error("Hash should change when merkle root changes")
	}
}

func TestValidateChain_InvalidMerkleRoot(t *testing.T) {
	genesis := NewGenesisBlock()
	block1 := NewBlock(genesis, []string{"tx1"})

	// Tamper with merkle root
	block1.Header.MerkleRoot = "invalid"
	// Don't recompute hash (so hash is correct but merkle root is wrong)

	chain := []*Block{genesis, block1}

	err := ValidateChain(chain)
	if err == nil {
		t.Error("Expected error for invalid merkle root")
	}
}

// Benchmark serialization
func BenchmarkSerialize(b *testing.B) {
	genesis := NewGenesisBlock()
	block := NewBlock(genesis, []string{
		"Transaction 1",
		"Transaction 2",
		"Transaction 3",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = block.Serialize()
	}
}

// Benchmark hash computation
func BenchmarkComputeHash(b *testing.B) {
	genesis := NewGenesisBlock()
	block := NewBlock(genesis, []string{
		"Transaction 1",
		"Transaction 2",
		"Transaction 3",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = block.ComputeHash()
	}
}

// Benchmark merkle root computation
func BenchmarkComputeMerkleRoot(b *testing.B) {
	transactions := []string{
		"Transaction 1",
		"Transaction 2",
		"Transaction 3",
		"Transaction 4",
		"Transaction 5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ComputeMerkleRoot(transactions)
	}
}

// Benchmark chain validation
func BenchmarkValidateChain(b *testing.B) {
	// Build a chain with 100 blocks
	chain := []*Block{NewGenesisBlock()}
	for i := 0; i < 99; i++ {
		block := NewBlock(chain[len(chain)-1], []string{"tx1", "tx2"})
		chain = append(chain, block)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateChain(chain)
	}
}

// Test that demonstrates the security of hash linking
func TestSecurityHashLinking(t *testing.T) {
	// Build a chain
	chain := []*Block{NewGenesisBlock()}
	for i := 0; i < 5; i++ {
		block := NewBlock(chain[len(chain)-1], []string{
			"Transaction " + string(rune('A'+i)),
		})
		chain = append(chain, block)
	}

	// Validate original chain
	if err := ValidateChain(chain); err != nil {
		t.Fatalf("Original chain should be valid: %v", err)
	}

	// Attempt 1: Modify transaction in middle block without updating hash
	tamperedChain1 := make([]*Block, len(chain))
	copy(tamperedChain1, chain)
	tamperedChain1[2].Transactions[0] = "Fraudulent transaction"

	if err := ValidateChain(tamperedChain1); err == nil {
		t.Error("Should detect tampered transaction (merkle root mismatch)")
	}

	// Attempt 2: Modify transaction and recompute that block's hash
	tamperedChain2 := make([]*Block, len(chain))
	for i := range chain {
		tamperedChain2[i] = &Block{
			Header:       chain[i].Header,
			Transactions: make([]string, len(chain[i].Transactions)),
			Hash:         chain[i].Hash,
		}
		copy(tamperedChain2[i].Transactions, chain[i].Transactions)
	}
	tamperedChain2[2].Transactions[0] = "Fraudulent transaction"
	tamperedChain2[2].Header.MerkleRoot = ComputeMerkleRoot(tamperedChain2[2].Transactions)
	tamperedChain2[2].Hash = tamperedChain2[2].ComputeHash()

	if err := ValidateChain(tamperedChain2); err == nil {
		t.Error("Should detect broken chain (next block's PrevHash doesn't match)")
	}
}

// Helper to verify SHA-256 properties
func TestSHA256Properties(t *testing.T) {
	data := []byte("test data")

	// Property 1: Deterministic
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(data)
	if hash1 != hash2 {
		t.Error("SHA-256 should be deterministic")
	}

	// Property 2: Fixed size (32 bytes)
	if len(hash1) != 32 {
		t.Errorf("SHA-256 should produce 32 bytes, got %d", len(hash1))
	}

	// Property 3: Avalanche effect
	data2 := []byte("test datb") // One character different
	hash3 := sha256.Sum256(data2)

	differentBytes := 0
	for i := 0; i < 32; i++ {
		if hash1[i] != hash3[i] {
			differentBytes++
		}
	}

	// At least 10 bytes should be different (avalanche effect)
	if differentBytes < 10 {
		t.Errorf("Expected significant difference due to avalanche effect, only %d bytes different", differentBytes)
	}
}
