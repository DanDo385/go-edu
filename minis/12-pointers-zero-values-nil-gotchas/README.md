# Project 12: Pointers, Zero Values, and Nil Gotchas

## What Is This Project About?

This project demystifies **pointers** in Go—one of the most powerful and dangerous features in systems programming. You'll learn:

1. **What pointers really are** (memory addresses, dereferencing, the & and * operators)
2. **Zero values for different types** (nil, 0, "", false, and composite types)
3. **The nil pointer dereference trap** (the most common runtime panic in Go)
4. **Pointer receivers vs value receivers** (when to use each)
5. **new vs make vs composite literals** (three ways to create values)
6. **Escape analysis for pointers** (stack vs heap allocation)

By the end, you'll understand when and why to use pointers, how to avoid nil panics, and how Go's memory model makes pointer arithmetic impossible (unlike C).

---

## The Fundamental Problem: Copying Is Expensive

### First Principles: Pass By Value

In Go, **everything is passed by value** by default. When you pass a variable to a function, Go **copies** it:

```go
type LargeStruct struct {
    data [1000000]int  // 8 MB of data!
}

func processData(s LargeStruct) {
    // s is a COPY of the original
    // This copied 8 MB of memory!
}

func main() {
    big := LargeStruct{}
    processData(big)  // Expensive copy!
}
```

**The problem:**
- Copying 8 MB is slow (CPU cache pollution, memory bandwidth)
- Modifications to `s` inside `processData` don't affect the original
- Wasted memory (two copies exist simultaneously)

**The solution:** Use a **pointer** to pass the memory address instead of copying the entire struct.

```go
func processData(s *LargeStruct) {
    // s is a pointer (just 8 bytes on 64-bit systems)
    // No copying of the struct!
}

func main() {
    big := LargeStruct{}
    processData(&big)  // Pass the address
}
```

---

## What Is a Pointer? (The Core Concept)

A **pointer** is a variable that stores a **memory address** of another variable.

### Memory Layout Example

```go
x := 42
p := &x  // p holds the address of x

fmt.Println(x)   // 42 (the value)
fmt.Println(p)   // 0xc0000140a8 (the address where x is stored)
fmt.Println(*p)  // 42 (dereference: follow the pointer to get the value)
```

**Memory diagram:**
```
Memory Address    Value
--------------    -----
0xc0000140a8 -->  42     (variable x)
0xc0000140b0 -->  0xc0000140a8  (variable p, stores address of x)
```

**Key operators:**
- `&x`: "Address of x" (gives you a pointer to x)
- `*p`: "Dereference p" (follows the pointer to get the value)

### Modifying Through Pointers

```go
x := 10
p := &x     // p points to x

*p = 20     // Modify x through the pointer

fmt.Println(x)  // 20 (x was changed!)
```

**What happened:**
1. `p` holds the address of `x`
2. `*p = 20` means "go to the address stored in p, and set the value to 20"
3. Since `p` points to `x`, this modifies `x`

---

## Zero Values: The Foundation of Safety

Go initializes **every variable** to a **zero value** (never uninitialized garbage like C).

### Zero Values for All Types

| Type              | Zero Value     | Safe to Use? |
|-------------------|----------------|--------------|
| `int`, `float64`  | `0`, `0.0`     | ✅ Yes       |
| `string`          | `""`           | ✅ Yes       |
| `bool`            | `false`        | ✅ Yes       |
| `[]int` (slice)   | `nil`          | ✅ Yes*      |
| `map[K]V`         | `nil`          | ⚠️ Read-only |
| `*T` (pointer)    | `nil`          | ❌ Dangerous |
| `chan T`          | `nil`          | ❌ Deadlock  |
| `func()`          | `nil`          | ❌ Panic     |
| `interface{}`     | `nil`          | ⚠️ Depends   |

**Key insights:**
- **Nil slices are safe**: `var s []int` (nil) can be used with `len()`, `append()`, iteration
- **Nil maps panic on write**: `var m map[string]int; m["key"] = 1` → panic!
- **Nil pointers panic on dereference**: `var p *int; *p = 1` → panic!

### Example: Nil Slice (Safe)

```go
var s []int  // nil slice

fmt.Println(len(s))      // 0 (safe)
fmt.Println(cap(s))      // 0 (safe)

for _, v := range s {    // No iterations (safe)
    fmt.Println(v)
}

s = append(s, 1)         // Works! (allocates backing array)
fmt.Println(s)           // [1]
```

### Example: Nil Map (Unsafe)

```go
var m map[string]int  // nil map

fmt.Println(m["key"])  // 0 (reads return zero value, safe)

m["key"] = 1           // PANIC: assignment to entry in nil map
```

**Fix:** Initialize with `make()`:
```go
m := make(map[string]int)
m["key"] = 1  // Works
```

---

## The Nil Pointer Dereference Trap

This is the **#1 runtime panic** in Go programs.

### The Classic Bug

```go
var p *int  // p is nil (zero value for pointer types)

*p = 42     // PANIC: invalid memory address or nil pointer dereference
```

**Why?**
`p` is `nil`, meaning it doesn't point to any valid memory address. Trying to dereference it is undefined behavior (in C, this could corrupt memory; Go panics instead).

### Real-World Example: Method on Nil Receiver

```go
type User struct {
    Name string
}

func (u *User) Greet() {
    fmt.Printf("Hello, %s!\n", u.Name)  // DANGER if u is nil!
}

func main() {
    var u *User  // nil pointer
    u.Greet()    // PANIC: nil pointer dereference
}
```

**Safe version:**
```go
func (u *User) Greet() {
    if u == nil {
        fmt.Println("Hello, guest!")
        return
    }
    fmt.Printf("Hello, %s!\n", u.Name)
}
```

### Defensive Nil Checks

Always check pointers before dereferencing:

```go
func processUser(u *User) {
    if u == nil {
        return  // or handle error
    }
    // Safe to use u now
    fmt.Println(u.Name)
}
```

---

## Pointer Receivers vs Value Receivers

When defining methods, you choose whether the receiver is a pointer or a value.

### Value Receiver (Copy)

```go
type Counter struct {
    count int
}

func (c Counter) Increment() {
    c.count++  // Modifies a COPY, not the original!
}

func main() {
    c := Counter{count: 0}
    c.Increment()
    fmt.Println(c.count)  // 0 (unchanged!)
}
```

**What happened?**
`Increment()` received a **copy** of `c`. The copy was modified, but the original remained unchanged.

### Pointer Receiver (Modify Original)

```go
func (c *Counter) Increment() {
    c.count++  // Modifies the original
}

func main() {
    c := Counter{count: 0}
    c.Increment()
    fmt.Println(c.count)  // 1 (changed!)
}
```

**Go's syntactic sugar:**
Even though `c` is a value and `Increment` expects `*Counter`, Go automatically does `(&c).Increment()` for you.

### When to Use Pointer Receivers

Use pointer receivers when:
1. **The method modifies the receiver**
2. **The receiver is large** (copying would be expensive)
3. **Consistency** (if one method uses a pointer receiver, use them for all methods)

### The Method Set Rule

**Critical for interfaces:**
- A type `T` has methods with receiver `T` and `*T` (automatic addressing)
- A type `*T` has **only** methods with receiver `*T`

**Example:**
```go
type Printer interface {
    Print()
}

type Doc struct{}

func (d Doc) Print() {
    fmt.Println("printing")
}

func main() {
    var p Printer

    p = Doc{}   // ✅ Doc implements Printer
    p = &Doc{}  // ✅ *Doc also implements Printer (methods of Doc are included)
}
```

**But if the method has a pointer receiver:**
```go
func (d *Doc) Print() {
    fmt.Println("printing")
}

func main() {
    var p Printer

    p = Doc{}   // ❌ COMPILE ERROR: Doc does NOT implement Printer
    p = &Doc{}  // ✅ *Doc implements Printer
}
```

---

## new vs make vs Composite Literals

Go has **three ways** to allocate memory. Each serves a different purpose.

### 1. Composite Literals (Most Common)

```go
// Structs
u := User{Name: "Alice"}        // Value
u := &User{Name: "Alice"}       // Pointer to struct

// Slices
s := []int{1, 2, 3}

// Maps
m := map[string]int{"a": 1}
```

**Use this when:** You know the initial values.

### 2. make() (For Slices, Maps, Channels ONLY)

```go
s := make([]int, 10)           // Slice with len=10, cap=10
m := make(map[string]int)      // Empty map (initialized, not nil)
ch := make(chan int, 5)        // Buffered channel
```

**Purpose:** Allocates and **initializes** the internal data structure.
- For slices: allocates backing array
- For maps: allocates hash table
- For channels: allocates channel buffer

**Use this when:** You need an empty but ready-to-use slice/map/channel.

### 3. new() (Allocates Zeroed Memory)

```go
p := new(int)        // Allocates an int, returns *int
fmt.Println(*p)      // 0 (zero value)

u := new(User)       // Allocates a User, returns *User
fmt.Println(u.Name)  // "" (zero value for string)
```

**What `new(T)` does:**
1. Allocates memory for a value of type `T`
2. Initializes it to the zero value
3. Returns a pointer `*T`

**Equivalent to:**
```go
var x T
p := &x
```

**Use this when:** You need a pointer to a zero-initialized value (rare; composite literals are usually better).

### Comparison Table

| Syntax               | Returns   | Use Case                          |
|----------------------|-----------|-----------------------------------|
| `T{}`                | `T`       | Value with initial values         |
| `&T{}`               | `*T`      | Pointer with initial values       |
| `new(T)`             | `*T`      | Pointer to zero value             |
| `make([]T, n)`       | `[]T`     | Initialized slice                 |
| `make(map[K]V)`      | `map[K]V` | Initialized map                   |
| `make(chan T, n)`    | `chan T`  | Initialized channel               |

---

## Escape Analysis: Stack vs Heap

Go's compiler decides whether to allocate on the **stack** or **heap** based on **escape analysis**.

### Stack Allocation (Fast)

```go
func foo() {
    x := 42  // Allocated on foo's stack frame
    // When foo returns, x is automatically destroyed
}
```

**Characteristics:**
- Very fast (just move the stack pointer)
- Automatically cleaned up when function returns
- Limited size (~1 MB per goroutine)

### Heap Allocation (Slower, GC'd)

```go
func makePointer() *int {
    x := 42
    return &x  // x ESCAPES to the heap (returned to caller)
}
```

**Characteristics:**
- Allocated in the heap (can be arbitrarily large)
- Survives after function returns
- Requires garbage collection (GC scans, marks, sweeps)

### How to Check (Compiler Flags)

```bash
go build -gcflags='-m' main.go
```

**Output:**
```
./main.go:5:2: moved to heap: x
./main.go:3:6: can inline makePointer
```

### Common Escape Scenarios

1. **Returning a pointer:** `return &x`
2. **Storing in a heap-allocated struct:** `globalVar.field = &x`
3. **Passing to interface:** `fmt.Println(x)` (interface wrapping may escape)
4. **Large variables:** Arrays/structs > ~64KB
5. **Closure captures:** Lambdas that outlive their scope

### Optimizing for Stack Allocation

```go
// HEAP (returns pointer)
func sumHeap(nums []int) *int {
    sum := 0
    for _, n := range nums {
        sum += n
    }
    return &sum  // Escapes!
}

// STACK (returns value)
func sumStack(nums []int) int {
    sum := 0
    for _, n := range nums {
        sum += n
    }
    return sum  // No escape, stays on stack
}
```

---

## Common Pitfalls

### Pitfall 1: Modifying Value Receiver

```go
func (c Counter) Increment() {
    c.count++  // BUG: Modifies copy
}
```

**Fix:** Use pointer receiver.

### Pitfall 2: Returning Pointer to Local Array

```go
func makeArray() *[3]int {
    arr := [3]int{1, 2, 3}
    return &arr  // OK in Go (escapes to heap), but consider returning slice
}
```

**Better:** Return a slice.

### Pitfall 3: Forgetting to Initialize Maps/Channels

```go
var m map[string]int
m["key"] = 1  // PANIC
```

**Fix:** `m := make(map[string]int)`.

### Pitfall 4: Nil Interface Gotcha

```go
var p *int  // nil pointer
var i interface{} = p  // i is NOT nil!

if i == nil {
    fmt.Println("nil")  // NOT printed!
}
```

**Why?**
An interface is `nil` only if **both** the type and value are `nil`. Here, type is `*int` (not nil), value is `nil`.

**Fix:** Check the value inside the interface or avoid wrapping nil pointers.

---

## How to Run

```bash
# Run the demonstration
cd minis/12-pointers-zero-values-nil-gotchas
go run cmd/pointers-demo/main.go

# Run exercises
cd exercise
go test -v

# Check escape analysis
go build -gcflags='-m' cmd/pointers-demo/main.go
```

---

## Key Takeaways

1. **Pointers store addresses**, not values (use `&` to get address, `*` to dereference)
2. **Zero values prevent uninitialized bugs** (but nil pointers/maps/channels are dangerous)
3. **Always check for nil** before dereferencing pointers
4. **Use pointer receivers** for methods that modify the receiver or for large types
5. **new allocates + zeros**, `make` allocates + initializes, composites do both with values
6. **Escape analysis decides stack vs heap** (prefer stack for performance)
7. **Go has no pointer arithmetic** (safe by design, unlike C)

---

## Connections to Other Projects

- **Project 11 (slices-internals)**: Slice headers contain pointers to backing arrays
- **Project 14 (methods-value-vs-pointer-receivers)**: Deep dive into receiver choice
- **Project 24 (sync-mutex-vs-rwmutex)**: Mutexes are passed by pointer (value receivers break locks)
- **Project 29 (escape-analysis-inlining)**: Detailed analysis of stack vs heap
- **Project 48 (reflection-introspection)**: Using reflect to inspect pointer types

---

## Stretch Goals

1. **Implement a simple linked list** using pointers
2. **Benchmark** value vs pointer receiver for large structs
3. **Write a function** that safely dereferences pointers with default values
4. **Create a nil-safe map wrapper** that auto-initializes on first write
