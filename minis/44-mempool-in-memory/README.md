# Project 44: Mempool In-Memory

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a blockchain node or cryptocurrency system. Thousands of transactions arrive every second from users around the world. You need to:

**❌ Naive approach:** Process transactions immediately in arrival order
- No prioritization: High-fee urgent transactions wait behind low-fee spam
- No validation: Invalid transactions consume processing resources
- No organization: Can't efficiently select best transactions for blocks
- No limits: Memory exhaustion from unlimited transaction storage

**✅ Better approach:** Use a **mempool** (memory pool)
- **Buffer incoming transactions** before they're included in blocks
- **Prioritize by fee**: Miners want high-fee transactions first
- **Validate before accepting**: Reject invalid transactions early
- **Bounded capacity**: Evict low-value transactions when full
- **Thread-safe access**: Multiple goroutines adding/removing transactions safely

This project teaches you how to build **production-grade in-memory transaction pools** that are:
- **Priority-aware**: Order transactions by fee, nonce, or custom criteria
- **Concurrency-safe**: Handle concurrent reads/writes without data races
- **Efficient**: Fast insertion, removal, and retrieval operations
- **Bounded**: Enforce size limits to prevent memory exhaustion
- **Flexible**: Support FIFO, priority queue, or hybrid strategies

### What You'll Learn

1. **Mempool pattern**: In-memory transaction buffering for blockchain systems
2. **Priority queues**: Heap-based data structures for efficient prioritization
3. **FIFO queues**: First-in-first-out ordering for fairness
4. **Concurrent data structures**: Thread-safe collections with mutexes
5. **Eviction policies**: LRU, lowest-fee, and size-based eviction
6. **Nonce ordering**: Per-account transaction sequencing

### The Challenge

Build a mempool system with:
- FIFO ordering (simple queue)
- Priority queue ordering (by fee)
- Per-account nonce ordering (sequential transactions)
- Concurrent access (thread-safe operations)
- Size limits with eviction policies
- Efficient transaction selection for block building

---

## 2. First Principles: Understanding Mempools

### What is a Mempool?

A **mempool** (memory pool) is an in-memory buffer that stores pending transactions before they're included in a block.

**Core concept:**
```
User Transactions → Mempool → Block Builder → Blockchain
                      ↓
            [Prioritized Buffer]
            - Validation
            - Ordering
            - Selection
```

**Why do we need mempools?**

1. **Buffering**: Transactions arrive faster than blocks are produced
2. **Prioritization**: Limited block space means we need to choose best transactions
3. **Validation**: Reject invalid transactions before expensive processing
4. **DoS protection**: Prevent spam from overwhelming the system

**Real-world examples:**
- Bitcoin Core: ~300 MB mempool, priority by fee
- Ethereum (Geth): Separate pools per account, ordered by nonce
- Custom blockchains: Application-specific ordering (time, stake, etc.)

### What is FIFO Ordering?

**FIFO** (First-In-First-Out) is the simplest ordering strategy: process transactions in arrival order.

**Analogy**: Like a queue at a coffee shop—first person in line gets served first.

**Data structure:**
```go
type FIFOMempool struct {
    mu           sync.RWMutex
    transactions []Transaction
    capacity     int
}
```

**Operations:**
- **Add**: Append to end of queue
- **Remove**: Take from front of queue
- **Peek**: Look at front without removing

**Time complexity:**
- Add: O(1)
- Remove: O(1)
- Contains: O(n)

**When to use FIFO:**
- Simple systems where all transactions are equal
- Fairness is more important than efficiency
- Low transaction volume

**Limitations:**
- No prioritization (spam blocks urgent transactions)
- Inefficient for high-value selection
- No way to prefer important transactions

### What is Priority Queue Ordering?

A **priority queue** is a data structure where elements are ordered by priority, not arrival time.

**Analogy**: Emergency room triage—patients with severe injuries go first, regardless of arrival time.

**Data structure**: Binary heap (efficient priority queue implementation)
```
        [Tx: fee=100]
       /              \
  [Tx: fee=50]    [Tx: fee=80]
   /         \
[Tx: fee=20] [Tx: fee=30]
```

**Heap property**: Parent has higher priority than children.

**Operations:**
- **Add**: Insert at end, bubble up until heap property satisfied
- **Remove**: Take root (highest priority), move last element to root, bubble down
- **Peek**: Look at root

**Time complexity:**
- Add: O(log n)
- Remove: O(log n)
- Peek: O(1)

**When to use priority queues:**
- Block space is limited (need best transactions)
- Users can specify priority (fees, tips)
- Miners/validators want to maximize revenue

**Implementation in Go:**
```go
type PriorityMempool struct {
    mu      sync.RWMutex
    heap    []*Transaction
    txIndex map[string]int  // txHash -> heap index
}

// Implement heap.Interface
func (m *PriorityMempool) Len() int
func (m *PriorityMempool) Less(i, j int) bool  // Compare priorities
func (m *PriorityMempool) Swap(i, j int)
func (m *PriorityMempool) Push(x interface{})
func (m *PriorityMempool) Pop() interface{}
```

### What is Nonce Ordering?

**Nonce** (number used once) is a counter that ensures transactions from the same account are processed in order.

**Problem without nonces:**
```
Alice sends:
  Tx1: Transfer $100 (account balance: $100)
  Tx2: Transfer $50

If Tx2 processes first:
  - Tx2 succeeds ($100 - $50 = $50)
  - Tx1 fails (only $50 remaining, can't transfer $100)
```

**Solution with nonces:**
```
Alice's account nonce: 5

Tx1: nonce=5, Transfer $100
Tx2: nonce=6, Transfer $50

Mempool ensures Tx1 processes before Tx2
```

**Data structure:**
```go
type NonceMempool struct {
    mu       sync.RWMutex
    accounts map[string]*AccountQueue  // address -> transactions
}

type AccountQueue struct {
    pendingNonce uint64
    txs          map[uint64]*Transaction  // nonce -> tx
}
```

**Rules:**
1. **Sequential processing**: Process nonce N before nonce N+1
2. **Gap handling**: If nonce N is missing, wait for it (don't process N+1)
3. **Replacement**: Nonce N can replace existing nonce N if fee is higher

**When to use nonce ordering:**
- Account-based blockchains (Ethereum, Polygon)
- Sequential operations required (prevent race conditions)
- State-dependent transactions

### What is Concurrency Safety?

**Concurrency safety** means multiple goroutines can access the mempool simultaneously without data races.

**The problem:**
```go
// ❌ UNSAFE
type UnsafeMempool struct {
    transactions []*Transaction
}

func (m *UnsafeMempool) Add(tx *Transaction) {
    m.transactions = append(m.transactions, tx)  // DATA RACE!
}

// Two goroutines calling Add() simultaneously can corrupt the slice
```

**The solution: Mutual exclusion (mutex)**
```go
// ✅ SAFE
type SafeMempool struct {
    mu           sync.RWMutex
    transactions []*Transaction
}

func (m *SafeMempool) Add(tx *Transaction) {
    m.mu.Lock()         // Acquire exclusive lock
    defer m.mu.Unlock() // Release lock when done

    m.transactions = append(m.transactions, tx)  // Safe!
}

func (m *SafeMempool) Get(hash string) *Transaction {
    m.mu.RLock()         // Acquire shared read lock
    defer m.mu.RUnlock() // Release lock when done

    for _, tx := range m.transactions {
        if tx.Hash == hash {
            return tx
        }
    }
    return nil
}
```

**Key concepts:**

1. **sync.Mutex**: Exclusive lock (only one goroutine at a time)
2. **sync.RWMutex**: Read-write lock (multiple readers OR one writer)
3. **Lock/Unlock**: Protect write operations
4. **RLock/RUnlock**: Protect read operations

**When to use RWMutex:**
- Many readers, few writers (common in mempools)
- Allows concurrent reads (better performance)
- Writers still get exclusive access

**Common patterns:**
```go
// Always use defer to ensure unlock
func (m *Mempool) Add(tx *Transaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    // ... even if panic occurs, unlock happens
}

// Read lock for reads
func (m *Mempool) Size() int {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return len(m.transactions)
}
```

---

## 3. Breaking Down the Solution

### Step 1: Define Transaction Structure

**What data do we need?**

```go
type Transaction struct {
    Hash      string    // Unique identifier
    From      string    // Sender address
    To        string    // Recipient address
    Value     uint64    // Amount to transfer
    Fee       uint64    // Transaction fee (for priority)
    Nonce     uint64    // Account nonce (for ordering)
    Timestamp time.Time // When transaction was created
}
```

**Why these fields?**
- **Hash**: Unique ID for deduplication and lookup
- **From/To**: Identify sender and recipient
- **Value**: Amount being transferred
- **Fee**: Priority metric (higher fee = higher priority)
- **Nonce**: Ordering within account
- **Timestamp**: FIFO tiebreaker, debugging

### Step 2: FIFO Mempool Implementation

**Data structure:**
```go
type FIFOMempool struct {
    mu       sync.RWMutex
    txs      []*Transaction
    txMap    map[string]*Transaction  // Hash -> Tx (for O(1) lookup)
    capacity int
}
```

**Add operation:**
```go
func (m *FIFOMempool) Add(tx *Transaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Check if already exists
    if _, exists := m.txMap[tx.Hash]; exists {
        return errors.New("transaction already in pool")
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
```

**Remove operation:**
```go
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
            // Remove by replacing with last element, then truncate
            m.txs[i] = m.txs[len(m.txs)-1]
            m.txs = m.txs[:len(m.txs)-1]
            break
        }
    }

    delete(m.txMap, hash)
    return tx, nil
}
```

**GetNext operation (FIFO):**
```go
func (m *FIFOMempool) GetNext() *Transaction {
    m.mu.RLock()
    defer m.mu.RUnlock()

    if len(m.txs) == 0 {
        return nil
    }

    return m.txs[0]  // First in = first out
}
```

### Step 3: Priority Queue Implementation

**Using Go's container/heap:**
```go
import "container/heap"

type PriorityMempool struct {
    mu      sync.RWMutex
    heap    *TxHeap
    txMap   map[string]int  // Hash -> heap index
}

type TxHeap []*Transaction

// Implement heap.Interface
func (h TxHeap) Len() int {
    return len(h)
}

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
```

**Add operation:**
```go
func (m *PriorityMempool) Add(tx *Transaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, exists := m.txMap[tx.Hash]; exists {
        return errors.New("transaction already in pool")
    }

    heap.Push(m.heap, tx)
    m.txMap[tx.Hash] = len(*m.heap) - 1

    return nil
}
```

**GetNext operation (highest priority):**
```go
func (m *PriorityMempool) GetNext() *Transaction {
    m.mu.RLock()
    defer m.mu.RUnlock()

    if m.heap.Len() == 0 {
        return nil
    }

    return (*m.heap)[0]  // Root = highest priority
}
```

**Remove operation:**
```go
func (m *PriorityMempool) Remove(hash string) (*Transaction, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    idx, exists := m.txMap[hash]
    if !exists {
        return nil, errors.New("transaction not found")
    }

    tx := heap.Remove(m.heap, idx).(*Transaction)
    delete(m.txMap, hash)

    return tx, nil
}
```

### Step 4: Nonce-Based Mempool

**Data structure:**
```go
type NonceMempool struct {
    mu       sync.RWMutex
    accounts map[string]*AccountQueue
}

type AccountQueue struct {
    address      string
    pendingNonce uint64
    txs          map[uint64]*Transaction  // nonce -> tx
}
```

**Add operation:**
```go
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
            return errors.New("transaction with same nonce has higher fee")
        }
    }

    queue.txs[tx.Nonce] = tx
    return nil
}
```

**GetNextForAccount operation:**
```go
func (m *NonceMempool) GetNextForAccount(address string) *Transaction {
    m.mu.RLock()
    defer m.mu.RUnlock()

    queue, exists := m.accounts[address]
    if !exists {
        return nil
    }

    // Return transaction with expected nonce
    return queue.txs[queue.pendingNonce]
}
```

**AdvanceNonce operation:**
```go
func (m *NonceMempool) AdvanceNonce(address string) {
    m.mu.Lock()
    defer m.mu.Unlock()

    if queue, exists := m.accounts[address]; exists {
        delete(queue.txs, queue.pendingNonce)
        queue.pendingNonce++
    }
}
```

### Step 5: Size Limits and Eviction

**Eviction strategies:**

1. **Reject new transactions** (simplest)
```go
func (m *Mempool) Add(tx *Transaction) error {
    if m.Size() >= m.capacity {
        return errors.New("mempool full")
    }
    // ... add transaction
}
```

2. **Evict lowest priority** (priority mempool)
```go
func (m *PriorityMempool) Add(tx *Transaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.heap.Len() >= m.capacity {
        // Get lowest priority transaction
        lowest := (*m.heap)[m.heap.Len()-1]

        // Only evict if new transaction has higher priority
        if tx.Fee <= lowest.Fee {
            return errors.New("transaction priority too low")
        }

        // Evict lowest priority
        heap.Remove(m.heap, m.heap.Len()-1)
        delete(m.txMap, lowest.Hash)
    }

    heap.Push(m.heap, tx)
    m.txMap[tx.Hash] = len(*m.heap) - 1
    return nil
}
```

3. **Evict oldest** (FIFO mempool)
```go
func (m *FIFOMempool) Add(tx *Transaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if len(m.txs) >= m.capacity {
        // Evict first transaction (oldest)
        oldest := m.txs[0]
        m.txs = m.txs[1:]
        delete(m.txMap, oldest.Hash)
    }

    m.txs = append(m.txs, tx)
    m.txMap[tx.Hash] = tx
    return nil
}
```

---

## 4. Complete Solution Walkthrough

### Full Implementation: Hybrid Mempool

**Combining priority queue with nonce ordering:**

```go
package exercise

import (
    "container/heap"
    "errors"
    "sync"
    "time"
)

type Transaction struct {
    Hash      string
    From      string
    To        string
    Value     uint64
    Fee       uint64
    Nonce     uint64
    Timestamp time.Time
}

// HybridMempool uses priority queue globally, but maintains nonce order per account
type HybridMempool struct {
    mu           sync.RWMutex
    accounts     map[string]*AccountQueue
    globalPriority *TxHeap
    txIndex      map[string]*Transaction
    capacity     int
}

type AccountQueue struct {
    address      string
    pendingNonce uint64
    txs          map[uint64]*Transaction
    ready        []*Transaction  // Txs with correct nonce
}

func NewHybridMempool(capacity int) *HybridMempool {
    h := &TxHeap{}
    heap.Init(h)

    return &HybridMempool{
        accounts:       make(map[string]*AccountQueue),
        globalPriority: h,
        txIndex:        make(map[string]*Transaction),
        capacity:       capacity,
    }
}

func (m *HybridMempool) Add(tx *Transaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Check if already exists
    if _, exists := m.txIndex[tx.Hash]; exists {
        return errors.New("transaction already exists")
    }

    // Check capacity
    if len(m.txIndex) >= m.capacity {
        // Evict lowest priority
        if m.globalPriority.Len() > 0 {
            lowest := (*m.globalPriority)[m.globalPriority.Len()-1]
            if tx.Fee <= lowest.Fee {
                return errors.New("mempool full, transaction priority too low")
            }
            m.evict(lowest)
        } else {
            return errors.New("mempool full")
        }
    }

    // Add to account queue
    queue := m.getOrCreateQueue(tx.From)
    if err := queue.add(tx); err != nil {
        return err
    }

    // If transaction has correct nonce, add to global priority queue
    if tx.Nonce == queue.pendingNonce {
        heap.Push(m.globalPriority, tx)
        queue.ready = append(queue.ready, tx)
    }

    m.txIndex[tx.Hash] = tx
    return nil
}

func (m *HybridMempool) Remove(hash string) (*Transaction, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    tx, exists := m.txIndex[hash]
    if !exists {
        return nil, errors.New("transaction not found")
    }

    delete(m.txIndex, hash)

    // Remove from account queue
    if queue, exists := m.accounts[tx.From]; exists {
        delete(queue.txs, tx.Nonce)
    }

    // Remove from global priority queue
    m.removeFromHeap(tx)

    return tx, nil
}

func (m *HybridMempool) GetNext() *Transaction {
    m.mu.RLock()
    defer m.mu.RUnlock()

    if m.globalPriority.Len() == 0 {
        return nil
    }

    return (*m.globalPriority)[0]
}

func (m *HybridMempool) Size() int {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return len(m.txIndex)
}

func (m *HybridMempool) getOrCreateQueue(address string) *AccountQueue {
    queue, exists := m.accounts[address]
    if !exists {
        queue = &AccountQueue{
            address:      address,
            pendingNonce: 0,
            txs:          make(map[uint64]*Transaction),
            ready:        make([]*Transaction, 0),
        }
        m.accounts[address] = queue
    }
    return queue
}

func (m *HybridMempool) evict(tx *Transaction) {
    delete(m.txIndex, tx.Hash)

    if queue, exists := m.accounts[tx.From]; exists {
        delete(queue.txs, tx.Nonce)
    }

    m.removeFromHeap(tx)
}

func (m *HybridMempool) removeFromHeap(tx *Transaction) {
    for i, htx := range *m.globalPriority {
        if htx.Hash == tx.Hash {
            heap.Remove(m.globalPriority, i)
            return
        }
    }
}

func (q *AccountQueue) add(tx *Transaction) error {
    if existing, exists := q.txs[tx.Nonce]; exists {
        if tx.Fee <= existing.Fee {
            return errors.New("transaction with same nonce has higher fee")
        }
    }

    q.txs[tx.Nonce] = tx
    return nil
}

// TxHeap implements heap.Interface
type TxHeap []*Transaction

func (h TxHeap) Len() int { return len(h) }

func (h TxHeap) Less(i, j int) bool {
    if h[i].Fee != h[j].Fee {
        return h[i].Fee > h[j].Fee
    }
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
```

---

## 5. Key Concepts Explained

### Concept 1: Heap vs Sorted List

**Why use a heap instead of a sorted list?**

| Operation | Sorted List | Binary Heap |
|-----------|-------------|-------------|
| Insert | O(n) | O(log n) |
| Remove max | O(1) | O(log n) |
| Get max | O(1) | O(1) |
| Space | O(n) | O(n) |

**When to use each:**
- **Sorted list**: Need to iterate in order, few insertions
- **Heap**: Frequent insertions/removals, only need max element

**Mempool use case:** Heap wins because we frequently add/remove transactions.

### Concept 2: Read-Write Mutex Benefits

**Performance comparison:**

```go
// With sync.Mutex (exclusive lock)
// 10 goroutines reading Size() sequentially: ~1000 μs

// With sync.RWMutex (shared read lock)
// 10 goroutines reading Size() concurrently: ~100 μs
```

**Rule of thumb:** Use RWMutex when reads outnumber writes 10:1 or more.

### Concept 3: Memory Efficiency

**Transaction size:**
```
Hash: 64 bytes (SHA-256)
From/To: 40 bytes each
Value/Fee/Nonce: 8 bytes each
Timestamp: 24 bytes
Total: ~200 bytes per transaction
```

**Mempool capacity:**
```
1,000 transactions = ~200 KB
10,000 transactions = ~2 MB
100,000 transactions = ~20 MB
```

**Production sizing:** Bitcoin Core uses ~300 MB mempool (~1.5M transactions).

### Concept 4: Nonce Gap Handling

**Problem:** What if nonce 5 arrives before nonce 4?

**Solutions:**

1. **Wait for gap to fill** (Ethereum approach)
```go
// Only process nonce 5 after nonce 4 arrives
queue.ready = transactions with contiguous nonces starting from pendingNonce
```

2. **Reject future nonces** (simple approach)
```go
if tx.Nonce > queue.pendingNonce + MaxNonceGap {
    return errors.New("nonce too far in future")
}
```

3. **Time-based expiry** (hybrid approach)
```go
if tx.Timestamp.Add(MaxWaitTime).Before(time.Now()) {
    evict(tx)  // Waited too long for gap to fill
}
```

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Transaction Validator

Validate transactions before adding to mempool.

```go
type Validator interface {
    Validate(tx *Transaction) error
}

type BasicValidator struct{}

func (v *BasicValidator) Validate(tx *Transaction) error {
    if tx.Hash == "" {
        return errors.New("missing hash")
    }
    if tx.From == "" || tx.To == "" {
        return errors.New("missing address")
    }
    if tx.Fee == 0 {
        return errors.New("fee must be > 0")
    }
    return nil
}

// Use in mempool
func (m *Mempool) Add(tx *Transaction) error {
    if err := m.validator.Validate(tx); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // ... add to mempool
}
```

### Pattern 2: Event Notifications

Notify listeners when transactions are added/removed.

```go
type MempoolListener interface {
    OnAdded(tx *Transaction)
    OnRemoved(tx *Transaction)
    OnEvicted(tx *Transaction)
}

type Mempool struct {
    // ...
    listeners []MempoolListener
}

func (m *Mempool) Add(tx *Transaction) error {
    // ... add transaction

    for _, listener := range m.listeners {
        listener.OnAdded(tx)
    }

    return nil
}
```

### Pattern 3: Batch Operations

Add/remove multiple transactions atomically.

```go
func (m *Mempool) AddBatch(txs []*Transaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Validate all first
    for _, tx := range txs {
        if err := m.validator.Validate(tx); err != nil {
            return err  // Reject entire batch
        }
    }

    // Add all
    for _, tx := range txs {
        m.unsafeAdd(tx)  // No locking (already locked)
    }

    return nil
}
```

### Pattern 4: Mempool Statistics

Track metrics for monitoring.

```go
type Stats struct {
    TotalTransactions int
    TotalFees         uint64
    OldestTransaction time.Time
    AverageFee        uint64
    HighestFee        uint64
    LowestFee         uint64
}

func (m *Mempool) GetStats() Stats {
    m.mu.RLock()
    defer m.mu.RUnlock()

    stats := Stats{}
    // ... calculate statistics

    return stats
}
```

### Pattern 5: Expiry/TTL

Remove old transactions automatically.

```go
type Mempool struct {
    // ...
    ttl time.Duration
}

func (m *Mempool) startExpiryMonitor() {
    ticker := time.NewTicker(1 * time.Minute)
    go func() {
        for range ticker.C {
            m.evictExpired()
        }
    }()
}

func (m *Mempool) evictExpired() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    for hash, tx := range m.txIndex {
        if now.Sub(tx.Timestamp) > m.ttl {
            m.unsafeRemove(hash)
        }
    }
}
```

---

## 7. Real-World Applications

### Bitcoin Core Mempool

**Features:**
- Priority by fee-per-byte
- Descendant/ancestor limits (prevent chain abuse)
- Replace-by-fee (RBF) support
- Size limit: ~300 MB

**Use case:** Miners select highest-fee transactions for blocks

### Ethereum (Geth) TxPool

**Features:**
- Per-account queues with nonce ordering
- Separate "pending" and "queued" pools
- Gas price threshold
- Account nonce gap handling

**Use case:** Validators build blocks with sequential account transactions

### Custom Blockchains

**Features:**
- Application-specific priority (stake, reputation, etc.)
- Multi-dimensional ordering (fee + time + stake)
- Custom validation rules

**Use case:** Specialized consensus mechanisms (PoS, PoA, etc.)

---

## 8. Common Mistakes to Avoid

### Mistake 1: Forgetting to Lock

```go
// ❌ WRONG
func (m *Mempool) Size() int {
    return len(m.txIndex)  // DATA RACE!
}

// ✅ CORRECT
func (m *Mempool) Size() int {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return len(m.txIndex)
}
```

### Mistake 2: Holding Lock Too Long

```go
// ❌ WRONG
func (m *Mempool) Add(tx *Transaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Expensive validation while holding lock!
    if err := expensiveValidation(tx); err != nil {
        return err
    }

    m.txIndex[tx.Hash] = tx
    return nil
}

// ✅ CORRECT
func (m *Mempool) Add(tx *Transaction) error {
    // Validate BEFORE acquiring lock
    if err := expensiveValidation(tx); err != nil {
        return err
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    m.txIndex[tx.Hash] = tx
    return nil
}
```

### Mistake 3: Not Updating Heap Index

```go
// ❌ WRONG
func (h TxHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
    // Forgot to update txMap indices!
}

// ✅ CORRECT
func (m *Mempool) Swap(i, j int) {
    h := m.heap
    h[i], h[j] = h[j], h[i]
    m.txMap[h[i].Hash] = i
    m.txMap[h[j].Hash] = j
}
```

### Mistake 4: Unbounded Growth

```go
// ❌ WRONG
func (m *Mempool) Add(tx *Transaction) error {
    m.txIndex[tx.Hash] = tx  // No size limit!
    return nil
}

// ✅ CORRECT
func (m *Mempool) Add(tx *Transaction) error {
    if len(m.txIndex) >= m.capacity {
        return m.evict(tx)
    }
    m.txIndex[tx.Hash] = tx
    return nil
}
```

---

## 9. Stretch Goals

### Goal 1: Implement LRU Eviction ⭐⭐

Evict least recently used transactions when mempool is full.

**Hint:** Combine map with doubly-linked list (like `container/list`).

### Goal 2: Add Transaction Dependencies ⭐⭐⭐

Support transactions that depend on other transactions (parent/child relationships).

**Hint:** Track dependency graph, only select child after parent.

### Goal 3: Implement Replace-by-Fee ⭐⭐

Allow replacing transaction with same nonce if fee is higher.

**Hint:** Compare fees, remove old transaction, add new one.

### Goal 4: Add Metrics and Monitoring ⭐⭐

Expose Prometheus metrics for mempool monitoring.

**Hint:** Track add/remove rates, size, fees, latencies.

### Goal 5: Implement Sharded Mempool ⭐⭐⭐

Partition mempool by transaction hash for better concurrency.

**Hint:** Multiple sub-mempools, each with own lock, shard by hash prefix.

---

## How to Run

```bash
# Run the demo
go run ./minis/44-mempool-in-memory/cmd/mempool-demo

# Run tests
go test ./minis/44-mempool-in-memory/exercise

# Run with race detector
go test -race ./minis/44-mempool-in-memory/exercise

# Benchmark
go test -bench=. ./minis/44-mempool-in-memory/exercise
```

---

## Summary

**What you learned:**
- Transaction pools buffer pending transactions
- FIFO queues provide fairness, priority queues optimize for value
- Nonce ordering ensures sequential account transactions
- RWMutex enables concurrent reads with exclusive writes
- Eviction policies prevent unbounded memory growth
- Heap data structure efficiently maintains priority ordering

**Why this matters:**
Mempools are critical infrastructure in blockchain systems, handling thousands of transactions per second with complex ordering requirements. The patterns you learned (priority queues, concurrent access, eviction policies) apply to many distributed systems beyond blockchains.

**Key takeaways:**
- Choose data structure based on access patterns
- Always protect shared state with mutexes
- Design for bounded resource usage
- Test concurrent code with race detector

Build robust transaction pools!
