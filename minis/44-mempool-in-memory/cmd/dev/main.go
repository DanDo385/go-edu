package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/example/go-10x-minis/minis/44-mempool-in-memory/exercise"
)

func main() {
	fmt.Println("=== Mempool In-Memory Demo ===\n")

	// Demo 1: FIFO Mempool
	demo1FIFOMempool()
	fmt.Println()

	// Demo 2: Priority Mempool
	demo2PriorityMempool()
	fmt.Println()

	// Demo 3: Nonce-Based Mempool
	demo3NonceMempool()
	fmt.Println()

	// Demo 4: Concurrent Access
	demo4ConcurrentAccess()
	fmt.Println()

	// Demo 5: Eviction Policies
	demo5EvictionPolicies()
}

// Demo 1: FIFO Mempool (First-In-First-Out)
func demo1FIFOMempool() {
	fmt.Println("--- Demo 1: FIFO Mempool ---")
	fmt.Println("Transactions are processed in arrival order (like a queue)")
	fmt.Println()

	mempool := exercise.NewFIFOMempool(5)

	// Create transactions with different fees
	txs := []*exercise.Transaction{
		createTx("Alice", "Bob", 100, 10, 0),
		createTx("Charlie", "Dave", 50, 50, 0),  // Higher fee but arrives later
		createTx("Eve", "Frank", 200, 5, 0),
	}

	// Add transactions
	for _, tx := range txs {
		err := mempool.Add(tx)
		if err != nil {
			log.Printf("Error adding tx: %v", err)
			continue
		}
		fmt.Printf("Added: %s -> %s (value=%d, fee=%d)\n",
			tx.From, tx.To, tx.Value, tx.Fee)
		time.Sleep(100 * time.Millisecond) // Simulate arrival time
	}

	fmt.Println()
	fmt.Printf("Mempool size: %d\n", mempool.Size())
	fmt.Println("\nProcessing order (FIFO):")

	// Process in FIFO order
	for i := 0; i < 3; i++ {
		tx := mempool.GetNext()
		if tx == nil {
			break
		}
		fmt.Printf("%d. %s -> %s (fee=%d) - First in = first out!\n",
			i+1, tx.From, tx.To, tx.Fee)
		mempool.Remove(tx.Hash)
	}

	fmt.Println("\nNote: Charlie's high-fee tx processed second, not first!")
}

// Demo 2: Priority Mempool (Highest Fee First)
func demo2PriorityMempool() {
	fmt.Println("--- Demo 2: Priority Mempool ---")
	fmt.Println("Transactions are processed by priority (highest fee first)")
	fmt.Println()

	mempool := exercise.NewPriorityMempool(5)

	// Create transactions with different fees
	txs := []*exercise.Transaction{
		createTx("Alice", "Bob", 100, 10, 0),
		createTx("Charlie", "Dave", 50, 50, 0),   // Highest fee
		createTx("Eve", "Frank", 200, 5, 0),      // Lowest fee
		createTx("Grace", "Henry", 75, 25, 0),    // Medium fee
	}

	// Add transactions in random order
	for _, tx := range txs {
		err := mempool.Add(tx)
		if err != nil {
			log.Printf("Error adding tx: %v", err)
			continue
		}
		fmt.Printf("Added: %s -> %s (fee=%d)\n", tx.From, tx.To, tx.Fee)
	}

	fmt.Println()
	fmt.Printf("Mempool size: %d\n", mempool.Size())
	fmt.Println("\nProcessing order (by priority):")

	// Process in priority order
	for i := 0; i < 4; i++ {
		tx := mempool.GetNext()
		if tx == nil {
			break
		}
		fmt.Printf("%d. %s -> %s (fee=%d) - Highest fee processed first!\n",
			i+1, tx.From, tx.To, tx.Fee)
		mempool.Remove(tx.Hash)
	}

	fmt.Println("\nNote: Processing order is by fee, not arrival time!")
}

// Demo 3: Nonce-Based Mempool (Account Ordering)
func demo3NonceMempool() {
	fmt.Println("--- Demo 3: Nonce-Based Mempool ---")
	fmt.Println("Transactions from same account must be processed in nonce order")
	fmt.Println()

	mempool := exercise.NewNonceMempool()

	// Alice sends multiple transactions (out of order)
	aliceTxs := []*exercise.Transaction{
		createTx("Alice", "Bob", 100, 10, 2),   // Nonce 2 (future)
		createTx("Alice", "Carol", 50, 15, 0),  // Nonce 0 (ready)
		createTx("Alice", "Dave", 75, 12, 1),   // Nonce 1 (next)
		createTx("Alice", "Eve", 25, 20, 3),    // Nonce 3 (future)
	}

	fmt.Println("Adding Alice's transactions (out of order):")
	for _, tx := range aliceTxs {
		err := mempool.Add(tx)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		fmt.Printf("  Nonce %d: Alice -> %s (value=%d)\n", tx.Nonce, tx.To, tx.Value)
	}

	fmt.Println("\nProcessing order (must follow nonce sequence):")

	for nonce := uint64(0); nonce < 4; nonce++ {
		tx := mempool.GetNextForAccount("Alice")
		if tx == nil {
			fmt.Printf("  Nonce %d: Waiting... (gap in sequence)\n", nonce)
			continue
		}
		fmt.Printf("  Nonce %d: Alice -> %s (value=%d) ✓\n", tx.Nonce, tx.To, tx.Value)
		mempool.AdvanceNonce("Alice")
	}

	fmt.Println("\nNote: Nonce 2 can't process until nonce 0 and 1 complete!")
}

// Demo 4: Concurrent Access
func demo4ConcurrentAccess() {
	fmt.Println("--- Demo 4: Concurrent Access ---")
	fmt.Println("Multiple goroutines safely accessing mempool concurrently")
	fmt.Println()

	mempool := exercise.NewPriorityMempool(100)

	var wg sync.WaitGroup
	numProducers := 5
	txsPerProducer := 10

	fmt.Printf("Starting %d goroutines, each adding %d transactions...\n",
		numProducers, txsPerProducer)

	start := time.Now()

	// Producer goroutines
	for i := 0; i < numProducers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < txsPerProducer; j++ {
				tx := createTx(
					fmt.Sprintf("Sender%d", id),
					fmt.Sprintf("Receiver%d", id),
					uint64(rand.Intn(1000)),
					uint64(rand.Intn(100)+1),
					uint64(j),
				)
				mempool.Add(tx)
			}
		}(i)
	}

	// Consumer goroutine
	wg.Add(1)
	consumedCount := 0
	go func() {
		defer wg.Done()
		for consumedCount < numProducers*txsPerProducer {
			if tx := mempool.GetNext(); tx != nil {
				mempool.Remove(tx.Hash)
				consumedCount++
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("\nCompleted in %v\n", elapsed)
	fmt.Printf("Processed %d transactions concurrently\n", numProducers*txsPerProducer)
	fmt.Printf("Final mempool size: %d\n", mempool.Size())
	fmt.Println("\nNo data races! All operations were thread-safe.")
}

// Demo 5: Eviction Policies
func demo5EvictionPolicies() {
	fmt.Println("--- Demo 5: Eviction Policies ---")
	fmt.Println("What happens when mempool reaches capacity?")
	fmt.Println()

	capacity := 5
	mempool := exercise.NewPriorityMempool(capacity)

	fmt.Printf("Mempool capacity: %d transactions\n\n", capacity)

	// Fill mempool to capacity
	fmt.Println("Filling mempool to capacity:")
	for i := 0; i < capacity; i++ {
		tx := createTx(
			fmt.Sprintf("Sender%d", i),
			fmt.Sprintf("Receiver%d", i),
			100,
			uint64((i+1)*10), // Increasing fees
			0,
		)
		mempool.Add(tx)
		fmt.Printf("  %d. Added tx with fee=%d (size: %d/%d)\n",
			i+1, tx.Fee, mempool.Size(), capacity)
	}

	fmt.Println("\nMempool is full! Trying to add more transactions...")

	// Try to add low-priority transaction
	lowPriorityTx := createTx("NewSender", "NewReceiver", 100, 5, 0)
	err := mempool.Add(lowPriorityTx)
	if err != nil {
		fmt.Printf("  ✗ Low-fee transaction (fee=5) rejected: %v\n", err)
	}

	// Try to add high-priority transaction
	highPriorityTx := createTx("VIP", "Receiver", 100, 100, 0)
	err = mempool.Add(highPriorityTx)
	if err == nil {
		fmt.Printf("  ✓ High-fee transaction (fee=100) accepted!\n")
		fmt.Printf("    Lowest-fee transaction was evicted.\n")
	}

	fmt.Println("\nEviction policy: Reject low-priority, evict when high-priority arrives")
}

// Helper functions

func createTx(from, to string, value, fee, nonce uint64) *exercise.Transaction {
	tx := &exercise.Transaction{
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

func hashTx(tx *exercise.Transaction) string {
	data := fmt.Sprintf("%s:%s:%d:%d:%d:%d",
		tx.From, tx.To, tx.Value, tx.Fee, tx.Nonce, tx.Timestamp.UnixNano())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes for brevity
}
