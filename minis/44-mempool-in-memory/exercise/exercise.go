//go:build !solution
// +build !solution

package exercise

import (
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

// FIFOMempool implements a first-in-first-out transaction pool.
// Transactions are processed in arrival order.
type FIFOMempool struct {
	// TODO: Add fields for:
	// - Mutex for thread safety
	// - Slice to store transactions in order
	// - Map for O(1) lookup by hash
	// - Capacity limit
}

// NewFIFOMempool creates a new FIFO mempool with the given capacity.
func NewFIFOMempool(capacity int) *FIFOMempool {
	// TODO: Initialize FIFOMempool
	return nil
}

// Add adds a transaction to the mempool.
// Returns error if transaction already exists or mempool is full.
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
func (m *FIFOMempool) GetNext() *Transaction {
	// TODO: Implement
	// 1. Read lock the mutex
	// 2. Return first transaction in slice
	return nil
}

// Size returns the current number of transactions in the mempool.
func (m *FIFOMempool) Size() int {
	// TODO: Implement with read lock
	return 0
}

// PriorityMempool implements a priority-based transaction pool.
// Transactions are ordered by fee (highest first).
type PriorityMempool struct {
	// TODO: Add fields for:
	// - Mutex for thread safety
	// - Heap to store transactions by priority
	// - Map for O(1) lookup by hash
	// - Capacity limit
}

// NewPriorityMempool creates a new priority mempool with the given capacity.
func NewPriorityMempool(capacity int) *PriorityMempool {
	// TODO: Initialize PriorityMempool
	// Hint: Use container/heap package
	return nil
}

// Add adds a transaction to the mempool.
// If mempool is full, evicts lowest priority transaction if new tx has higher priority.
func (m *PriorityMempool) Add(tx *Transaction) error {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Check if transaction already exists
	// 3. If at capacity, check if new tx has higher priority than lowest
	// 4. Add to heap and map
	return nil
}

// Remove removes a transaction from the mempool by hash.
func (m *PriorityMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Find transaction in map
	// 3. Remove from heap
	// 4. Remove from map
	return nil, nil
}

// GetNext returns the highest priority transaction without removing it.
func (m *PriorityMempool) GetNext() *Transaction {
	// TODO: Implement
	// Return root of heap (highest priority)
	return nil
}

// Size returns the current number of transactions in the mempool.
func (m *PriorityMempool) Size() int {
	// TODO: Implement with read lock
	return 0
}

// NonceMempool implements a nonce-based transaction pool.
// Transactions from the same account are ordered by nonce.
type NonceMempool struct {
	// TODO: Add fields for:
	// - Mutex for thread safety
	// - Map of account address to AccountQueue
}

// AccountQueue stores transactions for a single account, ordered by nonce.
type AccountQueue struct {
	// TODO: Add fields for:
	// - Account address
	// - Pending nonce (next expected nonce)
	// - Map of nonce to transaction
}

// NewNonceMempool creates a new nonce-based mempool.
func NewNonceMempool() *NonceMempool {
	// TODO: Initialize NonceMempool
	return nil
}

// Add adds a transaction to the mempool.
// If a transaction with the same nonce exists, replaces it only if fee is higher.
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
func (m *NonceMempool) GetNextForAccount(address string) *Transaction {
	// TODO: Implement
	// 1. Read lock the mutex
	// 2. Get account queue
	// 3. Return transaction with pending nonce
	return nil
}

// AdvanceNonce advances the pending nonce for the given account.
// This should be called after processing a transaction.
func (m *NonceMempool) AdvanceNonce(address string) {
	// TODO: Implement
	// 1. Lock the mutex
	// 2. Get account queue
	// 3. Remove transaction with current pending nonce
	// 4. Increment pending nonce
}

// Size returns the total number of transactions across all accounts.
func (m *NonceMempool) Size() int {
	// TODO: Implement with read lock
	return 0
}
