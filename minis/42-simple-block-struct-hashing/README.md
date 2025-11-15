# Project 42: Simple Block Struct Hashing

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a system that needs an **immutable audit log** - a record of events that can never be changed or tampered with once written. This is the core problem that blockchain technology solves.

**Examples needing immutable records:**
- **Financial transactions**: Bank transfers, stock trades, cryptocurrency payments
- **Medical records**: Patient history, prescriptions, lab results
- **Supply chain**: Product origin, shipping events, custody transfers
- **Voting systems**: Cast ballots, vote counts, election results

**Traditional approaches fail:**
- **Centralized database**: Admin can modify records
- **Append-only file**: File can be edited or corrupted
- **Signed documents**: Signatures can be forged or stripped

**Blockchain solution**: Link blocks together using cryptographic hashes, making tampering mathematically detectable.

This project teaches you how to build the fundamental building block (pun intended) of blockchain technology: **chained, hashed blocks** that form an immutable ledger.

### What You'll Learn

1. **Block structure**: Headers, data payload, metadata
2. **Cryptographic hashing**: SHA-256 for block fingerprinting
3. **Hash linking**: Previous block hash creates chain
4. **Serialization**: Converting structs to bytes for hashing
5. **Chain validation**: Detecting tampering by verifying hashes
6. **Timestamps**: Ordering blocks chronologically
7. **Merkle roots**: (Simplified) Data integrity within blocks

### The Challenge

Build a blockchain block system that:
- Creates blocks with headers and transaction data
- Links blocks using previous block's hash
- Serializes blocks deterministically for hashing
- Validates chain integrity (detects tampering)
- Handles timestamps correctly
- Uses SHA-256 for cryptographic security

---

## 2. First Principles: Understanding Blockchain Blocks

### What is a Block?

A **block** is a container for data with a cryptographic fingerprint (hash) that depends on:
1. The block's own data
2. The previous block's hash

**Analogy**: Imagine a journal where each new page references the previous page's content. If someone changes page 5, you'd know because page 6 references the old content of page 5.

**Block structure**:
```
┌─────────────────────────────────────┐
│           BLOCK HEADER              │
├─────────────────────────────────────┤
│ Index:           1                  │
│ Timestamp:       2025-11-15 10:30   │
│ Previous Hash:   0xabc123...        │
│ Merkle Root:     0xdef456...        │
│ Nonce:           0                  │
│ Hash:            0x789xyz...        │  ← Computed from all above
├─────────────────────────────────────┤
│           BLOCK DATA                │
├─────────────────────────────────────┤
│ Transaction 1: Alice → Bob $10      │
│ Transaction 2: Bob → Carol $5       │
│ Transaction 3: Carol → Dave $3      │
└─────────────────────────────────────┘
```

### What is a Block Header?

The **block header** contains metadata about the block. Think of it as the envelope of a letter - it tells you about the letter without revealing its contents.

**Essential header fields**:

1. **Index** (Block Number): Position in the chain (0, 1, 2, ...)
2. **Timestamp**: When the block was created (Unix timestamp)
3. **Previous Hash**: Hash of the previous block (the "link")
4. **Merkle Root**: Hash of all transactions (data integrity)
5. **Nonce**: Number used for proof-of-work (not used in this project, but important for mining)
6. **Hash**: This block's hash (computed from all above fields)

**Why separate header from data?**

- **Efficiency**: Can verify chain without downloading all transaction data
- **Light clients**: Mobile wallets can verify blocks using only headers
- **Storage**: Headers are ~80 bytes, blocks can be megabytes

### What is Hash Linking?

**Hash linking** creates a chain by including the previous block's hash in the current block.

**Visual example**:

```
Genesis Block (Block 0)
┌─────────────────────────┐
│ Index: 0                │
│ PrevHash: 0x000000...   │ ← No previous block
│ Data: "Genesis"         │
│ Hash: 0xabc123...       │ ← Computed from above
└─────────────────────────┘
            ↓
         (Hash: 0xabc123...)
            ↓
Block 1
┌─────────────────────────┐
│ Index: 1                │
│ PrevHash: 0xabc123...   │ ← Links to Block 0
│ Data: "Transaction 1"   │
│ Hash: 0xdef456...       │ ← Computed from above
└─────────────────────────┘
            ↓
         (Hash: 0xdef456...)
            ↓
Block 2
┌─────────────────────────┐
│ Index: 2                │
│ PrevHash: 0xdef456...   │ ← Links to Block 1
│ Data: "Transaction 2"   │
│ Hash: 0x789xyz...       │ ← Computed from above
└─────────────────────────┘
```

**What makes this tamper-proof?**

If someone changes Block 1's data:
1. Block 1's hash changes (hash depends on data)
2. Block 2's `PrevHash` field now points to wrong hash
3. Block 2's hash must be recomputed
4. Block 3's `PrevHash` field now points to wrong hash
5. ... cascade continues to the end of the chain

**Result**: Changing any historical block requires recomputing all subsequent blocks - computationally infeasible for long chains with proof-of-work.

### What is Serialization?

**Serialization** converts a struct into bytes in a consistent, deterministic way.

**Why needed?**

Cryptographic hash functions (like SHA-256) operate on bytes, not structs.

**Example**:

```go
type Block struct {
    Index     int
    Timestamp int64
    PrevHash  string
    Data      string
}

block := Block{
    Index:     1,
    Timestamp: 1731672600,
    PrevHash:  "abc123",
    Data:      "Transaction 1",
}

// Serialize to bytes
bytes := serialize(block)
// bytes = []byte{0x01, 0x00, 0x00, 0x00, 0x65, 0x3d, ...}

// Hash the bytes
hash := sha256.Sum256(bytes)
// hash = [32]byte{0x7a, 0x3b, 0x9f, ...}
```

**Critical requirement**: Serialization must be **deterministic**
- Same input always produces same bytes
- Order of fields matters
- No random padding or timestamps

**Common serialization methods**:
1. **JSON** (human-readable, but not deterministic due to whitespace/field order)
2. **Binary encoding** (compact, deterministic)
3. **Custom format** (full control, used in Bitcoin)

### What is a Timestamp?

A **timestamp** records when a block was created, proving the data existed at that time.

**Unix timestamp**: Seconds since January 1, 1970 00:00:00 UTC

```go
timestamp := time.Now().Unix()
// timestamp = 1731672600 (example: 2025-11-15 10:30:00 UTC)
```

**Why timestamps matter**:
- **Ordering**: Establishes chronological sequence
- **Proof of existence**: Data existed at this time
- **Validation**: Blocks shouldn't have timestamps in the future
- **Difficulty adjustment**: Used to calculate mining difficulty

**Timestamp validation rules**:
- Must be greater than previous block's timestamp
- Should not be more than 2 hours in the future (tolerance for clock drift)
- Cannot be in the distant past (prevents backdating)

### What is a Merkle Root?

A **Merkle root** is a single hash representing all transactions in a block.

**Merkle tree structure** (for 4 transactions):

```
         Root Hash (Merkle Root)
              /        \
            /            \
       Hash01           Hash23
        /  \             /  \
       /    \           /    \
    Hash0  Hash1    Hash2  Hash3
      |      |        |      |
    Tx0    Tx1      Tx2    Tx3
```

**How it's computed**:
1. Hash each transaction individually
2. Pair up hashes and hash the pairs
3. Repeat until one hash remains (the root)

**Why use Merkle trees?**
- **Efficiency**: Prove a transaction is in a block without downloading all transactions
- **Proof size**: O(log N) proof size instead of O(N)
- **Integrity**: Single hash represents all data

**Simplified version** (what we'll use):
Just hash all transaction data concatenated together:

```go
merkleRoot := sha256.Sum256([]byte(tx1 + tx2 + tx3 + ...))
```

---

## 3. Breaking Down the Solution

### Step 1: Define the Block Structure

**Block Header**:
```go
type BlockHeader struct {
    Index      int    // Block number (position in chain)
    Timestamp  int64  // Unix timestamp (seconds since epoch)
    PrevHash   string // Hash of previous block (hex string)
    MerkleRoot string // Hash of all transactions (hex string)
    Nonce      int    // For proof-of-work (0 for this project)
}
```

**Block**:
```go
type Block struct {
    Header       BlockHeader // Block metadata
    Transactions []string    // Transaction data
    Hash         string      // This block's hash (hex string)
}
```

**Why use strings for hashes?**
- Easy to print and debug
- Human-readable (hex encoding)
- Standard format: `"a1b2c3d4..."`

### Step 2: Serialize Block for Hashing

**Goal**: Convert block to bytes deterministically.

**Approach**: Concatenate all fields in fixed order:

```go
func (b *Block) Serialize() []byte {
    var buf bytes.Buffer

    // Write header fields
    binary.Write(&buf, binary.BigEndian, int64(b.Header.Index))
    binary.Write(&buf, binary.BigEndian, b.Header.Timestamp)
    buf.WriteString(b.Header.PrevHash)
    buf.WriteString(b.Header.MerkleRoot)
    binary.Write(&buf, binary.BigEndian, int64(b.Header.Nonce))

    // Write transactions
    for _, tx := range b.Transactions {
        buf.WriteString(tx)
    }

    return buf.Bytes()
}
```

**Key points**:
- Use `binary.BigEndian` for consistent byte order (big-endian is standard)
- Write fields in fixed order
- Include all fields that affect the hash

### Step 3: Compute Block Hash

**Algorithm**:
1. Serialize block to bytes
2. Compute SHA-256 hash
3. Encode hash as hex string

```go
func (b *Block) ComputeHash() string {
    data := b.Serialize()
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}
```

**Example output**:
```
"a1b2c3d4e5f6789012345678901234567890123456789012345678901234567890"
```

**Why SHA-256?**
- **Deterministic**: Same input always produces same hash
- **One-way**: Cannot reverse hash to get original data
- **Collision-resistant**: Nearly impossible to find two inputs with same hash
- **Fixed size**: Always 32 bytes (64 hex characters)
- **Fast**: Can compute millions of hashes per second

### Step 4: Compute Merkle Root

**Simplified approach** (concatenate all transactions and hash):

```go
func ComputeMerkleRoot(transactions []string) string {
    if len(transactions) == 0 {
        return ""
    }

    // Concatenate all transactions
    var data string
    for _, tx := range transactions {
        data += tx
    }

    // Hash the concatenated data
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

**Production version** would build a proper Merkle tree (see Project 40).

### Step 5: Create Genesis Block

The **genesis block** is the first block in the chain - it has no previous block.

```go
func NewGenesisBlock() *Block {
    block := &Block{
        Header: BlockHeader{
            Index:      0,
            Timestamp:  time.Now().Unix(),
            PrevHash:   "0000000000000000000000000000000000000000000000000000000000000000",
            MerkleRoot: "",
            Nonce:      0,
        },
        Transactions: []string{"Genesis Block"},
    }

    block.Header.MerkleRoot = ComputeMerkleRoot(block.Transactions)
    block.Hash = block.ComputeHash()

    return block
}
```

**Why all zeros for PrevHash?**
Convention - indicates no previous block exists.

### Step 6: Create New Block

**Algorithm**:
1. Create block with index = previous index + 1
2. Set previous hash = previous block's hash
3. Set timestamp = current time
4. Compute merkle root from transactions
5. Compute block hash

```go
func NewBlock(prevBlock *Block, transactions []string) *Block {
    block := &Block{
        Header: BlockHeader{
            Index:     prevBlock.Header.Index + 1,
            Timestamp: time.Now().Unix(),
            PrevHash:  prevBlock.Hash,
            Nonce:     0,
        },
        Transactions: transactions,
    }

    block.Header.MerkleRoot = ComputeMerkleRoot(block.Transactions)
    block.Hash = block.ComputeHash()

    return block
}
```

### Step 7: Validate Chain

**Validation checks**:

1. **Hash integrity**: Each block's stored hash matches its computed hash
2. **Hash linking**: Each block's PrevHash matches previous block's Hash
3. **Index continuity**: Indexes increment by 1
4. **Timestamp ordering**: Timestamps are non-decreasing
5. **Merkle root**: Matches computed root from transactions

```go
func ValidateChain(chain []*Block) error {
    if len(chain) == 0 {
        return errors.New("chain is empty")
    }

    // Validate genesis block
    if chain[0].Header.Index != 0 {
        return errors.New("genesis block must have index 0")
    }

    // Validate each block
    for i, block := range chain {
        // Check hash integrity
        if block.Hash != block.ComputeHash() {
            return fmt.Errorf("block %d: hash mismatch", i)
        }

        // Check hash linking (skip genesis)
        if i > 0 {
            if block.Header.PrevHash != chain[i-1].Hash {
                return fmt.Errorf("block %d: prev hash mismatch", i)
            }

            if block.Header.Index != chain[i-1].Header.Index + 1 {
                return fmt.Errorf("block %d: index not sequential", i)
            }

            if block.Header.Timestamp < chain[i-1].Header.Timestamp {
                return fmt.Errorf("block %d: timestamp out of order", i)
            }
        }

        // Check merkle root
        merkleRoot := ComputeMerkleRoot(block.Transactions)
        if block.Header.MerkleRoot != merkleRoot {
            return fmt.Errorf("block %d: merkle root mismatch", i)
        }
    }

    return nil
}
```

---

## 4. Complete Solution Walkthrough

### The Block Structure

```go
type BlockHeader struct {
    Index      int    // Position in chain
    Timestamp  int64  // Creation time (Unix timestamp)
    PrevHash   string // Previous block's hash
    MerkleRoot string // Hash of transactions
    Nonce      int    // For proof-of-work (unused here)
}

type Block struct {
    Header       BlockHeader
    Transactions []string
    Hash         string // This block's hash
}
```

**Design decisions**:

- **Why `int` for Index?** Blocks numbered 0, 1, 2, ... (int is sufficient)
- **Why `int64` for Timestamp?** Unix timestamps are int64 (seconds since 1970)
- **Why `string` for hashes?** Human-readable, easy to print/compare
- **Why `[]string` for Transactions?** Simple representation; production uses structs

### Serialization

```go
func (b *Block) Serialize() []byte {
    var buf bytes.Buffer

    // Serialize header
    binary.Write(&buf, binary.BigEndian, int64(b.Header.Index))
    binary.Write(&buf, binary.BigEndian, b.Header.Timestamp)
    buf.WriteString(b.Header.PrevHash)
    buf.WriteString(b.Header.MerkleRoot)
    binary.Write(&buf, binary.BigEndian, int64(b.Header.Nonce))

    // Serialize transactions
    for _, tx := range b.Transactions {
        buf.WriteString(tx)
    }

    return buf.Bytes()
}
```

**Why this approach?**

- **`bytes.Buffer`**: Efficient byte array builder
- **`binary.BigEndian`**: Standard byte order (network byte order)
- **Fixed order**: Same block always serializes to same bytes

**Alternative approaches**:
- **JSON**: `json.Marshal(b)` - not deterministic (field order varies)
- **Gob**: `gob.Encode(b)` - Go-specific, not cross-language
- **Protobuf**: `proto.Marshal(b)` - good for production, more complex

### Hash Computation

```go
func (b *Block) ComputeHash() string {
    data := b.Serialize()
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}
```

**Line-by-line**:

1. **Serialize**: Get deterministic byte representation
2. **SHA-256**: Compute 32-byte hash
3. **Hex encode**: Convert bytes to string (64 hex characters)

**Example**:
```
Input:  Block{Index: 1, Timestamp: 1731672600, ...}
Bytes:  [0x00, 0x00, 0x00, 0x01, 0x67, 0x3d, ...]
Hash:   [0xa1, 0xb2, 0xc3, 0xd4, ...]
Hex:    "a1b2c3d4e5f6..."
```

### Merkle Root

```go
func ComputeMerkleRoot(transactions []string) string {
    if len(transactions) == 0 {
        return ""
    }

    var data string
    for _, tx := range transactions {
        data += tx
    }

    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

**Simplified approach**: Just concatenate and hash all transactions.

**Production approach**: Build proper Merkle tree (see Project 40).

### Creating Blocks

```go
func NewGenesisBlock() *Block {
    block := &Block{
        Header: BlockHeader{
            Index:     0,
            Timestamp: time.Now().Unix(),
            PrevHash:  strings.Repeat("0", 64), // All zeros
            Nonce:     0,
        },
        Transactions: []string{"Genesis Block"},
    }

    block.Header.MerkleRoot = ComputeMerkleRoot(block.Transactions)
    block.Hash = block.ComputeHash()

    return block
}

func NewBlock(prevBlock *Block, transactions []string) *Block {
    block := &Block{
        Header: BlockHeader{
            Index:     prevBlock.Header.Index + 1,
            Timestamp: time.Now().Unix(),
            PrevHash:  prevBlock.Hash,
            Nonce:     0,
        },
        Transactions: transactions,
    }

    block.Header.MerkleRoot = ComputeMerkleRoot(block.Transactions)
    block.Hash = block.ComputeHash()

    return block
}
```

**Key patterns**:
- Genesis block has Index 0, PrevHash all zeros
- New blocks increment Index, link via PrevHash
- Always compute MerkleRoot before Hash (Hash depends on MerkleRoot)

---

## 5. Key Concepts Explained

### Concept 1: Cryptographic Hash Functions

A **cryptographic hash function** converts arbitrary data to a fixed-size fingerprint.

**Properties**:

1. **Deterministic**: Same input → same output
   ```go
   sha256.Sum256([]byte("hello")) // Always same result
   ```

2. **Fast**: Can hash gigabytes per second
   ```go
   // Hashing 1 MB takes ~1-2ms
   data := make([]byte, 1024*1024)
   hash := sha256.Sum256(data)
   ```

3. **One-way**: Cannot reverse
   ```go
   // Given hash, cannot find input
   hash := "a1b2c3..."  // No way to get original data
   ```

4. **Collision-resistant**: Hard to find two inputs with same hash
   ```go
   // Probability of collision: ~1 in 2^256
   // More atoms in universe: ~2^266
   ```

5. **Avalanche effect**: Small input change → completely different hash
   ```go
   sha256.Sum256([]byte("hello"))  // → "2cf24d..."
   sha256.Sum256([]byte("hallo"))  // → "d3751d..." (totally different)
   ```

**Common hash functions**:
- **SHA-256**: Used in Bitcoin, industry standard
- **SHA-3**: Newer, quantum-resistant alternative
- **Blake2**: Faster than SHA-256, used in some blockchains

### Concept 2: Hash Chains

A **hash chain** links data using hashes, creating tamper-evident history.

**Structure**:
```
Block 0: Hash0 = H(data0)
Block 1: Hash1 = H(data1 || Hash0)
Block 2: Hash2 = H(data2 || Hash1)
Block 3: Hash3 = H(data3 || Hash2)
```

**Tamper detection**:
```
Original chain:
  Block 0 → Block 1 → Block 2 → Block 3
  Hash0     Hash1     Hash2     Hash3

Attacker modifies Block 1:
  Block 0 → Block 1' → Block 2 → Block 3
  Hash0     Hash1'    Hash2     Hash3
                      ↑
                    Mismatch! (PrevHash ≠ Hash1')
```

**Applications beyond blockchain**:
- **Git commits**: Each commit hashes previous commit
- **Certificate transparency logs**: Tamper-evident certificate records
- **Audit logs**: Immutable event history

### Concept 3: Serialization and Determinism

**Why determinism matters**:

```go
// Non-deterministic (BAD)
type Block struct {
    Data      string
    Timestamp time.Time  // Includes nanoseconds, changes every call
}

func (b *Block) Hash() string {
    data, _ := json.Marshal(b)  // Field order not guaranteed
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}

// Calling Hash() twice on same block gives different results!
// Reason: Timestamp changes, JSON field order varies
```

```go
// Deterministic (GOOD)
type Block struct {
    Data      string
    Timestamp int64  // Fixed precision
}

func (b *Block) Serialize() []byte {
    var buf bytes.Buffer
    buf.WriteString(b.Data)
    binary.Write(&buf, binary.BigEndian, b.Timestamp)
    return buf.Bytes()  // Always same bytes for same block
}
```

**Deterministic serialization rules**:
1. **Fixed field order**: Always serialize in same sequence
2. **Fixed precision**: Use int64, not float64 or time.Time
3. **No randomness**: No UUIDs, random padding, etc.
4. **Canonical encoding**: One correct representation

### Concept 4: Timestamps in Distributed Systems

**Timestamp challenges**:

1. **Clock drift**: Different nodes have different times
   ```
   Node A: 10:30:15
   Node B: 10:30:23  (8 seconds ahead)
   ```

2. **Malicious timestamps**: Nodes can lie about time
   ```
   Block 5: Timestamp 2025-01-01 10:00
   Block 6: Timestamp 2024-01-01 10:00  ← Backdated!
   ```

**Solutions**:

1. **Median time**: Use median of previous 11 blocks (Bitcoin)
2. **Tolerance window**: Accept timestamps within ±2 hours
3. **Network time protocol**: Sync clocks using NTP
4. **Logical clocks**: Use sequence numbers instead of time

**Timestamp validation**:
```go
func ValidateTimestamp(block, prevBlock *Block) error {
    // Must be after previous block
    if block.Header.Timestamp <= prevBlock.Header.Timestamp {
        return errors.New("timestamp must increase")
    }

    // Cannot be too far in future (2 hours tolerance)
    now := time.Now().Unix()
    if block.Header.Timestamp > now + 2*3600 {
        return errors.New("timestamp too far in future")
    }

    return nil
}
```

### Concept 5: Immutability vs. Tamper-Evidence

**Immutability**: Cannot be changed.

**Tamper-evidence**: Changes are detectable.

**Blockchain is tamper-evident, not truly immutable:**
- You CAN change historical blocks
- But changes are DETECTABLE (hash chain breaks)
- With proof-of-work, changes are COMPUTATIONALLY INFEASIBLE

**Example**:
```go
// Attacker tries to modify block 100 in a 1000-block chain

// Step 1: Modify block 100's data
chain[100].Transactions[0] = "Fraudulent transaction"

// Step 2: Recompute block 100's hash
chain[100].Hash = chain[100].ComputeHash()

// Step 3: Fix block 101's PrevHash
chain[101].Header.PrevHash = chain[100].Hash
chain[101].Hash = chain[101].ComputeHash()

// Step 4-1000: Recompute all subsequent blocks...
// This takes time! With proof-of-work, this is infeasible.

// Meanwhile, honest nodes keep adding new blocks,
// so attacker can never catch up.
```

**51% attack**: If attacker controls >50% of network's computing power, they can rewrite history. This is why decentralization matters.

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Deterministic Serialization

```go
func SerializeStruct(s interface{}) ([]byte, error) {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    if err := enc.Encode(s); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
```

**Use case**: Any time you need consistent byte representation.

### Pattern 2: Hash-Based Integrity Check

```go
type Document struct {
    Content string
    Hash    string
}

func NewDocument(content string) *Document {
    doc := &Document{Content: content}
    doc.Hash = doc.ComputeHash()
    return doc
}

func (d *Document) ComputeHash() string {
    hash := sha256.Sum256([]byte(d.Content))
    return hex.EncodeToString(hash[:])
}

func (d *Document) Verify() bool {
    return d.Hash == d.ComputeHash()
}
```

**Use case**: File integrity, data verification, checksums.

### Pattern 3: Linked List with Hash Pointers

```go
type Node struct {
    Data     string
    PrevHash string
    Hash     string
}

type HashLinkedList struct {
    Nodes []*Node
}

func (ll *HashLinkedList) Append(data string) {
    var prevHash string
    if len(ll.Nodes) > 0 {
        prevHash = ll.Nodes[len(ll.Nodes)-1].Hash
    }

    node := &Node{Data: data, PrevHash: prevHash}
    node.Hash = computeHash(node.Data + node.PrevHash)

    ll.Nodes = append(ll.Nodes, node)
}
```

**Use case**: Audit logs, event sourcing, version control.

### Pattern 4: Merkle Tree Root

```go
func BuildMerkleRoot(data [][]byte) []byte {
    if len(data) == 0 {
        return nil
    }
    if len(data) == 1 {
        hash := sha256.Sum256(data[0])
        return hash[:]
    }

    // Pair up and hash
    var nextLevel [][]byte
    for i := 0; i < len(data); i += 2 {
        if i+1 < len(data) {
            combined := append(data[i], data[i+1]...)
            hash := sha256.Sum256(combined)
            nextLevel = append(nextLevel, hash[:])
        } else {
            // Odd number: duplicate last element
            hash := sha256.Sum256(data[i])
            nextLevel = append(nextLevel, hash[:])
        }
    }

    return BuildMerkleRoot(nextLevel)
}
```

**Use case**: Efficient data integrity, distributed databases, Git.

### Pattern 5: Chain Validation with Detailed Errors

```go
type ValidationError struct {
    BlockIndex int
    ErrorType  string
    Message    string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("Block %d (%s): %s", e.BlockIndex, e.ErrorType, e.Message)
}

func ValidateChainDetailed(chain []*Block) error {
    for i, block := range chain {
        if block.Hash != block.ComputeHash() {
            return &ValidationError{
                BlockIndex: i,
                ErrorType:  "HASH_MISMATCH",
                Message:    "stored hash doesn't match computed hash",
            }
        }

        if i > 0 && block.Header.PrevHash != chain[i-1].Hash {
            return &ValidationError{
                BlockIndex: i,
                ErrorType:  "LINK_BROKEN",
                Message:    "prev hash doesn't match previous block's hash",
            }
        }
    }

    return nil
}
```

---

## 7. Real-World Applications

### Bitcoin Blockchain

**Block structure** (simplified):

```go
type BitcoinBlock struct {
    Header BitcoinBlockHeader
    Transactions []Transaction
}

type BitcoinBlockHeader struct {
    Version       int32
    PrevBlockHash [32]byte
    MerkleRoot    [32]byte
    Timestamp     uint32
    Bits          uint32  // Difficulty target
    Nonce         uint32  // Proof-of-work
}
```

**Key differences from our implementation**:
- Fixed-size byte arrays instead of strings
- Difficulty and nonce for proof-of-work
- More compact binary serialization
- Full Merkle tree instead of simple hash

### Ethereum Blockchain

**Block structure**:

```go
type EthereumBlock struct {
    Header       EthereumBlockHeader
    Transactions []Transaction
    Uncles       []*BlockHeader  // For uncle block rewards
}

type EthereumBlockHeader struct {
    ParentHash   common.Hash
    UncleHash    common.Hash
    Coinbase     common.Address  // Miner's address
    Root         common.Hash     // State trie root
    TxHash       common.Hash     // Transaction trie root
    ReceiptHash  common.Hash     // Receipt trie root
    Bloom        Bloom           // Log bloom filter
    Difficulty   *big.Int
    Number       *big.Int
    GasLimit     uint64
    GasUsed      uint64
    Time         uint64
    Extra        []byte
    MixDigest    common.Hash
    Nonce        BlockNonce
}
```

**Differences**:
- State root (entire world state)
- Receipt root (transaction outcomes)
- Gas (transaction fees)
- Much more complex than Bitcoin

### Git Version Control

**Git commit structure** (similar to blockchain):

```go
type GitCommit struct {
    Tree      string   // Hash of file tree
    Parent    string   // Hash of parent commit
    Author    string
    Committer string
    Message   string
    Hash      string   // SHA-1 hash of above
}
```

**How it's blockchain-like**:
- Commits hash-linked (parent hash)
- Tampering changes all subsequent hashes
- Distributed consensus (multiple clones)

### Certificate Transparency

**CT Log Entry**:

```go
type CTLogEntry struct {
    Index       int64
    LeafInput   []byte  // Certificate data
    ExtraData   []byte  // Certificate chain
    Timestamp   uint64
    Hash        []byte  // Hash of this entry
}
```

**Usage**: Detect malicious SSL certificates by maintaining public, append-only log.

### Supply Chain Tracking

**Supply chain block**:

```go
type SupplyChainBlock struct {
    ProductID    string
    Event        string  // "Manufactured", "Shipped", "Delivered"
    Location     string
    Timestamp    int64
    PrevHash     string
    Hash         string
}
```

**Use case**: Track product from factory to customer, detect counterfeit goods.

---

## 8. Common Mistakes to Avoid

### Mistake 1: Non-Deterministic Serialization

**Wrong**:
```go
func (b *Block) Serialize() []byte {
    data, _ := json.Marshal(b)
    return data  // Field order not guaranteed!
}
```

**Result**: Same block produces different hashes.

**Right**:
```go
func (b *Block) Serialize() []byte {
    var buf bytes.Buffer
    binary.Write(&buf, binary.BigEndian, int64(b.Header.Index))
    // ... fixed order
    return buf.Bytes()
}
```

### Mistake 2: Hashing Before Setting MerkleRoot

**Wrong**:
```go
func NewBlock(prev *Block, txs []string) *Block {
    block := &Block{...}
    block.Hash = block.ComputeHash()  // MerkleRoot not set yet!
    block.Header.MerkleRoot = ComputeMerkleRoot(txs)
    return block
}
```

**Result**: Hash doesn't include merkle root.

**Right**:
```go
func NewBlock(prev *Block, txs []string) *Block {
    block := &Block{...}
    block.Header.MerkleRoot = ComputeMerkleRoot(txs)  // First
    block.Hash = block.ComputeHash()  // Then hash
    return block
}
```

### Mistake 3: Using Floating-Point Timestamps

**Wrong**:
```go
type Block struct {
    Timestamp float64  // Imprecise!
}
```

**Problem**: Floating-point has rounding errors, not deterministic.

**Right**:
```go
type Block struct {
    Timestamp int64  // Seconds since epoch (deterministic)
}
```

### Mistake 4: Forgetting to Validate Genesis Block

**Wrong**:
```go
func ValidateChain(chain []*Block) error {
    for i := 1; i < len(chain); i++ {  // Skips genesis!
        // validate...
    }
}
```

**Problem**: Genesis block could be invalid.

**Right**:
```go
func ValidateChain(chain []*Block) error {
    if len(chain) == 0 {
        return errors.New("empty chain")
    }

    // Validate genesis
    if chain[0].Header.Index != 0 {
        return errors.New("invalid genesis")
    }

    // Validate rest
    for i := 1; i < len(chain); i++ {
        // validate...
    }
}
```

### Mistake 5: Storing Hash in Serialization

**Wrong**:
```go
func (b *Block) Serialize() []byte {
    var buf bytes.Buffer
    buf.WriteString(b.Hash)  // Circular dependency!
    // ... other fields
    return buf.Bytes()
}
```

**Problem**: Hash depends on serialization, but serialization includes hash → infinite loop or wrong hash.

**Right**:
```go
func (b *Block) Serialize() []byte {
    var buf bytes.Buffer
    // Serialize header fields (NOT including hash)
    // Serialize transactions
    return buf.Bytes()
}
```

### Mistake 6: Comparing Byte Slices Incorrectly

**Wrong**:
```go
hash1 := sha256.Sum256(data1)
hash2 := sha256.Sum256(data2)

if hash1 == hash2 {  // Won't compile! [32]byte not comparable
    // ...
}
```

**Right**:
```go
hash1 := sha256.Sum256(data1)
hash2 := sha256.Sum256(data2)

if bytes.Equal(hash1[:], hash2[:]) {  // Convert to slice
    // ...
}

// Or compare hex strings
hex1 := hex.EncodeToString(hash1[:])
hex2 := hex.EncodeToString(hash2[:])
if hex1 == hex2 {
    // ...
}
```

---

## 9. Stretch Goals

### Goal 1: Add Block Size Limits

Prevent blocks from becoming too large.

**Hint**:
```go
const MaxBlockSize = 1024 * 1024 // 1 MB

func NewBlock(prev *Block, txs []string) (*Block, error) {
    block := &Block{...}

    if len(block.Serialize()) > MaxBlockSize {
        return nil, errors.New("block too large")
    }

    return block, nil
}
```

### Goal 2: Implement Block Height Queries

Find block by index efficiently.

**Hint**:
```go
type Blockchain struct {
    Blocks      []*Block
    BlockByHash map[string]*Block
    BlockByIndex map[int]*Block
}

func (bc *Blockchain) GetBlockByIndex(index int) *Block {
    return bc.BlockByIndex[index]
}
```

### Goal 3: Add Difficulty and Proof-of-Work

Require hash to start with N zeros (see Project 43).

**Hint**:
```go
func (b *Block) Mine(difficulty int) {
    target := strings.Repeat("0", difficulty)

    for {
        b.Header.Nonce++
        b.Hash = b.ComputeHash()

        if strings.HasPrefix(b.Hash, target) {
            break  // Found valid hash
        }
    }
}
```

### Goal 4: Persist Blockchain to Disk

Save and load chain from file.

**Hint**:
```go
func SaveChain(chain []*Block, filename string) error {
    data, err := json.MarshalIndent(chain, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}

func LoadChain(filename string) ([]*Block, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var chain []*Block
    err = json.Unmarshal(data, &chain)
    return chain, err
}
```

### Goal 5: Implement Proper Merkle Tree

Build full Merkle tree with proof generation (see Project 40).

---

## How to Run

```bash
# Run the demo
make run P=42-simple-block-struct-hashing

# Or directly
cd /home/user/go-edu/minis/42-simple-block-struct-hashing
go run cmd/block-demo/main.go

# Run tests
go test ./minis/42-simple-block-struct-hashing/exercise/...

# Run with verbose output
go test -v ./minis/42-simple-block-struct-hashing/exercise/...

# Run solution version (with build tag)
go test -tags=solution ./minis/42-simple-block-struct-hashing/exercise/...
```

---

## Summary

**What you learned**:
- Block structure: Headers + data + hash
- Hash linking: Creates tamper-evident chain
- Serialization: Deterministic byte encoding
- Cryptographic hashing: SHA-256 for integrity
- Chain validation: Detecting tampering
- Timestamps: Chronological ordering
- Merkle roots: Data integrity within blocks

**Why this matters**:
Blockchain is one of the most important innovations in distributed systems. Understanding how blocks work is fundamental to cryptocurrencies, supply chain, voting systems, and any application needing immutable records.

**Key insights**:
- Hashing makes tampering detectable (not impossible, but evident)
- Deterministic serialization is critical for consensus
- Hash chains create immutable history
- Validation must check multiple properties (hash, links, timestamps)

**Next steps**:
- Project 43: Proof-of-work and mining
- Project 44: Mempool and transaction ordering
- Project 45: P2P gossip network for block propagation

Build on! 🔗
