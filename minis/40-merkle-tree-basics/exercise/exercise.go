//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for understanding Merkle trees.
//
// LEARNING OBJECTIVES:
// - Build a Merkle tree from data blocks
// - Generate cryptographic proofs of inclusion
// - Verify proofs without accessing the full dataset
// - Understand blockchain-style data verification

package exercise

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
//
// REQUIREMENTS:
// - Hash each data block to create leaf nodes
// - Build parent nodes by hashing pairs of children
// - If odd number of nodes, duplicate the last one
// - Continue until only one hash remains (the root)
// - Store all levels for proof generation
//
// ALGORITHM:
// 1. Level 0 (leaves): Hash each data block
// 2. For each level with >1 node:
//    a. If odd count, duplicate last node
//    b. Pair nodes and hash: parent = hash(left + right)
//    c. Create next level with parent hashes
// 3. Root is the single hash at the top level
//
// EXAMPLE:
//   data = ["A", "B", "C", "D"]
//
//   Level 0 (leaves): [H(A), H(B), H(C), H(D)]
//   Level 1:          [H(H(A)+H(B)), H(H(C)+H(D))]
//   Level 2 (root):   [H(H(H(A)+H(B))+H(H(C)+H(D)))]
//
// EDGE CASES:
// - Empty data: Return tree with hash of empty byte slice as root
// - Single block: The hash of that block IS the root
//
// HINT: Use sha256.Sum256() for hashing. When hashing pairs, concatenate bytes.
func BuildMerkleTree(data [][]byte) *MerkleTree {
	// TODO: Implement this function
	//
	// Steps:
	// 1. Handle empty data case
	// 2. Create leaf nodes by hashing each data block
	// 3. Build tree level by level
	// 4. Handle odd number of nodes at each level
	// 5. Return MerkleTree with Root, Leaves, and Levels populated

	return nil
}

// GenerateProof creates a proof of inclusion for the data at the given index.
//
// REQUIREMENTS:
// - Collect all sibling hashes from leaf to root
// - Track whether each sibling is left or right
// - Return a proof that can reconstruct the root
//
// ALGORITHM:
// 1. Start at leaf index
// 2. For each level from bottom to top:
//    a. Find sibling index (if index is even, sibling is index+1; if odd, index-1)
//    b. Add sibling hash and position to proof
//    c. Move to parent index (index / 2)
// 3. Return the proof
//
// EXAMPLE:
//   Tree with 4 leaves [H1, H2, H3, H4], proving index 1 (H2):
//
//   Level 0: Need H1 (sibling of H2, IsLeft=true)
//   Level 1: Need H(H3+H4) (sibling of H(H1+H2), IsLeft=false)
//
//   Proof: [{H1, true}, {H(H3+H4), false}]
//
// EDGE CASES:
// - Invalid index: Return nil
// - Single leaf: Empty proof (data hash == root)
//
// HINT: Track whether sibling is left or right by comparing indices.
// The sibling is on the left if its index < current index.
func GenerateProof(tree *MerkleTree, index int) *MerkleProof {
	// TODO: Implement this function
	//
	// Steps:
	// 1. Validate index
	// 2. Initialize proof with LeafIndex
	// 3. Traverse from leaf to root
	// 4. At each level, find sibling and determine if it's left or right
	// 5. Add sibling to proof
	// 6. Move to parent level

	return nil
}

// VerifyProof verifies that data is in the tree by checking the proof against the root.
//
// REQUIREMENTS:
// - Hash the data to get the leaf hash
// - Combine with siblings according to their position
// - Compare final hash with expected root
// - Return true if they match, false otherwise
//
// ALGORITHM:
// 1. Compute hash of data
// 2. For each sibling in proof:
//    a. If sibling is left: hash = H(sibling + current)
//    b. If sibling is right: hash = H(current + sibling)
//    c. Update current hash
// 3. Compare final hash with root
//
// EXAMPLE:
//   Data: "B"
//   Proof: [{H(A), true}, {H(H(C)+H(D)), false}]
//   Root: H(H(H(A)+H(B))+H(H(C)+H(D)))
//
//   Step 1: current = H(B)
//   Step 2a: current = H(H(A) + H(B))  [sibling is left]
//   Step 2b: current = H(H(H(A)+H(B)) + H(H(C)+H(D)))  [sibling is right]
//   Step 3: current == root? → true
//
// SECURITY NOTE:
// Order matters! hash(A+B) ≠ hash(B+A), so we must respect IsLeft flag.
//
// HINT: Use a loop to combine current hash with each sibling.
func VerifyProof(data []byte, proof *MerkleProof, root []byte) bool {
	// TODO: Implement this function
	//
	// Steps:
	// 1. Handle nil proof
	// 2. Hash the data
	// 3. Iteratively combine with siblings
	// 4. Compare final hash with root

	return false
}

// GetMerkleRoot is a convenience function that builds a tree and returns the root.
//
// REQUIREMENTS:
// - Build a Merkle tree from the data
// - Return only the root hash
//
// HINT: Use BuildMerkleTree and extract the Root field.
func GetMerkleRoot(data [][]byte) []byte {
	// TODO: Implement this function
	//
	// This is a simple wrapper around BuildMerkleTree.

	return nil
}

// ========================================
// UTILITY FUNCTIONS (PROVIDED)
// ========================================

// hash computes SHA-256 hash of data
func hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// hashesEqual compares two hashes for equality
func hashesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
