//go:build !solution && !reference

package mempoolinmemory

import (
	"container/heap"
	"errors"
	"sync"
	"time"
)

type Transaction struct {
	Hash      string    // Unique transaction identifier
	From      string    // Sender address
	To        string    // Receiver address
	Value     uint64    // Amount to transfer
	Fee       uint64    // Transaction fee (used for prioritization)
	Nonce     uint64    // Account nonce (for ordering)
	Timestamp time.Time // When transaction was created
}

type FIFOMempool struct {
	mu       sync.RWMutex
	txs      []*Transaction
	txMap    map[string]*Transaction // Hash -> Transaction for O(1) lookup
	capacity int
}

type PriorityMempool struct {
	mu       sync.RWMutex
	heap     *TxHeap
	txMap    map[string]int // Hash -> heap index
	capacity int
}

type TxHeap []*Transaction

type NonceMempool struct {
	mu       sync.RWMutex
	accounts map[string]*AccountQueue
}

type AccountQueue struct {
	address      string
	pendingNonce uint64
	txs          map[uint64]*Transaction
}

// NewFIFOMempool implements the exercise.
//
// TODO: Implement this function
func NewFIFOMempool(capacity int) *FIFOMempool {
	// TODO: Implement
	return nil
}

// Add implements the exercise.
//
// TODO: Implement this function
func (m *FIFOMempool) Add(tx *Transaction) error {
	// TODO: Implement
	return nil
}

// Remove implements the exercise.
//
// TODO: Implement this function
func (m *FIFOMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement
	return nil, nil
}

// GetNext implements the exercise.
//
// TODO: Implement this function
func (m *FIFOMempool) GetNext() *Transaction {
	// TODO: Implement
	return nil
}

// Size implements the exercise.
//
// TODO: Implement this function
func (m *FIFOMempool) Size() int {
	// TODO: Implement
	return 0
}

// Len implements the exercise.
//
// TODO: Implement this function
func (h TxHeap) Len() int {
	// TODO: Implement
	return 0
}

// Less implements the exercise.
//
// TODO: Implement this function
func (h TxHeap) Less(i int, j int) bool {
	// TODO: Implement
	return false
}

// Swap implements the exercise.
//
// TODO: Implement this function
func (h TxHeap) Swap(i int, j int) {
	// TODO: Implement
}

// Push implements the exercise.
//
// TODO: Implement this function
func (h *TxHeap) Push(x interface{}) {
	// TODO: Implement
}

// Pop implements the exercise.
//
// TODO: Implement this function
func (h *TxHeap) Pop() interface{} {
	// TODO: Implement
	return nil
}

// NewPriorityMempool implements the exercise.
//
// TODO: Implement this function
func NewPriorityMempool(capacity int) *PriorityMempool {
	// TODO: Implement
	return nil
}

// Add implements the exercise.
//
// TODO: Implement this function
func (m *PriorityMempool) Add(tx *Transaction) error {
	// TODO: Implement
	return nil
}

// Remove implements the exercise.
//
// TODO: Implement this function
func (m *PriorityMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement
	return nil, nil
}

// GetNext implements the exercise.
//
// TODO: Implement this function
func (m *PriorityMempool) GetNext() *Transaction {
	// TODO: Implement
	return nil
}

// Size implements the exercise.
//
// TODO: Implement this function
func (m *PriorityMempool) Size() int {
	// TODO: Implement
	return 0
}

// NewNonceMempool implements the exercise.
//
// TODO: Implement this function
func NewNonceMempool() *NonceMempool {
	// TODO: Implement
	return nil
}

// Add implements the exercise.
//
// TODO: Implement this function
func (m *NonceMempool) Add(tx *Transaction) error {
	// TODO: Implement
	return nil
}

// GetNextForAccount implements the exercise.
//
// TODO: Implement this function
func (m *NonceMempool) GetNextForAccount(address string) *Transaction {
	// TODO: Implement
	return nil
}

// AdvanceNonce implements the exercise.
//
// TODO: Implement this function
func (m *NonceMempool) AdvanceNonce(address string) {
	// TODO: Implement
}

// Size implements the exercise.
//
// TODO: Implement this function
func (m *NonceMempool) Size() int {
	// TODO: Implement
	return 0
}
