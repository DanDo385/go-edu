# Project 29: Escape Analysis, Inlining & Compiler Optimizations

## What Is This Project About?

This project takes you **deep into Go's compiler optimizations**—the invisible machinery that makes your code faster without you changing a single line. You'll learn:

1. **Escape Analysis** (when variables live on stack vs heap)
2. **Stack vs Heap Allocation** (performance implications)
3. **Function Inlining** (eliminating function call overhead)
4. **Compiler Flags (gcflags)** (how to inspect compiler decisions)
5. **Optimization Techniques** (writing allocation-friendly code)

By the end, you'll understand **why some code is 10x faster** than seemingly equivalent alternatives, and how to write code that the compiler can optimize aggressively.

---

## The Fundamental Problem: Memory Management is Slow

### First Principles: Where Do Variables Live?

When you declare a variable in Go, it must be stored **somewhere in memory**. There are two main options:

1. **The Stack** (fast, automatic cleanup)
2. **The Heap** (slower, requires garbage collection)

**Example:**
```go
func main() {
    x := 42  // Where does x live? Stack or heap?
}
```

The answer: **It depends**. Go's compiler uses **escape analysis** to decide.

---

## What Is Escape Analysis?

**Escape analysis** is a compile-time optimization where the Go compiler determines whether a variable can safely live on the **stack** (the function's local memory) or must "escape" to the **heap** (long-lived memory shared across functions).

### The Core Question

> "Does this variable outlive the function that created it?"

- **NO** → Stack allocation (fast!)
- **YES** → Heap allocation (escapes)

### Why This Matters

**Stack allocation:**
- **Fast**: Allocating is just moving a stack pointer (nanoseconds)
- **No GC overhead**: When the function returns, the stack is automatically reclaimed
- **Better cache locality**: Stack memory is densely packed and hot in CPU cache

**Heap allocation:**
- **Slower**: Requires memory allocator to find free space
- **GC pressure**: Objects must be tracked and eventually garbage collected
- **Cache misses**: Heap objects are scattered in memory

**Performance difference:** Stack allocations can be **100-1000x faster** than heap allocations.

---

## Escape Analysis Examples

### Example 1: Stays on Stack

```go
func sum() int {
    nums := []int{1, 2, 3, 4, 5}
    total := 0
    for _, n := range nums {
        total += n
    }
    return total  // Only total (an int) escapes, not nums
}
```

**Analysis:**
- `nums` is created and used entirely within `sum()`
- After `sum()` returns, no one needs `nums` anymore
- **Result**: `nums` lives on the stack ✓

### Example 2: Escapes to Heap (Returned Pointer)

```go
func newCounter() *int {
    count := 0
    return &count  // Returning a pointer to a local variable
}
```

**Analysis:**
- The caller receives a pointer to `count`
- `count` must remain valid after `newCounter()` returns
- **Result**: `count` escapes to the heap ✗

**Memory diagram:**
```
Stack (after newCounter returns):
┌─────────────────┐
│ (stack unwound) │  ← count would be here, but it escaped
└─────────────────┘

Heap:
┌───────┐
│ count │ ← Allocated here, pointer returned to caller
│   0   │
└───────┘
```

### Example 3: Escapes to Heap (Too Large)

```go
func largeArray() {
    // Array larger than ~64KB escapes to heap
    var big [1000000]int  // 8MB on 64-bit systems
    big[0] = 1
}
```

**Analysis:**
- Stacks have limited size (typically 2-8KB initially, can grow to ~1GB)
- Very large allocations are placed on heap to avoid stack overflow
- **Result**: `big` escapes to the heap ✗

### Example 4: Escapes via Interface

```go
func printValue(x int) {
    fmt.Println(x)  // x escapes!
}
```

**Analysis:**
- `fmt.Println()` takes `interface{}` parameters
- Interfaces require a **type descriptor + pointer to value**
- To create the pointer, `x` must be placed somewhere addressable (heap)
- **Result**: `x` escapes to the heap ✗

**This is why high-performance code avoids `fmt` in hot paths.**

### Example 5: Escapes via Slice/Map Storage

```go
func storeInSlice() {
    cache := make([]interface{}, 10)
    value := 42
    cache[0] = value  // value escapes (copied to interface{})
}
```

**Analysis:**
- Storing a value in `interface{}` causes it to escape
- The interface needs a stable pointer to the value
- **Result**: `value` escapes to the heap ✗

---

## How to Check: Using -gcflags='-m'

Go provides a compiler flag to **see escape analysis decisions**:

```bash
go build -gcflags='-m' main.go
```

**Example code:**
```go
package main

func makeSlice() []int {
    s := []int{1, 2, 3}
    return s
}

func main() {
    _ = makeSlice()
}
```

**Output:**
```
./main.go:4:6: []int{...} escapes to heap
```

### More Verbose Analysis

Use `-m -m` (or more) for increasing detail:

```bash
go build -gcflags='-m -m' main.go
```

**Output:**
```
./main.go:4:6: []int{...} escapes to heap:
./main.go:4:6:   flow: ~r0 = &{storage for []int{...}}:
./main.go:4:6:     from []int{...} (spill) at ./main.go:4:6
./main.go:4:6:     from return []int{...} (return) at ./main.go:4:2
```

This shows the **data flow** that caused the escape.

---

## Stack vs Heap: Deep Dive

### The Stack: How It Works

The stack is a **contiguous block of memory** that grows and shrinks automatically as functions are called and return.

**Example:**
```go
func main() {           // Stack frame for main()
    x := 10
    y := add(x, 20)     // New stack frame for add()
    fmt.Println(y)
}

func add(a, b int) int {
    sum := a + b        // sum lives in add's stack frame
    return sum          // add's frame destroyed, sum value copied to caller
}
```

**Stack layout during `add()`:**
```
┌───────────────────┐ ← Top of stack (grows down)
│ add's frame:      │
│   sum = 30        │
│   b = 20          │
│   a = 10          │
│   return address  │
├───────────────────┤
│ main's frame:     │
│   y = (pending)   │
│   x = 10          │
└───────────────────┘ ← Bottom of stack
```

**After `add()` returns:**
```
┌───────────────────┐
│ main's frame:     │
│   y = 30          │ ← sum's value copied here
│   x = 10          │
└───────────────────┘
```

The space used by `add()` is instantly reclaimed (just move the stack pointer up).

### The Heap: How It Works

The heap is a **large pool of memory** managed by the allocator and garbage collector.

**Allocation flow:**
1. Request N bytes from the allocator
2. Allocator finds a free chunk (or requests more memory from OS)
3. Returns a pointer to the allocated memory
4. GC periodically scans for unreachable objects and frees them

**Why it's slower:**
- **Allocation**: Must search for free space (stack just increments a pointer)
- **GC pauses**: All goroutines pause while GC scans memory
- **Fragmentation**: Heap memory becomes fragmented over time

---

## Function Inlining: Eliminating Call Overhead

**Inlining** is an optimization where the compiler **replaces a function call with the function's body**.

### Without Inlining

```go
func add(a, b int) int {
    return a + b
}

func main() {
    x := add(3, 5)  // Function call overhead
}
```

**What happens at runtime:**
1. Push arguments onto stack
2. Jump to `add()`'s code
3. Execute `a + b`
4. Return result
5. Jump back to `main()`

**Cost**: ~10-20 nanoseconds (small but adds up in hot loops)

### With Inlining

The compiler rewrites `main()` as:

```go
func main() {
    x := 3 + 5  // Direct computation, no call overhead
}
```

**Benefits:**
- **No call overhead** (save ~10-20 ns per call)
- **Better optimization** (compiler can see more context)
- **Escape elimination** (values don't need to cross function boundaries)

### Inlining Rules

The Go compiler inlines functions that are:
1. **Small** (budget: ~80 "cost points")
2. **Simple** (no complex control flow)
3. **Not recursive**

**Example: Too large to inline**
```go
func complexFunc(a, b, c int) int {
    if a > b {
        for i := 0; i < c; i++ {
            a += b * c
        }
    } else {
        a -= b
    }
    return a
}
// Too complex: will NOT be inlined
```

**Checking inlining decisions:**
```bash
go build -gcflags='-m' main.go
```

**Output:**
```
./main.go:3:6: can inline add
./main.go:7:6: can inline main
./main.go:8:10: inlining call to add
```

---

## Compiler Flags (gcflags): Your X-Ray Vision

The `-gcflags` flag passes options to the Go compiler.

### Most Useful Flags

| Flag | Description | Example |
|------|-------------|---------|
| `-m` | Print escape analysis and inlining decisions | `go build -gcflags='-m'` |
| `-m -m` | More verbose escape analysis | `go build -gcflags='-m -m'` |
| `-l` | Disable inlining | `go build -gcflags='-l'` |
| `-N` | Disable optimizations | `go build -gcflags='-N'` |
| `-S` | Print assembly output | `go build -gcflags='-S'` |

### Example: Seeing Generated Assembly

```bash
go build -gcflags='-S' main.go 2>&1 | grep -A 10 'main.add'
```

**Output (excerpt):**
```assembly
"".add STEXT nosplit size=4 args=0x18 locals=0x0
    0x0000 00000 (main.go:3)    MOVQ    "".a+8(SP), AX
    0x0005 00005 (main.go:3)    ADDQ    "".b+16(SP), AX
    0x000a 00010 (main.go:3)    MOVQ    AX, "".~r0+24(SP)
    0x000f 00015 (main.go:3)    RET
```

This shows the actual machine instructions generated.

---

## Optimization Techniques

### Technique 1: Avoid Pointers When Not Needed

**Bad (causes escape):**
```go
func process() *Result {
    r := Result{Value: 42}
    return &r  // r escapes to heap
}
```

**Good (stays on stack):**
```go
func process() Result {
    r := Result{Value: 42}
    return r  // r copied to caller's stack
}
```

**When to use pointers:**
- Large structs (copying is expensive)
- Need to modify the original (not a copy)
- Interface implementations require pointer receivers

### Technique 2: Pre-Allocate Slices

**Bad (multiple allocations):**
```go
var results []int
for i := 0; i < 1000; i++ {
    results = append(results, i)  // Grows and reallocates
}
```

**Good (single allocation):**
```go
results := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    results = append(results, i)  // No reallocation
}
```

### Technique 3: Use Value Receivers for Small Structs

**Struct definition:**
```go
type Point struct {
    X, Y int
}
```

**Prefer value receiver (16 bytes):**
```go
func (p Point) Distance() float64 {
    return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}
```

**Avoid pointer receiver for small types:**
```go
func (p *Point) Distance() float64 {  // Less efficient for small structs
    return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}
```

**Guideline**: Use value receivers for structs ≤ 16-32 bytes.

### Technique 4: Avoid Interface{} in Hot Paths

**Bad (every value escapes):**
```go
func logValue(v interface{}) {
    fmt.Println(v)  // v must escape to heap
}

for i := 0; i < 1000; i++ {
    logValue(i)  // 1000 heap allocations!
}
```

**Good (use concrete types):**
```go
func logInt(v int) {
    // Custom formatting, no interface{}
}

for i := 0; i < 1000; i++ {
    logInt(i)  // No allocations
}
```

### Technique 5: Reuse Buffers

**Bad (allocate every time):**
```go
for _, item := range items {
    buf := make([]byte, 1024)
    process(item, buf)
}
```

**Good (reuse buffer):**
```go
buf := make([]byte, 1024)
for _, item := range items {
    clear(buf)  // Reset buffer (Go 1.21+)
    process(item, buf)
}
```

---

## Practical Example: Optimizing a Parser

### Naive Implementation

```go
type Token struct {
    Type  string
    Value string
}

func parseTokens(input string) []*Token {
    tokens := []*Token{}  // Slice of pointers
    for _, word := range strings.Fields(input) {
        token := &Token{  // Each token escapes to heap
            Type:  "WORD",
            Value: word,
        }
        tokens = append(tokens, token)
    }
    return tokens
}
```

**Problems:**
1. No pre-allocation (slice grows repeatedly)
2. Pointers to tokens (every token escapes to heap)
3. No capacity hint

### Optimized Implementation

```go
type Token struct {
    Type  string
    Value string
}

func parseTokens(input string) []Token {  // Slice of values, not pointers
    words := strings.Fields(input)
    tokens := make([]Token, 0, len(words))  // Pre-allocate

    for _, word := range words {
        tokens = append(tokens, Token{  // Value stays on stack
            Type:  "WORD",
            Value: word,
        })
    }
    return tokens  // Backing array escapes, but tokens themselves don't
}
```

**Improvements:**
1. Pre-allocated slice (no repeated growth)
2. Value semantics (tokens don't individually escape)
3. Single heap allocation for the backing array

**Performance gain:** ~5-10x faster, ~90% fewer allocations.

---

## Advanced: Bounds Check Elimination

The compiler can eliminate array/slice bounds checks if it can prove they're unnecessary.

### Example: Manual Bounds Check

```go
func sumSlice(s []int) int {
    if len(s) == 0 {
        return 0
    }

    total := s[0]  // Compiler knows: len > 0, so s[0] is safe
    for i := 1; i < len(s); i++ {
        total += s[i]  // Bounds check needed (i could be large)
    }
    return total
}
```

**Optimized with explicit length:**
```go
func sumSlice(s []int) int {
    n := len(s)
    if n == 0 {
        return 0
    }

    total := s[0]
    for i := 1; i < n; i++ {  // Compiler can prove: i < n == len(s)
        total += s[i]  // Bounds check eliminated!
    }
    return total
}
```

**Check with:**
```bash
go build -gcflags='-d=ssa/check_bce/debug=1' main.go
```

---

## Benchmarking Escape Analysis Impact

### Test Code

```go
// Escapes to heap
func allocateHeap() *int {
    x := 42
    return &x
}

// Stays on stack
func allocateStack() int {
    x := 42
    return x
}

func BenchmarkHeap(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = allocateHeap()
    }
}

func BenchmarkStack(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = allocateStack()
    }
}
```

### Running Benchmarks

```bash
go test -bench=. -benchmem
```

**Expected output:**
```
BenchmarkHeap-8     50000000    25.3 ns/op    8 B/op    1 allocs/op
BenchmarkStack-8    2000000000   0.5 ns/op    0 B/op    0 allocs/op
```

**Analysis:** Stack allocation is **50x faster** and allocates **0 bytes**.

---

## How to Run This Project

```bash
# Run the demonstration program
cd minis/29-escape-analysis-inlining
go run cmd/escape-demo/main.go

# Check escape analysis decisions
go build -gcflags='-m' cmd/escape-demo/main.go

# Check inlining decisions
go build -gcflags='-m -m' cmd/escape-demo/main.go

# Disable inlining to see the difference
go build -gcflags='-l -m' cmd/escape-demo/main.go

# See generated assembly
go build -gcflags='-S' cmd/escape-demo/main.go 2>&1 | less

# Run exercises
cd exercise
go test -v

# Benchmark exercises
go test -bench=. -benchmem
```

---

## Common Escape Scenarios (Quick Reference)

| Scenario | Escapes? | Why |
|----------|----------|-----|
| Return local variable by value | NO | Value copied to caller's stack |
| Return pointer to local variable | YES | Caller needs it after function returns |
| Store in global variable | YES | Globals live on heap |
| Store in interface{} | YES | Interface needs pointer to value |
| Pass to fmt.Println() | YES | fmt uses interface{} |
| Very large local variable | YES | Stack overflow risk |
| Closure captures variable by reference | YES | Closure outlives function |
| Slice returned from function | YES | Backing array must persist |
| Slice used only locally | NO | Compiler can stack-allocate |

---

## Key Takeaways

1. **Escape analysis** determines stack vs heap allocation at **compile time**
2. **Stack allocations** are ~100x faster than heap allocations
3. **Returning pointers** causes variables to escape to heap
4. **Interface{}** causes values to escape (avoid in hot paths)
5. **Function inlining** eliminates call overhead and enables more optimizations
6. **Use `-gcflags='-m'`** to see compiler decisions
7. **Pre-allocate slices** when size is known
8. **Prefer value receivers** for small structs
9. **Reuse buffers** instead of allocating new ones
10. **Profile before optimizing** (use Project 28's pprof techniques)

---

## Connections to Other Projects

- **Project 11 (slices-internals)**: You learned slice structure; now you know when they escape
- **Project 14 (methods-value-vs-pointer-receivers)**: Receiver choice affects escape behavior
- **Project 28 (pprof-cpu-mem-benchmarks)**: Use pprof to find allocations caused by escaping
- **Project 30+**: Apply these optimization techniques to real-world scenarios

---

## Exercises

See `exercise/` directory for hands-on optimization challenges:

1. **Fix Escapes**: Rewrite functions to avoid unnecessary heap allocations
2. **Inline Candidates**: Refactor code to enable inlining
3. **Buffer Reuse**: Optimize a parser with buffer pooling
4. **Benchmark Comparison**: Measure before/after optimization impact

---

## Further Reading

- [Go FAQ: Stack vs Heap](https://go.dev/doc/faq#stack_or_heap)
- [Go Compiler Optimizations](https://github.com/golang/go/wiki/CompilerOptimizations)
- [Escape Analysis in Go](https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-escape-analysis.html)
- [Go Performance Tips](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
