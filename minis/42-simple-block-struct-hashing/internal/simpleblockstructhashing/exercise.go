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

// Serialize implements the exercise.
//
// TODO: Implement this function
func (b *Block) Serialize() []byte {
	// TODO: Implement
	return nil
}

// ComputeHash implements the exercise.
//
// TODO: Implement this function
func (b *Block) ComputeHash() string {
	// TODO: Implement
	return ""
}

// ComputeMerkleRoot implements the exercise.
//
// TODO: Implement this function
func ComputeMerkleRoot(transactions []string) string {
	// TODO: Implement
	return ""
}

// NewGenesisBlock implements the exercise.
//
// TODO: Implement this function
func NewGenesisBlock() *Block {
	// TODO: Implement
	return nil
}

// NewBlock implements the exercise.
//
// TODO: Implement this function
func NewBlock(prevBlock *Block, transactions []string) *Block {
	// TODO: Implement
	return nil
}

// ValidateChain implements the exercise.
//
// TODO: Implement this function
func ValidateChain(chain []*Block) error {
	// TODO: Implement
	return nil
}
