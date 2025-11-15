//go:build !solution
// +build !solution

package exercise

// BlockHeader contains metadata about a block.
type BlockHeader struct {
	Index      int    // Block number (position in chain)
	Timestamp  int64  // Unix timestamp (seconds since epoch)
	PrevHash   string // Hash of previous block (hex string)
	MerkleRoot string // Hash of all transactions (hex string)
	Nonce      int    // For proof-of-work (unused in this project)
}

// Block represents a block in the blockchain.
type Block struct {
	Header       BlockHeader // Block metadata
	Transactions []string    // Transaction data
	Hash         string      // This block's hash (hex string)
}

// Serialize converts the block to a byte slice for hashing.
// The serialization must be deterministic (same block → same bytes).
//
// TODO: Implement serialization
// Hint: Use bytes.Buffer and binary.Write with binary.BigEndian
// Include: Index, Timestamp, PrevHash, MerkleRoot, Nonce, Transactions
func (b *Block) Serialize() []byte {
	// TODO: implement
	return nil
}

// ComputeHash calculates the SHA-256 hash of the block.
// Returns hex-encoded hash string.
//
// TODO: Implement hash computation
// Hint: Serialize the block, compute SHA-256, encode as hex
func (b *Block) ComputeHash() string {
	// TODO: implement
	return ""
}

// ComputeMerkleRoot calculates the merkle root of transactions.
// This is a simplified version - just hashes all transactions concatenated.
//
// TODO: Implement merkle root computation
// Hint: Concatenate all transactions, hash with SHA-256, encode as hex
func ComputeMerkleRoot(transactions []string) string {
	// TODO: implement
	return ""
}

// NewGenesisBlock creates the first block in the blockchain.
// The genesis block has Index 0 and PrevHash of all zeros.
//
// TODO: Implement genesis block creation
// Hint:
// 1. Create block with Index=0, Timestamp=now, PrevHash=all zeros
// 2. Add "Genesis Block" as the only transaction
// 3. Compute merkle root
// 4. Compute block hash
func NewGenesisBlock() *Block {
	// TODO: implement
	return nil
}

// NewBlock creates a new block linked to the previous block.
//
// TODO: Implement new block creation
// Hint:
// 1. Index = previous index + 1
// 2. Timestamp = current time
// 3. PrevHash = previous block's hash
// 4. Transactions = provided transactions
// 5. Compute merkle root
// 6. Compute block hash
func NewBlock(prevBlock *Block, transactions []string) *Block {
	// TODO: implement
	return nil
}

// ValidateChain validates the entire blockchain.
// Checks:
// - Genesis block has Index 0
// - Each block's hash is correct (matches computed hash)
// - Each block's PrevHash matches previous block's Hash
// - Indexes are sequential (0, 1, 2, ...)
// - Timestamps are non-decreasing
// - Merkle roots are correct
//
// TODO: Implement chain validation
// Returns nil if chain is valid, error describing the problem if invalid.
func ValidateChain(chain []*Block) error {
	// TODO: implement
	return nil
}
