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
func (b *Block) Serialize() []byte {
	// TODO: Implement this method.

	// For a hash to be consistent, the input data must be serialized into bytes in a deterministic way.
	// The order and format of the data must be the same every time.

	// Step 1: Create a `bytes.Buffer` to build the byte slice.
	// Step 2: Write the block header fields to the buffer.
	//   - Use `binary.Write` for fixed-size data like integers. It's important to specify a consistent byte order, like `binary.BigEndian`.
	//   - `binary.Write(&buf, binary.BigEndian, int64(b.Header.Index))`
	//   - Use `buf.WriteString` for string fields like `PrevHash` and `MerkleRoot`.
	// Step 3: Write the transaction data to the buffer.
	//   - Loop through the `Transactions` and write each one to the buffer.
	// Step 4: Return the buffer's contents: `buf.Bytes()`.
	return nil
}

// ComputeHash calculates the SHA-256 hash of the block.
func (b *Block) ComputeHash() string {
	// TODO: Implement this method.

	// Step 1: Get the deterministic byte representation of the block by calling `b.Serialize()`.
	// Step 2: Compute the SHA-256 hash of the serialized data.
	//   - `hash := sha256.Sum256(data)`
	//   - This returns a `[32]byte` array.
	// Step 3: Convert the hash array into a hex-encoded string.
	//   - `return hex.EncodeToString(hash[:])`
	//   - The `[:]` is important to convert the `[32]byte` array into a `[]byte` slice, which `EncodeToString` expects.
	return ""
}

// ComputeMerkleRoot calculates the merkle root of transactions.
func ComputeMerkleRoot(transactions []string) string {
	// TODO: Implement this simplified Merkle root calculation.

	// A real Merkle root is a tree of hashes. For this exercise, we will use a much simpler (but less efficient) approach:
	// just concatenating all transaction data and hashing the result.

	// Step 1: Handle the case of zero transactions. Return an empty string.
	// Step 2: Concatenate all transaction strings into a single string.
	//   - A `strings.Builder` is more efficient for this than `+=` in a loop.
	// Step 3: Compute the SHA-256 hash of the concatenated string's bytes.
	// Step 4: Return the hex-encoded hash.
	return ""
}

// NewGenesisBlock creates the first block in the blockchain.
func NewGenesisBlock() *Block {
	// TODO: Implement this function.

	// The genesis block is the first block in the chain and is created specially.

	// Step 1: Create a `Block` struct.
	// - `Header.Index` should be 0.
	// - `Header.Timestamp` should be the current time (`time.Now().Unix()`).
	// - `Header.PrevHash` should be a string of 64 zeros, as there is no previous block. `strings.Repeat("0", 64)` is good for this.
	// - `Transactions` should be a slice containing a single string, like `"Genesis Block"`.

	// Step 2: Calculate the Merkle root for the transactions and set it in the header.
	// - `block.Header.MerkleRoot = ComputeMerkleRoot(block.Transactions)`

	// Step 3: Calculate the block's own hash and set it. This must be done *after* all other header fields are set.
	// - `block.Hash = block.ComputeHash()`

	// Step 4: Return the block.
	return nil
}

// NewBlock creates a new block linked to the previous block.
func NewBlock(prevBlock *Block, transactions []string) *Block {
	// TODO: Implement this function.

	// This function creates a new block that cryptographically links to the previous one.

	// Step 1: Create a `Block` struct.
	// - `Header.Index` should be `prevBlock.Header.Index + 1`.
	// - `Header.Timestamp` should be the current time.
	// - `Header.PrevHash` should be the hash of the `prevBlock`. This is the crucial link in the chain.
	// - `Transactions` should be the provided `transactions`.

	// Step 2: Calculate the Merkle root and set it in the header.
	// Step 3: Calculate the block's own hash and set it.
	// Step 4: Return the block.
	return nil
}

// ValidateChain validates the entire blockchain.
func ValidateChain(chain []*Block) error {
	// TODO: Implement this function to validate the integrity of the entire chain.

	// Step 1: Handle an empty chain. It's not valid.

	// Step 2: Validate the Genesis Block.
	// - The first block (`chain[0]`) must have an `Index` of 0.

	// Step 3: Loop through the rest of the chain (from index 1).
	// - For each `block` at index `i`, you'll need the `prevBlock` at index `i-1`.

	// Step 4: Inside the loop, perform the following checks for each block:
	//   a. **Hash Integrity**: Re-compute the block's hash (`block.ComputeHash()`) and check if it matches the `block.Hash` field. If not, the block's data has been tampered with.
	//   b. **Merkle Root Integrity**: Re-compute the Merkle root (`ComputeMerkleRoot(block.Transactions)`) and check if it matches `block.Header.MerkleRoot`.
	//   c. **Chain Linking**: Check if `block.Header.PrevHash` is equal to `prevBlock.Hash`. This is the core of the chain's security.
	//   d. **Index Order**: Check if `block.Header.Index` is equal to `prevBlock.Header.Index + 1`.
	//   e. **Timestamp Order**: Check if the block's timestamp is greater than or equal to the previous block's timestamp.

	// Step 5: If any check fails, return a specific, descriptive `error` immediately.
	// - `return fmt.Errorf("block %d: hash mismatch", i)`

	// Step 6: If the loop completes without finding any errors, return `nil`.
	return nil
}
