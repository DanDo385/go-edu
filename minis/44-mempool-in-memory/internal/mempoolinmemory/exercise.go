//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package mempoolinmemory

import (
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
// TODO: implement NewFIFOMempool.
func NewFIFOMempool(capacity int) *FIFOMempool { panic("TODO: implement") }
// TODO: implement Add.
func (m *FIFOMempool) Add(tx *Transaction) error { panic("TODO: implement") }
// TODO: implement Remove.
func (m *FIFOMempool) Remove(hash string) (*Transaction, error) { panic("TODO: implement") }
// TODO: implement GetNext.
func (m *FIFOMempool) GetNext() *Transaction { panic("TODO: implement") }
// TODO: implement Size.
func (m *FIFOMempool) Size() int { panic("TODO: implement") }

type PriorityMempool struct {
	mu       sync.RWMutex
	heap     *TxHeap
	txMap    map[string]int // Hash -> heap index
	capacity int
}

type TxHeap []*Transaction
// TODO: implement Len.
func (h TxHeap) Len() int { panic("TODO: implement") }
// TODO: implement Less.
func (h TxHeap) Less(i, j int) bool { panic("TODO: implement") }
// TODO: implement Swap.
func (h TxHeap) Swap(i, j int) { panic("TODO: implement") }
// TODO: implement Push.
func (h *TxHeap) Push(x interface{}) { panic("TODO: implement") }
// TODO: implement Pop.
func (h *TxHeap) Pop() interface{} { panic("TODO: implement") }
// TODO: implement NewPriorityMempool.
func NewPriorityMempool(capacity int) *PriorityMempool { panic("TODO: implement") }
// TODO: implement Add.
func (m *PriorityMempool) Add(tx *Transaction) error { panic("TODO: implement") }
// TODO: implement Remove.
func (m *PriorityMempool) Remove(hash string) (*Transaction, error) { panic("TODO: implement") }
// TODO: implement GetNext.
func (m *PriorityMempool) GetNext() *Transaction { panic("TODO: implement") }
// TODO: implement Size.
func (m *PriorityMempool) Size() int { panic("TODO: implement") }

type NonceMempool struct {
	mu       sync.RWMutex
	accounts map[string]*AccountQueue
}

type AccountQueue struct {
	address      string
	pendingNonce uint64
	txs          map[uint64]*Transaction
}
// TODO: implement NewNonceMempool.
func NewNonceMempool() *NonceMempool { panic("TODO: implement") }
// TODO: implement Add.
func (m *NonceMempool) Add(tx *Transaction) error { panic("TODO: implement") }
// TODO: implement GetNextForAccount.
func (m *NonceMempool) GetNextForAccount(address string) *Transaction { panic("TODO: implement") }
// TODO: implement AdvanceNonce.
func (m *NonceMempool) AdvanceNonce(address string) { panic("TODO: implement") }
// TODO: implement Size.
func (m *NonceMempool) Size() int { panic("TODO: implement") }
