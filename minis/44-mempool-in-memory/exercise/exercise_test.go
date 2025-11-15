package exercise

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Helper function to create test transactions
func createTestTx(from, to string, value, fee, nonce uint64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Value:     value,
		Fee:       fee,
		Nonce:     nonce,
		Timestamp: time.Now(),
	}
	tx.Hash = hashTx(tx)
	return tx
}

func hashTx(tx *Transaction) string {
	data := fmt.Sprintf("%s:%s:%d:%d:%d:%d",
		tx.From, tx.To, tx.Value, tx.Fee, tx.Nonce, tx.Timestamp.UnixNano())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// ============================================================================
// FIFO Mempool Tests
// ============================================================================

func TestFIFOMempool_Basic(t *testing.T) {
	mempool := NewFIFOMempool(5)

	// Test empty mempool
	if mempool.Size() != 0 {
		t.Errorf("Expected size 0, got %d", mempool.Size())
	}

	if mempool.GetNext() != nil {
		t.Error("Expected nil for empty mempool")
	}

	// Add transaction
	tx1 := createTestTx("Alice", "Bob", 100, 10, 0)
	err := mempool.Add(tx1)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	if mempool.Size() != 1 {
		t.Errorf("Expected size 1, got %d", mempool.Size())
	}

	// Get next should return first transaction
	next := mempool.GetNext()
	if next == nil || next.Hash != tx1.Hash {
		t.Error("GetNext should return first transaction")
	}

	// Remove transaction
	removed, err := mempool.Remove(tx1.Hash)
	if err != nil {
		t.Fatalf("Failed to remove transaction: %v", err)
	}
	if removed.Hash != tx1.Hash {
		t.Error("Removed wrong transaction")
	}

	if mempool.Size() != 0 {
		t.Errorf("Expected size 0 after removal, got %d", mempool.Size())
	}
}

func TestFIFOMempool_FIFOOrder(t *testing.T) {
	mempool := NewFIFOMempool(10)

	// Add transactions in specific order
	txs := []*Transaction{
		createTestTx("Alice", "Bob", 100, 10, 0),
		createTestTx("Charlie", "Dave", 50, 50, 0), // Higher fee
		createTestTx("Eve", "Frank", 200, 5, 0),    // Lower fee
	}

	for _, tx := range txs {
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
		err := mempool.Add(tx)
		if err != nil {
			t.Fatalf("Failed to add transaction: %v", err)
		}
	}

	// Should get transactions in FIFO order (not by fee)
	for i, expectedTx := range txs {
		next := mempool.GetNext()
		if next == nil {
			t.Fatalf("Expected transaction %d, got nil", i)
		}
		if next.Hash != expectedTx.Hash {
			t.Errorf("Transaction %d: expected %s, got %s", i, expectedTx.From, next.From)
		}
		mempool.Remove(next.Hash)
	}
}

func TestFIFOMempool_Capacity(t *testing.T) {
	capacity := 3
	mempool := NewFIFOMempool(capacity)

	// Fill to capacity
	for i := 0; i < capacity; i++ {
		tx := createTestTx(fmt.Sprintf("Sender%d", i), "Receiver", 100, 10, 0)
		err := mempool.Add(tx)
		if err != nil {
			t.Fatalf("Failed to add transaction %d: %v", i, err)
		}
	}

	// Try to add one more
	tx := createTestTx("Extra", "Receiver", 100, 100, 0)
	err := mempool.Add(tx)
	if err == nil {
		t.Error("Expected error when adding to full mempool")
	}
}

func TestFIFOMempool_Duplicate(t *testing.T) {
	mempool := NewFIFOMempool(5)

	tx := createTestTx("Alice", "Bob", 100, 10, 0)
	err := mempool.Add(tx)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	// Try to add same transaction again
	err = mempool.Add(tx)
	if err == nil {
		t.Error("Expected error when adding duplicate transaction")
	}
}

// ============================================================================
// Priority Mempool Tests
// ============================================================================

func TestPriorityMempool_Basic(t *testing.T) {
	mempool := NewPriorityMempool(5)

	// Test empty mempool
	if mempool.Size() != 0 {
		t.Errorf("Expected size 0, got %d", mempool.Size())
	}

	if mempool.GetNext() != nil {
		t.Error("Expected nil for empty mempool")
	}

	// Add transaction
	tx1 := createTestTx("Alice", "Bob", 100, 10, 0)
	err := mempool.Add(tx1)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	if mempool.Size() != 1 {
		t.Errorf("Expected size 1, got %d", mempool.Size())
	}

	// Get next should return highest priority
	next := mempool.GetNext()
	if next == nil || next.Hash != tx1.Hash {
		t.Error("GetNext should return transaction")
	}
}

func TestPriorityMempool_PriorityOrder(t *testing.T) {
	mempool := NewPriorityMempool(10)

	// Add transactions with different fees
	txs := []*Transaction{
		createTestTx("Alice", "Bob", 100, 10, 0),
		createTestTx("Charlie", "Dave", 50, 50, 0),  // Highest fee
		createTestTx("Eve", "Frank", 200, 5, 0),     // Lowest fee
		createTestTx("Grace", "Henry", 75, 25, 0),   // Medium fee
	}

	for _, tx := range txs {
		err := mempool.Add(tx)
		if err != nil {
			t.Fatalf("Failed to add transaction: %v", err)
		}
	}

	// Should get transactions in priority order (by fee)
	expectedOrder := []string{"Charlie", "Grace", "Alice", "Eve"}
	for i, expectedSender := range expectedOrder {
		next := mempool.GetNext()
		if next == nil {
			t.Fatalf("Expected transaction %d, got nil", i)
		}
		if next.From != expectedSender {
			t.Errorf("Transaction %d: expected from %s, got %s", i, expectedSender, next.From)
		}
		mempool.Remove(next.Hash)
	}
}

func TestPriorityMempool_Eviction(t *testing.T) {
	capacity := 3
	mempool := NewPriorityMempool(capacity)

	// Fill with low-fee transactions
	for i := 0; i < capacity; i++ {
		tx := createTestTx(fmt.Sprintf("Sender%d", i), "Receiver", 100, uint64(i+1), 0)
		err := mempool.Add(tx)
		if err != nil {
			t.Fatalf("Failed to add transaction %d: %v", i, err)
		}
	}

	// Try to add very low fee transaction (should be rejected)
	lowFeeTx := createTestTx("LowFee", "Receiver", 100, 0, 0)
	err := mempool.Add(lowFeeTx)
	if err == nil {
		t.Error("Expected error when adding low-fee transaction to full mempool")
	}

	// Add high fee transaction (should evict lowest)
	highFeeTx := createTestTx("HighFee", "Receiver", 100, 100, 0)
	err = mempool.Add(highFeeTx)
	if err != nil {
		t.Fatalf("Failed to add high-fee transaction: %v", err)
	}

	// Verify size is still at capacity
	if mempool.Size() != capacity {
		t.Errorf("Expected size %d, got %d", capacity, mempool.Size())
	}

	// Verify highest fee is at top
	next := mempool.GetNext()
	if next == nil || next.Fee != 100 {
		t.Error("Expected high-fee transaction at top")
	}
}

// ============================================================================
// Nonce Mempool Tests
// ============================================================================

func TestNonceMempool_Basic(t *testing.T) {
	mempool := NewNonceMempool()

	// Test empty mempool
	if mempool.Size() != 0 {
		t.Errorf("Expected size 0, got %d", mempool.Size())
	}

	// Add transaction
	tx1 := createTestTx("Alice", "Bob", 100, 10, 0)
	err := mempool.Add(tx1)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	if mempool.Size() != 1 {
		t.Errorf("Expected size 1, got %d", mempool.Size())
	}

	// Get next for account
	next := mempool.GetNextForAccount("Alice")
	if next == nil || next.Hash != tx1.Hash {
		t.Error("GetNextForAccount should return transaction")
	}

	// Advance nonce
	mempool.AdvanceNonce("Alice")

	// After advancing, should not find transaction
	next = mempool.GetNextForAccount("Alice")
	if next != nil {
		t.Error("Expected nil after advancing nonce")
	}
}

func TestNonceMempool_NonceOrdering(t *testing.T) {
	mempool := NewNonceMempool()

	// Add transactions with different nonces (out of order)
	txs := []*Transaction{
		createTestTx("Alice", "Bob", 100, 10, 2),   // Future nonce
		createTestTx("Alice", "Carol", 50, 15, 0),  // Current nonce
		createTestTx("Alice", "Dave", 75, 12, 1),   // Next nonce
		createTestTx("Alice", "Eve", 25, 20, 3),    // Future nonce
	}

	for _, tx := range txs {
		err := mempool.Add(tx)
		if err != nil {
			t.Fatalf("Failed to add transaction: %v", err)
		}
	}

	// Process in nonce order
	expectedReceivers := []string{"Carol", "Dave", "Bob", "Eve"}
	for i, expectedReceiver := range expectedReceivers {
		next := mempool.GetNextForAccount("Alice")
		if next == nil {
			t.Fatalf("Expected transaction %d, got nil", i)
		}
		if next.To != expectedReceiver {
			t.Errorf("Transaction %d: expected to %s, got %s", i, expectedReceiver, next.To)
		}
		mempool.AdvanceNonce("Alice")
	}
}

func TestNonceMempool_Replacement(t *testing.T) {
	mempool := NewNonceMempool()

	// Add transaction with nonce 0
	tx1 := createTestTx("Alice", "Bob", 100, 10, 0)
	err := mempool.Add(tx1)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	// Try to replace with lower fee (should fail)
	tx2 := createTestTx("Alice", "Carol", 100, 5, 0)
	err = mempool.Add(tx2)
	if err == nil {
		t.Error("Expected error when replacing with lower fee")
	}

	// Replace with higher fee (should succeed)
	tx3 := createTestTx("Alice", "Dave", 100, 20, 0)
	err = mempool.Add(tx3)
	if err != nil {
		t.Fatalf("Failed to replace transaction: %v", err)
	}

	// Verify replaced transaction
	next := mempool.GetNextForAccount("Alice")
	if next == nil || next.To != "Dave" {
		t.Error("Expected replaced transaction")
	}
}

func TestNonceMempool_MultipleAccounts(t *testing.T) {
	mempool := NewNonceMempool()

	// Add transactions for different accounts
	accounts := []string{"Alice", "Bob", "Charlie"}
	for _, account := range accounts {
		for nonce := uint64(0); nonce < 3; nonce++ {
			tx := createTestTx(account, "Receiver", 100, 10, nonce)
			err := mempool.Add(tx)
			if err != nil {
				t.Fatalf("Failed to add transaction: %v", err)
			}
		}
	}

	// Verify size
	expectedSize := len(accounts) * 3
	if mempool.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, mempool.Size())
	}

	// Process each account independently
	for _, account := range accounts {
		for nonce := uint64(0); nonce < 3; nonce++ {
			next := mempool.GetNextForAccount(account)
			if next == nil {
				t.Fatalf("Expected transaction for %s nonce %d", account, nonce)
			}
			if next.From != account || next.Nonce != nonce {
				t.Errorf("Expected %s nonce %d, got %s nonce %d",
					account, nonce, next.From, next.Nonce)
			}
			mempool.AdvanceNonce(account)
		}
	}
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestFIFOMempool_Concurrent(t *testing.T) {
	mempool := NewFIFOMempool(1000)
	var wg sync.WaitGroup

	// Multiple goroutines adding transactions
	numGoroutines := 10
	txsPerGoroutine := 20

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < txsPerGoroutine; j++ {
				tx := createTestTx(
					fmt.Sprintf("Sender%d", id),
					fmt.Sprintf("Receiver%d", j),
					100, 10, uint64(j),
				)
				mempool.Add(tx)
			}
		}(i)
	}

	wg.Wait()

	// Verify all transactions were added (within capacity)
	size := mempool.Size()
	expected := numGoroutines * txsPerGoroutine
	if size != expected {
		t.Errorf("Expected size %d, got %d", expected, size)
	}
}

func TestPriorityMempool_Concurrent(t *testing.T) {
	mempool := NewPriorityMempool(1000)
	var wg sync.WaitGroup

	numGoroutines := 10
	txsPerGoroutine := 20

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < txsPerGoroutine; j++ {
				tx := createTestTx(
					fmt.Sprintf("Sender%d", id),
					fmt.Sprintf("Receiver%d", j),
					100, uint64(j+1), uint64(j),
				)
				mempool.Add(tx)
			}
		}(i)
	}

	wg.Wait()

	size := mempool.Size()
	expected := numGoroutines * txsPerGoroutine
	if size != expected {
		t.Errorf("Expected size %d, got %d", expected, size)
	}
}

func TestNonceMempool_Concurrent(t *testing.T) {
	mempool := NewNonceMempool()
	var wg sync.WaitGroup

	numAccounts := 10
	txsPerAccount := 20

	for i := 0; i < numAccounts; i++ {
		wg.Add(1)
		go func(accountID int) {
			defer wg.Done()
			account := fmt.Sprintf("Account%d", accountID)
			for nonce := uint64(0); nonce < uint64(txsPerAccount); nonce++ {
				tx := createTestTx(account, "Receiver", 100, 10, nonce)
				mempool.Add(tx)
			}
		}(i)
	}

	wg.Wait()

	size := mempool.Size()
	expected := numAccounts * txsPerAccount
	if size != expected {
		t.Errorf("Expected size %d, got %d", expected, size)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkFIFOMempool_Add(b *testing.B) {
	mempool := NewFIFOMempool(b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tx := createTestTx(fmt.Sprintf("Sender%d", i), "Receiver", 100, 10, 0)
		mempool.Add(tx)
	}
}

func BenchmarkPriorityMempool_Add(b *testing.B) {
	mempool := NewPriorityMempool(b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tx := createTestTx(fmt.Sprintf("Sender%d", i), "Receiver", 100, uint64(i), 0)
		mempool.Add(tx)
	}
}

func BenchmarkNonceMempool_Add(b *testing.B) {
	mempool := NewNonceMempool()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tx := createTestTx("Alice", "Receiver", 100, 10, uint64(i))
		mempool.Add(tx)
	}
}
