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

func (b *Block) Serialize() []byte {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func (b *Block) ComputeHash() string {
	// TODO: Implement this function
	panic("not implemented")
}

func ComputeMerkleRoot(transactions []string) string {
	// TODO: Implement this function
	panic("not implemented")
}

func NewGenesisBlock() *Block {
	// TODO: Implement this function
	panic("not implemented")
}

func NewBlock(prevBlock *Block, transactions []string) *Block {
	// TODO: Implement this function
	panic("not implemented")
}

func ValidateChain(chain []*Block) error {
	// TODO: Implement this function
	panic("not implemented")
}
