//go:build !solution && !reference

package mempoolinmemory

import (
	"container/heap"
	"errors"
	"sync"
	"time"
)

func NewFIFOMempool(capacity int) *FIFOMempool {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func (m *FIFOMempool) Add(tx *Transaction) error {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *FIFOMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *FIFOMempool) GetNext() *Transaction {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *FIFOMempool) Size() int {
	// TODO: Implement this function
	panic("not implemented")
}

func (h TxHeap) Len() int {
	// TODO: Implement this function
	panic("not implemented")
}

func (h TxHeap) Less(i, j int) bool {
	// TODO: Implement this function
	panic("not implemented")
}

func (h TxHeap) Swap(i, j int) {
	// TODO: Implement this function
	panic("not implemented")
}

func (h *TxHeap) Push(x interface{}) {
	// TODO: Implement this function
	panic("not implemented")
}

func (h *TxHeap) Pop() interface{} {
	// TODO: Implement this function
	panic("not implemented")
}

func NewPriorityMempool(capacity int) *PriorityMempool {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *PriorityMempool) Add(tx *Transaction) error {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *PriorityMempool) Remove(hash string) (*Transaction, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *PriorityMempool) GetNext() *Transaction {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *PriorityMempool) Size() int {
	// TODO: Implement this function
	panic("not implemented")
}

func NewNonceMempool() *NonceMempool {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *NonceMempool) Add(tx *Transaction) error {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *NonceMempool) GetNextForAccount(address string) *Transaction {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *NonceMempool) AdvanceNonce(address string) {
	// TODO: Implement this function
	panic("not implemented")
}

func (m *NonceMempool) Size() int {
	// TODO: Implement this function
	panic("not implemented")
}
