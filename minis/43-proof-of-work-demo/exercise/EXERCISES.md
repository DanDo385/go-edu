# Proof of Work Exercises

## Overview

In these exercises, you'll implement a complete Proof of Work mining system from scratch. You'll learn how to:
- Calculate cryptographic hashes
- Mine blocks by finding valid nonces
- Build and validate blockchains
- Adjust mining difficulty
- Measure mining performance

## Exercise 1: Calculate Block Hash ⭐

**File**: `exercise.go` → `CalculateBlockHash()`

Implement a function that calculates the SHA-256 hash of a block.

**Requirements**:
- Concatenate all block fields in order: Index, Timestamp, Data, PrevHash, Nonce
- Compute SHA-256 hash of the concatenated string
- Return the hash as a hexadecimal string

**Example**:
```go
block := Block{
    Index:     1,
    Timestamp: 1609459200,
    Data:      "Hello, blockchain!",
    PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
    Nonce:     42,
}

hash := CalculateBlockHash(block)
// hash should be a 64-character hexadecimal string
// hash should change completely if any field changes (avalanche effect)
```

**Hints**:
- Use `crypto/sha256` package
- Use `fmt.Sprintf()` to concatenate fields
- Use `encoding/hex` to convert bytes to hex string
- The hash is 64 characters long (32 bytes * 2 hex chars per byte)

**Test**: Run `go test -run TestCalculateBlockHash` to verify your implementation.

---

## Exercise 2: Validate Proof of Work ⭐

**File**: `exercise.go` → `IsValidProof()`

Implement a function that checks if a block's hash meets the proof of work difficulty requirement.

**Requirements**:
- Check if the hash starts with the required number of leading zeros
- Difficulty 1 = "0..." (1 zero), Difficulty 5 = "00000..." (5 zeros)
- Return true if valid, false otherwise

**Example**:
```go
// Valid proof
hash := "0000abc123..."  // 4 leading zeros
difficulty := 4
isValid := IsValidProof(hash, difficulty)  // Should return true

// Invalid proof
hash := "000abc123..."   // Only 3 leading zeros
difficulty := 4
isValid := IsValidProof(hash, difficulty)  // Should return false
```

**Hints**:
- Use `strings.HasPrefix()` to check the prefix
- Use `strings.Repeat()` to create the target string of zeros

**Test**: Run `go test -run TestIsValidProof` to verify your implementation.

---

## Exercise 3: Mine a Block ⭐⭐

**File**: `exercise.go` → `MineBlock()`

Implement the core mining algorithm that finds a valid nonce for a block.

**Requirements**:
- Start with nonce = 0
- Calculate the hash with the current nonce
- Check if hash meets difficulty requirement
- If not, increment nonce and try again
- Continue until a valid hash is found
- Update the block's Nonce and Hash fields
- Return the number of attempts it took

**Example**:
```go
block := Block{
    Index:     1,
    Timestamp: 1609459200,
    Data:      "Transaction data",
    PrevHash:  "000...",
    Nonce:     0,
}

attempts := MineBlock(&block, 4)
// block.Nonce should now contain the valid nonce
// block.Hash should start with "0000"
// attempts should be > 0
```

**Hints**:
- Use a loop that continues until proof is valid
- Increment nonce on each iteration
- Update both Nonce AND Hash fields of the block
- Consider: what if nonce overflows? (Advanced: update timestamp)

**Test**: Run `go test -run TestMineBlock` to verify your implementation.

---

## Exercise 4: Validate Blockchain ⭐⭐

**File**: `exercise.go` → `ValidateChain()`

Implement a function that validates an entire blockchain.

**Requirements**:
Check each block in the chain:
1. Hash is correct (stored hash matches calculated hash)
2. Proof of work is valid (hash meets difficulty requirement)
3. Links to previous block (PrevHash matches previous block's Hash)

**Example**:
```go
chain := []Block{
    {Index: 0, Hash: "abc...", ...},  // Genesis block
    {Index: 1, Hash: "def...", PrevHash: "abc...", ...},
    {Index: 2, Hash: "ghi...", PrevHash: "def...", ...},
}

isValid := ValidateChain(chain, 4)  // Should return true

// If we tamper with block 1:
chain[1].Data = "HACKED"
isValid = ValidateChain(chain, 4)  // Should return false
```

**Hints**:
- Loop through blocks starting at index 1 (skip genesis)
- For each block, verify: hash correctness, proof of work, and linkage
- Return false immediately if any check fails
- Return true only if all blocks pass all checks

**Test**: Run `go test -run TestValidateChain` to verify your implementation.

---

## Exercise 5: Adjust Difficulty ⭐⭐⭐

**File**: `exercise.go` → `AdjustDifficulty()`

Implement a difficulty adjustment algorithm that maintains consistent block times.

**Requirements**:
- Calculate actual time taken for the last N blocks
- Calculate expected time (N blocks * target block time)
- If mining is too fast (actual < expected/2), increase difficulty
- If mining is too slow (actual > expected*2), decrease difficulty
- Otherwise, keep difficulty unchanged
- Return the new difficulty

**Example**:
```go
chain := []Block{
    {Index: 0, Timestamp: 1000, ...},
    {Index: 1, Timestamp: 1002, ...},  // 2 seconds
    {Index: 2, Timestamp: 1004, ...},  // 2 seconds
    {Index: 3, Timestamp: 1006, ...},  // 2 seconds
}

targetBlockTime := 10  // Want 10 seconds per block
currentDifficulty := 4

newDifficulty := AdjustDifficulty(chain, targetBlockTime, currentDifficulty)
// Mining is too fast (2 sec vs 10 sec target), so increase difficulty
// newDifficulty should be 5
```

**Hints**:
- Need at least 2 blocks to calculate time
- Use timestamps of first and last block in the window
- Number of intervals = number of blocks - 1
- Use thresholds to avoid reacting to small variances
- Minimum difficulty should be 1

**Test**: Run `go test -run TestAdjustDifficulty` to verify your implementation.

---

## Exercise 6: Calculate Mining Probability ⭐⭐⭐

**File**: `exercise.go` → `MiningProbability()`

Implement a function that calculates the probability of finding a block within a certain time.

**Requirements**:
- Given: hash rate (hashes/second), difficulty, and time (seconds)
- Calculate: probability of finding a valid block within that time
- Use the formula: P = 1 - e^(-λt), where λ = hashRate / expectedAttempts
- Expected attempts = 16^difficulty (for hexadecimal zeros)

**Example**:
```go
hashRate := 1000.0      // 1000 hashes per second
difficulty := 3         // Expected attempts = 16^3 = 4096
timeSeconds := 10.0     // 10 seconds

probability := MiningProbability(hashRate, difficulty, timeSeconds)
// λ = 1000 / 4096 ≈ 0.244
// P = 1 - e^(-0.244 * 10) ≈ 0.914 (91.4% chance)
```

**Hints**:
- Use `math.Pow(16, difficulty)` for expected attempts
- Use `math.Exp(x)` for e^x
- λ (lambda) = hashRate / expectedAttempts
- The formula models mining as a Poisson process

**Test**: Run `go test -run TestMiningProbability` to verify your implementation.

---

## Running the Tests

```bash
# Run all tests
go test -v

# Run a specific test
go test -run TestCalculateBlockHash

# Run tests with coverage
go test -cover

# Run benchmarks
go test -bench=. -benchmem

# Run with race detector (if using concurrency)
go test -race
```

---

## Stretch Goals

Once you've completed all exercises, try these challenges:

### Challenge 1: Implement a Merkle Root ⭐⭐⭐

Instead of hashing all transactions directly, build a Merkle tree and use the root hash.

**Benefits**:
- Efficiently prove a transaction is in a block (O(log n) instead of O(n))
- Light clients can verify transactions without downloading entire blocks

```go
func BuildMerkleRoot(transactions []string) string {
    // Build binary tree of hashes
    // Pair-wise hash until single root remains
}
```

### Challenge 2: Simulate a 51% Attack ⭐⭐⭐

Simulate a scenario where an attacker tries to rewrite blockchain history.

```go
func Simulate51PercentAttack(honestChain []Block, attackerHashRate, honestHashRate float64) {
    // Attacker mines a private chain
    // Honest miners mine on the public chain
    // Race to see which chain grows longer
}
```

### Challenge 3: Implement Mining Pools ⭐⭐⭐⭐

Simulate multiple miners working together in a pool.

```go
type Miner struct {
    ID       string
    HashRate float64
}

func SimulatePool(miners []Miner, difficulty int, blocks int) map[string]float64 {
    // Distribute work among miners
    // Track shares (near-valid hashes)
    // Distribute rewards proportionally
}
```

### Challenge 4: Optimize Mining with Goroutines ⭐⭐⭐⭐

Use multiple goroutines to mine in parallel.

```go
func MineBlockParallel(block *Block, difficulty int, workers int) int {
    // Split nonce space among workers
    // First worker to find valid nonce wins
    // Use channels to communicate results
}
```

**Hint**: Be careful with shared state! Each worker needs its own copy of the block.

---

## Common Pitfalls

### Pitfall 1: Forgetting to Update Hash After Mining

```go
// ❌ WRONG
block.Nonce = validNonce
// Forgot to update block.Hash!

// ✅ CORRECT
block.Nonce = validNonce
block.Hash = CalculateBlockHash(block)
```

### Pitfall 2: Off-by-One in Validation

```go
// ❌ WRONG
for i := 0; i < len(chain); i++ {
    if chain[i].PrevHash != chain[i-1].Hash {  // Crash on i=0!
        return false
    }
}

// ✅ CORRECT
for i := 1; i < len(chain); i++ {  // Start at 1
    if chain[i].PrevHash != chain[i-1].Hash {
        return false
    }
}
```

### Pitfall 3: Integer Overflow in Nonce

```go
// ❌ WRONG
for {
    block.Nonce++  // Can overflow!
    if isValid(block) {
        return
    }
}

// ✅ CORRECT
for {
    block.Nonce++
    if block.Nonce == math.MaxInt {
        block.Timestamp = time.Now().Unix()
        block.Nonce = 0
    }
    if isValid(block) {
        return
    }
}
```

### Pitfall 4: Not Handling Empty Chain

```go
// ❌ WRONG
func AdjustDifficulty(chain []Block) int {
    lastBlock := chain[len(chain)-1]  // Crash if chain is empty!
    // ...
}

// ✅ CORRECT
func AdjustDifficulty(chain []Block) int {
    if len(chain) < 2 {
        return currentDifficulty  // Need at least 2 blocks
    }
    // ...
}
```

---

## Learning Resources

- **Bitcoin Whitepaper**: The original paper by Satoshi Nakamoto
- **Mastering Bitcoin**: Andreas Antonopoulos (free online)
- **Hash Functions**: Understanding SHA-256 and cryptographic properties
- **Byzantine Generals Problem**: The theoretical foundation of distributed consensus
- **Proof of Stake**: Alternative consensus mechanism used by Ethereum

---

## Next Steps

After mastering Proof of Work, explore:
- **Digital Signatures**: How transactions are authenticated (ECDSA)
- **Merkle Trees**: Efficient verification of transaction inclusion
- **UTXO Model**: Bitcoin's approach to tracking ownership
- **Smart Contracts**: Programmable blockchains (Ethereum)
- **Layer 2 Solutions**: Lightning Network, rollups, state channels
- **Proof of Stake**: Ethereum's current consensus mechanism

Good luck! ⛏️
