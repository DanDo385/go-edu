//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package simpleblockstructhashing

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
// TODO: implement Serialize.
func (b *Block) Serialize() []byte { panic("TODO: implement") }
// TODO: implement ComputeHash.
func (b *Block) ComputeHash() string { panic("TODO: implement") }
// TODO: implement ComputeMerkleRoot.
func ComputeMerkleRoot(transactions []string) string { panic("TODO: implement") }
// TODO: implement NewGenesisBlock.
func NewGenesisBlock() *Block { panic("TODO: implement") }
// TODO: implement NewBlock.
func NewBlock(prevBlock *Block, transactions []string) *Block { panic("TODO: implement") }
// TODO: implement ValidateChain.
func ValidateChain(chain []*Block) error { panic("TODO: implement") }
