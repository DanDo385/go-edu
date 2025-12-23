//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "container/heap" for priority queue implementation
// - "errors" for error handling
// - "sync" for RWMutex (thread-safe concurrent access)
// - "time" for transaction timestamps
//
// import (
//     "container/heap"
//     "errors"
//     "sync"
//     "time"
// )

// ============================================================================
// IN-MEMORY TRANSACTION MEMPOOL: Managing Pending Blockchain Transactions
// ============================================================================
//
// A mempool (memory pool) is a critical component of blockchain systems:
// - Stores pending transactions before they're included in a block
// - Implements eviction policies when full (FIFO, priority-based, or nonce-based)
// - Must be thread-safe for concurrent reads/writes
// - Balances between O(1) lookups and ordered processing
//
// Why mempools are important:
// - Bitcoin/Ethereum nodes maintain thousands of pending transactions
// - Miners select transactions based on fees (priority)
// - Network propagates transactions between nodes
// - Must prevent duplicate transactions and handle replacements
//
// Three Implementation Patterns:
// 1. FIFO (First-In-First-Out): Simple, fair, but ignores transaction fees
// 2. Priority: Fee-based ordering, maximizes miner revenue
// 3. Nonce-based: Maintains sequential ordering per account
//
// Memory Management Considerations:
// - Each transaction is 100-500 bytes (hash, addresses, value, fee, nonce)
// - Typical mempool: 10,000-50,000 transactions = 5-25 MB
// - Map overhead: ~48 bytes per entry (key + value + hash table metadata)
// - Slice overhead: 24 bytes (pointer, length, capacity)
// - Mutex: ~16 bytes (platform-dependent)
//
// ============================================================================

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash      string    // Unique transaction identifier (32-byte hex string)
	From      string    // Sender address (20-byte hex string)
	To        string    // Receiver address (20-byte hex string)
	Value     uint64    // Amount to transfer (in smallest unit, e.g., wei)
	Fee       uint64    // Transaction fee (used for prioritization)
	Nonce     uint64    // Account nonce (for ordering, prevents replay attacks)
	Timestamp time.Time // When transaction was created
}

// ============================================================================
// Exercise 1: FIFO Mempool (First-In-First-Out)
// ============================================================================

// FIFOMempool implements a first-in-first-out transaction pool.
type FIFOMempool struct {
	mu       sync.RWMutex
	txs      []*Transaction
	txMap    map[string]*Transaction
	capacity int
}

// NewFIFOMempool creates a new FIFO mempool with the given capacity.
func NewFIFOMempool(capacity int) *FIFOMempool {
	// TODO: Implement this function.
	// - Initialize the `FIFOMempool` struct.
	// - `txs` slice should be created with a capacity of `capacity` to reduce re-allocations.
	// - `txMap` should be initialized as an empty map.
	return &FIFOMempool{
		txs:      make([]*Transaction, 0, capacity),
		txMap:    make(map[string]*Transaction),
		capacity: capacity,
	}
}

// Add adds a transaction to the mempool.
func (m *FIFOMempool) Add(tx *Transaction) error {
	// TODO: Implement this thread-safe add operation.

	// Step 1: Acquire a write lock since you're modifying the data structures.
	// - `m.mu.Lock()`
	// - `defer m.mu.Unlock()`

	// Step 2: Check if the transaction already exists using the `txMap`.
	// - This provides an O(1) lookup.
	// - If it exists, return an error.

	// Step 3: Check if the mempool is at capacity.
	// - `if len(m.txs) >= m.capacity`
	// - If full, return an error.

	// Step 4: Add the new transaction to both the slice and the map.
	// - `m.txs = append(m.txs, tx)`
	// - `m.txMap[tx.Hash] = tx`
	return nil
}

// Remove removes a transaction from the mempool by hash.
func (m *FIFOMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement this thread-safe remove operation.

	// Step 1: Acquire a write lock.
	// Step 2: Check if the transaction exists in the `txMap`. If not, return an error.
	// Step 3: Find the transaction in the `txs` slice. This will be an O(n) linear scan.
	// Step 4: Once found at index `i`, remove it from the slice.
	// - The `append(m.txs[:i], m.txs[i+1:]...)` trick is efficient for this.
	// Step 5: Delete the transaction from the `txMap`.
	// Step 6: Return the removed transaction.
	return nil, nil
}

// GetNext returns the next transaction to process (oldest) without removing it.
func (m *FIFOMempool) GetNext() *Transaction {
	// TODO: Implement this thread-safe read operation.
	// - Use a read lock (`m.mu.RLock()`) because you are not modifying data.
	// - If the `txs` slice is empty, return `nil`.
	// - Otherwise, return the first element (`m.txs[0]`).
	return nil
}

// Size returns the current number of transactions in the mempool.
func (m *FIFOMempool) Size() int {
	// TODO: Implement this thread-safe size check.
	// - Use a read lock.
	// - Return the length of the `txs` slice.
	return 0
}

// ============================================================================
// Exercise 2: Priority Mempool (Fee-Based Ordering)
// ============================================================================

// PriorityMempool implements a priority-based transaction pool using a heap.
type PriorityMempool struct {
	mu       sync.RWMutex
	heap     *TxHeap
	txMap    map[string]int
	capacity int
}

// TxHeap implements heap.Interface for a priority queue of transactions.
type TxHeap []*Transaction

func (h TxHeap) Len() int { return len(h) }

func (h TxHeap) Less(i, j int) bool {
	// This defines the priority. We want a "max-heap" based on the fee.
	// `Less` should return true if `i` has higher priority than `j`.
	if h[i].Fee != h[j].Fee {
		return h[i].Fee > h[j].Fee
	}
	// If fees are equal, the earlier transaction has higher priority.
	return h[i].Timestamp.Before(h[j].Timestamp)
}

func (h TxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *TxHeap) Push(x interface{}) {
	*h = append(*h, x.(*Transaction))
}

func (h *TxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	tx := old[n-1]
	*h = old[0 : n-1]
	return tx
}

// NewPriorityMempool creates a new priority mempool.
func NewPriorityMempool(capacity int) *PriorityMempool {
	// TODO: Implement this function.
	// - Create a new `TxHeap`.
	// - Initialize it as a heap using `heap.Init(h)`.
	// - Return a new `PriorityMempool` with the heap, an initialized map, and the capacity.
	return nil
}

// Add adds a transaction to the mempool.
func (m *PriorityMempool) Add(tx *Transaction) error {
	// TODO: Implement this thread-safe add operation with eviction logic.

	// Step 1: Acquire a write lock.
	// Step 2: Check for duplicates in `txMap`.
	// Step 3: If the mempool is at capacity:
	//   - Look at the lowest-priority item without removing it. In a max-heap, this is one of the leaves. For simplicity, you can peek at the last element if your heap implementation allows, though it's not guaranteed to be the absolute minimum. A better way is to check `(*m.heap)[0]` if it were a min-heap. For a max-heap, you might need to search the last `n/2` elements. A simpler approach for this exercise is to just check against the fee of the element at the end of the slice.
	//   - If the new transaction's fee is not higher than the lowest fee in the pool, reject it.
	//   - If it is higher, evict the lowest-priority item using `heap.Pop()`. Don't forget to remove it from `txMap`.
	// Step 4: Add the new transaction to the heap using `heap.Push()`.
	// Step 5: Update the `txMap` to store the index of the new item in the heap.
	return nil
}

// Remove removes a transaction from the mempool by hash.
func (m *PriorityMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement this thread-safe remove operation.

	// Removing an arbitrary element from a heap is more complex than a slice.
	// Step 1: Acquire a write lock.
	// Step 2: Look up the transaction's *index* in the `txMap`. If it doesn't exist, return an error.
	// Step 3: Use `heap.Remove(m.heap, index)` to remove the element. This function will maintain the heap property.
	// Step 4: Delete the hash from `txMap`.
	// Step 5: The `heap.Remove` operation might have swapped elements, invalidating other indices in your `txMap`. You need to re-index the map or find a way to update it. For this exercise, a simple but slow solution is to rebuild the index map after a removal. A better solution involves updating the index of the swapped element.
	return nil, nil
}

// GetNext returns the highest priority transaction without removing it.
func (m *PriorityMempool) GetNext() *Transaction {
	// TODO: Implement this thread-safe peek operation.
	// - Use a read lock.
	// - If the heap is empty, return `nil`.
	// - The highest-priority item is always at the root of the heap, which is index 0 of the slice. Return `(*m.heap)[0]`.
	return nil
}

// Size returns the current number of transactions in the mempool.
func (m *PriorityMempool) Size() int {
	// TODO: Implement this thread-safe size check.
	// - Use a read lock.
	// - Return `m.heap.Len()`.
	return 0
}

// ============================================================================
// Exercise 3: Nonce-Based Mempool (Account Ordering)
// ============================================================================

// NonceMempool implements a nonce-based transaction pool.
type NonceMempool struct {
	mu       sync.RWMutex
	accounts map[string]*AccountQueue
}

// AccountQueue stores transactions for a single account, ordered by nonce.
type AccountQueue struct {
	address      string
	pendingNonce uint64
	txs          map[uint64]*Transaction
}

// NewNonceMempool creates a new nonce-based mempool.
func NewNonceMempool() *NonceMempool {
	// TODO: Implement this function.
	// - Initialize the `NonceMempool` with an empty map of accounts.
	return nil
}

// Add adds a transaction to the mempool.
func (m *NonceMempool) Add(tx *Transaction) error {
	// TODO: Implement this thread-safe add operation with replacement logic.

	// Step 1: Acquire a write lock.
	// Step 2: Get the `AccountQueue` for the transaction's sender (`tx.From`).
	// Step 3: If the queue doesn't exist, create a new one and add it to the `m.accounts` map.
	// Step 4: Check if a transaction with the same nonce already exists in the queue.
	//   - `if existing, exists := queue.txs[tx.Nonce]; exists`
	//   - If it exists, only replace it if the new transaction has a strictly higher fee (`tx.Fee > existing.Fee`). This is "Replace-by-Fee". If not, return an error.
	// Step 5: Add the new transaction to the `queue.txs` map, keyed by its nonce.
	return nil
}

// GetNextForAccount returns the next transaction for the given account.
func (m *NonceMempool) GetNextForAccount(address string) *Transaction {
	// TODO: Implement this thread-safe read operation.
	// - Use a read lock.
	// - Get the `AccountQueue` for the given `address`.
	// - If the queue exists, return the transaction from its `txs` map that corresponds to the `queue.pendingNonce`.
	// - If the queue or the specific transaction doesn't exist, return `nil`.
	return nil
}

// AdvanceNonce advances the pending nonce for the given account.
func (m *NonceMempool) AdvanceNonce(address string) {
	// TODO: Implement this thread-safe method.
	// - This should be called after a transaction for an account has been successfully processed and included in a block.
	// - Use a write lock.
	// - Get the `AccountQueue` for the address.
	// - If it exists, delete the transaction corresponding to the `pendingNonce` and then increment `pendingNonce`.
}

// Size returns the total number of transactions across all accounts.
func (m *NonceMempool) Size() int {
	// TODO: Implement this thread-safe size check.
	// - Use a read lock.
	// - Iterate through all the `AccountQueue`s in the `m.accounts` map and sum the lengths of their `txs` maps.
	return 0
}

// ============================================================================
// After implementing all functions:
// - Run: go test -v ./...
// - Run: go test -race ./... (CRITICAL: check for data races!)
// - Compare with solution.go to see detailed implementations
// - Experiment: Try adding 10,000 transactions and measure memory usage
// - Benchmark: Compare FIFO vs Priority for different workloads
//
// Performance Comparison (10,000 transactions):
// FIFO:     Add: O(1), Remove: O(n), GetNext: O(1), Memory: ~500 KB
// Priority: Add: O(log n), Remove: O(n), GetNext: O(1), Memory: ~600 KB
// Nonce:    Add: O(1), Remove: O(1), GetNext: O(1), Memory: ~800 KB
//
// When to use each:
// - FIFO: Simple, fair, but ignores fees (not used in production)
// - Priority: Maximizes miner revenue, common in Bitcoin/Ethereum
// - Nonce: Account-based ordering, prevents replay attacks (Ethereum)
// ============================================================================
