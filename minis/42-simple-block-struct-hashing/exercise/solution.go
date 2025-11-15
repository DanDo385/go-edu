//go:build solution
// +build solution

/*
Problem: Implement a simple blockchain with blocks, hashing, and validation

Requirements:
1. Block structure with header and transactions
2. Cryptographic hash linking between blocks
3. Deterministic serialization for hashing
4. Merkle root for transaction integrity
5. Chain validation to detect tampering

Data Structure:
- Block: Header + Transactions + Hash
- BlockHeader: Index, Timestamp, PrevHash, MerkleRoot, Nonce
- Chain: Linked list of blocks via hash pointers

Time/Space Complexity:
- NewBlock: O(n) where n = number of transactions (hash each)
- ValidateChain: O(m*n) where m = chain length, n = avg transactions
- Space: O(m*n) for storing entire chain

Why Go is well-suited:
- crypto/sha256: Built-in cryptographic hashing
- encoding/hex: Easy hex encoding for hashes
- bytes.Buffer: Efficient serialization
- encoding/binary: Deterministic binary encoding
- Strong typing: Prevents serialization mistakes

Compared to other languages:
- Python: hashlib for hashing, but slower and less type-safe
- Rust: Similar primitives, but more complex lifetimes
- JavaScript: Web Crypto API, but less suitable for backend
- C++: Fast but error-prone manual memory management

Blockchain concepts:
- Hash chain: Each block hashes previous block's hash
- Tamper evidence: Changing any block breaks chain
- Genesis block: First block with no predecessor
- Merkle root: Commit to all transactions with one hash
- Serialization: Must be deterministic for consensus
*/

package exercise

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

// BlockHeader contains metadata about a block.
type BlockHeader struct {
	Index      int    // Block number (position in chain)
	Timestamp  int64  // Unix timestamp (seconds since epoch)
	PrevHash   string // Hash of previous block (hex string)
	MerkleRoot string // Hash of all transactions (hex string)
	Nonce      int    // For proof-of-work (unused in this project)
}

// Block represents a block in the blockchain.
type Block struct {
	Header       BlockHeader // Block metadata
	Transactions []string    // Transaction data
	Hash         string      // This block's hash (hex string)
}

// Serialize converts the block to a byte slice for hashing.
// The serialization must be deterministic (same block → same bytes).
//
// Go Concepts:
// - bytes.Buffer: Efficient byte array builder
// - binary.Write: Deterministic binary encoding
// - binary.BigEndian: Standard byte order (network byte order)
func (b *Block) Serialize() []byte {
	var buf bytes.Buffer

	// Serialize header fields (NOT including the block hash itself)
	binary.Write(&buf, binary.BigEndian, int64(b.Header.Index))
	binary.Write(&buf, binary.BigEndian, b.Header.Timestamp)
	buf.WriteString(b.Header.PrevHash)
	buf.WriteString(b.Header.MerkleRoot)
	binary.Write(&buf, binary.BigEndian, int64(b.Header.Nonce))

	// Serialize transactions
	for _, tx := range b.Transactions {
		buf.WriteString(tx)
	}

	return buf.Bytes()
}

// ComputeHash calculates the SHA-256 hash of the block.
// Returns hex-encoded hash string.
//
// Go Concepts:
// - sha256.Sum256: Cryptographic hash function (returns [32]byte)
// - hex.EncodeToString: Convert bytes to hex string
// - [:]  slice operator to convert array to slice
func (b *Block) ComputeHash() string {
	data := b.Serialize()
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// ComputeMerkleRoot calculates the merkle root of transactions.
// This is a simplified version - just hashes all transactions concatenated.
//
// Production version would build a proper Merkle tree (see Project 40).
//
// Go Concepts:
// - String concatenation: Simple but inefficient for many strings
// - Better alternative: Use strings.Builder or bytes.Buffer
func ComputeMerkleRoot(transactions []string) string {
	if len(transactions) == 0 {
		return ""
	}

	// Concatenate all transactions
	var data string
	for _, tx := range transactions {
		data += tx
	}

	// Hash the concatenated data
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// NewGenesisBlock creates the first block in the blockchain.
// The genesis block has Index 0 and PrevHash of all zeros.
//
// Go Concepts:
// - strings.Repeat: Create string of repeated characters
// - time.Now().Unix(): Get current Unix timestamp
// - Pointer return: Return *Block instead of Block for consistency
func NewGenesisBlock() *Block {
	block := &Block{
		Header: BlockHeader{
			Index:     0,
			Timestamp: time.Now().Unix(),
			PrevHash:  strings.Repeat("0", 64), // 64 zeros (SHA-256 is 32 bytes = 64 hex chars)
			Nonce:     0,
		},
		Transactions: []string{"Genesis Block"},
	}

	// Compute merkle root from transactions
	block.Header.MerkleRoot = ComputeMerkleRoot(block.Transactions)

	// Compute block hash (must be done AFTER merkle root is set)
	block.Hash = block.ComputeHash()

	return block
}

// NewBlock creates a new block linked to the previous block.
//
// Go Concepts:
// - Struct initialization with field names
// - Method chaining: block.Header.Index
// - Pointer receiver: Modifying prevBlock would require pointer
func NewBlock(prevBlock *Block, transactions []string) *Block {
	block := &Block{
		Header: BlockHeader{
			Index:     prevBlock.Header.Index + 1,
			Timestamp: time.Now().Unix(),
			PrevHash:  prevBlock.Hash, // Hash linking!
			Nonce:     0,
		},
		Transactions: transactions,
	}

	// Compute merkle root from transactions
	block.Header.MerkleRoot = ComputeMerkleRoot(block.Transactions)

	// Compute block hash (must be done AFTER merkle root is set)
	block.Hash = block.ComputeHash()

	return block
}

// ValidateChain validates the entire blockchain.
// Checks:
// - Genesis block has Index 0
// - Each block's hash is correct (matches computed hash)
// - Each block's PrevHash matches previous block's Hash
// - Indexes are sequential (0, 1, 2, ...)
// - Timestamps are non-decreasing
// - Merkle roots are correct
//
// Go Concepts:
// - Error handling: Return nil for success, error for failure
// - fmt.Errorf: Format error messages with context
// - Early return: Return immediately on first error found
func ValidateChain(chain []*Block) error {
	if len(chain) == 0 {
		return errors.New("chain is empty")
	}

	// Validate genesis block
	if chain[0].Header.Index != 0 {
		return errors.New("genesis block must have index 0")
	}

	// Validate each block
	for i, block := range chain {
		// Check 1: Hash integrity (stored hash matches computed hash)
		computedHash := block.ComputeHash()
		if block.Hash != computedHash {
			return fmt.Errorf("block %d: hash mismatch (stored=%s, computed=%s)",
				i, block.Hash[:16], computedHash[:16])
		}

		// Check 2: Merkle root integrity
		computedMerkle := ComputeMerkleRoot(block.Transactions)
		if block.Header.MerkleRoot != computedMerkle {
			return fmt.Errorf("block %d: merkle root mismatch (stored=%s, computed=%s)",
				i, block.Header.MerkleRoot[:16], computedMerkle[:16])
		}

		// Checks for non-genesis blocks
		if i > 0 {
			prevBlock := chain[i-1]

			// Check 3: Hash linking (current PrevHash matches previous Hash)
			if block.Header.PrevHash != prevBlock.Hash {
				return fmt.Errorf("block %d: prev hash mismatch (block.PrevHash=%s, prev.Hash=%s)",
					i, block.Header.PrevHash[:16], prevBlock.Hash[:16])
			}

			// Check 4: Sequential indexes
			if block.Header.Index != prevBlock.Header.Index+1 {
				return fmt.Errorf("block %d: index not sequential (expected=%d, got=%d)",
					i, prevBlock.Header.Index+1, block.Header.Index)
			}

			// Check 5: Timestamp ordering (non-decreasing)
			if block.Header.Timestamp < prevBlock.Header.Timestamp {
				return fmt.Errorf("block %d: timestamp out of order (prev=%d, current=%d)",
					i, prevBlock.Header.Timestamp, block.Header.Timestamp)
			}
		}
	}

	return nil
}

/*
Alternatives & Trade-offs:

1. Merkle Tree Implementation:
   Current: Simple concatenation and hash
   Alternative: Proper binary Merkle tree
   Pros: Enables efficient proofs (SPV)
   Cons: More complex implementation
   When to use: Production blockchains (see Project 40)

2. Serialization Format:
   Current: Custom binary with bytes.Buffer
   Alternative: JSON, Protobuf, or RLP (Ethereum)
   Pros: JSON is human-readable; Protobuf is compact
   Cons: JSON not deterministic; Protobuf adds dependency
   When to use: JSON for debugging, Protobuf for production

3. Hash Function:
   Current: SHA-256
   Alternative: SHA-3, Blake2, or other
   Pros: SHA-3 is quantum-resistant; Blake2 is faster
   Cons: SHA-256 is most widely adopted
   When to use: SHA-256 for compatibility; Blake2 for speed

4. Timestamp Validation:
   Current: Just check non-decreasing
   Alternative: Use median-time-past (Bitcoin) or network time
   Pros: More resistant to timestamp manipulation
   Cons: More complex; requires knowing network time
   When to use: Production blockchains with consensus

5. Block Storage:
   Current: In-memory slice
   Alternative: Database (LevelDB, BoltDB) or file system
   Pros: Persistent; handles large chains
   Cons: More complex; I/O overhead
   When to use: Production systems

Real-world blockchain differences:

Bitcoin:
- Uses double SHA-256 (SHA-256 of SHA-256)
- Merkle tree for transactions (enables SPV wallets)
- Proof-of-work in block hash (leading zeros)
- More compact binary serialization
- Stores UTXO set separately

Ethereum:
- Uses Keccak-256 (SHA-3 variant)
- Three Merkle trees: transactions, receipts, state
- Account-based model (not UTXO)
- RLP serialization
- Stores entire world state

Our implementation:
- Single SHA-256
- Simplified merkle root (just concatenate + hash)
- No proof-of-work (Nonce unused)
- Simple binary serialization
- In-memory storage only

Go vs X for blockchain:

Go vs Rust:
  Rust: More performance, memory safety guarantees
  Go: Simpler, faster development, better tooling
  Example: Rust is used for high-performance nodes (Parity)
           Go is used for most blockchain infrastructure

Go vs C++:
  C++: Maximum performance, used in Bitcoin Core
  Go: Safer, more productive, used in Ethereum, Hyperledger
  Example: C++ for consensus-critical code
           Go for everything else

Go vs JavaScript:
  JS: Good for web3 frontends, not suitable for nodes
  Go: Production-grade blockchain nodes
  Example: JS for dApps, Go for blockchain backend

Why Go won for blockchain infrastructure:
1. Fast compilation → quick iteration
2. Strong typing → fewer bugs
3. Built-in concurrency → handle P2P networking
4. Good standard library → crypto, networking, encoding
5. Static binary → easy deployment
6. Growing ecosystem → many blockchain projects use Go
*/
