# Project 14: Methods - Value vs Pointer Receivers

## What Is This Project About?

This project provides a **deep, comprehensive understanding** of method receivers in Go—one of the most critical decisions you'll make when designing types. You'll learn:

1. **What method receivers are** (value vs pointer, syntactic sugar, method sets)
2. **When to use each type** (mutability, performance, semantics)
3. **How receivers affect interface satisfaction** (the method set rules)
4. **Performance implications** (copying overhead, escape analysis, memory)
5. **Common pitfalls** (mixing receivers, nil receivers, interface confusion)
6. **Best practices** (consistency, choosing the right receiver)

By the end, you'll know exactly when to write `func (t T)` vs `func (t *T)` and why it matters for correctness, performance, and API design.

---

## The Fundamental Problem: Modifying State vs Reading State

### First Principles: What Is a Method?

A **method** is a function with a special **receiver** parameter that comes before the function name:

```go
// Function (no receiver)
func Increment(c Counter) {
    c.count++
}

// Method (with receiver)
func (c Counter) Increment() {
    c.count++
}
```

The receiver (`c Counter`) is like a special first parameter that determines:
1. **What type** the method belongs to
2. **Whether modifications** affect the original value
3. **Which interfaces** the type can satisfy

---

## Value Receivers: The Copy Semantics

### What Is a Value Receiver?

When you use a **value receiver**, the method receives a **copy** of the value:

```go
type Counter struct {
    count int
}

// VALUE RECEIVER: c is a COPY of the original Counter
func (c Counter) Increment() {
    c.count++  // Modifies the COPY, not the original!
}

func main() {
    c := Counter{count: 0}
    c.Increment()
    fmt.Println(c.count)  // 0 (unchanged!)
}
```

**Memory diagram:**
```
Stack (main):
  c: Counter{count: 0}

Stack (Increment):
  c: Counter{count: 0}  ← COPY
     └─ Modified to {count: 1}

After Increment returns:
  main's c: still {count: 0}
```

**Why this happens:**
- Go passes everything by value (makes a copy)
- The method modifies its local copy
- The copy is discarded when the method returns
- The original remains unchanged

### When Value Receivers Work Well

Value receivers are perfect for:

1. **Read-only operations** (no mutation needed):
```go
type Point struct {
    X, Y float64
}

// Read-only: just calculate, don't modify
func (p Point) Distance() float64 {
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}
```

2. **Small types** (copying is cheap):
```go
type UUID [16]byte

// Small (16 bytes), copying is fast
func (u UUID) String() string {
    return fmt.Sprintf("%x-%x-%x-%x-%x",
        u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
```

3. **Value semantics** (immutable data):
```go
type Temperature float64

// Temperature values are conceptually immutable
func (t Temperature) Celsius() float64 {
    return float64(t)
}

func (t Temperature) Fahrenheit() float64 {
    return float64(t)*9/5 + 32
}
```

---

## Pointer Receivers: The Reference Semantics

### What Is a Pointer Receiver?

When you use a **pointer receiver**, the method receives a **pointer** to the original value:

```go
// POINTER RECEIVER: c is a POINTER to the original Counter
func (c *Counter) Increment() {
    c.count++  // Modifies the ORIGINAL through the pointer!
}

func main() {
    c := Counter{count: 0}
    c.Increment()  // Go automatically does (&c).Increment()
    fmt.Println(c.count)  // 1 (changed!)
}
```

**Memory diagram:**
```
Stack (main):
  c: Counter{count: 0}
     ↑
     │ (pointer passed)
Stack (Increment):
  c: *Counter = &main.c
     └─ *c.count++ modifies main's c directly
```

**Why this works:**
- The method receives a pointer (just 8 bytes on 64-bit systems)
- Dereferencing the pointer accesses the original value
- Modifications persist after the method returns

### Go's Syntactic Sugar

Go automatically handles addressing/dereferencing for you:

```go
c := Counter{count: 0}

// All these are equivalent:
c.Increment()         // Go does (&c).Increment() automatically
(&c).Increment()      // Explicit (same result)

p := &c
p.Increment()         // Works with pointer variable too
(*p).Increment()      // Explicit dereference (same result)
```

**The rule:**
- If the receiver is `*T` and you call with `T`, Go takes the address automatically (`&`)
- If the receiver is `T` and you call with `*T`, Go dereferences automatically (`*`)

---

## The Method Set Rules (Critical for Interfaces)

This is where things get subtle and **critically important** for interfaces.

### The Core Rule

**For a type `T`:**
- Values of type `T` can call methods with receiver `T` or `*T`
- Values of type `*T` can call methods with receiver `T` or `*T`

**But for interfaces:**
- Type `T` implements interfaces requiring methods with receiver `T` only
- Type `*T` implements interfaces requiring methods with receiver `T` **or** `*T`

### Example: Interface Satisfaction

```go
type Printer interface {
    Print()
}

type Document struct {
    content string
}

// Pointer receiver
func (d *Document) Print() {
    fmt.Println(d.content)
}

func main() {
    var p Printer

    // This FAILS:
    d := Document{content: "hello"}
    p = d  // ❌ COMPILE ERROR: Document does not implement Printer
           // (Print has pointer receiver)

    // This WORKS:
    p = &d  // ✅ *Document implements Printer
    p.Print()  // "hello"
}
```

**Why this matters:**
```go
// You can CALL the method on a value (Go takes address):
d := Document{content: "hello"}
d.Print()  // ✅ Works (Go does (&d).Print())

// But you CANNOT assign to interface:
var p Printer = d  // ❌ Compile error!
```

**The reason:**
- Go can auto-address `d` to call the method because `d` is addressable
- But Go cannot auto-address when assigning to interface (interface might store unaddressable values)

### Method Set Table

| Receiver Type | Can be called on `T`? | Can be called on `*T`? | `T` satisfies interface? | `*T` satisfies interface? |
|---------------|----------------------|------------------------|-------------------------|--------------------------|
| `(t T)`       | ✅ Yes               | ✅ Yes (auto deref)    | ✅ Yes                  | ✅ Yes                   |
| `(t *T)`      | ✅ Yes (auto addr)   | ✅ Yes                 | ❌ No                   | ✅ Yes                   |

---

## When to Use Value vs Pointer Receivers

### Use Value Receivers When:

1. **The type is small** (a few words or less):
```go
type Point struct { X, Y float64 }  // 16 bytes, cheap to copy
func (p Point) Distance() float64 { ... }
```

2. **The type is naturally immutable**:
```go
type Temperature float64
type UUID [16]byte
```

3. **You never need to modify the receiver**:
```go
type Color struct { R, G, B uint8 }
func (c Color) Hex() string { ... }  // Read-only
```

### Use Pointer Receivers When:

1. **The method modifies the receiver**:
```go
type Counter struct { count int }
func (c *Counter) Increment() {
    c.count++  // Needs to modify
}
```

2. **The type is large** (copying would be expensive):
```go
type Image struct {
    pixels [1920][1080][3]uint8  // ~6MB!
}
func (i *Image) SetPixel(x, y int, color [3]uint8) {
    i.pixels[x][y] = color
}
```

3. **Consistency** (if any method needs a pointer, use pointers for all):
```go
type DB struct { conn *sql.Conn }

// If one method uses *DB, use *DB for all
func (db *DB) Query(q string) { ... }
func (db *DB) Close() { ... }
```

4. **The zero value is not useful** (needs initialization):
```go
type File struct { fd int }  // fd must be initialized via Open()
func (f *File) Read(buf []byte) (int, error) { ... }
```

### The Consistency Rule (Most Important!)

**Always be consistent:** If one method has a pointer receiver, all methods should have pointer receivers.

**Why?**
- Predictable interface satisfaction
- Avoids confusion about which methods modify
- Clearer API

```go
// ❌ BAD: Mixing receivers
type User struct { name string }
func (u User) Name() string { return u.name }      // Value
func (u *User) SetName(n string) { u.name = n }    // Pointer

// ✅ GOOD: Consistent receivers
func (u *User) Name() string { return u.name }     // Pointer
func (u *User) SetName(n string) { u.name = n }    // Pointer
```

---

## Performance Implications

### Copying Overhead

**Small types (≤ 3 words = 24 bytes on 64-bit):**
- Value receiver: ~1-2 ns (register or stack copy)
- Pointer receiver: ~1-2 ns (pass pointer)
- **Difference: Negligible**

**Large types (> 64 bytes):**
```go
type LargeStruct struct {
    data [1000]int  // 8000 bytes
}

// Value receiver: Copies 8000 bytes every call!
func (l LargeStruct) Process() { ... }

// Pointer receiver: Copies 8 bytes (the pointer)
func (l *LargeStruct) Process() { ... }
```

**Benchmark result:**
```
BenchmarkValueReceiver-8     1000000    1200 ns/op   8000 B/op
BenchmarkPointerReceiver-8  10000000     120 ns/op      0 B/op
```

**10x faster with pointer receiver!**

### Escape Analysis Impact

**Value receivers encourage stack allocation:**
```go
func (p Point) Distance() float64 { ... }

func main() {
    p := Point{X: 3, Y: 4}  // Stack allocated
    d := p.Distance()       // No heap escape
}
```

**Pointer receivers may force heap allocation:**
```go
func (p *Point) Distance() float64 { ... }

func main() {
    p := Point{X: 3, Y: 4}  // Might escape to heap (if passed to interface)
    d := p.Distance()       // Compiler might heap-allocate p
}
```

**Check with:**
```bash
go build -gcflags='-m' main.go
```

---

## Nil Receivers: A Special Case

Methods can be called on nil pointers! This enables nil-safe patterns.

### Example: Nil-Safe Methods

```go
type Tree struct {
    value int
    left  *Tree
    right *Tree
}

// This method works even if t is nil!
func (t *Tree) Sum() int {
    if t == nil {
        return 0  // Nil tree has sum 0
    }
    return t.value + t.left.Sum() + t.right.Sum()
}

func main() {
    var t *Tree  // nil
    fmt.Println(t.Sum())  // 0 (doesn't panic!)
}
```

**When is this useful?**
- Optional values (nil means "not set")
- Recursive data structures (nil terminates recursion)
- Lazy initialization

**Warning:**
```go
func (t *Tree) SetValue(v int) {
    t.value = v  // PANIC if t is nil!
}
```

Always check for nil if you dereference fields!

---

## Common Pitfalls

### Pitfall 1: Value Receiver Doesn't Modify

```go
type Counter struct { count int }

func (c Counter) Increment() {
    c.count++  // BUG: Modifies copy
}

func main() {
    c := Counter{count: 0}
    c.Increment()
    fmt.Println(c.count)  // 0 (unchanged!)
}
```

**Fix:** Use pointer receiver.

### Pitfall 2: Mixing Receivers Breaks Interfaces

```go
type Stringer interface {
    String() string
}

type Person struct { name string }

func (p Person) String() string { return p.name }
func (p *Person) SetName(n string) { p.name = n }

func main() {
    var s Stringer = Person{name: "Alice"}  // ✅ Works
    var s2 Stringer = &Person{name: "Bob"}  // ✅ Works

    // But now Person has inconsistent receivers!
}
```

**Fix:** Use all pointer receivers for consistency.

### Pitfall 3: Forgetting Interface Method Sets

```go
type Printer interface {
    Print()
}

type Doc struct{}
func (d *Doc) Print() { fmt.Println("printing") }

func PrintAll(docs []Printer) {
    for _, d := range docs {
        d.Print()
    }
}

func main() {
    docs := []Printer{
        Doc{},   // ❌ COMPILE ERROR: Doc doesn't implement Printer
        &Doc{},  // ✅ Works
    }
    PrintAll(docs)
}
```

**Remember:** Only `*Doc` implements `Printer` because `Print` has a pointer receiver.

### Pitfall 4: Taking Address of Map Element

```go
type Counter struct { count int }
func (c *Counter) Inc() { c.count++ }

func main() {
    m := map[string]Counter{
        "a": {count: 0},
    }

    m["a"].Inc()  // ❌ COMPILE ERROR: cannot take address of map element
}
```

**Why?**
- Map elements are not addressable (map might rehash/relocate)
- You cannot call pointer-receiver methods directly

**Fix:**
```go
// Option 1: Extract, modify, store back
c := m["a"]
c.Inc()
m["a"] = c

// Option 2: Use *Counter in the map
m := map[string]*Counter{
    "a": &Counter{count: 0},
}
m["a"].Inc()  // ✅ Works
```

---

## Decision Tree: Choosing the Right Receiver

```
START: Writing a method for type T

│
├─ Does the method modify the receiver?
│  ├─ YES → Use pointer receiver (*T)
│  └─ NO ↓
│
├─ Is T large (> 64 bytes)?
│  ├─ YES → Use pointer receiver (*T)
│  └─ NO ↓
│
├─ Does T contain a mutex (sync.Mutex, etc)?
│  ├─ YES → Use pointer receiver (*T) [MUST!]
│  └─ NO ↓
│
├─ Do other methods already use pointer receivers?
│  ├─ YES → Use pointer receiver (*T) [consistency]
│  └─ NO ↓
│
├─ Is T a map, slice, or interface?
│  ├─ YES → Use value receiver (T) [already reference-like]
│  └─ NO ↓
│
├─ Does T represent a value type (numbers, IDs, etc)?
│  ├─ YES → Use value receiver (T)
│  └─ NO → Use pointer receiver (*T) [default for structs]
```

---

## How to Run

```bash
# Run the demonstration program
cd minis/14-methods-value-vs-pointer-receivers
go run cmd/methods-demo/main.go

# Run tests
cd exercise
go test -v

# Check which variables escape to heap
go build -gcflags='-m' cmd/methods-demo/main.go

# Benchmark value vs pointer receivers
go test -bench=. -benchmem
```

---

## Expected Output (Demo Program)

```
=== Value Receiver: No Modification ===
Before: Counter{count: 0}
After:  Counter{count: 0}  ← Unchanged!

=== Pointer Receiver: Modification Works ===
Before: Counter{count: 0}
After:  Counter{count: 5}  ← Changed!

=== Interface Satisfaction ===
*Document implements Printer ✓
Document does NOT implement Printer ✗

=== Performance: Large Struct ===
Value receiver:   1200 ns/op, 8000 B/op
Pointer receiver:  120 ns/op,    0 B/op
Speedup: 10x faster!
```

---

## Key Takeaways

1. **Value receivers copy the value** (modifications don't persist)
2. **Pointer receivers modify the original** (use for mutations or large types)
3. **Method sets affect interface satisfaction** (only `*T` has pointer-receiver methods)
4. **Consistency is critical** (if one method uses `*T`, all should)
5. **Small types can use value receivers** (< 64 bytes, immutable data)
6. **Large types should use pointer receivers** (avoid copying overhead)
7. **Nil receiver checks enable nil-safe APIs** (check before dereferencing)
8. **Go auto-addresses/dereferences** (syntactic sugar for method calls)

---

## Connections to Other Projects

- **Project 12 (pointers-zero-values-nil-gotchas)**: Foundation for understanding receivers
- **Project 11 (slices-internals)**: Slices use value receivers (already reference-like)
- **Project 24 (sync-mutex-vs-rwmutex)**: Mutexes REQUIRE pointer receivers
- **Project 29 (escape-analysis-inlining)**: How receivers affect heap allocation
- **Project 48 (reflection-introspection)**: Inspecting method sets at runtime
- **Project 62 (interface-design-best-practices)**: Designing interfaces with receivers in mind

---

## Stretch Goals

1. **Implement a binary tree** with nil-safe pointer-receiver methods
2. **Benchmark** value vs pointer receivers for different struct sizes (16B, 64B, 1KB, 1MB)
3. **Create a type** that uses value receivers but still allows mutation (using slices/maps internally)
4. **Write a code analyzer** that detects inconsistent receiver types in a package
5. **Experiment with escape analysis** to see when value receivers prevent heap allocation
