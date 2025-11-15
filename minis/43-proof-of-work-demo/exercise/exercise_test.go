package exercise

import (
	"math"
	"strings"
	"testing"
)

// Test Exercise 1: Calculate Block Hash
func TestCalculateBlockHash(t *testing.T) {
	tests := []struct {
		name  string
		block Block
	}{
		{
			name: "simple block",
			block: Block{
				Index:     1,
				Timestamp: 1609459200,
				Data:      "Hello, blockchain!",
				PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
				Nonce:     42,
			},
		},
		{
			name: "genesis block",
			block: Block{
				Index:     0,
				Timestamp: 1000000000,
				Data:      "Genesis Block",
				PrevHash:  "0",
				Nonce:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := CalculateBlockHash(tt.block)

			// Check hash is 64 characters (32 bytes * 2 hex chars)
			if len(hash) != 64 {
				t.Errorf("Hash length = %d, want 64", len(hash))
			}

			// Check hash is hexadecimal
			for _, char := range hash {
				if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
					t.Errorf("Hash contains non-hex character: %c", char)
				}
			}

			// Check determinism (same input -> same output)
			hash2 := CalculateBlockHash(tt.block)
			if hash != hash2 {
				t.Error("Hash is not deterministic")
			}

			// Check avalanche effect (small change -> completely different hash)
			modifiedBlock := tt.block
			modifiedBlock.Data = tt.block.Data + "!"
			differentHash := CalculateBlockHash(modifiedBlock)
			if hash == differentHash {
				t.Error("Avalanche effect not working (same hash for different data)")
			}

			// Verify against solution
			expected := CalculateBlockHashSolution(tt.block)
			if hash != expected {
				t.Errorf("Hash = %s, want %s", hash, expected)
			}
		})
	}
}

// Test Exercise 2: Validate Proof of Work
func TestIsValidProof(t *testing.T) {
	tests := []struct {
		name       string
		hash       string
		difficulty int
		want       bool
	}{
		{
			name:       "valid with 1 zero",
			hash:       "0abc123456789",
			difficulty: 1,
			want:       true,
		},
		{
			name:       "valid with 4 zeros",
			hash:       "0000abc123456",
			difficulty: 4,
			want:       true,
		},
		{
			name:       "invalid - not enough zeros",
			hash:       "000abc123456",
			difficulty: 4,
			want:       false,
		},
		{
			name:       "invalid - no zeros",
			hash:       "abc123456789",
			difficulty: 1,
			want:       false,
		},
		{
			name:       "valid with 6 zeros",
			hash:       "000000abc123",
			difficulty: 6,
			want:       true,
		},
		{
			name:       "difficulty 0",
			hash:       "anything",
			difficulty: 0,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidProof(tt.hash, tt.difficulty)
			if got != tt.want {
				t.Errorf("IsValidProof(%q, %d) = %v, want %v",
					tt.hash, tt.difficulty, got, tt.want)
			}

			// Verify against solution
			expected := IsValidProofSolution(tt.hash, tt.difficulty)
			if got != expected {
				t.Errorf("Result doesn't match solution: got %v, want %v", got, expected)
			}
		})
	}
}

// Test Exercise 3: Mine a Block
func TestMineBlock(t *testing.T) {
	tests := []struct {
		name       string
		block      Block
		difficulty int
		maxTime    int // maximum expected attempts
	}{
		{
			name: "easy difficulty",
			block: Block{
				Index:     1,
				Timestamp: 1609459200,
				Data:      "Test transaction",
				PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
				Nonce:     0,
			},
			difficulty: 2,
			maxTime:    1000,
		},
		{
			name: "medium difficulty",
			block: Block{
				Index:     2,
				Timestamp: 1609459300,
				Data:      "Another transaction",
				PrevHash:  "00abc123",
				Nonce:     0,
			},
			difficulty: 3,
			maxTime:    10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the test case
			block := tt.block

			attempts := MineBlock(&block, tt.difficulty)

			// Check that attempts is positive
			if attempts <= 0 {
				t.Errorf("Attempts = %d, want > 0", attempts)
			}

			// Check that hash is valid
			if !IsValidProof(block.Hash, tt.difficulty) {
				t.Errorf("Mined hash %s is not valid for difficulty %d",
					block.Hash, tt.difficulty)
			}

			// Check that hash matches calculated hash
			calculatedHash := CalculateBlockHash(block)
			if block.Hash != calculatedHash {
				t.Errorf("Stored hash doesn't match calculated hash")
			}

			// Check that nonce was updated
			if block.Nonce == 0 && attempts > 1 {
				t.Error("Nonce should have been updated")
			}

			t.Logf("Mined with difficulty %d in %d attempts (nonce=%d)",
				tt.difficulty, attempts, block.Nonce)
		})
	}
}

// Test Exercise 4: Validate Blockchain
func TestValidateChain(t *testing.T) {
	difficulty := 3

	// Create a valid chain
	validChain := CreateTestChain(5, difficulty)

	tests := []struct {
		name  string
		chain []Block
		want  bool
	}{
		{
			name:  "valid chain",
			chain: validChain,
			want:  true,
		},
		{
			name:  "empty chain",
			chain: []Block{},
			want:  true,
		},
		{
			name: "single block valid",
			chain: []Block{
				{
					Index:     0,
					Timestamp: 1000,
					Data:      "Genesis",
					PrevHash:  "0",
					Hash:      "",
				},
			},
			want: false, // Invalid because hash is empty
		},
	}

	// Add test for tampered data
	tamperedChain := make([]Block, len(validChain))
	copy(tamperedChain, validChain)
	tamperedChain[2].Data = "TAMPERED"
	tests = append(tests, struct {
		name  string
		chain []Block
		want  bool
	}{
		name:  "tampered data",
		chain: tamperedChain,
		want:  false,
	})

	// Add test for broken link
	brokenLinkChain := make([]Block, len(validChain))
	copy(brokenLinkChain, validChain)
	brokenLinkChain[2].PrevHash = "wronghash"
	tests = append(tests, struct {
		name  string
		chain []Block
		want  bool
	}{
		name:  "broken link",
		chain: brokenLinkChain,
		want:  false,
	})

	// Add test for invalid proof of work
	invalidProofChain := make([]Block, len(validChain))
	copy(invalidProofChain, validChain)
	invalidProofChain[2].Hash = "ffffffffff" // No leading zeros
	tests = append(tests, struct {
		name  string
		chain []Block
		want  bool
	}{
		name:  "invalid proof of work",
		chain: invalidProofChain,
		want:  false,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateChain(tt.chain, difficulty)
			if got != tt.want {
				t.Errorf("ValidateChain() = %v, want %v", got, tt.want)
			}

			// Verify against solution
			expected := ValidateChainSolution(tt.chain, difficulty)
			if got != expected {
				t.Errorf("Result doesn't match solution: got %v, want %v", got, expected)
			}
		})
	}
}

// Test Exercise 5: Adjust Difficulty
func TestAdjustDifficulty(t *testing.T) {
	tests := []struct {
		name              string
		targetBlockTime   int64
		currentDifficulty int
		blocks            []Block
		wantHigher        bool  // true if difficulty should increase
		wantLower         bool  // true if difficulty should decrease
		wantSame          bool  // true if difficulty should stay the same
	}{
		{
			name:              "mining too fast",
			targetBlockTime:   10,
			currentDifficulty: 3,
			blocks: []Block{
				{Timestamp: 1000},
				{Timestamp: 1002}, // 2 seconds (too fast)
				{Timestamp: 1004}, // 2 seconds
				{Timestamp: 1006}, // 2 seconds
				{Timestamp: 1008}, // 2 seconds
				{Timestamp: 1010}, // 2 seconds
				{Timestamp: 1012}, // 2 seconds
				{Timestamp: 1014}, // 2 seconds
				{Timestamp: 1016}, // 2 seconds
				{Timestamp: 1018}, // 2 seconds
				{Timestamp: 1020}, // 2 seconds - total 20 sec for 10 blocks
			},
			wantHigher: true,
		},
		{
			name:              "mining too slow",
			targetBlockTime:   10,
			currentDifficulty: 3,
			blocks: []Block{
				{Timestamp: 1000},
				{Timestamp: 1025}, // 25 seconds (too slow)
				{Timestamp: 1050}, // 25 seconds
				{Timestamp: 1075}, // 25 seconds
				{Timestamp: 1100}, // 25 seconds
				{Timestamp: 1125}, // 25 seconds
				{Timestamp: 1150}, // 25 seconds
				{Timestamp: 1175}, // 25 seconds
				{Timestamp: 1200}, // 25 seconds
				{Timestamp: 1225}, // 25 seconds
				{Timestamp: 1250}, // 25 seconds - total 250 sec for 10 blocks
			},
			wantLower: true,
		},
		{
			name:              "mining at target rate",
			targetBlockTime:   10,
			currentDifficulty: 3,
			blocks: []Block{
				{Timestamp: 1000},
				{Timestamp: 1010}, // 10 seconds
				{Timestamp: 1020}, // 10 seconds
				{Timestamp: 1030}, // 10 seconds
				{Timestamp: 1040}, // 10 seconds
				{Timestamp: 1050}, // 10 seconds
				{Timestamp: 1060}, // 10 seconds
				{Timestamp: 1070}, // 10 seconds
				{Timestamp: 1080}, // 10 seconds
				{Timestamp: 1090}, // 10 seconds
				{Timestamp: 1100}, // 10 seconds - total 100 sec for 10 blocks
			},
			wantSame: true,
		},
		{
			name:              "not enough blocks",
			targetBlockTime:   10,
			currentDifficulty: 3,
			blocks: []Block{
				{Timestamp: 1000},
			},
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AdjustDifficulty(tt.blocks, tt.targetBlockTime, tt.currentDifficulty)

			if tt.wantHigher && got <= tt.currentDifficulty {
				t.Errorf("AdjustDifficulty() = %d, want > %d (should increase)",
					got, tt.currentDifficulty)
			}
			if tt.wantLower && got >= tt.currentDifficulty {
				t.Errorf("AdjustDifficulty() = %d, want < %d (should decrease)",
					got, tt.currentDifficulty)
			}
			if tt.wantSame && got != tt.currentDifficulty {
				t.Errorf("AdjustDifficulty() = %d, want = %d (should stay same)",
					got, tt.currentDifficulty)
			}

			// Verify difficulty doesn't go below 1
			if got < 1 {
				t.Errorf("Difficulty = %d, should never be < 1", got)
			}

			// Verify against solution
			expected := AdjustDifficultySolution(tt.blocks, tt.targetBlockTime, tt.currentDifficulty)
			if got != expected {
				t.Errorf("Result doesn't match solution: got %d, want %d", got, expected)
			}

			t.Logf("Difficulty adjusted from %d to %d", tt.currentDifficulty, got)
		})
	}
}

// Test Exercise 6: Mining Probability
func TestMiningProbability(t *testing.T) {
	tests := []struct {
		name        string
		hashRate    float64
		difficulty  int
		timeSeconds float64
		wantMin     float64 // minimum expected probability
		wantMax     float64 // maximum expected probability
	}{
		{
			name:        "low probability",
			hashRate:    100,
			difficulty:  5,
			timeSeconds: 1,
			wantMin:     0.0,
			wantMax:     0.02,
		},
		{
			name:        "medium probability",
			hashRate:    1000,
			difficulty:  3,
			timeSeconds: 4,
			wantMin:     0.5,
			wantMax:     1.0,
		},
		{
			name:        "high probability",
			hashRate:    10000,
			difficulty:  2,
			timeSeconds: 10,
			wantMin:     0.99,
			wantMax:     1.0,
		},
		{
			name:        "zero time",
			hashRate:    1000,
			difficulty:  3,
			timeSeconds: 0,
			wantMin:     0.0,
			wantMax:     0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MiningProbability(tt.hashRate, tt.difficulty, tt.timeSeconds)

			// Check probability is in valid range [0, 1]
			if got < 0 || got > 1 {
				t.Errorf("Probability = %f, want in range [0, 1]", got)
			}

			// Check probability is in expected range
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("Probability = %f, want in range [%f, %f]",
					got, tt.wantMin, tt.wantMax)
			}

			// Verify against solution
			expected := MiningProbabilitySolution(tt.hashRate, tt.difficulty, tt.timeSeconds)
			if math.Abs(got-expected) > 0.0001 {
				t.Errorf("Result doesn't match solution: got %f, want %f", got, expected)
			}

			t.Logf("Hash rate: %.0f H/s, Difficulty: %d, Time: %.1f s → Probability: %.4f",
				tt.hashRate, tt.difficulty, tt.timeSeconds, got)
		})
	}
}

// Benchmark for mining at different difficulties
func BenchmarkMining(b *testing.B) {
	difficulties := []int{1, 2, 3, 4}

	for _, difficulty := range difficulties {
		b.Run(string(rune('0'+difficulty))+" zeros", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				block := Block{
					Index:     1,
					Timestamp: 1609459200,
					Data:      "Benchmark block",
					PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
					Nonce:     0,
				}
				MineBlock(&block, difficulty)
			}
		})
	}
}

// Benchmark hash calculation
func BenchmarkCalculateHash(b *testing.B) {
	block := Block{
		Index:     1,
		Timestamp: 1609459200,
		Data:      "Benchmark block",
		PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
		Nonce:     12345,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalculateBlockHash(block)
	}
}

// Test that demonstrates the avalanche effect
func TestAvalancheEffect(t *testing.T) {
	block := Block{
		Index:     1,
		Timestamp: 1609459200,
		Data:      "Original data",
		PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
		Nonce:     42,
	}

	originalHash := CalculateBlockHash(block)

	// Change just one character
	block.Data = "Original datb" // Changed last 'a' to 'b'
	modifiedHash := CalculateBlockHash(block)

	// Count how many characters are different
	differences := 0
	for i := 0; i < len(originalHash); i++ {
		if originalHash[i] != modifiedHash[i] {
			differences++
		}
	}

	// Avalanche effect means >50% of bits should change
	// For 64 hex characters, we expect ~32 to be different
	if differences < 20 {
		t.Errorf("Avalanche effect too weak: only %d/64 characters changed", differences)
	}

	t.Logf("Avalanche effect: %d/64 characters changed (%.1f%%)",
		differences, float64(differences)/64*100)
}

// Test difficulty scaling
func TestDifficultyScaling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping difficulty scaling test in short mode")
	}

	block := Block{
		Index:     1,
		Timestamp: 1609459200,
		Data:      "Test block",
		PrevHash:  "0000000000000000000000000000000000000000000000000000000000000000",
	}

	var previousAttempts int
	for difficulty := 1; difficulty <= 4; difficulty++ {
		testBlock := block
		attempts := MineBlock(&testBlock, difficulty)

		t.Logf("Difficulty %d: %d attempts (hash: %s...)",
			difficulty, attempts, testBlock.Hash[:16])

		// Each increment should roughly multiply attempts by 16
		if difficulty > 1 && previousAttempts > 0 {
			ratio := float64(attempts) / float64(previousAttempts)
			// Allow wide range because of randomness
			if ratio < 2 || ratio > 128 {
				t.Logf("Warning: Unexpected scaling ratio: %.2fx (expected ~16x)",
					ratio)
			}
		}

		previousAttempts = attempts
	}
}

// Helper function to verify chain integrity
func TestChainIntegrity(t *testing.T) {
	chain := CreateTestChain(5, 3)

	// Verify each block links correctly
	for i := 1; i < len(chain); i++ {
		if chain[i].PrevHash != chain[i-1].Hash {
			t.Errorf("Block %d doesn't link to block %d", i, i-1)
		}

		if chain[i].Index != i {
			t.Errorf("Block index mismatch: got %d, want %d", chain[i].Index, i)
		}

		// Verify hash has correct number of leading zeros
		if !strings.HasPrefix(chain[i].Hash, "000") {
			t.Errorf("Block %d hash doesn't have 3 leading zeros: %s",
				i, chain[i].Hash)
		}
	}

	t.Logf("Chain integrity verified for %d blocks", len(chain))
}
