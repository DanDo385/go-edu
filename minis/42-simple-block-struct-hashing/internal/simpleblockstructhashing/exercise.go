//go:build !solution && !reference

package simpleblockstructhashing



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
	// TODO: Implement this function
	panic("unimplemented")
}

// ComputeHash calculates the SHA-256 hash of the block.
// Returns hex-encoded hash string.
//
// Go Concepts:
// - sha256.Sum256: Cryptographic hash function (returns [32]byte)
// - hex.EncodeToString: Convert bytes to hex string
// - [:]  slice operator to convert array to slice
func (b *Block) ComputeHash() string {
	// TODO: Implement this function
	panic("unimplemented")
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
	// TODO: Implement this function
	panic("unimplemented")
}

// NewGenesisBlock creates the first block in the blockchain.
// The genesis block has Index 0 and PrevHash of all zeros.
//
// Go Concepts:
// - strings.Repeat: Create string of repeated characters
// - time.Now().Unix(): Get current Unix timestamp
// - Pointer return: Return *Block instead of Block for consistency
func NewGenesisBlock() *Block {
	// TODO: Implement this function
	panic("unimplemented")
}

// NewBlock creates a new block linked to the previous block.
//
// Go Concepts:
// - Struct initialization with field names
// - Method chaining: block.Header.Index
// - Pointer receiver: Modifying prevBlock would require pointer
func NewBlock(prevBlock *Block, transactions []string) *Block {
	// TODO: Implement this function
	panic("unimplemented")
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
	// TODO: Implement this function
	panic("unimplemented")
}


