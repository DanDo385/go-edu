# Project 13: Interfaces, Duck Typing & Method Sets

## What Is This Project About?

This project demystifies **interfaces** in Go—one of the most elegant and powerful features in the language. You'll learn:

1. **What interfaces really are** (implicit implementation, method sets, runtime polymorphism)
2. **Duck typing in Go** ("if it walks like a duck and quacks like a duck...")
3. **The empty interface** (`interface{}` and `any`)
4. **Type assertions and type switches** (extracting concrete types)
5. **Nil interface gotchas** (the two-part nil problem)
6. **Interface satisfaction rules** (pointer vs value receivers)
7. **Real-world interface patterns** (io.Reader, error, Stringer)

By the end, you'll understand how Go achieves polymorphism without inheritance, why interfaces enable loose coupling, and how to avoid common interface pitfalls.

---

## The Fundamental Problem: Rigid Type Systems

### First Principles: The Coupling Problem

In traditional object-oriented languages, you create tightly coupled hierarchies:

```java
// Java-style rigid hierarchy
class Animal { }
class Dog extends Animal { }
class Cat extends Animal { }

// If you later want to add a Robot that barks,
// you can't make it extend Animal (it's not an animal!)
```

**The problems:**
- Inheritance creates tight coupling between types
- You must know all types at design time
- Difficult to make unrelated types compatible
- Forces awkward "is-a" relationships

### Real-World Scenario

Imagine building a logging system. You want to log to:
- Files
- Network sockets
- Cloud storage
- stdout/stderr
- Test buffers

In a language with classes, you'd create a `Logger` base class and force every output destination to inherit from it. But what if you want to use an existing library's `FileWriter` that doesn't inherit from your `Logger`?

**This is why interfaces exist.**

---

## What Is an Interface? (The Core Concept)

An **interface** in Go is a **contract** that specifies behavior (methods) without specifying implementation. Any type that implements all the methods of an interface **automatically** satisfies that interface—no explicit declaration needed.

### The Interface Declaration

```go
// An interface is a collection of method signatures
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**This says:** "A Writer is anything that has a `Write` method with this exact signature."

### Implicit Implementation (Duck Typing)

**Go's revolutionary idea:** You don't declare "I implement Writer." If your type has the right methods, it **is** a Writer.

```go
// File implements Writer (implicitly)
type File struct {
    name string
}

func (f *File) Write(p []byte) (n int, err error) {
    // Write bytes to file
    return len(p), nil
}

// Network implements Writer (implicitly)
type Network struct {
    addr string
}

func (n *Network) Write(p []byte) (n int, err error) {
    // Send bytes over network
    return len(p), nil
}

// Both File and Network are Writers, no declaration needed!
```

**Memory layout:**

An interface value is a **two-word structure**:

```
Interface value:
┌────────────────┬─────────────────┐
│ Type pointer   │ Value pointer   │
│ (*File)        │ 0x400020 ────┐  │
└────────────────┴──────────────┼──┘
                                │
                                └──> Actual File struct:
                                     ┌──────────────┐
                                     │ name: "log"  │
                                     └──────────────┘
```

**Key insight:** An interface stores **both** the type information and the value. This allows runtime polymorphism.

---

## Duck Typing: "If It Quacks Like a Duck..."

The famous duck test: "If it walks like a duck and quacks like a duck, it's a duck."

### Go's Version

```go
type Quacker interface {
    Quack() string
}

// Duck is obviously a Quacker
type Duck struct{}

func (d Duck) Quack() string {
    return "Quack!"
}

// Robot happens to quack too
type Robot struct{}

func (r Robot) Quack() string {
    return "Mechanical quack!"
}

// Both are Quackers, even though they're completely unrelated types
func MakeItQuack(q Quacker) {
    fmt.Println(q.Quack())
}

func main() {
    MakeItQuack(Duck{})   // Works!
    MakeItQuack(Robot{})  // Also works!
}
```

**Why this matters:**
- You can make unrelated types compatible
- Library code can work with types it never knew about
- Decouples interface definition from implementation

---

## The Empty Interface: `interface{}` and `any`

An interface with **zero methods** is satisfied by **every type**.

### The Universal Type

```go
var x interface{}  // Can hold ANY value

x = 42
x = "hello"
x = []int{1, 2, 3}
x = struct{ Name string }{"Alice"}
```

**Since Go 1.18:**
```go
var x any  // Alias for interface{}, more readable
```

### Why Empty Interfaces Exist

Before generics (Go 1.18), the empty interface was the only way to write functions that accept any type:

```go
func Print(value interface{}) {
    fmt.Println(value)
}

Print(42)
Print("hello")
Print([]int{1, 2, 3})
```

**The tradeoff:**
- **Flexibility:** Accepts any type
- **Type safety:** Lost at compile time (must use type assertions at runtime)

### Modern Alternative: Generics

```go
// Go 1.18+: Use generics for type-safe generic code
func Print[T any](value T) {
    fmt.Println(value)
}
```

---

## Type Assertions: Extracting Concrete Types

When you have an interface value, you can extract the underlying concrete type using **type assertions**.

### Syntax

```go
var w Writer = &File{name: "log.txt"}

// Type assertion: "w is a *File"
f := w.(*File)  // Panics if w is not a *File

// Safe type assertion with boolean check
f, ok := w.(*File)
if ok {
    fmt.Println("It's a file:", f.name)
} else {
    fmt.Println("Not a file")
}
```

### Type Switches: Pattern Matching

For checking multiple types, use a **type switch**:

```go
func Describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    case *File:
        fmt.Printf("File: %s\n", v.name)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

**Performance note:** Type assertions and switches are fast (single pointer comparison + vtable lookup).

---

## Nil Interface Gotchas: The Two-Part Nil

This is the **most confusing** aspect of interfaces in Go.

### The Problem

```go
var p *int = nil
var i interface{} = p

fmt.Println(p == nil)  // true
fmt.Println(i == nil)  // false ← SURPRISE!
```

**Why?**

An interface is `nil` **only if both** the type and value are `nil`:

```
Nil interface:
┌──────┬──────┐
│ nil  │ nil  │  ← BOTH parts are nil
└──────┴──────┘

Non-nil interface with nil value:
┌──────┬──────┐
│ *int │ nil  │  ← Type is set, value is nil
└──────┴──────┘
```

### Real-World Bug

```go
func GetUser(id int) *User {
    if id == 0 {
        return nil  // Return nil pointer
    }
    return &User{ID: id}
}

func Process() error {
    user := GetUser(0)  // user is nil
    return user         // BUG: Returns non-nil error!
}

func main() {
    err := Process()
    if err != nil {
        // This always executes, even though user was nil!
        fmt.Println("Error:", err)
    }
}
```

**Why?** Returning `user` (a `*User = nil`) as an `error` interface creates an interface with `(type=*User, value=nil)`, which is **not** `nil`.

### The Fix

```go
func Process() error {
    user := GetUser(0)
    if user == nil {
        return nil  // Return typed nil
    }
    return user
}
```

**Rule of thumb:** Never return a typed nil pointer as an interface. Return the interface-typed `nil` directly.

---

## Interface Satisfaction Rules: Method Sets

Whether a type satisfies an interface depends on its **method set**.

### The Rule

- A type `T` has methods with receiver `T` (value receiver)
- A type `*T` has methods with receiver `T` **and** `*T` (pointer receiver)

**Critical consequence:** If an interface method has a pointer receiver, only the pointer type satisfies the interface.

### Example 1: Value Receiver

```go
type Stringer interface {
    String() string
}

type Person struct {
    Name string
}

// Value receiver
func (p Person) String() string {
    return p.Name
}

func main() {
    var s Stringer

    s = Person{Name: "Alice"}   // ✅ Person implements Stringer
    s = &Person{Name: "Bob"}    // ✅ *Person also implements Stringer
}
```

**Why both work?** Methods with value receivers are in the method set of both `T` and `*T`.

### Example 2: Pointer Receiver

```go
type Incrementer interface {
    Increment()
}

type Counter struct {
    count int
}

// Pointer receiver (must modify the original)
func (c *Counter) Increment() {
    c.count++
}

func main() {
    var inc Incrementer

    inc = Counter{count: 0}      // ❌ COMPILE ERROR: Counter doesn't implement Incrementer
    inc = &Counter{count: 0}     // ✅ *Counter implements Incrementer
}
```

**Why?** Methods with pointer receivers are **only** in the method set of `*T`, not `T`.

**Intuition:** Go can't automatically take the address of a value in all contexts (e.g., map values, interface values), so it can't auto-promote `T` to `*T` for pointer receiver methods.

### The Addressability Exception

```go
c := Counter{count: 0}
c.Increment()  // ✅ Works! Go does (&c).Increment()
```

**Why?** When you call a method on a **variable**, Go can take its address. But when assigning to an interface, the value might not be addressable.

---

## Real-World Interface Patterns

Go's standard library uses interfaces extensively. Understanding these patterns is key to idiomatic Go.

### Pattern 1: io.Reader and io.Writer

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**Used by:** Files, network connections, buffers, compression, encryption, HTTP bodies, etc.

**Why powerful:** You can compose readers/writers (e.g., `gzip.NewReader(file)` wraps a file reader with decompression).

### Pattern 2: error Interface

```go
type error interface {
    Error() string
}
```

**Why it's an interface:** Allows custom error types with additional context.

```go
type MyError struct {
    Code int
    Msg  string
}

func (e MyError) Error() string {
    return fmt.Sprintf("error %d: %s", e.Code, e.Msg)
}
```

### Pattern 3: fmt.Stringer

```go
type Stringer interface {
    String() string
}
```

**Used by:** `fmt.Println`, `fmt.Printf` with `%v` and `%s`.

```go
type Point struct {
    X, Y int
}

func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func main() {
    fmt.Println(Point{3, 4})  // Calls String(), prints: (3, 4)
}
```

### Pattern 4: Interface Composition

```go
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

**Interfaces can embed other interfaces**, creating larger contracts from smaller ones.

---

## Performance Implications

### Interface Call Overhead

**Direct call (non-interface):**
```go
f := &File{name: "log"}
f.Write(data)  // Direct method call (inlinable)
```

**Interface call:**
```go
var w Writer = &File{name: "log"}
w.Write(data)  // Indirect call through vtable (not inlinable)
```

**Cost:** ~1-2 nanoseconds overhead for the indirection. Usually negligible, but matters in ultra-hot loops.

### Escape Analysis

```go
func Direct() {
    f := File{name: "log"}  // Stack-allocated
}

func ViaInterface() {
    var w Writer = &File{name: "log"}  // Escapes to heap!
}
```

**Why?** Assigning to an interface often forces heap allocation because the compiler can't determine the size at compile time.

### When to Avoid Interfaces

- Ultra-hot loops (billions of iterations)
- Embedded systems with tight memory constraints
- When you have only one implementation (interfaces add complexity for no benefit)

### When to Use Interfaces

- Abstracting I/O (files, network, etc.)
- Testing (mock implementations)
- Plugin architectures
- Decoupling packages

---

## Common Pitfalls and How to Avoid Them

### Pitfall 1: Returning Typed Nil

```go
// WRONG
func GetUser() *User {
    return nil  // Returns typed nil
}

func main() {
    if GetUser() != nil {
        // Might execute even though user is nil!
    }
}
```

**Fix:** Always return interface-typed nil directly.

### Pitfall 2: Comparing Interface Values

```go
var a, b Writer
a = &File{name: "log"}
b = &File{name: "log"}

fmt.Println(a == b)  // false (different pointers!)
```

**Why?** Interface comparison compares **both** type and value. Two different File pointers are not equal.

### Pitfall 3: Modifying Through Interface

```go
type Setter interface {
    Set(int)
}

type Value struct {
    x int
}

func (v Value) Set(x int) {  // Value receiver!
    v.x = x  // Modifies a copy
}

func main() {
    var s Setter = Value{}
    s.Set(42)  // Does nothing (modifies copy)
}
```

**Fix:** Use pointer receiver for methods that modify the receiver.

### Pitfall 4: Over-Engineering with Interfaces

```go
// WRONG: Premature abstraction
type UserServiceInterface interface {
    GetUser(id int) *User
    CreateUser(name string) *User
    // ... 20 more methods
}

// Only one implementation exists!
```

**Fix:** "Accept interfaces, return structs." Only define interfaces when you have multiple implementations or need to decouple for testing.

---

## How to Run

```bash
# Run the demonstration program
cd minis/13-interfaces-duck-typing
go run cmd/interfaces-demo/main.go

# Run tests
cd exercise
go test -v

# Check which types implement which interfaces
go doc io.Writer
```

---

## Expected Output (Demo Program)

```
=== Interface Basics ===
File writes: Hello, World!
Network writes: Hello, World!

=== Type Assertions ===
It's a file: log.txt

=== Type Switches ===
Integer: 42
String: hello
File: data.txt

=== Nil Interface Gotcha ===
Pointer is nil: true
Interface is nil: false  ← GOTCHA!

=== Empty Interface ===
Value: 42, Type: int
Value: hello, Type: string
```

---

## Key Takeaways

1. **Interfaces are implicit** (duck typing: if it has the methods, it's the interface)
2. **Interfaces are two words** (type pointer + value pointer)
3. **Empty interface accepts anything** (but loses type safety)
4. **Type assertions extract concrete types** (use `v, ok := i.(T)` for safety)
5. **Nil interfaces are tricky** (both type and value must be nil)
6. **Method set rules matter** (pointer receivers → only `*T` implements the interface)
7. **Accept interfaces, return structs** (keep APIs flexible, implementations concrete)

---

## Connections to Other Projects

- **Project 12 (pointers-zero-values-nil-gotchas)**: Understanding nil is critical for interfaces
- **Project 14 (methods-value-vs-pointer-receivers)**: Method sets determine interface satisfaction
- **Project 15 (error-wrapping-sentinel-errors)**: The error interface is Go's error handling foundation
- **Project 19 (channels-basics)**: Channels are interfaces to goroutine communication
- **Project 48 (reflection-introspection)**: Reflection uses interfaces extensively

---

## Stretch Goals

1. **Implement a plugin system** using interfaces for extensibility
2. **Create a custom error type** with structured fields
3. **Build a mock HTTP server** using interface-based dependency injection
4. **Benchmark** interface calls vs direct calls in a hot loop
5. **Write a type-safe generic cache** using interfaces (pre-generics style)
