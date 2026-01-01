//go:build !solution && !reference

package merkletreebasics

import (
	"crypto/sha256"
)

func BuildMerkleTree(data [][]byte) *MerkleTree {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func GenerateProof(tree *MerkleTree, index int) *MerkleProof {
	// TODO: Implement this function
	panic("not implemented")
}

func VerifyProof(data []byte, proof *MerkleProof, root []byte) bool {
	// TODO: Implement this function
	panic("not implemented")
}

func GetMerkleRoot(data [][]byte) []byte {
	// TODO: Implement this function
	panic("not implemented")
}

func hash(data []byte) []byte {
	// TODO: Implement this function
	panic("not implemented")
}

func hashesEqual(a, b []byte) bool {
	// TODO: Implement this function
	panic("not implemented")
}
