//go:build !solution && !reference

package simpleblockstructhashing

/*
Problem: Implement a simple blockchain with blocks, hashing, and validation
Requirements:
1. Block structure with header and transactions
2. Cryptographic hash linking between blocks
3. Deterministic serialization for hashing
Time/Space Complexity:
- NewBlock: O(n) where n = number of transactions (hash each)
- ValidateChain: O(m*n) where m = chain length, n = avg transactions
- Space: O(m*n) for storing entire chain
*/

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

type BlockHeader struct {
	Index      int    // Block number (position in chain)
	Timestamp  int64  // Unix timestamp (seconds since epoch)
	PrevHash   string // Hash of previous block (hex string)
	MerkleRoot string // Hash of all transactions (hex string)
	Nonce      int    // For proof-of-work (unused in this project)
}

type Block struct {
	Header       BlockHeader // Block metadata
	Transactions []string    // Transaction data
	Hash         string      // This block's hash (hex string)
}

// Serialize - TODO: implement this function
func (b *Block) Serialize() []byte {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// ComputeHash - TODO: implement this function
func (b *Block) ComputeHash() string {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// ComputeMerkleRoot - TODO: implement this function
func ComputeMerkleRoot(transactions []string) string {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return ""
}

// NewGenesisBlock - TODO: implement this function
func NewGenesisBlock() *Block {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// NewBlock - TODO: implement this function
func NewBlock(prevBlock *Block, transactions []string) *Block {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// ValidateChain - TODO: implement this function
func ValidateChain(chain []*Block) error {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

