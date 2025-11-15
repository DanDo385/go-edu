//go:build solution
// +build solution

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
func BuildMerkleTree(data [][]byte) *MerkleTree {
	// MICRO-COMMENT: Handle empty data case
	if len(data) == 0 {
		emptyHash := hash([]byte{})
		return &MerkleTree{
			Root:   emptyHash,
			Leaves: [][]byte{},
			Levels: [][][]byte{{emptyHash}},
		}
	}

	// MACRO-COMMENT: Build tree bottom-up, level by level
	// Level 0: Hash all data blocks to create leaves
	// Each subsequent level: Hash pairs of children to create parents
	// Continue until only one hash remains (the root)

	// MICRO-COMMENT: Create leaf nodes by hashing each data block
	leaves := make([][]byte, len(data))
	for i, d := range data {
		leaves[i] = hash(d)
	}

	// MICRO-COMMENT: Initialize levels with the leaf level
	levels := [][][]byte{leaves}
	currentLevel := leaves

	// MICRO-COMMENT: Build tree level by level until we reach the root
	for len(currentLevel) > 1 {
		nextLevel := [][]byte{}

		// MICRO-COMMENT: If odd number of nodes, duplicate the last one
		// This is the standard approach used in Bitcoin
		if len(currentLevel)%2 == 1 {
			currentLevel = append(currentLevel, currentLevel[len(currentLevel)-1])
		}

		// MICRO-COMMENT: Pair nodes and hash them to create parent nodes
		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			right := currentLevel[i+1]

			// MICRO-COMMENT: Parent hash = hash(left || right)
			// We concatenate the bytes and hash them
			parentHash := hash(append(left, right...))
			nextLevel = append(nextLevel, parentHash)
		}

		// MICRO-COMMENT: Store this level and move up the tree
		levels = append(levels, nextLevel)
		currentLevel = nextLevel
	}

	return &MerkleTree{
		Root:   currentLevel[0],
		Leaves: leaves,
		Levels: levels,
	}
}

// GenerateProof creates a proof of inclusion for the data at the given index.
func GenerateProof(tree *MerkleTree, index int) *MerkleProof {
	// MICRO-COMMENT: Validate index bounds
	if index < 0 || index >= len(tree.Leaves) {
		return nil
	}

	// MICRO-COMMENT: Initialize proof structure
	proof := &MerkleProof{
		LeafIndex: index,
		Siblings:  []ProofNode{},
	}

	// MACRO-COMMENT: Traverse from leaf to root, collecting sibling hashes
	// At each level, we need the sibling of the current node to compute
	// the parent hash. Track whether the sibling is left or right.

	currentIndex := index

	// MICRO-COMMENT: Iterate through all levels except the root
	for level := 0; level < len(tree.Levels)-1; level++ {
		currentLevelNodes := tree.Levels[level]

		// MICRO-COMMENT: Handle odd number of nodes (duplicated last node)
		// We need to duplicate here to correctly compute sibling indices
		if len(currentLevelNodes)%2 == 1 {
			currentLevelNodes = append(currentLevelNodes, currentLevelNodes[len(currentLevelNodes)-1])
		}

		// MICRO-COMMENT: Find sibling index
		// If current is even (left child), sibling is index+1 (right)
		// If current is odd (right child), sibling is index-1 (left)
		var siblingIndex int
		if currentIndex%2 == 0 {
			siblingIndex = currentIndex + 1
		} else {
			siblingIndex = currentIndex - 1
		}

		// MICRO-COMMENT: Add sibling to proof
		siblingHash := currentLevelNodes[siblingIndex]
		isLeft := siblingIndex < currentIndex

		proof.Siblings = append(proof.Siblings, ProofNode{
			Hash:   siblingHash,
			IsLeft: isLeft,
		})

		// MICRO-COMMENT: Move to parent index for next level
		// Parent index is always currentIndex / 2 (integer division)
		currentIndex = currentIndex / 2
	}

	return proof
}

// VerifyProof verifies that data is in the tree by checking the proof against the root.
func VerifyProof(data []byte, proof *MerkleProof, root []byte) bool {
	// MICRO-COMMENT: Handle nil proof
	if proof == nil {
		return false
	}

	// MACRO-COMMENT: Reconstruct the root hash from the data and proof
	// Start with the hash of the data, then iteratively combine with
	// siblings according to their position (left or right)

	// MICRO-COMMENT: Start with hash of the data (leaf hash)
	currentHash := hash(data)

	// MICRO-COMMENT: Combine with each sibling to move up the tree
	for _, sibling := range proof.Siblings {
		if sibling.IsLeft {
			// MICRO-COMMENT: Sibling goes on the left
			currentHash = hash(append(sibling.Hash, currentHash...))
		} else {
			// MICRO-COMMENT: Sibling goes on the right
			currentHash = hash(append(currentHash, sibling.Hash...))
		}
	}

	// MICRO-COMMENT: Compare computed root with expected root
	return hashesEqual(currentHash, root)
}

// GetMerkleRoot is a convenience function that builds a tree and returns the root.
func GetMerkleRoot(data [][]byte) []byte {
	// MICRO-COMMENT: Build tree and extract root
	tree := BuildMerkleTree(data)
	return tree.Root
}

// ========================================
// UTILITY FUNCTIONS
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
