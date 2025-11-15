// Package main demonstrates error handling patterns in Go.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program provides comprehensive demonstrations of:
// 1. The error interface and how errors work in Go
// 2. Creating errors with errors.New and fmt.Errorf
// 3. Sentinel errors (predefined error values)
// 4. Error wrapping with %w (preserving error chains)
// 5. Error inspection with errors.Is and errors.As
// 6. Custom error types with additional data
// 7. Error handling patterns and best practices
//
// PHILOSOPHY:
// In Go, errors are VALUES, not exceptions. This means:
// - Errors are returned explicitly (not thrown)
// - You handle errors where they occur (no try/catch scanning)
// - Errors can be inspected, wrapped, and transformed like any other value
//
// THE ERROR INTERFACE:
// type error interface {
//     Error() string
// }
// That's it! Any type with an Error() method is an error.

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// ============================================================================
// SECTION 1: Sentinel Errors (Predefined Error Values)
// ============================================================================

// MACRO-COMMENT: What Are Sentinel Errors?
// Sentinel errors are predefined error variables that represent specific conditions.
// They're called "sentinel" because they act as guards or markers.
//
// WHY USE THEM:
// 1. Callers can check for specific error conditions
// 2. Consistent error values across your package
// 3. Self-documenting (the name tells you what went wrong)
//
// NAMING CONVENTION:
// - Start with "Err" prefix: ErrNotFound, ErrInvalidInput, etc.
// - Use package-level variables (var, not const)
// - Create with errors.New (not fmt.Errorf)

var (
	// ErrNotFound indicates a resource doesn't exist.
	// MICRO-COMMENT: Common pattern in database/cache operations
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists indicates a resource already exists.
	// MICRO-COMMENT: Common pattern in creation operations
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrInvalidInput indicates the input doesn't meet requirements.
	// MICRO-COMMENT: Common pattern in validation
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized indicates insufficient permissions.
	// MICRO-COMMENT: Common pattern in authentication/authorization
	ErrUnauthorized = errors.New("unauthorized access")
)

// getUser simulates a database lookup that returns a sentinel error.
//
// MACRO-COMMENT: Returning Sentinel Errors
// When a specific condition occurs (user not found), we return the
// sentinel error directly. Callers can then use == or errors.Is to check.
func getUser(id int) (string, error) {
	// MICRO-COMMENT: Simulate database lookup
	users := map[int]string{
		1: "Alice",
		2: "Bob",
		3: "Charlie",
	}

	user, exists := users[id]
	if !exists {
		// MICRO-COMMENT: Return the sentinel error directly
		// This allows callers to check: if err == ErrNotFound { ... }
		return "", ErrNotFound
	}

	return user, nil
}

// demonstrateSentinelErrors shows how to use and check sentinel errors.
func demonstrateSentinelErrors() {
	fmt.Println("=== Sentinel Errors ===")

	// MICRO-COMMENT: Try to get a user that doesn't exist
	_, err := getUser(999)

	// MACRO-COMMENT: Checking Sentinel Errors (Two Ways)
	// Way 1: Direct comparison with ==
	// This works because sentinel errors are package-level variables
	// with stable memory addresses (pointer equality).
	if err == ErrNotFound {
		fmt.Println("User not found (checked with ==)")
	}

	// Way 2: Using errors.Is (preferred, explained in detail later)
	// This is better because it works even if the error is wrapped.
	if errors.Is(err, ErrNotFound) {
		fmt.Println("User not found (checked with errors.Is)")
	}

	fmt.Printf("Error message: %v\n", err)
	fmt.Println()
}

// ============================================================================
// SECTION 2: Error Wrapping (Adding Context to Errors)
// ============================================================================

// MACRO-COMMENT: Why Wrap Errors?
// When an error occurs deep in the call stack, the error message alone
// often doesn't tell you WHERE it happened or WHAT was being attempted.
//
// Example without wrapping:
//   "no such file or directory"
//   - Which file? For what purpose?
//
// Example with wrapping:
//   "load config: read file config.json: no such file or directory"
//   - Clear chain of what happened!
//
// THE %w VERB:
// fmt.Errorf("context: %w", err) does two things:
// 1. Includes the error message (like %v)
// 2. Preserves the error for errors.Is and errors.As

// readFile simulates reading a file and wraps any error.
func readFile(path string) ([]byte, error) {
	// MICRO-COMMENT: Attempt to read the file
	data, err := os.ReadFile(path)
	if err != nil {
		// MACRO-COMMENT: Wrapping with %w
		// This creates a NEW error that:
		// - Contains the message: "read file <path>: <original message>"
		// - Wraps the original error (accessible via errors.Unwrap)
		// - Preserves the error chain for errors.Is and errors.As
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	return data, nil
}

// loadConfig tries to load a config file, adding more context.
func loadConfig(path string) (string, error) {
	// MICRO-COMMENT: Read the file (which may wrap an error)
	data, err := readFile(path)
	if err != nil {
		// MACRO-COMMENT: Wrapping Again (Multi-Level Context)
		// Each layer adds context:
		//   loadConfig: "load config: ..."
		//   readFile: "read file config.json: ..."
		//   os.ReadFile: "open config.json: no such file or directory"
		//
		// Final message: "load config: read file config.json: open config.json: no such file or directory"
		return "", fmt.Errorf("load config: %w", err)
	}

	return string(data), nil
}

// demonstrateErrorWrapping shows error wrapping in action.
func demonstrateErrorWrapping() {
	fmt.Println("=== Error Wrapping ===")

	// MICRO-COMMENT: Try to load a config file that doesn't exist
	_, err := loadConfig("nonexistent.json")

	// MICRO-COMMENT: Print the error (shows the full chain)
	fmt.Printf("Error: %v\n", err)
	// Output: "load config: read file nonexistent.json: open nonexistent.json: no such file or directory"

	// MACRO-COMMENT: Unwrapping Errors
	// Even though the error is wrapped multiple times, we can still
	// check if the ORIGINAL error is os.ErrNotExist.
	//
	// errors.Is walks the error chain:
	//   1. Check if err == os.ErrNotExist? No
	//   2. Unwrap err, check again? No
	//   3. Unwrap again, check again? No
	//   4. Unwrap again, check again? Yes! (os.ErrNotExist)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Root cause: file does not exist (found via errors.Is)")
	}

	// MACRO-COMMENT: Contrast with %v (NO WRAPPING)
	// If we had used %v instead of %w:
	//   fmt.Errorf("load config: %v", err)
	//
	// The message would look the same, BUT:
	// - errors.Is would NOT find os.ErrNotExist (chain is broken)
	// - errors.As would NOT work (can't access original error)
	//
	// RULE: Always use %w when wrapping errors!

	fmt.Println()
}

// ============================================================================
// SECTION 3: Custom Error Types (Structured Error Data)
// ============================================================================

// MACRO-COMMENT: Why Custom Error Types?
// Sometimes you need more than just a string. Custom error types let you:
// 1. Attach structured data (field name, timestamp, error code, etc.)
// 2. Implement additional behavior (retry logic, logging, etc.)
// 3. Let callers extract and inspect the data
//
// PATTERN:
// 1. Define a struct with error data
// 2. Implement the Error() method
// 3. Optionally implement Unwrap() if wrapping another error

// ValidationError represents a validation failure with details.
//
// MICRO-COMMENT: This error type carries structured data:
// - Field: which field failed validation
// - Value: the invalid value
// - Reason: why it's invalid
type ValidationError struct {
	Field  string      // The field that failed validation
	Value  interface{} // The invalid value
	Reason string      // Why it's invalid
}

// Error implements the error interface.
//
// MACRO-COMMENT: The Error() Method
// This method must return a human-readable string describing the error.
// It's what gets printed when you fmt.Println(err) or err.Error().
//
// BEST PRACTICE: Include all relevant details in the message.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field %s: %s (value: %v)",
		e.Field, e.Reason, e.Value)
}

// validateAge checks if an age is valid.
func validateAge(age int) error {
	if age < 0 {
		// MICRO-COMMENT: Return a custom error with structured data
		return ValidationError{
			Field:  "age",
			Value:  age,
			Reason: "must be non-negative",
		}
	}
	if age > 150 {
		return ValidationError{
			Field:  "age",
			Value:  age,
			Reason: "must be 150 or less",
		}
	}
	return nil
}

// demonstrateCustomErrors shows how to create and use custom error types.
func demonstrateCustomErrors() {
	fmt.Println("=== Custom Error Types ===")

	// MICRO-COMMENT: Try to validate an invalid age
	err := validateAge(-5)

	// MACRO-COMMENT: Using Custom Errors (Two Approaches)
	//
	// Approach 1: Just check if there's an error
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Approach 2: Extract the custom error type to access fields
	var ve ValidationError
	if errors.As(err, &ve) {
		// MICRO-COMMENT: Now we have access to all fields!
		fmt.Printf("Field: %s\n", ve.Field)
		fmt.Printf("Value: %v\n", ve.Value)
		fmt.Printf("Reason: %s\n", ve.Reason)
	}

	fmt.Println()
}

// ============================================================================
// SECTION 4: Wrapping Custom Errors
// ============================================================================

// DatabaseError wraps another error and adds database-specific context.
//
// MACRO-COMMENT: Wrapping Pattern for Custom Errors
// When your custom error wraps another error:
// 1. Store the wrapped error in a field
// 2. Implement Error() to include both messages
// 3. Implement Unwrap() to return the wrapped error
//
// This allows errors.Is and errors.As to traverse the chain.
type DatabaseError struct {
	Operation string // What operation failed (SELECT, INSERT, etc.)
	Table     string // Which table was involved
	Err       error  // The underlying error
}

// Error implements the error interface.
func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s on %s: %v",
		e.Operation, e.Table, e.Err)
}

// Unwrap returns the wrapped error.
//
// MACRO-COMMENT: The Unwrap() Method
// This method is recognized by errors.Is and errors.As.
// It allows them to traverse the error chain.
//
// Without Unwrap():
//   errors.Is(dbErr, ErrNotFound) → false (can't see wrapped error)
//
// With Unwrap():
//   errors.Is(dbErr, ErrNotFound) → true (finds it in the chain)
func (e DatabaseError) Unwrap() error {
	return e.Err
}

// queryUser simulates a database query that returns a wrapped custom error.
func queryUser(id int) (string, error) {
	// MICRO-COMMENT: Simulate query failure
	_, err := getUser(id)
	if err != nil {
		// MACRO-COMMENT: Wrap with DatabaseError
		// This adds database-specific context while preserving
		// the original error (ErrNotFound) for inspection.
		return "", DatabaseError{
			Operation: "SELECT",
			Table:     "users",
			Err:       err,
		}
	}

	return "user data", nil
}

// demonstrateWrappedCustomErrors shows wrapped custom errors.
func demonstrateWrappedCustomErrors() {
	fmt.Println("=== Wrapped Custom Errors ===")

	// MICRO-COMMENT: Query a user that doesn't exist
	_, err := queryUser(999)

	// MICRO-COMMENT: Print the full error message
	fmt.Printf("Error: %v\n", err)

	// MACRO-COMMENT: Extract the DatabaseError
	var dbErr DatabaseError
	if errors.As(err, &dbErr) {
		fmt.Printf("Operation: %s\n", dbErr.Operation)
		fmt.Printf("Table: %s\n", dbErr.Table)

		// MACRO-COMMENT: Check the wrapped error
		// Even though the error is wrapped in DatabaseError,
		// errors.Is can still find ErrNotFound because
		// DatabaseError implements Unwrap().
		if errors.Is(dbErr, ErrNotFound) {
			fmt.Println("Root cause: resource not found")
		}
	}

	fmt.Println()
}

// ============================================================================
// SECTION 5: errors.Is (Checking Error Identity in Chains)
// ============================================================================

// MACRO-COMMENT: What Is errors.Is?
// errors.Is(err, target) checks if ANY error in the chain matches target.
//
// HOW IT WORKS:
// 1. Compare err == target
// 2. If not equal, check if err has an Unwrap() method
// 3. If yes, call Unwrap() and repeat from step 1
// 4. If no more errors to unwrap, return false
//
// WHY USE IT:
// - Works with wrapped errors (== doesn't)
// - Future-proof (if someone wraps your error later)
// - Idiomatic Go

func demonstrateErrorsIs() {
	fmt.Println("=== errors.Is (Checking Error Identity) ===")

	// MICRO-COMMENT: Create a chain of wrapped errors
	baseErr := ErrNotFound
	wrapped1 := fmt.Errorf("level 1: %w", baseErr)
	wrapped2 := fmt.Errorf("level 2: %w", wrapped1)
	wrapped3 := fmt.Errorf("level 3: %w", wrapped2)

	fmt.Printf("Error chain: %v\n", wrapped3)

	// MACRO-COMMENT: Direct Comparison Fails
	// wrapped3 is NOT the same pointer as ErrNotFound
	isEqual := (wrapped3 == ErrNotFound)
	fmt.Printf("wrapped3 == ErrNotFound: %v (pointer comparison)\n", isEqual)

	// MACRO-COMMENT: errors.Is Succeeds
	// errors.Is walks the chain and finds ErrNotFound
	isFound := errors.Is(wrapped3, ErrNotFound)
	fmt.Printf("errors.Is(wrapped3, ErrNotFound): %v (chain traversal)\n", isFound)

	// MICRO-COMMENT: errors.Is works at any level
	fmt.Printf("errors.Is(wrapped3, wrapped1): %v\n", errors.Is(wrapped3, wrapped1))
	fmt.Printf("errors.Is(wrapped2, baseErr): %v\n", errors.Is(wrapped2, baseErr))

	// MACRO-COMMENT: Standard Library Sentinel Errors
	// The standard library defines many sentinel errors:
	// - io.EOF (end of file)
	// - os.ErrNotExist (file doesn't exist)
	// - context.Canceled (context was canceled)
	// You should always use errors.Is to check for them.

	fmt.Println()
}

// ============================================================================
// SECTION 6: errors.As (Extracting Error Types from Chains)
// ============================================================================

// MACRO-COMMENT: What Is errors.As?
// errors.As(err, &target) finds the first error in the chain that
// matches target's type and assigns it to target.
//
// HOW IT WORKS:
// 1. Check if err can be assigned to *target (type assertion)
// 2. If yes, assign and return true
// 3. If no, check if err has Unwrap()
// 4. If yes, call Unwrap() and repeat from step 1
// 5. If no more errors, return false
//
// WHY USE IT:
// - Extract custom error types from wrapped errors
// - Access fields and methods on the original error
// - Type-safe alternative to type assertions

// TimeoutError represents a timeout with context.
type TimeoutError struct {
	Operation string
	Duration  time.Duration
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("timeout after %v during %s", e.Duration, e.Operation)
}

// demonstrateErrorsAs shows how to extract error types.
func demonstrateErrorsAs() {
	fmt.Println("=== errors.As (Extracting Error Types) ===")

	// MICRO-COMMENT: Create a wrapped timeout error
	baseErr := TimeoutError{
		Operation: "database query",
		Duration:  5 * time.Second,
	}
	wrapped := fmt.Errorf("failed to fetch user: %w", baseErr)

	fmt.Printf("Error: %v\n", wrapped)

	// MACRO-COMMENT: Extract the TimeoutError
	// errors.As takes a POINTER to the target type
	var te TimeoutError
	if errors.As(wrapped, &te) {
		// MICRO-COMMENT: Now we have the original TimeoutError!
		fmt.Printf("Timeout operation: %s\n", te.Operation)
		fmt.Printf("Timeout duration: %v\n", te.Duration)

		// MACRO-COMMENT: Handle Timeout-Specific Logic
		// Based on the timeout duration, we might decide to retry
		if te.Duration < 10*time.Second {
			fmt.Println("Timeout was short, could retry")
		}
	}

	// MACRO-COMMENT: errors.As vs Type Assertion
	// Type assertion only works on the immediate error:
	//   te, ok := wrapped.(TimeoutError)  // ❌ Fails (wrapped is *wrapError, not TimeoutError)
	//
	// errors.As walks the chain:
	//   errors.As(wrapped, &te)  // ✅ Works (finds TimeoutError in chain)

	fmt.Println()
}

// ============================================================================
// SECTION 7: Multi-Error Handling (Collecting Multiple Errors)
// ============================================================================

// MACRO-COMMENT: Multi-Error Pattern
// Sometimes you want to process multiple items and collect all errors,
// not just stop at the first one.
//
// USE CASES:
// - Batch processing (process all items, report all failures)
// - Validation (check all fields, report all violations)
// - Concurrent operations (collect errors from goroutines)

// MultiError holds multiple errors.
type MultiError struct {
	Errors []error
}

// Error implements the error interface.
func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred: %v (and %d more)",
		len(m.Errors), m.Errors[0], len(m.Errors)-1)
}

// Unwrap returns all errors for Go 1.20+ multi-error support.
//
// MACRO-COMMENT: Go 1.20+ Multi-Error Unwrapping
// Since Go 1.20, Unwrap() can return []error (not just error).
// This allows errors.Is and errors.As to check ALL wrapped errors.
func (m MultiError) Unwrap() []error {
	return m.Errors
}

// processItems simulates processing multiple items.
func processItems(items []int) error {
	var multi MultiError

	// MICRO-COMMENT: Process each item, collect errors
	for _, item := range items {
		if err := processItem(item); err != nil {
			multi.Errors = append(multi.Errors, err)
		}
	}

	// MICRO-COMMENT: Return MultiError only if there are errors
	if len(multi.Errors) > 0 {
		return multi
	}
	return nil
}

// processItem simulates processing a single item.
func processItem(item int) error {
	if item < 0 {
		return ValidationError{
			Field:  "item",
			Value:  item,
			Reason: "must be non-negative",
		}
	}
	if item > 100 {
		return ValidationError{
			Field:  "item",
			Value:  item,
			Reason: "must be 100 or less",
		}
	}
	return nil
}

// demonstrateMultiError shows multi-error handling.
func demonstrateMultiError() {
	fmt.Println("=== Multi-Error Handling ===")

	// MICRO-COMMENT: Process items with some invalid values
	items := []int{10, -5, 50, 200, 75}
	err := processItems(items)

	if err != nil {
		fmt.Printf("Error: %v\n", err)

		// MACRO-COMMENT: Extract MultiError to Access All Errors
		var multi MultiError
		if errors.As(err, &multi) {
			fmt.Printf("Total errors: %d\n", len(multi.Errors))
			for i, e := range multi.Errors {
				fmt.Printf("  [%d] %v\n", i+1, e)
			}
		}

		// MACRO-COMMENT: Check If Any Error Is a Specific Type
		// With Go 1.20+, errors.Is checks ALL errors in the slice
		for _, item := range items {
			expectedErr := ValidationError{Field: "item", Value: item}
			if errors.As(err, &expectedErr) {
				fmt.Printf("Found validation error for item %d\n", item)
				break
			}
		}
	}

	fmt.Println()
}

// ============================================================================
// SECTION 8: Error Handling Patterns
// ============================================================================

// MACRO-COMMENT: Common Error Handling Patterns
// 1. Guard clauses (return early)
// 2. Named return values for cleanup
// 3. Defer with error checking
// 4. Retry on temporary errors

// parseAndValidate demonstrates guard clauses (early returns).
//
// MICRO-COMMENT: Guard Clause Pattern
// Check for errors immediately and return early.
// This keeps the happy path unindented and easy to follow.
func parseAndValidate(input string) (int, error) {
	// MICRO-COMMENT: Guard clause 1: Empty input
	if input == "" {
		return 0, ErrInvalidInput
	}

	// MICRO-COMMENT: Guard clause 2: Parse error
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("parse input: %w", err)
	}

	// MICRO-COMMENT: Guard clause 3: Validation error
	if err := validateAge(value); err != nil {
		return 0, fmt.Errorf("validate: %w", err)
	}

	// MICRO-COMMENT: Happy path (unindented)
	return value, nil
}

// writeToFile demonstrates named return values for cleanup.
//
// MACRO-COMMENT: Named Return Values for Cleanup
// Use a named error return value to capture cleanup errors.
func writeToFile(path string, data []byte) (err error) {
	// MICRO-COMMENT: Create the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	// MACRO-COMMENT: Defer Cleanup with Error Checking
	// The defer captures and can modify the named return value (err).
	// If the write succeeds but close fails, we return the close error.
	defer func() {
		closeErr := file.Close()
		if err == nil {
			// MICRO-COMMENT: Only set err if no previous error
			err = closeErr
		}
		// MICRO-COMMENT: If both write and close fail, we return the write error
		// (the first error is usually more important)
	}()

	// MICRO-COMMENT: Write the data
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// demonstrateErrorPatterns shows common error handling patterns.
func demonstrateErrorPatterns() {
	fmt.Println("=== Error Handling Patterns ===")

	// MICRO-COMMENT: Guard clauses
	value, err := parseAndValidate("25")
	if err != nil {
		fmt.Printf("Parse failed: %v\n", err)
	} else {
		fmt.Printf("Parse succeeded: %d\n", value)
	}

	// MICRO-COMMENT: Named return with cleanup
	err = writeToFile("/tmp/test.txt", []byte("hello"))
	if err != nil {
		fmt.Printf("Write failed: %v\n", err)
	} else {
		fmt.Println("Write succeeded")
	}

	fmt.Println()
}

// ============================================================================
// SECTION 9: When to Panic vs Return Error
// ============================================================================

// MACRO-COMMENT: Panic vs Error - The Decision Tree
//
// USE ERROR when:
// - The error is expected/possible (file not found, network timeout)
// - The caller can recover or handle it
// - The error is part of normal operation
//
// USE PANIC when:
// - The error represents a programmer mistake (bug)
// - The error is unrecoverable (can't continue)
// - During initialization (before main() starts)
//
// EXAMPLES:
// - File doesn't exist → ERROR (caller might try another path)
// - Out of bounds slice access → PANIC (programmer bug)
// - Failed to parse config at startup → PANIC (can't run without config)

// mustLoadConfig panics if config loading fails (unrecoverable).
func mustLoadConfig(path string) string {
	cfg, err := loadConfig(path)
	if err != nil {
		// MACRO-COMMENT: Panic for Unrecoverable Errors
		// If the config file is required for the app to function,
		// panicking is acceptable (especially during initialization).
		panic(fmt.Sprintf("FATAL: Cannot load config: %v", err))
	}
	return cfg
}

// safeDivide returns an error for invalid input (recoverable).
func safeDivide(a, b int) (int, error) {
	if b == 0 {
		// MACRO-COMMENT: Return Error for Expected Conditions
		// Division by zero is an expected error condition.
		// The caller can handle it (e.g., show error message).
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// demonstratePanicVsError shows when to panic vs return error.
func demonstratePanicVsError() {
	fmt.Println("=== Panic vs Error ===")

	// MICRO-COMMENT: Error (expected condition, caller handles)
	result, err := safeDivide(10, 0)
	if err != nil {
		fmt.Printf("Division failed (handled gracefully): %v\n", err)
	} else {
		fmt.Printf("Result: %d\n", result)
	}

	// MICRO-COMMENT: Panic (unrecoverable, commented out to not crash)
	// mustLoadConfig("required.conf")  // Would panic if file doesn't exist

	fmt.Println("Use errors for expected failures, panic for unrecoverable states")
	fmt.Println()
}

// ============================================================================
// SECTION 10: Standard Library Error Examples
// ============================================================================

// demonstrateStdlibErrors shows common standard library errors.
func demonstrateStdlibErrors() {
	fmt.Println("=== Standard Library Errors ===")

	// MACRO-COMMENT: io.EOF (End of File)
	// io.EOF is a sentinel error indicating the end of input.
	// It's NOT an error condition—it's expected!
	fmt.Printf("io.EOF: %v\n", io.EOF)
	fmt.Printf("Is io.EOF an error? %v (yes, but a special one)\n", io.EOF != nil)

	// MACRO-COMMENT: os.ErrNotExist
	// Returned when a file doesn't exist.
	_, err := os.Open("nonexistent.txt")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File doesn't exist (checked with errors.Is)")
	}

	// MACRO-COMMENT: os.ErrPermission
	// Returned when you don't have permission.
	// (We can't easily demonstrate this, but you'd check the same way)

	// MACRO-COMMENT: Pathological Error (Should Never Happen)
	// Some errors indicate bugs if they occur.
	// Example: os.ErrInvalid (invalid argument to a function)

	fmt.Println()
}

// ============================================================================
// MAIN: Run All Demonstrations
// ============================================================================

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  Error Handling in Go                                    ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	demonstrateSentinelErrors()
	demonstrateErrorWrapping()
	demonstrateCustomErrors()
	demonstrateWrappedCustomErrors()
	demonstrateErrorsIs()
	demonstrateErrorsAs()
	demonstrateMultiError()
	demonstrateErrorPatterns()
	demonstratePanicVsError()
	demonstrateStdlibErrors()

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  Key Insights                                             ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Println("║  1. Errors are VALUES (not exceptions)                   ║")
	fmt.Println("║  2. Use errors.New for simple sentinel errors            ║")
	fmt.Println("║  3. Use fmt.Errorf with %w to wrap errors                ║")
	fmt.Println("║  4. Use errors.Is to check for specific errors           ║")
	fmt.Println("║  5. Use errors.As to extract custom error types          ║")
	fmt.Println("║  6. Custom errors carry structured data                  ║")
	fmt.Println("║  7. Implement Unwrap() to support error chains           ║")
	fmt.Println("║  8. Return errors (don't panic) unless unrecoverable     ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
}
