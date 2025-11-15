# Project 46: Generics and Parallel Map-Reduce

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a data processing system that needs to:
- Transform millions of records (map)
- Aggregate results (reduce)
- Work with different data types (users, transactions, logs)
- Process data in parallel for speed

**Without generics**, you'd write:
```go
// Separate function for each type
func MapInts(data []int, fn func(int) int) []int { ... }
func MapStrings(data []string, fn func(string) string) []string { ... }
func ReduceInts(data []int, fn func(int, int) int) int { ... }
// ... hundreds more functions
```

**With generics**, you write once:
```go
// Works for ANY type
func Map[T, U any](data []T, fn func(T) U) []U { ... }
func Reduce[T, U any](data []T, initial U, fn func(U, T) U) U { ... }
```

This project teaches you:
1. **Go generics** from first principles
2. **Type parameters and constraints**
3. **Parallel map-reduce engine** for high-performance data processing
4. **Real-world applications** of functional programming patterns

### What You'll Learn

1. **Generics fundamentals**: Type parameters, constraints, instantiation
2. **Type constraints**: any, comparable, custom interfaces, type unions
3. **Generic functions**: Map, Filter, Reduce, FlatMap
4. **Generic data structures**: Optional[T], Result[T, E], Pair[A, B]
5. **Parallel processing**: Worker pools, goroutines, channel coordination
6. **Map-Reduce pattern**: Distributed computation model
7. **Performance**: Benchmarking generic vs non-generic code

### The Challenge

Build a generic map-reduce engine that:
- Processes data in parallel using worker pools
- Works with any data type
- Provides common functional programming primitives
- Achieves near-linear speedup with multiple cores
- Is type-safe at compile time
- Has clean, reusable APIs

---

## 2. First Principles: Understanding Generics

### What Problem Do Generics Solve?

**The Fundamental Dilemma**: You want to write reusable code that works with different types, but Go is statically typed.

#### Option 1: Code Duplication (Type-Specific)

```go
func SumInts(nums []int) int {
    sum := 0
    for _, n := range nums {
        sum += n
    }
    return sum
}

func SumFloats(nums []float64) float64 {
    sum := 0.0
    for _, n := range nums {
        sum += n
    }
    return sum
}

// Need separate functions for every type!
```

**Problems**:
- Code duplication (violates DRY principle)
- More code to maintain and test
- Can't abstract common patterns

#### Option 2: Interface{} (Type Erasure)

```go
func Sum(nums []interface{}) interface{} {
    sum := 0
    for _, n := range nums {
        sum += n.(int)  // Runtime type assertion - can panic!
    }
    return sum
}
```

**Problems**:
- Loses type safety (runtime panics possible)
- Requires type assertions everywhere
- Poor performance (interface allocations, type switches)
- Can't use operators (+, -, *, /) on interface{}

#### Option 3: Generics (Type Parameters) ✅

```go
func Sum[T int | float64](nums []T) T {
    var sum T
    for _, n := range nums {
        sum += n  // Compiler knows T supports +
    }
    return sum
}

// Usage:
ints := Sum[int]([]int{1, 2, 3})        // = 6
floats := Sum[float64]([]float64{1.5, 2.5})  // = 4.0
```

**Benefits**:
- ✅ Write once, use for multiple types
- ✅ Type-safe at compile time (no runtime panics)
- ✅ No performance penalty (monomorphization)
- ✅ Clean, readable code

### How Do Generics Work?

**Type Parameters** are placeholders for concrete types:

```go
// [T any] is a type parameter
//  T is the parameter name
//  any is the constraint (any type allowed)
func Identity[T any](val T) T {
    return val
}

// Instantiation (compiler generates specific versions):
x := Identity[int](42)        // T = int
y := Identity[string]("hello") // T = string
```

**Behind the scenes**, Go uses **monomorphization**:
- Compiler generates a separate function for each type used
- `Identity[int]` becomes `Identity_int`
- `Identity[string]` becomes `Identity_string`
- No runtime overhead (unlike Java/C# generics with type erasure)

### Type Constraints

**Constraints** limit which types can be used:

#### 1. `any` Constraint (No Restrictions)

```go
// T can be ANY type
func First[T any](slice []T) T {
    return slice[0]
}
```

#### 2. `comparable` Constraint (Supports == and !=)

```go
// T must support equality comparison
func Contains[T comparable](slice []T, val T) bool {
    for _, item := range slice {
        if item == val {  // Requires T to be comparable
            return true
        }
    }
    return false
}

// Works: int, string, bool, pointers, structs with comparable fields
// Doesn't work: slices, maps, functions
```

#### 3. Type Union (One of Several Types)

```go
// T can be int OR float64
func Sum[T int | float64](nums []T) T {
    var sum T
    for _, n := range nums {
        sum += n  // Both int and float64 support +
    }
    return sum
}
```

#### 4. Interface Constraint (Must Implement Methods)

```go
// T must implement String() method
type Stringer interface {
    String() string
}

func PrintAll[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())  // OK: T implements Stringer
    }
}
```

#### 5. Approximate Type Constraint (~T)

```go
// T can be any type with underlying type int
type MyInt int

type Integer interface {
    ~int  // Matches int, MyInt, and any other type based on int
}

func Double[T Integer](n T) T {
    return n * 2
}

var x MyInt = 5
y := Double(x)  // Works! MyInt has underlying type int
```

#### 6. Combined Constraints

```go
// T must be comparable AND implement Stringer
type ComparableStringer interface {
    comparable
    Stringer
}

func FindFirst[T ComparableStringer](slice []T, val T) (T, bool) {
    for _, item := range slice {
        if item == val {
            return item, true
        }
    }
    var zero T
    return zero, false
}
```

---

## 3. Map-Reduce: The Fundamental Pattern

### What is Map-Reduce?

**Map-Reduce** is a programming model for processing large datasets by:
1. **Map**: Transform each element independently
2. **Reduce**: Combine all elements into a single result

**Origin**: Introduced by Google (2004) for distributed computation across thousands of machines.

**Why it matters**: Embarrassingly parallel - perfect for multi-core CPUs and distributed systems.

### The Map Operation

**Map** applies a function to each element, producing a new collection:

```
Input:  [1, 2, 3, 4]
Map:    x => x * 2
Output: [2, 4, 6, 8]
```

**Generic implementation**:
```go
func Map[T, U any](data []T, fn func(T) U) []U {
    result := make([]U, len(data))
    for i, item := range data {
        result[i] = fn(item)
    }
    return result
}

// Usage:
numbers := []int{1, 2, 3, 4}
doubled := Map(numbers, func(x int) int { return x * 2 })
// [2, 4, 6, 8]

strings := Map(numbers, func(x int) string { return fmt.Sprintf("#%d", x) })
// ["#1", "#2", "#3", "#4"]
```

**Key insight**: Type parameter `U` can differ from `T` - you can transform types!

### The Reduce Operation

**Reduce** combines all elements into a single value:

```
Input:  [1, 2, 3, 4]
Reduce: (acc, x) => acc + x, initial: 0
Steps:  0 + 1 = 1
        1 + 2 = 3
        3 + 3 = 6
        6 + 4 = 10
Output: 10
```

**Generic implementation**:
```go
func Reduce[T, U any](data []T, initial U, fn func(U, T) U) U {
    acc := initial
    for _, item := range data {
        acc = fn(acc, item)
    }
    return acc
}

// Usage:
numbers := []int{1, 2, 3, 4}

sum := Reduce(numbers, 0, func(acc, x int) int { return acc + x })
// 10

product := Reduce(numbers, 1, func(acc, x int) int { return acc * x })
// 24

concat := Reduce(numbers, "", func(acc string, x int) string {
    return acc + fmt.Sprintf("%d", x)
})
// "1234"
```

### Parallel Map-Reduce

**The Power**: Map operations are independent - can run in parallel!

**Sequential map** (1 core):
```
[1, 2, 3, 4, 5, 6, 7, 8]
 ↓  ↓  ↓  ↓  ↓  ↓  ↓  ↓
[2, 4, 6, 8, 10, 12, 14, 16]
Time: 8 units
```

**Parallel map** (4 cores):
```
Core 1: [1, 2] → [2, 4]
Core 2: [3, 4] → [6, 8]
Core 3: [5, 6] → [10, 12]
Core 4: [7, 8] → [14, 16]
Time: 2 units (4x speedup!)
```

**Implementation strategy**:
1. Split data into chunks (one per worker)
2. Process each chunk in a goroutine
3. Collect results via channels
4. Combine results

---

## 4. Building the Generic Map-Reduce Engine

### Design: Type Parameters and Constraints

```go
// MapFunc transforms T → U
type MapFunc[T, U any] func(T) U

// ReduceFunc combines U + T → U
type ReduceFunc[T, U any] func(accumulator U, item T) U

// FilterFunc decides whether to keep T
type FilterFunc[T any] func(T) bool
```

### Sequential Implementations

#### Generic Map

```go
func Map[T, U any](data []T, fn func(T) U) []U {
    result := make([]U, len(data))
    for i, item := range data {
        result[i] = fn(item)
    }
    return result
}
```

**Time complexity**: O(n)
**Space complexity**: O(n)

#### Generic Filter

```go
func Filter[T any](data []T, predicate func(T) bool) []T {
    result := make([]T, 0, len(data))
    for _, item := range data {
        if predicate(item) {
            result = append(result, item)
        }
    }
    return result
}
```

**Example**:
```go
numbers := []int{1, 2, 3, 4, 5, 6}
evens := Filter(numbers, func(x int) bool { return x%2 == 0 })
// [2, 4, 6]
```

#### Generic Reduce

```go
func Reduce[T, U any](data []T, initial U, fn func(U, T) U) U {
    acc := initial
    for _, item := range data {
        acc = fn(acc, item)
    }
    return acc
}
```

#### Generic FlatMap

```go
func FlatMap[T, U any](data []T, fn func(T) []U) []U {
    result := make([]U, 0, len(data))
    for _, item := range data {
        result = append(result, fn(item)...)
    }
    return result
}
```

**Example**:
```go
words := []string{"hello", "world"}
chars := FlatMap(words, func(s string) []rune {
    return []rune(s)
})
// ['h', 'e', 'l', 'l', 'o', 'w', 'o', 'r', 'l', 'd']
```

### Parallel Map Implementation

**Strategy**: Split work across goroutines using worker pool pattern.

```go
func ParallelMap[T, U any](data []T, fn func(T) U, numWorkers int) []U {
    n := len(data)
    result := make([]U, n)

    // Channel for work items
    jobs := make(chan int, n)

    // Start workers
    var wg sync.WaitGroup
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := range jobs {
                result[i] = fn(data[i])
            }
        }()
    }

    // Send work
    for i := 0; i < n; i++ {
        jobs <- i
    }
    close(jobs)

    // Wait for completion
    wg.Wait()

    return result
}
```

**Key decisions**:
1. **Pre-allocate result slice**: Avoid synchronization when writing results
2. **Index-based approach**: Each worker writes to different indices (no races)
3. **Channel for work distribution**: Fair distribution across workers
4. **WaitGroup for synchronization**: Know when all workers are done

### Parallel Reduce Implementation

**Challenge**: Reduce is inherently sequential (depends on previous results).

**Solution**: Two-phase approach:
1. **Map phase**: Split data, reduce each chunk in parallel
2. **Reduce phase**: Combine chunk results sequentially

```go
func ParallelReduce[T, U any](data []T, initial U, fn func(U, T) U, numWorkers int) U {
    n := len(data)
    if n == 0 {
        return initial
    }

    // Calculate chunk size
    chunkSize := (n + numWorkers - 1) / numWorkers

    // Channel for partial results
    partials := make(chan U, numWorkers)

    // Process chunks in parallel
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if end > n {
            end = n
        }
        if start >= n {
            break
        }

        wg.Add(1)
        go func(chunk []T) {
            defer wg.Done()
            acc := initial
            for _, item := range chunk {
                acc = fn(acc, item)
            }
            partials <- acc
        }(data[start:end])
    }

    // Close channel when all workers done
    go func() {
        wg.Wait()
        close(partials)
    }()

    // Combine partial results
    acc := initial
    for partial := range partials {
        acc = fn(acc, partial)
    }

    return acc
}
```

**Important note**: This only works for **associative** operations:
- ✅ Addition: `(a + b) + c = a + (b + c)`
- ✅ Multiplication: `(a * b) * c = a * (b * c)`
- ❌ Subtraction: `(a - b) - c ≠ a - (b - c)`

---

## 5. Advanced Generic Patterns

### Pattern 1: Optional[T] (Type-Safe Null)

```go
type Optional[T any] struct {
    value T
    valid bool
}

func Some[T any](val T) Optional[T] {
    return Optional[T]{value: val, valid: true}
}

func None[T any]() Optional[T] {
    return Optional[T]{valid: false}
}

func (o Optional[T]) Get() (T, bool) {
    return o.value, o.valid
}

func (o Optional[T]) OrElse(defaultVal T) T {
    if o.valid {
        return o.value
    }
    return defaultVal
}

func (o Optional[T]) Map[U any](fn func(T) U) Optional[U] {
    if !o.valid {
        return None[U]()
    }
    return Some(fn(o.value))
}
```

**Usage**:
```go
func FindUser(id int) Optional[User] {
    user, found := db.Get(id)
    if !found {
        return None[User]()
    }
    return Some(user)
}

// Chaining operations
userName := FindUser(123).
    Map(func(u User) string { return u.Name }).
    OrElse("Anonymous")
```

### Pattern 2: Result[T, E] (Error Handling)

```go
type Result[T, E any] struct {
    value T
    err   E
    ok    bool
}

func Ok[T, E any](val T) Result[T, E] {
    return Result[T, E]{value: val, ok: true}
}

func Err[T, E any](err E) Result[T, E] {
    return Result[T, E]{err: err, ok: false}
}

func (r Result[T, E]) Unwrap() (T, E, bool) {
    return r.value, r.err, r.ok
}

func (r Result[T, E]) Map[U any](fn func(T) U) Result[U, E] {
    if !r.ok {
        return Err[U, E](r.err)
    }
    return Ok[U, E](fn(r.value))
}
```

**Usage**:
```go
func Divide(a, b float64) Result[float64, string] {
    if b == 0 {
        return Err[float64, string]("division by zero")
    }
    return Ok[float64, string](a / b)
}

result := Divide(10, 2)
if val, _, ok := result.Unwrap(); ok {
    fmt.Println("Result:", val)  // 5.0
}
```

### Pattern 3: Pair[A, B] (Tuples)

```go
type Pair[A, B any] struct {
    First  A
    Second B
}

func MakePair[A, B any](a A, b B) Pair[A, B] {
    return Pair[A, B]{First: a, Second: b}
}

func (p Pair[A, B]) Swap() Pair[B, A] {
    return Pair[B, A]{First: p.Second, Second: p.First}
}
```

**Usage**:
```go
// Group by key
type User struct {
    ID   int
    Name string
}

users := []User{{1, "Alice"}, {2, "Bob"}}
pairs := Map(users, func(u User) Pair[int, string] {
    return MakePair(u.ID, u.Name)
})
// [{1, "Alice"}, {2, "Bob"}]
```

### Pattern 4: Generic Stack

```go
type Stack[T any] struct {
    items []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}
```

### Pattern 5: Generic Constraint Helpers

```go
// Number constraint (built-in in Go 1.21+)
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64
}

func Max[T Number](a, b T) T {
    if a > b {
        return a
    }
    return b
}

func Sum[T Number](nums []T) T {
    var sum T
    for _, n := range nums {
        sum += n
    }
    return sum
}

// Ordered constraint
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64 | ~string
}

func Min[T Ordered](slice []T) T {
    if len(slice) == 0 {
        var zero T
        return zero
    }
    min := slice[0]
    for _, item := range slice[1:] {
        if item < min {
            min = item
        }
    }
    return min
}
```

---

## 6. Real-World Applications

### Application 1: Data Processing Pipeline

**Scenario**: Process millions of log entries.

```go
type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
}

logs := []LogEntry{ /* millions of entries */ }

// Filter errors, extract messages, count words
errorMessages := logs |>
    Filter(func(log LogEntry) bool { return log.Level == "ERROR" }) |>
    Map(func(log LogEntry) string { return log.Message }) |>
    FlatMap(func(msg string) []string { return strings.Fields(msg) }) |>
    Reduce(make(map[string]int), func(counts map[string]int, word string) map[string]int {
        counts[word]++
        return counts
    })
```

With parallel processing:
```go
// 8x faster on 8 cores
errorMessages := ParallelMap(logs, extractWords, 8)
wordCounts := ParallelReduce(errorMessages, initialMap, countWords, 8)
```

### Application 2: ETL (Extract, Transform, Load)

```go
type RawRecord struct {
    ID   string
    Data string
}

type CleanRecord struct {
    ID    int
    Value float64
}

func ProcessETL(raw []RawRecord) []CleanRecord {
    return raw |>
        // Filter valid records
        Filter(func(r RawRecord) bool {
            return r.ID != "" && r.Data != ""
        }) |>
        // Transform to clean format
        Map(func(r RawRecord) Optional[CleanRecord] {
            id, err1 := strconv.Atoi(r.ID)
            val, err2 := strconv.ParseFloat(r.Data, 64)
            if err1 != nil || err2 != nil {
                return None[CleanRecord]()
            }
            return Some(CleanRecord{ID: id, Value: val})
        }) |>
        // Filter out failed parses
        Filter(func(opt Optional[CleanRecord]) bool {
            _, ok := opt.Get()
            return ok
        }) |>
        // Extract values
        Map(func(opt Optional[CleanRecord]) CleanRecord {
            val, _ := opt.Get()
            return val
        })
}
```

### Application 3: Image Processing

```go
type Pixel struct {
    R, G, B uint8
}

type Image [][]Pixel

func Brighten(img Image, factor float64) Image {
    return ParallelMap(img, func(row []Pixel) []Pixel {
        return Map(row, func(p Pixel) Pixel {
            return Pixel{
                R: uint8(math.Min(float64(p.R)*factor, 255)),
                G: uint8(math.Min(float64(p.G)*factor, 255)),
                B: uint8(math.Min(float64(p.B)*factor, 255)),
            }
        })
    }, runtime.NumCPU())
}
```

### Application 4: Financial Calculations

```go
type Transaction struct {
    Amount float64
    Type   string
}

func CalculateBalance(transactions []Transaction) float64 {
    return ParallelReduce(
        transactions,
        0.0,
        func(balance float64, tx Transaction) float64 {
            if tx.Type == "credit" {
                return balance + tx.Amount
            }
            return balance - tx.Amount
        },
        8,
    )
}
```

---

## 7. Performance Considerations

### When to Use Parallel Processing

**Parallel map is worth it when**:
- Data size is large (>10,000 items typically)
- Operation is CPU-intensive (not just memory access)
- You have multiple cores available

**Benchmark example**:
```go
func BenchmarkMapSequential(b *testing.B) {
    data := generateData(1000000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Map(data, expensiveFunc)
    }
}

func BenchmarkMapParallel(b *testing.B) {
    data := generateData(1000000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ParallelMap(data, expensiveFunc, 8)
    }
}
```

**Results** (8 cores, CPU-heavy operation):
```
BenchmarkMapSequential-8    10    120ms/op
BenchmarkMapParallel-8      80     15ms/op
```
8x speedup!

### Generic Performance

**Key insight**: Generics have **zero runtime overhead** in Go.

Why? **Monomorphization** - compiler generates specialized code for each type:
```go
// You write:
func Sum[T int | float64](nums []T) T { ... }

// Compiler generates (conceptually):
func Sum_int(nums []int) int { ... }
func Sum_float64(nums []float64) float64 { ... }
```

**Benchmark**:
```go
func BenchmarkSumGeneric(b *testing.B) {
    nums := []int{1, 2, 3, 4, 5}
    for i := 0; i < b.N; i++ {
        Sum(nums)  // Generic version
    }
}

func BenchmarkSumSpecific(b *testing.B) {
    nums := []int{1, 2, 3, 4, 5}
    for i := 0; i < b.N; i++ {
        SumInt(nums)  // Hand-written int version
    }
}
```

**Results**: Identical performance!

---

## 8. Common Patterns and Idioms

### Pattern 1: Method Chaining (Fluent API)

```go
type Stream[T any] struct {
    data []T
}

func NewStream[T any](data []T) Stream[T] {
    return Stream[T]{data: data}
}

func (s Stream[T]) Map[U any](fn func(T) U) Stream[U] {
    return Stream[U]{data: Map(s.data, fn)}
}

func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
    return Stream[T]{data: Filter(s.data, predicate)}
}

func (s Stream[T]) Reduce[U any](initial U, fn func(U, T) U) U {
    return Reduce(s.data, initial, fn)
}

// Usage:
result := NewStream([]int{1, 2, 3, 4, 5}).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * 2 }).
    Reduce(0, func(acc, x int) int { return acc + x })
```

### Pattern 2: Generic Wrapper Types

```go
type Box[T any] struct {
    Value T
}

func NewBox[T any](val T) Box[T] {
    return Box[T]{Value: val}
}

func (b Box[T]) Map[U any](fn func(T) U) Box[U] {
    return Box[U]{Value: fn(b.Value)}
}
```

### Pattern 3: Type Inference

```go
// Explicit type parameters (verbose)
result := Map[int, string](numbers, toString)

// Type inference (cleaner)
result := Map(numbers, toString)  // Compiler infers int → string

// Type inference from return type
var result []string = Map(numbers, toString)
```

---

## 9. Common Mistakes to Avoid

### Mistake 1: Using `any` When You Need Constraints

**❌ Wrong**:
```go
func Sum[T any](nums []T) T {
    var sum T
    for _, n := range nums {
        sum += n  // ERROR: + not defined for any
    }
    return sum
}
```

**✅ Correct**:
```go
func Sum[T int | float64](nums []T) T {
    var sum T
    for _, n := range nums {
        sum += n  // OK: both int and float64 support +
    }
    return sum
}
```

### Mistake 2: Forgetting Type Parameters in Methods

**❌ Wrong**:
```go
type Stack[T any] struct {
    items []T
}

func (s *Stack) Push(item T) {  // ERROR: T not in scope
    s.items = append(s.items, item)
}
```

**✅ Correct**:
```go
func (s *Stack[T]) Push(item T) {  // Must include [T]
    s.items = append(s.items, item)
}
```

### Mistake 3: Type Assertions on Generic Types

**❌ Wrong**:
```go
func GetInt[T any](val T) int {
    return val.(int)  // ERROR: T might not be int
}
```

**✅ Correct**:
```go
func GetInt[T int | int64](val T) int {
    return int(val)  // OK: conversion, not assertion
}
```

### Mistake 4: Parallel Reduce Without Associative Operation

**❌ Wrong**:
```go
// Subtraction is NOT associative
// (a - b) - c ≠ a - (b - c)
result := ParallelReduce(nums, 0, func(acc, x int) int {
    return acc - x  // Results will be wrong!
}, 4)
```

**✅ Correct**:
```go
// Addition IS associative
result := ParallelReduce(nums, 0, func(acc, x int) int {
    return acc + x  // Safe for parallel
}, 4)
```

---

## 10. Exercises

The `exercise/` directory contains progressive challenges:

### Exercise 1: Implement Basic Generics ⭐
Implement `Map`, `Filter`, `Reduce` for any type.

### Exercise 2: Generic Data Structures ⭐⭐
Implement `Optional[T]`, `Result[T, E]`, `Pair[A, B]`.

### Exercise 3: Parallel Map-Reduce ⭐⭐⭐
Implement `ParallelMap` and `ParallelReduce` with worker pools.

### Exercise 4: Stream API ⭐⭐⭐
Build a fluent API for chaining operations.

### Exercise 5: Real-World Application ⭐⭐⭐⭐
Build a log analyzer that processes millions of entries in parallel.

---

## How to Run

```bash
# Run the demo
make run P=46-generics-map-reduce

# Run tests
go test ./minis/46-generics-map-reduce/...

# Run benchmarks
go test -bench=. -benchmem ./minis/46-generics-map-reduce/exercise/

# Test with different worker counts
go test -bench=BenchmarkParallel -benchmem \
    -cpuprofile=cpu.prof ./minis/46-generics-map-reduce/exercise/

# View CPU profile
go tool pprof cpu.prof
```

---

## Summary

**What you learned**:
- ✅ Go generics fundamentals (type parameters, constraints)
- ✅ Generic functions (Map, Filter, Reduce, FlatMap)
- ✅ Generic data structures (Optional, Result, Pair, Stack)
- ✅ Type constraints (any, comparable, unions, interfaces)
- ✅ Parallel map-reduce for high performance
- ✅ Real-world functional programming patterns
- ✅ Performance characteristics (monomorphization)

**Why this matters**:
Generics enable writing reusable, type-safe code without sacrificing performance. Map-reduce patterns are fundamental to modern data processing, from ETL pipelines to distributed systems like Hadoop and Spark.

**Key insights**:
- Generics eliminate code duplication while maintaining type safety
- Type constraints make APIs both flexible and safe
- Parallel map-reduce achieves near-linear speedup on multi-core systems
- Go's monomorphization means zero runtime overhead
- Functional patterns (map/filter/reduce) lead to cleaner, more composable code

**Next steps**:
- Project 47: Plugin systems with hot reload
- Project 48: Reflection and introspection
- Advanced: Build a generic database query builder
- Advanced: Implement a distributed map-reduce framework

Happy coding! 🚀
