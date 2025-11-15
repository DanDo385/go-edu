//go:build solution
// +build solution

package exercise

import (
	"container/heap"
	"errors"
	"sync"
	"time"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash      string    // Unique transaction identifier
	From      string    // Sender address
	To        string    // Receiver address
	Value     uint64    // Amount to transfer
	Fee       uint64    // Transaction fee (used for prioritization)
	Nonce     uint64    // Account nonce (for ordering)
	Timestamp time.Time // When transaction was created
}

// ============================================================================
// FIFO Mempool Implementation
// ============================================================================

// FIFOMempool implements a first-in-first-out transaction pool.
type FIFOMempool struct {
	mu       sync.RWMutex
	txs      []*Transaction
	txMap    map[string]*Transaction // Hash -> Transaction for O(1) lookup
	capacity int
}

// NewFIFOMempool creates a new FIFO mempool with the given capacity.
func NewFIFOMempool(capacity int) *FIFOMempool {
	return &FIFOMempool{
		txs:      make([]*Transaction, 0, capacity),
		txMap:    make(map[string]*Transaction),
		capacity: capacity,
	}
}

// Add adds a transaction to the mempool.
func (m *FIFOMempool) Add(tx *Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already exists
	if _, exists := m.txMap[tx.Hash]; exists {
		return errors.New("transaction already exists")
	}

	// Check capacity
	if len(m.txs) >= m.capacity {
		return errors.New("mempool full")
	}

	// Add to slice and map
	m.txs = append(m.txs, tx)
	m.txMap[tx.Hash] = tx

	return nil
}

// Remove removes a transaction from the mempool by hash.
func (m *FIFOMempool) Remove(hash string) (*Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tx, exists := m.txMap[hash]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	// Find and remove from slice
	for i, t := range m.txs {
		if t.Hash == hash {
			// Remove by shifting elements
			m.txs = append(m.txs[:i], m.txs[i+1:]...)
			break
		}
	}

	delete(m.txMap, hash)
	return tx, nil
}

// GetNext returns the next transaction to process (oldest).
func (m *FIFOMempool) GetNext() *Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.txs) == 0 {
		return nil
	}

	return m.txs[0]
}

// Size returns the current number of transactions in the mempool.
func (m *FIFOMempool) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.txs)
}

// ============================================================================
// Priority Mempool Implementation
// ============================================================================

// PriorityMempool implements a priority-based transaction pool using a heap.
type PriorityMempool struct {
	mu       sync.RWMutex
	heap     *TxHeap
	txMap    map[string]int // Hash -> heap index
	capacity int
}

// TxHeap implements heap.Interface for priority queue.
type TxHeap []*Transaction

func (h TxHeap) Len() int { return len(h) }

func (h TxHeap) Less(i, j int) bool {
	// Higher fee = higher priority
	if h[i].Fee != h[j].Fee {
		return h[i].Fee > h[j].Fee
	}
	// Tiebreaker: earlier timestamp
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

// NewPriorityMempool creates a new priority mempool with the given capacity.
func NewPriorityMempool(capacity int) *PriorityMempool {
	h := &TxHeap{}
	heap.Init(h)

	return &PriorityMempool{
		heap:     h,
		txMap:    make(map[string]int),
		capacity: capacity,
	}
}

// Add adds a transaction to the mempool.
func (m *PriorityMempool) Add(tx *Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already exists
	if _, exists := m.txMap[tx.Hash]; exists {
		return errors.New("transaction already exists")
	}

	// If at capacity, check if we should evict lowest priority
	if m.heap.Len() >= m.capacity {
		lowest := (*m.heap)[m.heap.Len()-1]
		if tx.Fee <= lowest.Fee {
			return errors.New("mempool full, transaction priority too low")
		}

		// Evict lowest priority
		evicted := heap.Remove(m.heap, m.heap.Len()-1).(*Transaction)
		delete(m.txMap, evicted.Hash)
	}

	// Add to heap
	heap.Push(m.heap, tx)
	m.txMap[tx.Hash] = m.heap.Len() - 1

	return nil
}

// Remove removes a transaction from the mempool by hash.
func (m *PriorityMempool) Remove(hash string) (*Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.txMap[hash]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	// Find current index (may have changed due to heap operations)
	actualIdx := -1
	for i, tx := range *m.heap {
		if tx.Hash == hash {
			actualIdx = i
			break
		}
	}

	if actualIdx == -1 {
		return nil, errors.New("transaction not found in heap")
	}

	tx := heap.Remove(m.heap, actualIdx).(*Transaction)
	delete(m.txMap, hash)

	return tx, nil
}

// GetNext returns the highest priority transaction without removing it.
func (m *PriorityMempool) GetNext() *Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.heap.Len() == 0 {
		return nil
	}

	return (*m.heap)[0]
}

// Size returns the current number of transactions in the mempool.
func (m *PriorityMempool) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.heap.Len()
}

// ============================================================================
// Nonce Mempool Implementation
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
	return &NonceMempool{
		accounts: make(map[string]*AccountQueue),
	}
}

// Add adds a transaction to the mempool.
func (m *NonceMempool) Add(tx *Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Get or create account queue
	queue, exists := m.accounts[tx.From]
	if !exists {
		queue = &AccountQueue{
			address:      tx.From,
			pendingNonce: 0,
			txs:          make(map[uint64]*Transaction),
		}
		m.accounts[tx.From] = queue
	}

	// Check if nonce already exists
	if existing, exists := queue.txs[tx.Nonce]; exists {
		// Replace only if higher fee
		if tx.Fee <= existing.Fee {
			return errors.New("transaction with same nonce has higher or equal fee")
		}
	}

	queue.txs[tx.Nonce] = tx
	return nil
}

// GetNextForAccount returns the next transaction for the given account.
func (m *NonceMempool) GetNextForAccount(address string) *Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	queue, exists := m.accounts[address]
	if !exists {
		return nil
	}

	return queue.txs[queue.pendingNonce]
}

// AdvanceNonce advances the pending nonce for the given account.
func (m *NonceMempool) AdvanceNonce(address string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if queue, exists := m.accounts[address]; exists {
		delete(queue.txs, queue.pendingNonce)
		queue.pendingNonce++
	}
}

// Size returns the total number of transactions across all accounts.
func (m *NonceMempool) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, queue := range m.accounts {
		count += len(queue.txs)
	}
	return count
}
