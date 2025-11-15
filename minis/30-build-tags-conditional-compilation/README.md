# Project 30: Build Tags & Conditional Compilation

## What Is This Project About?

Imagine you're building a cross-platform application that needs to:
- Use different system calls on Linux, macOS, and Windows
- Include debug logging only in development builds
- Enable premium features only when a license is activated
- Optimize code differently for ARM vs x86 processors

This project teaches you how to use **build tags** (also called build constraints) to conditionally compile different code based on the operating system, architecture, or custom flags you define.

You'll learn:
1. **Build Tag Syntax**: Both legacy (`// +build`) and modern (`//go:build`) formats
2. **Platform-Specific Code**: Writing OS and architecture-specific implementations
3. **Feature Flags**: Enabling/disabling features at compile time
4. **Conditional Compilation**: How Go's build system decides what code to include

## The Fundamental Problem: One Codebase, Many Environments

### First Principles: Why Do We Need Conditional Compilation?

When you write code, you're creating instructions for a computer. But different computers have different capabilities:

- **Operating Systems**: Linux, macOS, Windows have different APIs for file systems, processes, networking
- **Architectures**: ARM chips (phones, M1 Macs) vs x86/AMD64 (traditional PCs) have different instruction sets
- **Environments**: Development vs production may need different logging, debugging, or monitoring
- **Features**: Free vs paid tiers may unlock different functionality

**The Challenge**: You want ONE codebase that works everywhere, but you need DIFFERENT code to run in different situations.

**The Solution**: Build tags let you mark files or code blocks to be compiled only under certain conditions.

### How Other Languages Handle This

- **C/C++**: Use preprocessor directives (`#ifdef`, `#ifndef`)
- **Python**: Check at runtime (`if platform.system() == 'Linux'`)
- **Java**: Separate modules or reflection-based checks

**Go's Approach**: Compile-time decisions using build tags. Code that doesn't match your build conditions doesn't even get compiled into your binary.

**Benefits**:
- **No runtime overhead**: Unused code isn't in the binary
- **Smaller binaries**: Only include what you need
- **Type safety**: Conditional code is still checked by the compiler
- **Explicit**: You can see exactly what code runs where

## Build Tags Syntax: Two Formats

Go supports two syntaxes for build constraints. The modern `//go:build` format was introduced in Go 1.17 and is now preferred.

### Modern Syntax: `//go:build`

```go
//go:build linux && amd64
```

**Rules**:
- Must be the FIRST line of the file (or after package comments)
- Must have EXACTLY ONE space after `//` (no extra spaces)
- Uses boolean operators: `&&` (AND), `||` (OR), `!` (NOT)
- Parentheses supported for grouping: `(linux || darwin) && !386`

### Legacy Syntax: `// +build`

```go
// +build linux,amd64
```

**Rules**:
- Must be near the top of the file, before the package declaration
- Followed by a blank line before package
- Uses different operators: `,` (AND), ` ` (OR), `!` (NOT)
- Multiple lines are OR'd together

**Example equivalences**:
```go
// Modern: //go:build linux && amd64
// Legacy: // +build linux,amd64

// Modern: //go:build linux || darwin
// Legacy: // +build linux darwin

// Modern: //go:build !windows
// Legacy: // +build !windows
```

### Using Both (For Compatibility)

For Go 1.16 and earlier compatibility, you can use both:

```go
//go:build linux && amd64
// +build linux,amd64

package main
```

The `go fmt` tool can automatically add the legacy format when you write the modern one.

## Understanding Build Tag Conditions

### Predefined Build Tags

Go automatically defines these tags based on your build environment:

**Operating Systems**:
- `linux`, `darwin` (macOS), `windows`, `freebsd`, `openbsd`, `netbsd`, `dragonfly`, `solaris`, `android`, `ios`

**Architectures**:
- `amd64` (x86-64), `386` (x86-32), `arm`, `arm64`, `ppc64`, `ppc64le`, `mips`, `mipsle`, `mips64`, `mips64le`, `s390x`, `wasm`

**Compiler**:
- `gc` (standard Go compiler), `gccgo` (GCC-based Go compiler)

**CGO**:
- `cgo` (enabled when using C interop)

**Go Version** (Go 1.21+):
- `go1.21`, `go1.22`, etc. (requires that version or newer)

**Runtime**:
- `unix` (Linux, macOS, BSD, etc.)

### Custom Build Tags

You can define your own tags for features:

```go
//go:build premium

package features

func EnableAdvancedAnalytics() {
    // Only compiled when building with -tags=premium
}
```

Build with:
```bash
go build -tags=premium
go build -tags="premium,debug"  # Multiple tags
```

## Pattern 1: Platform-Specific Implementations

The most common use case: different implementations for different operating systems.

### Problem: Getting the Username

Each OS stores the current user differently:
- **Unix/Linux**: Environment variable `$USER` or system calls
- **Windows**: Different environment variable `%USERNAME%` or Windows API

### Solution: Separate Files with Build Tags

**File: `user_unix.go`**
```go
//go:build unix

package main

import "os"

func GetUsername() string {
    return os.Getenv("USER")
}
```

**File: `user_windows.go`**
```go
//go:build windows

package main

import "os"

func GetUsername() string {
    return os.Getenv("USERNAME")
}
```

**File: `main.go`** (no build tags)
```go
package main

import "fmt"

func main() {
    // This calls the right GetUsername based on OS
    fmt.Println("Current user:", GetUsername())
}
```

When you run `go build`:
- On Linux/macOS: Only `user_unix.go` is compiled
- On Windows: Only `user_windows.go` is compiled
- Both define the same function signature, so `main.go` compiles on all platforms

### The File Naming Convention

Go also supports a convention where the OS/architecture is in the filename:

**Pattern**: `filename_GOOS_GOARCH.go`

Examples:
- `network_linux.go` - Only on Linux
- `network_darwin.go` - Only on macOS
- `network_windows.go` - Only on Windows
- `math_amd64.go` - Only on amd64 architecture
- `math_arm64.go` - Only on arm64 architecture
- `crypto_linux_amd64.go` - Only on Linux + amd64

**When to use naming vs explicit tags**:
- **Naming**: Simple single-OS or single-arch conditions
- **Explicit tags**: Complex conditions with AND/OR/NOT

## Pattern 2: Feature Flags

Enable or disable features at compile time using custom tags.

### Problem: Debug Logging

You want detailed debug logs during development but not in production.

### Solution: Debug Tag

**File: `logging.go`** (always compiled)
```go
package main

// Production logging - always available
func LogInfo(msg string) {
    println("[INFO]", msg)
}

func LogError(msg string) {
    println("[ERROR]", msg)
}
```

**File: `logging_debug.go`**
```go
//go:build debug

package main

import "fmt"

// Debug logging - only when built with -tags=debug
func LogDebug(msg string) {
    fmt.Printf("[DEBUG] %s\n", msg)
}

func EnableVerboseLogging() {
    println("Verbose logging enabled")
}
```

**File: `logging_nodebug.go`**
```go
//go:build !debug

package main

// Stub implementations when debug is disabled
func LogDebug(msg string) {
    // No-op in production
}

func EnableVerboseLogging() {
    // No-op in production
}
```

**Usage**:
```bash
# Production build - no debug code
go build -o myapp

# Debug build - includes debug code
go build -tags=debug -o myapp-debug
```

### Pattern 3: Architecture-Specific Optimizations

Different CPU architectures have different capabilities.

### Problem: Fast Math Operations

Modern CPUs have SIMD instructions (Single Instruction, Multiple Data) for parallel processing, but they're architecture-specific.

**File: `math_amd64.go`**
```go
//go:build amd64

package compute

// Uses AVX2 instructions available on modern x86-64 CPUs
func VectorSum(data []float64) float64 {
    // Optimized implementation using assembly or intrinsics
    // (simplified example)
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    return sum
}
```

**File: `math_arm64.go`**
```go
//go:build arm64

package compute

// Uses NEON instructions available on ARM64
func VectorSum(data []float64) float64 {
    // Optimized for ARM architecture
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    return sum
}
```

**File: `math_generic.go`**
```go
//go:build !amd64 && !arm64

package compute

// Fallback for other architectures
func VectorSum(data []float64) float64 {
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    return sum
}
```

## Pattern 4: Development vs Production

### Problem: Different Configurations

Development needs easy debugging, production needs performance.

**File: `config_dev.go`**
```go
//go:build dev

package config

const (
    APIEndpoint = "http://localhost:8080"
    LogLevel    = "DEBUG"
    CacheSize   = 10 // Small cache for testing
)
```

**File: `config_prod.go`**
```go
//go:build !dev

package config

const (
    APIEndpoint = "https://api.production.com"
    LogLevel    = "INFO"
    CacheSize   = 10000 // Large cache for performance
)
```

Build:
```bash
go build -tags=dev          # Development build
go build                    # Production build (default)
```

## Complex Build Constraints

### Boolean Logic

You can combine conditions with boolean operators:

```go
// Linux OR macOS
//go:build linux || darwin

// Linux AND 64-bit
//go:build linux && amd64

// NOT Windows
//go:build !windows

// (Linux OR macOS) AND 64-bit AND NOT CGO
//go:build (linux || darwin) && amd64 && !cgo

// Debug mode OR test coverage
//go:build debug || testcover
```

### Version Constraints

Require a minimum Go version:

```go
//go:build go1.21

package main

// This code uses features from Go 1.21+
import "slices"

func Example() {
    nums := []int{3, 1, 4, 1, 5}
    slices.Sort(nums) // Available in Go 1.21+
}
```

## Best Practices

### 1. Use Modern `//go:build` Syntax

The `//go:build` format is clearer and required for new code.

```go
// ✓ Good
//go:build linux && (amd64 || arm64)

// ✗ Avoid (legacy)
// +build linux
// +build amd64 arm64
```

### 2. Keep Platform-Specific Code Minimal

Abstract away platform differences with interfaces:

```go
// platform.go
type FileWatcher interface {
    Watch(path string) error
}

func NewFileWatcher() FileWatcher {
    return newPlatformWatcher()
}

// platform_linux.go
//go:build linux
func newPlatformWatcher() FileWatcher {
    return &linuxWatcher{}
}

// platform_windows.go
//go:build windows
func newPlatformWatcher() FileWatcher {
    return &windowsWatcher{}
}
```

### 3. Provide Fallback Implementations

Always have a fallback for unsupported platforms:

```go
//go:build !linux && !darwin && !windows

package main

func GetUsername() string {
    return "unknown"
}
```

### 4. Document Your Build Tags

Make it clear what tags are available:

```go
// Package features provides optional functionality.
//
// Build tags:
//   - premium: Enables premium features
//   - debug: Enables debug logging
//   - metrics: Enables Prometheus metrics
package features
```

### 5. Test All Combinations

Test your code with different build tags:

```bash
# Test default build
go test

# Test with debug enabled
go test -tags=debug

# Test with multiple tags
go test -tags="debug,premium"

# Test on different platforms (requires access or CI)
GOOS=windows go test
GOOS=darwin go test
GOOS=linux go test
```

### 6. Use `go list` to See What Gets Compiled

```bash
# List all Go files
go list -f '{{.GoFiles}}'

# List files for specific build tags
go list -tags=debug -f '{{.GoFiles}}'

# List files for different OS
GOOS=windows go list -f '{{.GoFiles}}'
```

## Common Pitfalls

### 1. Forgetting the Blank Line (Legacy Syntax)

```go
// ✗ Wrong
// +build linux
package main

// ✓ Correct
// +build linux

package main
```

### 2. Wrong Position for `//go:build`

```go
// ✗ Wrong - after package declaration
package main

//go:build linux

// ✓ Correct - before package declaration
//go:build linux

package main
```

### 3. Mismatched Function Signatures

All platform-specific implementations must have identical signatures:

```go
// ✗ Wrong - different signatures
// file_linux.go
func ReadFile() string { ... }

// file_windows.go
func ReadFile() (string, error) { ... } // Different!

// ✓ Correct - same signatures
// file_linux.go
func ReadFile() (string, error) { ... }

// file_windows.go
func ReadFile() (string, error) { ... }
```

### 4. Build Tag Typos

Build tags are NOT checked at compile time (except predefined ones):

```go
//go:build premiun  // Typo! This will never match

// Build with:
go build -tags=premium  // Won't include this file!
```

### 5. Overlapping Conditions

```go
// file1.go
//go:build linux

// file2.go
//go:build unix  // This includes linux!

// Both files will be compiled on Linux - might cause conflicts
```

## Real-World Examples

### 1. Standard Library: `os` Package

The `os` package has tons of platform-specific code:
- `os/file_unix.go` - Unix file operations
- `os/file_windows.go` - Windows file operations
- `os/exec_unix.go` - Unix process execution
- `os/exec_windows.go` - Windows process execution

### 2. Cross-Platform System Information

```go
//go:build linux
func GetCPUCount() int {
    // Read from /proc/cpuinfo
}

//go:build darwin
func GetCPUCount() int {
    // Use sysctl command
}

//go:build windows
func GetCPUCount() int {
    // Use WMI or environment variable
}
```

### 3. Embedded Devices

```go
//go:build !linux || (linux && !arm && !arm64)
// Build for everything except Linux ARM

//go:build tinygo
// Special code for TinyGo compiler (embedded devices)
```

## Build Tags vs Runtime Checks

### When to Use Build Tags

```go
//go:build windows

// At compile time, only Windows code is included
func OpenFile() { /* Windows-specific */ }
```

**Pros**:
- No runtime overhead
- Smaller binaries
- Compile-time safety

**Use when**: Platform differences are fundamental (system calls, APIs)

### When to Use Runtime Checks

```go
func OpenFile() {
    if runtime.GOOS == "windows" {
        // Windows path
    } else {
        // Unix path
    }
}
```

**Pros**:
- Single binary for all platforms
- Dynamic behavior

**Use when**: Small differences that don't warrant separate files

## How to Run This Project

```bash
# Run default build (uses your current OS)
cd minis/30-build-tags-conditional-compilation
go run cmd/demo/main.go

# Build for different platforms
GOOS=linux go build -o demo-linux cmd/demo/main.go
GOOS=windows go build -o demo-windows.exe cmd/demo/main.go
GOOS=darwin go build -o demo-macos cmd/demo/main.go

# Build with custom tags
go build -tags=debug cmd/demo/main.go
go build -tags="debug,premium" cmd/demo/main.go

# See what files would be compiled
go list -f '{{.GoFiles}}'
go list -tags=debug -f '{{.GoFiles}}'

# Run tests
cd exercise
go test
go test -tags=debug

# Cross-compile for different architectures
GOOS=linux GOARCH=arm64 go build cmd/demo/main.go
GOOS=windows GOARCH=amd64 go build cmd/demo/main.go
```

## Exercises

The `exercise/` directory contains hands-on practice:

1. **Exercise 1**: Implement platform-specific system information functions
2. **Exercise 2**: Create feature flags for different product tiers
3. **Exercise 3**: Write debug vs production logging implementations
4. **Exercise 4**: Handle architecture-specific optimizations

```bash
cd exercise
# Implement functions in exercise.go
go test                    # Test default build
go test -tags=debug       # Test with debug tag
go test -tags=premium     # Test with premium tag
```

## Key Takeaways

1. **Build tags control compilation**: Code with non-matching tags isn't compiled
2. **Two syntaxes exist**: Modern `//go:build` is preferred over legacy `// +build`
3. **Platform-specific code**: Use file naming or explicit tags for OS/arch-specific implementations
4. **Feature flags**: Custom tags enable/disable features at compile time
5. **Boolean logic**: Combine conditions with `&&`, `||`, `!`
6. **Keep it simple**: Minimize platform-specific code, use interfaces to abstract
7. **Test everything**: Verify all tag combinations work correctly

## Further Reading

- [Official Go Build Constraints Documentation](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Go 1.17 `//go:build` Proposal](https://go.dev/design/draft-gobuild)
- [Cross-compilation Guide](https://pkg.go.dev/cmd/go#hdr-Environment_variables)
- [Standard Library Examples](https://cs.opensource.google/go/go/+/refs/tags/go1.21.0:src/os/)

## Stretch Goals

1. **CI/CD Integration**: Set up GitHub Actions to build for multiple platforms
2. **Build Tag Linter**: Write a tool to verify all tag combinations compile
3. **Dynamic Plugin Loading**: Use tags to conditionally compile plugin systems
4. **Embedded Assets**: Use tags to include/exclude large assets from binaries
5. **Performance Comparison**: Benchmark build tag approach vs runtime checks
