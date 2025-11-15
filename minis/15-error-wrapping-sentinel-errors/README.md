# Project 15: Error Wrapping and Sentinel Errors

## What Is This Project About?

This project provides a **deep, comprehensive understanding** of error handling in Go—one of the most distinctive and powerful features of the language. You'll learn:

1. **What errors are in Go** (the error interface, first principles)
2. **How to create errors** (errors.New, fmt.Errorf, custom types)
3. **Sentinel errors** (predefined error values for comparison)
4. **Error wrapping** (adding context while preserving the original error)
5. **Error inspection** (errors.Is, errors.As for checking error types)
6. **Custom error types** (when and how to create them)
7. **Error handling philosophy** (when to wrap, when to return, when to handle)

By the end, you'll understand Go's error handling philosophy and be able to design robust error handling for production applications.

---

## The Fundamental Problem: What Is an Error?

### First Principles: Errors as Values

In many languages, errors are **exceptions**—special control flow that "throws" up the call stack:

```python
# Python: Exceptions
try:
    result = risky_operation()
except ValueError as e:
    handle_error(e)
```

Go takes a different approach: **errors are just values**. They're returned like any other return value:

```go
// Go: Errors as values
result, err := riskyOperation()
if err != nil {
    handleError(err)
}
```

**Why this matters:**
- **Explicit**: You can't ignore errors (the compiler forces you to handle or explicitly ignore them)
- **Composable**: Errors are first-class values you can pass around, transform, and inspect
- **Predictable**: No hidden control flow (no try/catch scanning)
- **Local**: Error handling happens exactly where the error occurs

### The Error Interface

In Go, an error is **any type** that implements the `error` interface:

```go
type error interface {
    Error() string
}
```

That's it! Any type with an `Error() string` method is an error.

**Example:**
```go
type MyError struct {
    When time.Time
    What string
}

func (e MyError) Error() string {
    return fmt.Sprintf("error at %v: %s", e.When, e.What)
}

// Now MyError is an error!
var err error = MyError{When: time.Now(), What: "something broke"}
fmt.Println(err.Error())  // "error at 2025-11-15 10:30:00: something broke"
```

**Key insight:** The error interface is **intentionally simple**. This allows maximum flexibility in how errors are created, stored, and inspected.

---

## Creating Errors: Three Main Patterns

### Pattern 1: errors.New (Simple Static Errors)

For **constant error messages**, use `errors.New`:

```go
import "errors"

var ErrNotFound = errors.New("resource not found")
var ErrUnauthorized = errors.New("unauthorized access")

func GetUser(id int) (*User, error) {
    if id < 0 {
        return nil, ErrNotFound  // Return the sentinel error
    }
    // ...
}
```

**When to use:**
- Fixed error messages
- Errors you want to compare by identity (sentinel errors)
- Package-level error definitions

**Under the hood:**
```go
// Simplified implementation
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}

func New(text string) error {
    return &errorString{text}
}
```

### Pattern 2: fmt.Errorf (Dynamic Error Messages)

For **formatted error messages with context**, use `fmt.Errorf`:

```go
import "fmt"

func ReadFile(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        // Add context: which file failed?
        return nil, fmt.Errorf("failed to open %s: %w", path, err)
    }
    // ...
}
```

**When to use:**
- Error messages with dynamic values (file names, user IDs, etc.)
- Adding context to errors from other functions
- Error wrapping (explained below)

**Format verbs:**
- `%v`: Default format
- `%s`: String representation
- `%q`: Quoted string
- `%w`: **Wrap the error** (special, explained next)

### Pattern 3: Custom Error Types

For **errors that need additional data**, create a custom type:

```go
type ValidationError struct {
    Field string
    Value interface{}
    Reason string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed on field %s: %s (value: %v)",
        e.Field, e.Reason, e.Value)
}

func ValidateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field: "age",
            Value: age,
            Reason: "must be non-negative",
        }
    }
    return nil
}
```

**When to use:**
- Errors that carry structured data (not just strings)
- Errors that need special handling in calling code
- Errors that implement additional interfaces

---

## Sentinel Errors: Predefined Error Values

### What Are Sentinel Errors?

**Sentinel errors** are predefined error variables that can be compared by **identity** (pointer equality):

```go
// Standard library examples
var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

// Usage
func Get(key string) (string, error) {
    // ...
    return "", ErrNotFound
}

func main() {
    val, err := Get("foo")
    if err == ErrNotFound {  // Compare by identity
        fmt.Println("Key doesn't exist")
    }
}
```

### When to Use Sentinel Errors

**Good use cases:**
1. **Package-level errors** that callers need to check:
   ```go
   // io package
   var EOF = errors.New("EOF")

   // os package
   var ErrNotExist = errors.New("file does not exist")
   ```

2. **Distinct error conditions** in your API:
   ```go
   var (
       ErrInvalidInput = errors.New("invalid input")
       ErrTimeout = errors.New("operation timed out")
       ErrCanceled = errors.New("operation canceled")
   )
   ```

### Naming Convention

Sentinel errors should be:
- **Package-level variables** (exported or unexported)
- **Named with `Err` prefix**: `ErrNotFound`, `ErrInvalidInput`
- **Created with `errors.New`** (not `fmt.Errorf`)

### The Problem with Sentinel Errors

**Sentinel errors break when wrapped:**

```go
var ErrNotFound = errors.New("not found")

func Get(key string) error {
    // ...
    // BUG: Wrapping breaks identity comparison!
    return fmt.Errorf("get %s: %w", key, ErrNotFound)
}

func main() {
    err := Get("foo")
    if err == ErrNotFound {  // ❌ FALSE! The error is wrapped
        // This code never runs
    }
}
```

**Solution:** Use `errors.Is` (explained next).

---

## Error Wrapping: Adding Context to Errors

### Why Wrap Errors?

When an error happens deep in the call stack, the error message alone often doesn't provide enough context:

```go
// Low-level function
func readConfig(path string) ([]byte, error) {
    return os.ReadFile(path)
}

func main() {
    _, err := readConfig("/etc/app.conf")
    fmt.Println(err)
    // "open /etc/app.conf: no such file or directory"
    // BUT: Was this the config file? A template? A log file?
}
```

**Error wrapping** lets you add context at each layer:

```go
func readConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }
    return data, nil
}

func loadApp() error {
    _, err := readConfig("/etc/app.conf")
    if err != nil {
        return fmt.Errorf("load app: %w", err)
    }
    return nil
}

func main() {
    err := loadApp()
    fmt.Println(err)
    // "load app: read config: open /etc/app.conf: no such file or directory"
    // Clear chain of what happened!
}
```

### The %w Verb: Wrapping Errors

**Key syntax:** Use `%w` in `fmt.Errorf` to wrap an error:

```go
err := someOperation()
if err != nil {
    // %w preserves the original error for inspection
    return fmt.Errorf("operation failed: %w", err)
}
```

**What %w does:**
1. **Includes the error message** (like `%v`)
2. **Preserves the error chain** (for `errors.Is` and `errors.As`)
3. **Allows unwrapping** (to get the original error)

**Contrast with %v:**
```go
// %v: Includes message but breaks the chain
return fmt.Errorf("operation failed: %v", err)  // ❌ Can't unwrap!

// %w: Includes message AND preserves the chain
return fmt.Errorf("operation failed: %w", err)  // ✅ Can unwrap!
```

### Multiple Wrapping (Go 1.20+)

Since Go 1.20, you can wrap **multiple errors** in one:

```go
err1 := errors.New("first error")
err2 := errors.New("second error")

err := fmt.Errorf("both failed: %w and %w", err1, err2)

// Check for either error
errors.Is(err, err1)  // true
errors.Is(err, err2)  // true
```

---

## Error Inspection: errors.Is and errors.As

### errors.Is: Checking Error Identity

**Problem:** Direct comparison (`==`) breaks with wrapped errors.

**Solution:** `errors.Is` checks the entire error chain:

```go
var ErrNotFound = errors.New("not found")

err := fmt.Errorf("get user: %w", ErrNotFound)

// Direct comparison fails
err == ErrNotFound  // ❌ false (err is wrapped)

// errors.Is succeeds
errors.Is(err, ErrNotFound)  // ✅ true (finds it in the chain)
```

**How errors.Is works:**
1. Compares `err` with `target`
2. If not equal, unwraps `err` and tries again
3. Repeats until it finds a match or reaches the end of the chain

**Diagram:**
```
err = "get user: database: not found"
      ↓ (unwrap)
      "database: not found"
      ↓ (unwrap)
      ErrNotFound ← Match!
```

### errors.As: Extracting Error Types

**Problem:** You need to access fields on a custom error type, but it's wrapped.

**Solution:** `errors.As` finds and extracts the first error of a given type:

```go
type ValidationError struct {
    Field string
    Value interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("invalid %s: %v", e.Field, e.Value)
}

// Function that wraps the error
func SaveUser(age int) error {
    err := ValidateAge(age)
    if err != nil {
        return fmt.Errorf("save user: %w", err)
    }
    return nil
}

// Caller extracts the ValidationError
func main() {
    err := SaveUser(-5)

    var ve ValidationError
    if errors.As(err, &ve) {  // ✅ Finds ValidationError in the chain
        fmt.Printf("Validation failed on field: %s\n", ve.Field)
        fmt.Printf("Invalid value: %v\n", ve.Value)
    }
}
```

**How errors.As works:**
1. Checks if `err` can be assigned to `*target`
2. If not, unwraps `err` and tries again
3. If found, assigns the error to `*target` and returns `true`

**Important:** The second parameter is a **pointer** to the error type:
```go
var ve ValidationError
errors.As(err, &ve)  // ✅ Correct (pointer)
errors.As(err, ve)   // ❌ Compile error!
```

---

## Custom Error Types: When and How

### When to Create Custom Error Types

Create a custom error type when:

1. **You need structured data** beyond a string:
   ```go
   type HTTPError struct {
       StatusCode int
       Body string
   }
   ```

2. **Callers need to handle the error specially**:
   ```go
   type TemporaryError interface {
       error
       Temporary() bool
   }

   if te, ok := err.(TemporaryError); ok && te.Temporary() {
       retry()
   }
   ```

3. **You want to implement additional interfaces**:
   ```go
   type FileError struct {
       Op   string
       Path string
       Err  error
   }

   func (e FileError) Unwrap() error { return e.Err }
   ```

### How to Create Custom Error Types

**Step 1: Define the struct**
```go
type QueryError struct {
    Query string
    Err   error  // Wrapped error (optional)
}
```

**Step 2: Implement Error() method**
```go
func (e QueryError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("query failed: %s: %v", e.Query, e.Err)
    }
    return fmt.Sprintf("query failed: %s", e.Query)
}
```

**Step 3 (optional): Implement Unwrap() for wrapping**
```go
func (e QueryError) Unwrap() error {
    return e.Err
}
```

**Usage:**
```go
err := db.Query("SELECT * FROM users")
if err != nil {
    return QueryError{Query: "SELECT * FROM users", Err: err}
}
```

### Custom Error Types: Best Practices

1. **Use value types** (not pointers) for error structs:
   ```go
   return ValidationError{...}  // ✅ Good
   return &ValidationError{...} // ❌ Causes pointer comparisons issues
   ```

2. **Implement Unwrap()** if your error wraps another:
   ```go
   func (e MyError) Unwrap() error { return e.wrappedErr }
   ```

3. **Make fields exported** if callers need to inspect them:
   ```go
   type PathError struct {
       Path string  // Exported: callers can read it
       err  error   // Unexported: internal
   }
   ```

4. **Don't use `Error()` for control flow**:
   ```go
   // ❌ Bad: Parsing error strings
   if strings.Contains(err.Error(), "not found") { ... }

   // ✅ Good: Use errors.Is or errors.As
   if errors.Is(err, ErrNotFound) { ... }
   ```

---

## Error Handling Philosophy

### The Four Options for Handling Errors

When you receive an error, you have **four choices**:

#### 1. Return the Error (Most Common)

Pass the error up the call stack, optionally adding context:

```go
func readConfig(path string) (Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        // Add context and return
        return Config{}, fmt.Errorf("read config %s: %w", path, err)
    }
    // ...
}
```

**When to use:** Almost always! Let callers decide how to handle it.

#### 2. Handle the Error Locally

If you can recover, handle it and don't return:

```go
func loadConfigWithFallback(path string) Config {
    cfg, err := readConfig(path)
    if err != nil {
        // Handle: Use default config
        log.Printf("Config load failed, using defaults: %v", err)
        return defaultConfig()
    }
    return cfg
}
```

**When to use:** You have a reasonable fallback or recovery.

#### 3. Retry on Transient Errors

Some errors are temporary (network glitches, resource contention):

```go
func fetchWithRetry(url string) ([]byte, error) {
    var lastErr error
    for i := 0; i < 3; i++ {
        data, err := fetch(url)
        if err == nil {
            return data, nil
        }

        // Check if error is retryable
        var netErr net.Error
        if errors.As(err, &netErr) && netErr.Temporary() {
            lastErr = err
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }

        // Non-retryable error
        return nil, err
    }
    return nil, fmt.Errorf("failed after 3 retries: %w", lastErr)
}
```

**When to use:** Network operations, database queries, external services.

#### 4. Panic (Rare!)

Only panic if the error represents **programmer error** or **unrecoverable state**:

```go
func mustLoadConfig(path string) Config {
    cfg, err := readConfig(path)
    if err != nil {
        // Config is required for the app to function
        panic(fmt.Sprintf("FATAL: Cannot load config: %v", err))
    }
    return cfg
}
```

**When to use:**
- Initialization code (before `main` starts)
- Critical resources that the app can't run without
- Programming errors (bugs, not runtime errors)

**Guideline:** If the caller can handle it, return an error. If not, panic.

---

## Error Handling Patterns

### Pattern 1: Guard Clauses

Check for errors immediately and return early:

```go
// ✅ Good: Guard clauses
func process(input string) error {
    if input == "" {
        return errors.New("input is empty")
    }

    data, err := parse(input)
    if err != nil {
        return fmt.Errorf("parse: %w", err)
    }

    if err := validate(data); err != nil {
        return fmt.Errorf("validate: %w", err)
    }

    return save(data)
}

// ❌ Bad: Nested if statements
func process(input string) error {
    if input != "" {
        data, err := parse(input)
        if err == nil {
            if err := validate(data); err == nil {
                return save(data)
            } else {
                return err
            }
        } else {
            return err
        }
    }
    return errors.New("input is empty")
}
```

### Pattern 2: Error Variables for Cleanup

Use named error variable to handle cleanup:

```go
func processFile(path string) (err error) {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open: %w", err)
    }
    defer func() {
        closeErr := file.Close()
        if err == nil {
            err = closeErr  // Only set if no other error
        }
    }()

    // Process file...
    return nil
}
```

### Pattern 3: Multi-Error Collection

Collect multiple errors when processing batches:

```go
type MultiError struct {
    Errors []error
}

func (m MultiError) Error() string {
    return fmt.Sprintf("%d errors occurred", len(m.Errors))
}

func processFiles(paths []string) error {
    var multi MultiError

    for _, path := range paths {
        if err := processFile(path); err != nil {
            multi.Errors = append(multi.Errors, err)
        }
    }

    if len(multi.Errors) > 0 {
        return multi
    }
    return nil
}
```

---

## Common Mistakes to Avoid

### Mistake 1: Ignoring Errors

```go
// ❌ Bad: Silently ignoring errors
data, _ := os.ReadFile(path)

// ✅ Good: Handle or explicitly ignore with comment
data, err := os.ReadFile(path)
if err != nil {
    log.Printf("Failed to read %s, using defaults: %v", path, err)
    data = defaultData
}

// ✅ Acceptable: When you truly don't care
_ = file.Close()  // Cleanup, error doesn't matter
```

### Mistake 2: Using %v Instead of %w

```go
// ❌ Bad: Breaks error chain
return fmt.Errorf("failed to process: %v", err)

// ✅ Good: Preserves error chain
return fmt.Errorf("failed to process: %w", err)
```

### Mistake 3: Wrapping Too Much

```go
// ❌ Bad: Every layer wraps
func a() error {
    return fmt.Errorf("a: %w", b())
}
func b() error {
    return fmt.Errorf("b: %w", c())
}
func c() error {
    return fmt.Errorf("c: %w", errors.New("base"))
}
// Result: "a: b: c: base" (redundant!)

// ✅ Good: Only add meaningful context
func a() error {
    return b()  // No context to add
}
func b() error {
    if err := c(); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}
```

### Mistake 4: Comparing Wrapped Errors with ==

```go
// ❌ Bad: Direct comparison fails with wrapped errors
if err == ErrNotFound { ... }

// ✅ Good: Use errors.Is
if errors.Is(err, ErrNotFound) { ... }
```

---

## How to Run

```bash
# Run the demonstration program
cd minis/15-error-wrapping-sentinel-errors
go run cmd/errors-demo/main.go

# Run tests
cd exercise
go test -v

# Run specific tests
go test -run TestSentinelErrors -v
go test -run TestErrorWrapping -v
```

---

## Expected Output (Demo Program)

```
=== Sentinel Errors ===
Error: resource not found
Is ErrNotFound? true

=== Error Wrapping ===
Error: load config: read file: open config.txt: no such file or directory
Is os.ErrNotExist? true

=== Custom Error Types ===
Error: validation failed on field age: must be positive (value: -5)
Field: age
Value: -5

=== errors.Is and errors.As ===
Found ErrNotFound in chain: true
Extracted ValidationError: age

=== Multi-Error Handling ===
Processing file1.txt: success
Processing file2.txt: error
Processing file3.txt: success
Total errors: 1
```

---

## Key Takeaways

1. **Errors are values** (not exceptions)
2. **Use errors.New** for simple sentinel errors
3. **Use fmt.Errorf with %w** to wrap errors and add context
4. **Use errors.Is** to check for specific errors in a chain
5. **Use errors.As** to extract custom error types
6. **Create custom error types** when you need structured data
7. **Wrap errors** to add context, but don't over-wrap
8. **Return errors** (don't panic) unless it's unrecoverable
9. **Handle errors early** with guard clauses
10. **Never ignore errors** without a good reason

---

## Connections to Other Projects

- **Project 13 (interfaces-duck-typing)**: The error interface is a fundamental interface
- **Project 24 (sync-mutex-vs-rwmutex)**: Thread-safe error handling
- **Project 28 (context-cancellation-timeouts)**: Context.Err() returns sentinel errors
- **Project 47 (testing-table-driven-subtests)**: Testing error cases
- **Project 59 (logging-structured-json)**: Logging errors with structured data

---

## Stretch Goals

1. **Implement a retry wrapper** that handles temporary errors automatically
2. **Create a custom error type** with stack traces (inspired by `github.com/pkg/errors`)
3. **Build an error aggregator** that collects errors from concurrent goroutines
4. **Write an error analyzer** that extracts all errors in a chain
5. **Implement error metrics** (count errors by type for monitoring)
