//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package merkletreebasics

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
// TODO: implement BuildMerkleTree.
func BuildMerkleTree(data [][]byte) *MerkleTree { panic("TODO: implement") }
// TODO: implement GenerateProof.
func GenerateProof(tree *MerkleTree, index int) *MerkleProof { panic("TODO: implement") }
// TODO: implement VerifyProof.
func VerifyProof(data []byte, proof *MerkleProof, root []byte) bool { panic("TODO: implement") }
// TODO: implement GetMerkleRoot.
func GetMerkleRoot(data [][]byte) []byte { panic("TODO: implement") }
// TODO: implement hash.
func hash(data []byte) []byte { panic("TODO: implement") }
// TODO: implement hashesEqual.
func hashesEqual(a, b []byte) bool { panic("TODO: implement") }
