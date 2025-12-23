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
// Transactions are processed in arrival order (fair but inefficient).
//
// TODO: Add fields for:
// - Mutex for thread safety (sync.RWMutex preferred for read-heavy workloads)
// - Slice to store transactions in arrival order (maintains insertion order)
// - Map for O(1) lookup by hash (avoids duplicate insertions)
// - Capacity limit (prevents unbounded growth)
//
// Memory layout example (capacity=1000):
// - txs: []Transaction with cap=1000 → 24 bytes (slice header)
// - txMap: map[string]*Transaction → 48 bytes * 1000 entries = 48 KB
// - mu: sync.RWMutex → 16 bytes
// - Total overhead: ~50 KB + transaction data
//
// type FIFOMempool struct {
//     mu       sync.RWMutex               // Protects txs and txMap
//     txs      []*Transaction             // Ordered list (insertion order)
//     txMap    map[string]*Transaction    // Hash → Transaction (fast lookup)
//     capacity int                        // Maximum number of transactions
// }

// NewFIFOMempool creates a new FIFO mempool with the given capacity.
// TODO: Initialize FIFOMempool
//
// Memory allocation strategy:
// - Pre-allocate slice with capacity to avoid repeated reallocations
// - Initialize map with expected size hint (if Go supported it)
// - Return pointer to avoid copying the large struct
//
// func NewFIFOMempool(capacity int) *FIFOMempool {
//     return &FIFOMempool{
//         txs:      make([]*Transaction, 0, capacity),  // Pre-allocate slice
//         txMap:    make(map[string]*Transaction),      // Empty map
//         capacity: capacity,
//     }
// }
func NewFIFOMempool(capacity int) *FIFOMempool {
	// TODO: Implement
	return nil
}

// Add adds a transaction to the mempool.
// Returns error if transaction already exists or mempool is full.
//
// TODO: Implement thread-safe add operation
//
// Algorithm:
// 1. Acquire write lock (mu.Lock())
// 2. Check if hash already exists in txMap → return error
// 3. Check if len(txs) >= capacity → return error
// 4. Append to txs slice
// 5. Add to txMap for O(1) lookup
// 6. Release lock (mu.Unlock())
//
// Why use both slice and map?
// - Slice: Maintains insertion order for FIFO processing (O(n) for search)
// - Map: Fast duplicate detection (O(1) lookup by hash)
// - Trade-off: Extra memory for better time complexity
//
// func (m *FIFOMempool) Add(tx *Transaction) error {
//     m.mu.Lock()
//     defer m.mu.Unlock()
//
//     // Check for duplicates
//     if _, exists := m.txMap[tx.Hash]; exists {
//         return errors.New("transaction already exists")
//     }
//
//     // Check capacity
//     if len(m.txs) >= m.capacity {
//         return errors.New("mempool full")
//     }
//
//     // Add to both slice and map
//     m.txs = append(m.txs, tx)
//     m.txMap[tx.Hash] = tx
//
//     return nil
// }
func (m *FIFOMempool) Add(tx *Transaction) error {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Check if transaction already exists
	// 3. Check if mempool is at capacity
	// 4. Add to slice and map
	return nil
}

// Remove removes a transaction from the mempool by hash.
// Returns the removed transaction or error if not found.
//
// TODO: Implement thread-safe remove operation
//
// Algorithm:
// 1. Acquire write lock
// 2. Look up transaction in map → return error if not found
// 3. Find transaction in slice (linear search, O(n))
// 4. Remove from slice using append trick: txs = append(txs[:i], txs[i+1:]...)
// 5. Delete from map
// 6. Return removed transaction
//
// Slice removal performance:
// - Worst case: O(n) for find + O(n) for shift = O(n)
// - Memory: No reallocation, just shifts elements left
// - Alternative: Use linked list for O(1) removal, but more complex
//
// func (m *FIFOMempool) Remove(hash string) (*Transaction, error) {
//     m.mu.Lock()
//     defer m.mu.Unlock()
//
//     // Check if exists
//     tx, exists := m.txMap[hash]
//     if !exists {
//         return nil, errors.New("transaction not found")
//     }
//
//     // Find and remove from slice
//     for i, t := range m.txs {
//         if t.Hash == hash {
//             m.txs = append(m.txs[:i], m.txs[i+1:]...)
//             break
//         }
//     }
//
//     // Remove from map
//     delete(m.txMap, hash)
//     return tx, nil
// }
func (m *FIFOMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Find transaction in map
	// 3. Remove from slice (maintain order)
	// 4. Remove from map
	return nil, nil
}

// GetNext returns the next transaction to process (oldest) without removing it.
// Returns nil if mempool is empty.
//
// TODO: Implement thread-safe read operation
//
// Why use RLock instead of Lock?
// - Multiple goroutines can read simultaneously (read lock is shared)
// - Writers are blocked while readers are active
// - Better performance for read-heavy workloads
//
// func (m *FIFOMempool) GetNext() *Transaction {
//     m.mu.RLock()
//     defer m.mu.RUnlock()
//
//     if len(m.txs) == 0 {
//         return nil
//     }
//
//     return m.txs[0]  // First transaction (oldest)
// }
func (m *FIFOMempool) GetNext() *Transaction {
	// TODO: Implement
	// 1. Read lock the mutex
	// 2. Return first transaction in slice
	return nil
}

// Size returns the current number of transactions in the mempool.
//
// TODO: Implement thread-safe size check
//
// func (m *FIFOMempool) Size() int {
//     m.mu.RLock()
//     defer m.mu.RUnlock()
//     return len(m.txs)
// }
func (m *FIFOMempool) Size() int {
	// TODO: Implement with read lock
	return 0
}

// ============================================================================
// Exercise 2: Priority Mempool (Fee-Based Ordering)
// ============================================================================

// PriorityMempool implements a priority-based transaction pool.
// Transactions are ordered by fee (highest first), using a heap data structure.
//
// TODO: Add fields for:
// - Mutex for thread safety
// - Heap to store transactions by priority (container/heap)
// - Map for O(1) lookup by hash
// - Capacity limit
//
// Why use a heap?
// - Heap: O(log n) insert/remove, O(1) peek at max
// - Sorted slice: O(n) insert, O(1) peek, O(n) remove
// - Unsorted slice: O(1) insert, O(n) peek/remove
// - Trade-off: Heap is best for priority queue operations
//
// Heap memory layout:
// - Heap is backed by a slice: []*Transaction
// - Parent at index i, children at 2i+1 and 2i+2
// - Max-heap property: parent.Fee >= children.Fee
//
// type PriorityMempool struct {
//     mu       sync.RWMutex
//     heap     *TxHeap                    // Priority queue (max-heap by fee)
//     txMap    map[string]int             // Hash → heap index (for O(1) removal)
//     capacity int
// }
//
// TxHeap implements heap.Interface for transactions.
// type TxHeap []*Transaction
//
// func (h TxHeap) Len() int { return len(h) }
// func (h TxHeap) Less(i, j int) bool {
//     // Max-heap: higher fee = higher priority
//     if h[i].Fee != h[j].Fee {
//         return h[i].Fee > h[j].Fee
//     }
//     // Tiebreaker: earlier timestamp
//     return h[i].Timestamp.Before(h[j].Timestamp)
// }
// func (h TxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
// func (h *TxHeap) Push(x interface{}) { *h = append(*h, x.(*Transaction)) }
// func (h *TxHeap) Pop() interface{} {
//     old := *h
//     n := len(old)
//     tx := old[n-1]
//     *h = old[0 : n-1]
//     return tx
// }

// NewPriorityMempool creates a new priority mempool with the given capacity.
//
// TODO: Initialize PriorityMempool
// Hint: Use container/heap package and heap.Init()
//
// func NewPriorityMempool(capacity int) *PriorityMempool {
//     h := &TxHeap{}
//     heap.Init(h)
//     return &PriorityMempool{
//         heap:     h,
//         txMap:    make(map[string]int),
//         capacity: capacity,
//     }
// }
func NewPriorityMempool(capacity int) *PriorityMempool {
	// TODO: Initialize PriorityMempool
	// Hint: Use container/heap package
	return nil
}

// Add adds a transaction to the mempool.
// If mempool is full, evicts lowest priority transaction if new tx has higher priority.
//
// TODO: Implement priority-based add with eviction
//
// Algorithm:
// 1. Lock mutex
// 2. Check if already exists → return error
// 3. If at capacity:
//    a. Peek at lowest priority tx (last element in heap)
//    b. If new tx fee <= lowest fee → return error (too low priority)
//    c. Otherwise, evict lowest priority: heap.Remove(heap, len(heap)-1)
// 4. Add to heap: heap.Push(heap, tx)
// 5. Update map with heap index
//
// Why evict lowest instead of rejecting?
// - Allows higher-fee transactions to replace lower-fee ones
// - Maximizes potential miner revenue
// - Common in Bitcoin/Ethereum mempools
//
// func (m *PriorityMempool) Add(tx *Transaction) error {
//     m.mu.Lock()
//     defer m.mu.Unlock()
//
//     if _, exists := m.txMap[tx.Hash]; exists {
//         return errors.New("transaction already exists")
//     }
//
//     // If at capacity, check if we should evict
//     if m.heap.Len() >= m.capacity {
//         lowest := (*m.heap)[m.heap.Len()-1]
//         if tx.Fee <= lowest.Fee {
//             return errors.New("mempool full, transaction priority too low")
//         }
//         // Evict lowest priority
//         evicted := heap.Remove(m.heap, m.heap.Len()-1).(*Transaction)
//         delete(m.txMap, evicted.Hash)
//     }
//
//     // Add to heap
//     heap.Push(m.heap, tx)
//     m.txMap[tx.Hash] = m.heap.Len() - 1
//
//     return nil
// }
func (m *PriorityMempool) Add(tx *Transaction) error {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Check if transaction already exists
	// 3. If at capacity, check if new tx has higher priority than lowest
	// 4. Add to heap and map
	return nil
}

// Remove removes a transaction from the mempool by hash.
//
// TODO: Implement heap-based removal
//
// Heap removal is tricky:
// - heap.Remove(h, i) removes element at index i
// - After removal, heap property is restored (O(log n))
// - But our map tracks indices, which become stale after removal
// - Solution: Linear search to find current index (O(n))
//
// func (m *PriorityMempool) Remove(hash string) (*Transaction, error) {
//     m.mu.Lock()
//     defer m.mu.Unlock()
//
//     _, exists := m.txMap[hash]
//     if !exists {
//         return nil, errors.New("transaction not found")
//     }
//
//     // Find current index (indices may have shifted)
//     actualIdx := -1
//     for i, tx := range *m.heap {
//         if tx.Hash == hash {
//             actualIdx = i
//             break
//         }
//     }
//
//     if actualIdx == -1 {
//         return nil, errors.New("transaction not found in heap")
//     }
//
//     tx := heap.Remove(m.heap, actualIdx).(*Transaction)
//     delete(m.txMap, hash)
//
//     return tx, nil
// }
func (m *PriorityMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Find transaction in map
	// 3. Remove from heap
	// 4. Remove from map
	return nil, nil
}

// GetNext returns the highest priority transaction without removing it.
//
// TODO: Implement priority-based peek
//
// func (m *PriorityMempool) GetNext() *Transaction {
//     m.mu.RLock()
//     defer m.mu.RUnlock()
//
//     if m.heap.Len() == 0 {
//         return nil
//     }
//
//     return (*m.heap)[0]  // Root of heap (highest priority)
// }
func (m *PriorityMempool) GetNext() *Transaction {
	// TODO: Implement
	// Return root of heap (highest priority)
	return nil
}

// Size returns the current number of transactions in the mempool.
//
// TODO: Implement thread-safe size check
//
// func (m *PriorityMempool) Size() int {
//     m.mu.RLock()
//     defer m.mu.RUnlock()
//     return m.heap.Len()
// }
func (m *PriorityMempool) Size() int {
	// TODO: Implement with read lock
	return 0
}

// ============================================================================
// Exercise 3: Nonce-Based Mempool (Account Ordering)
// ============================================================================

// NonceMempool implements a nonce-based transaction pool.
// Transactions from the same account are ordered by nonce (sequence number).
//
// TODO: Add fields for:
// - Mutex for thread safety
// - Map of account address to AccountQueue
//
// Why nonce-based ordering?
// - Prevents replay attacks (can't re-broadcast old transactions)
// - Ensures transactions are processed in order per account
// - Common in account-based blockchains (Ethereum)
//
// Nonce example:
// - Alice sends 3 transactions: nonce 0, 1, 2
// - Must be processed in order: 0 → 1 → 2
// - If nonce 1 is missing, nonce 2 waits (pending)
//
// type NonceMempool struct {
//     mu       sync.RWMutex
//     accounts map[string]*AccountQueue  // Address → AccountQueue
// }

// AccountQueue stores transactions for a single account, ordered by nonce.
//
// TODO: Add fields for:
// - Account address
// - Pending nonce (next expected nonce)
// - Map of nonce to transaction
//
// type AccountQueue struct {
//     address      string
//     pendingNonce uint64                   // Next nonce to process
//     txs          map[uint64]*Transaction  // Nonce → Transaction
// }

// NewNonceMempool creates a new nonce-based mempool.
//
// TODO: Initialize NonceMempool
//
// func NewNonceMempool() *NonceMempool {
//     return &NonceMempool{
//         accounts: make(map[string]*AccountQueue),
//     }
// }
func NewNonceMempool() *NonceMempool {
	// TODO: Initialize NonceMempool
	return nil
}

// Add adds a transaction to the mempool.
// If a transaction with the same nonce exists, replaces it only if fee is higher.
//
// TODO: Implement nonce-based add with replacement
//
// Algorithm:
// 1. Lock mutex
// 2. Get or create AccountQueue for tx.From
// 3. Check if transaction with same nonce exists
// 4. If exists and new fee <= old fee → return error
// 5. Add/replace transaction in queue
//
// Why allow replacement?
// - User may want to speed up transaction by increasing fee
// - Common pattern: "replace-by-fee" (RBF)
//
// func (m *NonceMempool) Add(tx *Transaction) error {
//     m.mu.Lock()
//     defer m.mu.Unlock()
//
//     // Get or create account queue
//     queue, exists := m.accounts[tx.From]
//     if !exists {
//         queue = &AccountQueue{
//             address:      tx.From,
//             pendingNonce: 0,
//             txs:          make(map[uint64]*Transaction),
//         }
//         m.accounts[tx.From] = queue
//     }
//
//     // Check for replacement (higher fee)
//     if existing, exists := queue.txs[tx.Nonce]; exists {
//         if tx.Fee <= existing.Fee {
//             return errors.New("transaction with same nonce has higher or equal fee")
//         }
//     }
//
//     queue.txs[tx.Nonce] = tx
//     return nil
// }
func (m *NonceMempool) Add(tx *Transaction) error {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Get or create account queue for tx.From
	// 3. Check if nonce already exists
	// 4. If exists, only replace if higher fee
	// 5. Add to account queue
	return nil
}

// GetNextForAccount returns the next transaction for the given account
// (with nonce equal to pending nonce), or nil if not available.
//
// TODO: Implement nonce-based next transaction lookup
//
// func (m *NonceMempool) GetNextForAccount(address string) *Transaction {
//     m.mu.RLock()
//     defer m.mu.RUnlock()
//
//     queue, exists := m.accounts[address]
//     if !exists {
//         return nil
//     }
//
//     return queue.txs[queue.pendingNonce]
// }
func (m *NonceMempool) GetNextForAccount(address string) *Transaction {
	// TODO: Implement
	// 1. Read lock the mutex
	// 2. Get account queue
	// 3. Return transaction with pending nonce
	return nil
}

// AdvanceNonce advances the pending nonce for the given account.
// This should be called after processing a transaction.
//
// TODO: Implement nonce advancement
//
// Algorithm:
// 1. Lock mutex
// 2. Get account queue
// 3. Remove transaction with current pending nonce
// 4. Increment pending nonce
//
// func (m *NonceMempool) AdvanceNonce(address string) {
//     m.mu.Lock()
//     defer m.mu.Unlock()
//
//     if queue, exists := m.accounts[address]; exists {
//         delete(queue.txs, queue.pendingNonce)
//         queue.pendingNonce++
//     }
// }
func (m *NonceMempool) AdvanceNonce(address string) {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Get account queue
	// 3. Remove transaction with current pending nonce
	// 4. Increment pending nonce
}

// Size returns the total number of transactions across all accounts.
//
// TODO: Implement total transaction count
//
// func (m *NonceMempool) Size() int {
//     m.mu.RLock()
//     defer m.mu.RUnlock()
//
//     count := 0
//     for _, queue := range m.accounts {
//         count += len(queue.txs)
//     }
//     return count
// }
func (m *NonceMempool) Size() int {
	// TODO: Implement with read lock
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
