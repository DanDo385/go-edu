//go:build !solution && !reference

package merkletreebasics

import "crypto/sha256"

type MerkleTree struct {
	Root   []byte     // The root hash representing the entire tree
	Leaves [][]byte   // All leaf hashes (hashes of original data)
	Levels [][][]byte // All levels of the tree (for proof generation)
}

type ProofNode struct {
	Hash   []byte // Hash value of the sibling node
	IsLeft bool   // True if this hash should be on the left when combining
}

type MerkleProof struct {
	LeafIndex int         // Index of the data in the original dataset
	Siblings  []ProofNode // Sibling hashes needed to reconstruct the root
}

// BuildMerkleTree implements the exercise.
//
// TODO: Implement this function
func BuildMerkleTree(data [][]byte) *MerkleTree {
	// TODO: Implement
	return nil
}

// GenerateProof implements the exercise.
//
// TODO: Implement this function
func GenerateProof(tree *MerkleTree, index int) *MerkleProof {
	// TODO: Implement
	return nil
}

// VerifyProof implements the exercise.
//
// TODO: Implement this function
func VerifyProof(data []byte, proof *MerkleProof, root []byte) bool {
	// TODO: Implement
	return false
}

// GetMerkleRoot implements the exercise.
//
// TODO: Implement this function
func GetMerkleRoot(data [][]byte) []byte {
	// TODO: Implement
	return nil
}
