//go:build !solution && !reference

package merkletreebasics

import (
	"crypto/sha256"
)

// MerkleTree represents a complete Merkle tree structure
type MerkleTree struct {
	Root   []byte     // The root hash representing the entire tree
	Leaves [][]byte   // All leaf hashes (hashes of original data)
	Levels [][][]byte // All levels of the tree (for proof generation)
}

// ProofNode represents one element in a Merkle proof
type ProofNode struct {
	Hash   []byte // Hash value of the sibling node
	IsLeft bool   // True if this hash should be on the left when combining
}

// MerkleProof represents a proof of inclusion for a specific data block
type MerkleProof struct {
	LeafIndex int         // Index of the data in the original dataset
	Siblings  []ProofNode // Sibling hashes needed to reconstruct the root
}

// BuildMerkleTree constructs a Merkle tree from data blocks.
func BuildMerkleTree(data [][]byte) *MerkleTree {
	// TODO: Implement this function
	panic("unimplemented")
}

// GenerateProof creates a proof of inclusion for the data at the given index.
func GenerateProof(tree *MerkleTree, index int) *MerkleProof {
	// TODO: Implement this function
	panic("unimplemented")
}

// VerifyProof verifies that data is in the tree by checking the proof against the root.
func VerifyProof(data []byte, proof *MerkleProof, root []byte) bool {
	// TODO: Implement this function
	panic("unimplemented")
}

// GetMerkleRoot is a convenience function that builds a tree and returns the root.
func GetMerkleRoot(data [][]byte) []byte {
	// TODO: Implement this function
	panic("unimplemented")
}

// ========================================
// UTILITY FUNCTIONS
// ========================================

// hash computes SHA-256 hash of data
func hash(data []byte) []byte {
	// TODO: Implement this function
	panic("unimplemented")
}

// hashesEqual compares two hashes for equality
func hashesEqual(a, b []byte) bool {
	// TODO: Implement this function
	panic("unimplemented")
}
