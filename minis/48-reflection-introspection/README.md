# Project 48: Reflection & Introspection

## What Is This Project About?

This project teaches you **reflection**—Go's ability to inspect and manipulate types and values at runtime. You'll learn:

1. **What reflection is** (examining types and values at runtime)
2. **The reflect package** (Type and Value, the core abstractions)
3. **Type inspection** (Kind, Name, Fields, Methods)
4. **Struct tags** (metadata for serialization, validation, ORMs)
5. **Dynamic method calls** (calling methods by name at runtime)
6. **Value modification** (changing values through reflection)
7. **Unsafe pitfalls** (when reflection causes panics, performance issues)
8. **When NOT to use reflection** (alternatives and best practices)

By the end, you'll understand how JSON marshaling works, how ORMs inspect structs, how testing frameworks discover tests, and when reflection is worth the complexity.

---

## The Fundamental Problem: Static Types vs Runtime Flexibility

### First Principles: The Type System Trade-Off

Go is **statically typed**: the compiler knows the exact type of every variable at compile time. This gives you:

- **Safety**: Type errors caught before your code runs
- **Performance**: No runtime type checks needed (in most cases)
- **Clarity**: Reading code, you know exactly what types you're working with

**But** this creates a problem: What if you need to work with types you don't know at compile time?

### Real-World Scenarios

**Scenario 1: JSON Unmarshaling**

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// How does json.Unmarshal know which fields to set?
var user User
json.Unmarshal(data, &user)
```

The `json.Unmarshal` function doesn't know about your `User` type—it's in a different package! Yet it can:
1. Discover the struct has `Name` and `Email` fields
2. Read the `json:"..."` tags
3. Set those fields from JSON data

**How?** Reflection.

**Scenario 2: ORM Database Mapping**

```go
type Product struct {
    ID    int    `db:"id"`
    Name  string `db:"product_name"`
    Price float64 `db:"price"`
}

// How does the ORM know how to build this SQL query?
db.Query("SELECT * FROM products").Scan(&product)
```

The ORM uses reflection to:
1. Inspect the struct fields
2. Read the `db:"..."` tags
3. Map database columns to struct fields

**Scenario 3: Test Discovery**

```go
// How does `go test` find all Test* functions?
func TestUserCreation(t *testing.T) { ... }
func TestUserValidation(t *testing.T) { ... }
```

The testing framework uses reflection to:
1. Find all functions starting with "Test"
2. Check if they have the right signature `func(*testing.T)`
3. Call them dynamically

**The pattern:** Reflection lets library code work with YOUR types without knowing about them at compile time.

---

## What Is Reflection? (The Core Concept)

**Reflection** is the ability of a program to examine its own structure at runtime.

### The Three Pillars of Reflection

1. **Type Inspection**: "What type is this value?"
2. **Value Inspection**: "What's inside this value?"
3. **Value Modification**: "Can I change this value?"

### Go's Reflection Model

Go's reflection is built on two key types in the `reflect` package:

```go
type Type interface { ... }   // Represents a Go type
type Value struct { ... }     // Represents a Go value
```

**Analogy**: Think of reflection like looking at yourself in a mirror:
- `reflect.Type` tells you about your structure (height, weight, eye color)
- `reflect.Value` tells you about your current state (what you're wearing, your expression)

### Getting Type and Value

```go
import "reflect"

x := 42
t := reflect.TypeOf(x)   // Type: int
v := reflect.ValueOf(x)  // Value: 42

fmt.Println(t)           // "int"
fmt.Println(v)           // "42"
fmt.Println(v.Kind())    // "int"
```

**Key insight**: `TypeOf` and `ValueOf` are your entry points into reflection. Everything else builds on these.

---

## Type Inspection: Understanding Types

### Type vs Kind

This is the most confusing distinction in reflection.

**Type** = The exact type as written in code
**Kind** = The underlying category of type

```go
type MyInt int

x := MyInt(42)

t := reflect.TypeOf(x)
fmt.Println(t)       // "main.MyInt" (the Type)
fmt.Println(t.Kind()) // "int" (the Kind)
```

**Why the distinction?**

- **Type** is specific: `MyInt`, `time.Duration`, `os.File`
- **Kind** is general: `int`, `struct`, `ptr`, `slice`, `map`, etc.

**The kinds** (from `reflect` package):

```
Invalid, Bool, Int, Int8, Int16, Int32, Int64,
Uint, Uint8, Uint16, Uint32, Uint64, Uintptr,
Float32, Float64, Complex64, Complex128,
Array, Chan, Func, Interface, Map, Pointer, Slice, String, Struct, UnsafePointer
```

### Inspecting Structs

Reflection shines when working with structs:

```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alice", Age: 30}
t := reflect.TypeOf(p)

fmt.Println(t.Kind())        // "struct"
fmt.Println(t.NumField())    // 2

// Iterate over fields
for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    fmt.Printf("Field %d: %s (%s)\n", i, field.Name, field.Type)
}

// Output:
// Field 0: Name (string)
// Field 1: Age (int)
```

**Field metadata** you can access:
- `Name`: Field name ("Name", "Age")
- `Type`: Field type (string, int)
- `Tag`: Struct tag (for metadata)
- `Offset`: Memory offset (for unsafe operations)
- `Anonymous`: Is it an embedded field?

### Inspecting Methods

```go
type Counter struct {
    count int
}

func (c *Counter) Increment() { c.count++ }
func (c *Counter) Get() int { return c.count }

c := &Counter{}
t := reflect.TypeOf(c)

fmt.Println(t.NumMethod())   // 2

for i := 0; i < t.NumMethod(); i++ {
    method := t.Method(i)
    fmt.Println(method.Name)
}

// Output: Get, Increment (alphabetical order)
```

---

## Struct Tags: Metadata for Your Structs

### What Are Struct Tags?

Struct tags are **string literals** attached to struct fields, containing metadata for other code to read.

```go
type User struct {
    Name  string `json:"name" db:"user_name" validate:"required"`
    Email string `json:"email" db:"email" validate:"email"`
    Age   int    `json:"age,omitempty" db:"age" validate:"min=0,max=120"`
}
```

**Format**: `` `key:"value" key2:"value2"` ``

**Common uses**:
- `json:"..."`: JSON marshaling/unmarshaling
- `xml:"..."`: XML marshaling
- `db:"..."`: Database column mapping
- `validate:"..."`: Validation rules
- `yaml:"..."`: YAML marshaling

### Reading Struct Tags

```go
type User struct {
    Name string `json:"name" db:"user_name"`
}

u := User{}
t := reflect.TypeOf(u)
field, _ := t.FieldByName("Name")

jsonTag := field.Tag.Get("json")  // "name"
dbTag := field.Tag.Get("db")      // "user_name"
missing := field.Tag.Get("xml")   // "" (empty string if not present)
```

### Tag Format Rules

Tags use the convention `key:"value"` with spaces between pairs:

```go
// GOOD
`json:"name" db:"user_name"`

// BAD (no space)
`json:"name"db:"user_name"`

// GOOD (multiple options)
`json:"name,omitempty"`

// GOOD (multiple tags)
`json:"name" validate:"required,min=1"`
```

**Parsing complex tags**:

```go
tag := field.Tag.Get("json")  // "name,omitempty"

// For comma-separated options, you need to parse manually
parts := strings.Split(tag, ",")
name := parts[0]              // "name"
omitempty := len(parts) > 1 && parts[1] == "omitempty"
```

### How JSON Uses Tags

Here's simplified pseudocode for how `json.Marshal` works:

```go
func Marshal(v interface{}) ([]byte, error) {
    val := reflect.ValueOf(v)
    typ := reflect.TypeOf(v)

    if typ.Kind() != reflect.Struct {
        // handle other types...
    }

    result := make(map[string]interface{})

    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        fieldValue := val.Field(i)

        // Read the json tag
        jsonName := field.Tag.Get("json")
        if jsonName == "" {
            jsonName = field.Name  // Use field name if no tag
        }

        // Store in result map
        result[jsonName] = fieldValue.Interface()
    }

    return encodeJSON(result), nil
}
```

---

## Value Inspection and Modification

### Getting Values

```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Bob", Age: 25}
v := reflect.ValueOf(p)

nameField := v.FieldByName("Name")
fmt.Println(nameField.String())  // "Bob"

ageField := v.FieldByName("Age")
fmt.Println(ageField.Int())      // 25
```

**Type-specific methods**:
- `v.Bool()` for booleans
- `v.Int()` for int types
- `v.Uint()` for uint types
- `v.Float()` for float types
- `v.String()` for strings
- `v.Interface()` for getting as `interface{}`

### Setting Values (The Tricky Part)

**The Rule**: You can only set a value if it's **addressable** and **exported**.

```go
// WRONG: Passing by value
p := Person{Name: "Alice"}
v := reflect.ValueOf(p)
nameField := v.FieldByName("Name")
nameField.SetString("Bob")  // PANIC: reflect: reflect.Value.SetString using unaddressable value

// RIGHT: Passing a pointer
p := Person{Name: "Alice"}
v := reflect.ValueOf(&p).Elem()  // Get pointer, then dereference
nameField := v.FieldByName("Name")
nameField.SetString("Bob")  // ✅ Works!
fmt.Println(p.Name)  // "Bob"
```

**Why the pointer?**

When you pass by value, reflection gets a **copy**. Modifying the copy wouldn't affect the original.
When you pass a pointer, reflection can follow it to the actual data.

**The `.Elem()` method**: For pointers, `.Elem()` gives you the value being pointed to.

```go
x := 42
v := reflect.ValueOf(&x)  // v is reflect.Value of *int
fmt.Println(v.Kind())     // "ptr"

elem := v.Elem()          // elem is reflect.Value of int
fmt.Println(elem.Kind())  // "int"
fmt.Println(elem.Int())   // 42

elem.SetInt(100)          // Modify through reflection
fmt.Println(x)            // 100
```

### Checking If You Can Set

```go
v := reflect.ValueOf(x)

if v.CanSet() {
    v.SetInt(100)  // Safe to set
} else {
    // Cannot set (not addressable or unexported)
}
```

### Unexported Fields

```go
type Person struct {
    Name string
    age  int    // Unexported!
}

p := Person{Name: "Alice", age: 30}
v := reflect.ValueOf(&p).Elem()

nameField := v.FieldByName("Name")
nameField.SetString("Bob")  // ✅ Works (exported)

ageField := v.FieldByName("age")
ageField.SetInt(40)  // ❌ PANIC (unexported)
```

**Why?** Go's visibility rules apply to reflection too. You can't access unexported fields from outside the package, even with reflection (unless you use `unsafe`, which is... unsafe).

---

## Dynamic Method Calls

### Calling Methods by Name

```go
type Calculator struct{}

func (c Calculator) Add(a, b int) int {
    return a + b
}

calc := Calculator{}
v := reflect.ValueOf(calc)

// Get the method by name
method := v.MethodByName("Add")

// Prepare arguments (must be []reflect.Value)
args := []reflect.Value{
    reflect.ValueOf(5),
    reflect.ValueOf(3),
}

// Call the method
results := method.Call(args)

// Extract result (also []reflect.Value)
sum := results[0].Int()  // 8
```

**Use cases**:
- RPC frameworks (call methods by name from network requests)
- Plugin systems (call methods on dynamically loaded code)
- Test frameworks (discover and run test functions)

### Checking Method Signatures

```go
m := v.MethodByName("Add")
if !m.IsValid() {
    // Method doesn't exist
}

t := m.Type()
fmt.Println(t.NumIn())   // Number of input parameters
fmt.Println(t.NumOut())  // Number of return values

// Check parameter types
if t.NumIn() == 2 && t.In(0).Kind() == reflect.Int {
    // First parameter is an int
}
```

---

## Performance Implications: The Cost of Flexibility

### Reflection is Slow

**Direct call**:
```go
x := 42
y := x + 10  // ~1 nanosecond
```

**Reflection**:
```go
v := reflect.ValueOf(42)
result := v.Int() + 10  // ~100-200 nanoseconds
```

**Why the difference?**

1. **Type checking at runtime**: Reflection must verify types dynamically
2. **Interface conversions**: Values are wrapped in interfaces
3. **No inlining**: Compiler can't optimize reflection code
4. **Heap allocations**: Reflection often forces allocations

### Benchmark Comparison

```
BenchmarkDirectCall-8        1000000000    0.5 ns/op
BenchmarkReflectionCall-8      10000000  150.0 ns/op
```

**~300x slower** for simple operations.

### When Reflection is Acceptable

Despite the overhead, reflection is used extensively because:

1. **Amortized cost**: JSON marshaling is I/O bound (network/disk), not CPU bound
2. **Startup-time only**: Many frameworks use reflection once at startup, then cache results
3. **Rare operations**: If you're only doing it once per request, 150ns doesn't matter
4. **No alternative**: Some problems can't be solved without reflection

### Optimization Strategies

**1. Cache reflection results**

```go
// SLOW: Reflect every time
func GetField(obj interface{}, name string) interface{} {
    v := reflect.ValueOf(obj)
    return v.FieldByName(name).Interface()
}

// FAST: Cache the field index
var nameFieldIndex int
func init() {
    t := reflect.TypeOf(User{})
    for i := 0; i < t.NumField(); i++ {
        if t.Field(i).Name == "Name" {
            nameFieldIndex = i
            break
        }
    }
}

func GetName(obj User) string {
    v := reflect.ValueOf(obj)
    return v.Field(nameFieldIndex).String()
}
```

**2. Use code generation instead**

Many projects use tools to generate code at build time:
- `stringer`: Generate `String()` methods for enums
- `mockgen`: Generate mock implementations
- `protoc-gen-go`: Generate Go code from Protobuf definitions

**Generated code is fast** (no reflection at runtime) but less flexible.

---

## Unsafe Pitfalls: When Reflection Panics

### Panic 1: Setting Unaddressable Values

```go
x := 42
v := reflect.ValueOf(x)  // Value, not pointer!
v.SetInt(100)            // PANIC: reflect.Value.SetInt using unaddressable value
```

**Fix**: Pass a pointer and use `.Elem()`

```go
x := 42
v := reflect.ValueOf(&x).Elem()
v.SetInt(100)  // ✅ Works
```

### Panic 2: Wrong Type Assertion

```go
v := reflect.ValueOf("hello")
x := v.Int()  // PANIC: reflect: call of reflect.Value.Int on string Value
```

**Fix**: Check the kind first

```go
if v.Kind() == reflect.Int {
    x := v.Int()
}
```

### Panic 3: Calling Methods with Wrong Arguments

```go
// Method signature: Add(int, int) int
method := v.MethodByName("Add")
args := []reflect.Value{reflect.ValueOf("wrong")}  // Wrong type!
method.Call(args)  // PANIC: reflect: Call using string as type int
```

**Fix**: Validate argument types

```go
mt := method.Type()
if mt.NumIn() != len(args) {
    return errors.New("wrong number of arguments")
}

for i := 0; i < mt.NumIn(); i++ {
    if mt.In(i) != args[i].Type() {
        return fmt.Errorf("argument %d has wrong type", i)
    }
}

method.Call(args)  // ✅ Safe
```

### Panic 4: Nil Pointer Dereference

```go
var p *Person = nil
v := reflect.ValueOf(p).Elem()  // PANIC: reflect: call of reflect.Value.Elem on zero Value
```

**Fix**: Check for nil first

```go
v := reflect.ValueOf(p)
if v.IsNil() {
    // Handle nil case
}
```

---

## When NOT to Use Reflection

### Antipattern 1: Configuration/Flags

```go
// WRONG: Reflection for simple config
type Config struct {
    Port int
    Host string
}

func SetConfig(cfg *Config, key string, value interface{}) {
    v := reflect.ValueOf(cfg).Elem()
    field := v.FieldByName(key)
    // ... complex reflection code ...
}

// RIGHT: Simple, direct code
type Config struct {
    Port int
    Host string
}

func (c *Config) Set(key, value string) error {
    switch key {
    case "Port":
        c.Port, _ = strconv.Atoi(value)
    case "Host":
        c.Host = value
    default:
        return fmt.Errorf("unknown key: %s", key)
    }
    return nil
}
```

**Why direct is better**: Clearer, faster, type-safe, easier to debug.

### Antipattern 2: Type Switches Instead of Reflection

```go
// WRONG: Reflection for type handling
func Process(val interface{}) {
    v := reflect.ValueOf(val)
    switch v.Kind() {
    case reflect.Int:
        fmt.Println("Integer:", v.Int())
    case reflect.String:
        fmt.Println("String:", v.String())
    }
}

// RIGHT: Type switch
func Process(val interface{}) {
    switch v := val.(type) {
    case int:
        fmt.Println("Integer:", v)
    case string:
        fmt.Println("String:", v)
    }
}
```

**Why type switch is better**: ~100x faster, more idiomatic, type-safe.

### Antipattern 3: Generics Instead of Reflection (Go 1.18+)

```go
// WRONG (pre-Go 1.18): Reflection for generic functions
func Max(a, b interface{}) interface{} {
    va := reflect.ValueOf(a)
    vb := reflect.ValueOf(b)
    if va.Int() > vb.Int() {
        return a
    }
    return b
}

// RIGHT (Go 1.18+): Use generics
func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

**Why generics are better**: Type-safe, faster, clearer.

### When Reflection IS Appropriate

✅ **Good use cases**:
- Serialization/deserialization (JSON, XML, protobuf)
- ORM database mapping
- Dependency injection frameworks
- Testing frameworks (test discovery)
- RPC frameworks (calling methods by name)
- CLI flag parsing (mapping strings to struct fields)
- Validation frameworks (checking struct constraints)

**Common pattern**: The library doesn't know your types, but needs to work with them.

---

## How to Run

```bash
# Run the demonstration program
cd minis/48-reflection-introspection
go run cmd/reflect-demo/main.go

# Run tests
cd exercise
go test -v

# Run benchmarks to see reflection overhead
go test -bench=. -benchmem

# See what the compiler knows about types
go build -gcflags='-m' cmd/reflect-demo/main.go
```

---

## Expected Output (Demo Program)

```
=== Type Inspection ===
Type: main.Person
Kind: struct
Number of fields: 2
Field 0: Name (string)
Field 1: Age (int)

=== Struct Tags ===
Name field: json:"name" db:"user_name" validate:"required"
JSON tag: name
DB tag: user_name
Validate tag: required

=== Value Inspection ===
Name: Alice
Age: 30

=== Value Modification ===
Before: {Alice 30}
After: {Bob 25}

=== Dynamic Method Calls ===
Methods: Get, Increment
Calling Increment...
Calling Get...
Result: 5

=== Performance Warning ===
Direct field access: 1000000000 iterations in 0.5s
Reflection field access: 10000000 iterations in 1.5s
Reflection is ~150x slower
```

---

## Key Takeaways

1. **Reflection lets you inspect types and values at runtime** (examining structure dynamically)
2. **`reflect.Type` describes the type, `reflect.Value` holds the value** (two sides of the same coin)
3. **Type vs Kind**: Type is specific (`MyInt`), Kind is general (`int`)
4. **Struct tags are metadata** (read with `field.Tag.Get("key")`)
5. **Setting values requires addressability** (pass pointers, use `.Elem()`)
6. **Reflection is 100-300x slower** (use only when necessary)
7. **Prefer alternatives**: Type switches, generics, code generation
8. **Reflection shines for libraries** (JSON, ORMs, testing frameworks)
9. **Always validate before calling** (check types, nil, addressability)
10. **Cache reflection results** (don't reflect in hot loops)

---

## Connections to Other Projects

- **Project 13 (interfaces-duck-typing)**: Reflection uses interfaces extensively
- **Project 04 (jsonl-log-filter)**: JSON marshaling uses reflection under the hood
- **Project 38 (config-loader-env-yaml)**: Config loaders often use reflection
- **Project 46 (generics-map-reduce)**: Generics often replace reflection
- **Project 15 (error-wrapping-sentinel-errors)**: Reflection can inspect error types

---

## Stretch Goals

1. **Build a JSON marshaler** from scratch using reflection
2. **Create a struct validator** that reads `validate:"..."` tags
3. **Write a simple ORM** that maps structs to SQL
4. **Benchmark reflection vs generics** for common operations
5. **Build a dependency injection container** using reflection
6. **Create a command dispatcher** that calls methods by name
7. **Implement a deep equality checker** (like `reflect.DeepEqual`)

---

## Further Reading

- [The Laws of Reflection (Go Blog)](https://go.dev/blog/laws-of-reflection)
- [reflect package documentation](https://pkg.go.dev/reflect)
- [How to use struct tags in Go](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go)
- [Reflection is never clear](https://rakyll.org/typesystem/) by Jaana Dogan
- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/gopherchina-2019.html#reflection)
