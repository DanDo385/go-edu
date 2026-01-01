package proofofworkdemo

import (
	"math"
)

// Block represents a single block in the blockchain
type Block struct {
	Index     int
	Timestamp int64
	Data      string
	PrevHash  string
	Nonce     int
	Hash      string
}

// Exercise 1: Calculate Block Hash
func CalculateBlockHash(b Block) string {
	// TODO: Implement this function.

	// The hash of a block is a cryptographic fingerprint of its contents. It must be deterministic.

	// Step 1: Create a single string by concatenating all the important block fields.
	// - The order must be consistent.
	// - Use `fmt.Sprintf` to combine `Index`, `Timestamp`, `Data`, `PrevHash`, and `Nonce`.
	// - `record := fmt.Sprintf("%d%d%s%s%d", ...)`

	// Step 2: Compute the SHA-256 hash of the string.
	// - The hash function takes a byte slice, so convert the string: `[]byte(record)`.
	// - `hash := sha256.Sum256([]byte(record))`

	// Step 3: Convert the hash (which is a `[32]byte` array) to a hexadecimal string.
	// - The `encoding/hex` package is perfect for this.
	// - `return hex.EncodeToString(hash[:])`

	return ""
}

// Exercise 2: Validate Proof of Work
func IsValidProof(hash string, difficulty int) bool {
	// TODO: Implement this function.

	// Proof-of-work in this system means finding a hash that starts with a certain number of zeros.
	// The `difficulty` determines how many zeros are required.

	// Step 1: Create the target prefix.
	// - This is a string consisting of `difficulty` zeros.
	// - `target := strings.Repeat("0", difficulty)`

	// Step 2: Check if the hash has the required prefix.
	// - `return strings.HasPrefix(hash, target)`

	return false
}

// Exercise 3: Mine a Block
func MineBlock(block *Block, difficulty int) int {
	// TODO: Implement this function.

	// "Mining" is the process of finding a `Nonce` that results in a block hash meeting the `difficulty` requirement.
	// This is a brute-force search.

	// Step 1: Use an infinite `for` loop that starts with `attempts := 0`.
	// - `for attempts := 0; ; attempts++ { ... }`

	// Step 2: Inside the loop, set the block's nonce.
	// - `block.Nonce = attempts`

	// Step 3: Calculate the hash for the block with this new nonce.
	// - `hash := CalculateBlockHash(*block)`
	// - Note: `CalculateBlockHash` takes a `Block` by value, so we dereference the pointer.

	// Step 4: Check if the new hash is a valid proof.
	// - `if IsValidProof(hash, difficulty) { ... }`

	// Step 5: If the proof is valid:
	//   - You've found the block!
	//   - Set the block's `Hash` field: `block.Hash = hash`.
	//   - Return the number of `attempts` it took.

	return 0
}

// Exercise 4: Validate Blockchain
func ValidateChain(chain []Block, difficulty int) bool {
	// TODO: Implement this function.

	// This function verifies the integrity of the entire chain.

	// Step 1: Handle edge cases.
	// - If the chain is empty, it's trivially valid.
	// - If the chain has only one block (the genesis block), calculate its hash and check its proof-of-work.

	// Step 2: Loop through the chain from the second block (`i = 1`).
	// - `for i := 1; i < len(chain); i++ { ... }`

	// Step 3: Inside the loop, for each `current` block and its `previous` block:
	//   a. **Verify Hash Integrity**: Check if `current.Hash` is equal to the result of `CalculateBlockHash(current)`. If not, the block's data was tampered with. Return `false`.
	//   b. **Verify Proof of Work**: Check if `IsValidProof(current.Hash, difficulty)` is `true`. If not, the block is invalid. Return `false`.
	//   c. **Verify Chain Link**: Check if `current.PrevHash` is equal to `previous.Hash`. This is the most critical check. If they don't match, the chain is broken. Return `false`.

	// Step 4: If the loop completes without finding any issues, the chain is valid. Return `true`.

	return false
}

// Exercise 5: Adjust Difficulty
func AdjustDifficulty(chain []Block, targetBlockTime int64, currentDifficulty int) int {
	// TODO: Implement this function.

	// The goal is to keep the average block time close to `targetBlockTime`.

	// Step 1: Check if there are enough blocks to make an adjustment.
	// - If the chain has fewer than, say, 10 blocks, it's too early to tell. Return the `currentDifficulty`.

	// Step 2: Define a window of recent blocks to analyze.
	// - For example, the last 10 blocks.

	// Step 3: Calculate the actual time it took to mine those blocks.
	// - `actualTime := chain[len(chain)-1].Timestamp - chain[len(chain)-10].Timestamp`

	// Step 4: Calculate the expected time for that many blocks.
	// - `expectedTime := targetBlockTime * 9` (9 intervals between 10 blocks).

	// Step 5: Compare and adjust.
	// - If `actualTime` is much faster than `expectedTime` (e.g., `actualTime < expectedTime / 2`), increase the difficulty.
	// - If `actualTime` is much slower (e.g., `actualTime > expectedTime * 2`), decrease the difficulty.
	// - Make sure the difficulty never goes below 1.
	// - Otherwise, keep the difficulty the same.

	return currentDifficulty
}

// Exercise 6: Calculate Mining Probability
func MiningProbability(hashRate float64, difficulty int, timeSeconds float64) float64 {
	// TODO: Implement this function.

	// This function models the block finding process as a Poisson process.

	// Step 1: Calculate the size of the search space (the expected number of hashes you need to try to find one block).
	// - Since our hash is hexadecimal, each character has 16 possibilities.
	// - `expectedAttempts := math.Pow(16, float64(difficulty))`

	// Step 2: Calculate lambda (λ), the average number of blocks you expect to find per second.
	// - `lambda := hashRate / expectedAttempts`

	// Step 3: Use the Poisson probability formula to find the probability of finding *at least one* block in the given time.
	// - The probability of finding *zero* blocks is `e^(-λ * t)`.
	// - Therefore, the probability of finding at least one is `1 - e^(-λ * t)`.
	// - `probability := 1 - math.Exp(-lambda * timeSeconds)`

	// Step 4: Return the probability.
	return 0.0
}

// Helper function for testing - creates a simple blockchain
func CreateTestChain(numBlocks int, difficulty int) []Block {
	chain := make([]Block, numBlocks)

	// Genesis block
	chain[0] = Block{
		Index:     0,
		Timestamp: 1000,
		Data:      "Genesis Block",
		PrevHash:  "0",
		Nonce:     0,
	}
	MineBlock(&chain[0], difficulty)

	// Subsequent blocks
	for i := 1; i < numBlocks; i++ {
		chain[i] = Block{
			Index:     i,
			Timestamp: chain[i-1].Timestamp + 10, // 10 second intervals
			Data:      "Block " + string(rune('0'+i)),
			PrevHash:  chain[i-1].Hash,
			Nonce:     0,
		}
		MineBlock(&chain[i], difficulty)
	}

	return chain
}

// STRETCH GOAL: Merkle Root
func BuildMerkleRoot(transactions []string) string {
	// TODO: (Stretch goal) Implement a proper Merkle tree.

	// Step 1: Handle edge cases (zero or one transaction).
	// Step 2: Create the "leaf" nodes by hashing each individual transaction.
	// Step 3: In a loop, build the tree level by level until only one hash (the root) remains.
	//   - `for len(nodes) > 1 { ... }`
	//   - Inside the loop, take pairs of hashes, concatenate them, and hash the result to create the parent nodes for the next level.
	//   - If you have an odd number of nodes at any level, duplicate the last one and hash it with itself.
	// Step 4: Return the final root hash.

	return ""
}

// STRETCH GOAL: Mining Pool Share Calculation
func CalculatePoolRewards(miners []Miner, blockReward float64) map[string]float64 {
	// TODO: (Stretch goal) Implement proportional reward distribution.

	// Step 1: Calculate the total number of shares contributed by all miners in the pool.
	// Step 2: Create a map to store the rewards for each miner.
	// Step 3: Loop through each miner.
	//   - Calculate their proportion of the total shares: `float64(miner.Shares) / float64(totalShares)`.
	//   - Calculate their reward: `proportion * blockReward`.
	//   - Store it in the rewards map.
	// Step 4: Return the map.

	return nil
}

// STRETCH GOAL: Estimate Attack Cost
func EstimateAttackCost(
	numBlocks int,
	difficulty int,
	honestHashRate float64,
	attackerHashRate float64,
	electricityCostPerKWh float64,
	minerWattage float64,
) float64 {
	// TODO: (Stretch goal) Implement the attack cost estimation.

	// This is a simplified model. It assumes the attacker must re-mine N blocks plus any blocks the honest network finds during the attack.

	// Step 1: Calculate the expected number of hashes per block for the given difficulty.
	// Step 2: Calculate the time it would take the attacker to mine `numBlocks`.
	// Step 3: Calculate how many blocks the honest network would find in that time.
	// Step 4: The attacker must mine the original `numBlocks` plus the new honest blocks. Calculate the total time this would take the attacker.
	// Step 5: Convert this total time into energy usage (kWh) based on the `minerWattage`.
	// Step 6: Multiply the energy usage by the `electricityCostPerKWh` to get the final cost.

	return 0.0
}

// STRETCH GOAL: Selfish Mining Simulation
func SimulateSelfishMining(
	honestHashRate float64,
	selfishHashRate float64,
	numBlocks int,
	difficulty int,
) (honestBlocks int, selfishBlocks int) {
	// TODO: (Stretch goal) Implement a simplified selfish mining simulation.

	// This is a complex simulation. A simplified model would be:
	// - In a loop, decide who finds the next block based on their proportion of the total hash rate.
	// - If the honest network finds a block, increment their public chain length.
	// - If the selfish miner finds a block, increment their *private* chain length.
	// - Implement the selfish miner's strategy:
	//   - If their private chain is longer than the public chain by 2 or more, they can broadcast it, invalidating the public chain's work.
	//   - If the public chain catches up, they are forced to broadcast their chain in a race.
	// - Keep track of how many blocks each side successfully gets into the "real" chain over a simulation of `numBlocks` total blocks.
	return 0, 0
}

// Additional utility functions you might need:

// CountLeadingZeros counts the number of leading zero characters in a hex string
func CountLeadingZeros(hash string) int {
	count := 0
	for _, char := range hash {
		if char == '0' {
			count++
		} else {
			break
		}
	}
	return count
}

// ExpectedMiningTime calculates the expected time to mine a block
func ExpectedMiningTime(hashRate float64, difficulty int) float64 {
	expectedAttempts := math.Pow(16, float64(difficulty))
	return expectedAttempts / hashRate
}
