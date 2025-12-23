//go:build !solution
// +build !solution

package exercise

import (
	"errors"
	"fmt"
)

// ============================================================================ 
// Sentinel Errors
// ============================================================================

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrInvalidUserID = errors.New("invalid user ID")
)

// FindUser simulates finding a user by ID.
func FindUser(id int) (string, error) {
	// TODO: Implement this function.

	// This function introduces "sentinel errors".
	// A sentinel error is a specific, pre-declared error value that you can check for by identity.

	// Step 1: Validate the input.
	// - If `id` is less than or equal to 0, return `ErrInvalidUserID`.

	// Step 2: Simulate a "not found" condition.
	// - If `id` is greater than 1000, return `ErrUserNotFound`.

	// Step 3: Handle the success case.
	// - Return a simulated username and a `nil` error.

	// Memory & Language Comparison:
	// - Go: Sentinel errors are a common pattern for simple, fixed error conditions. They are compared by value (`err == ErrUserNotFound`).
	// - C: This is similar to returning specific integer error codes (e.g., -1 for not found).
	// - Java/C#: This is somewhat similar to throwing specific, well-known exception types (e.g., `FileNotFoundException`).
	// - Rust: This is strongly discouraged. Rust prefers using enums for error types (e.g., `enum UserError { NotFound, InvalidId }`), which is more type-safe.
	return "", nil
}

// CreateUser simulates creating a new user.
func CreateUser(username string) error {
	// TODO: Implement this function.

	// Step 1: Validate the input.
	// - If `username` is empty, return `ErrInvalidUserID`.

	// Step 2: Simulate a "conflict" or "already exists" condition.
	// - If the username is "admin" or "root", return `ErrUserExists`.

	// Step 3: Handle the success case.
	// - Return `nil`.
	return nil
}

// ReadConfig simulates reading a configuration file.
func ReadConfig(id int) (string, error) {
	// TODO: Implement this function.

	// This function demonstrates error wrapping.
	// Error wrapping allows you to add context to an error without losing the original error.

	// Step 1: Call `FindUser(id)`.
	// Step 2: If `FindUser` returns an error, wrap it with additional context.
	// - Use `fmt.Errorf("read config for user %d: %w", id, err)`.
	// - The `%w` verb is the key part. It tells `fmt.Errorf` to create a new error that "wraps" the original `err`.
	// - This creates a chain of errors, allowing you to inspect the underlying cause later.

	// Step 3: If there's no error, return the username.
	return "", nil
}

// LoadUserData simulates loading user data from multiple sources.
func LoadUserData(id int) (string, error) {
	// TODO: Implement this function.

	// This function shows how error chains can be built up.

	// Step 1: Call `ReadConfig(id)`.
	// Step 2: If `ReadConfig` returns an error, wrap it *again*.
	// - `return "", fmt.Errorf("load user data: %w", err)`
	// - Now you have a three-level error chain: `LoadUserData` error -> `ReadConfig` error -> `FindUser` error.

	// Step 3: If there's no error, return the username.

	// Memory & Language Comparison:
	// - Go: The `%w` verb and the `errors` package provide a standardized way to handle error chains.
	// - Java/C#: This is the default behavior. When you catch an exception and throw a new one, you can pass the original exception as the "cause" (`new Exception("message", cause)`), creating a similar chain.
	// - Python: You can achieve a similar effect with `raise NewException("message") from original_exception`.
	// - Rust: Error handling libraries like `anyhow` or `thiserror` provide powerful ways to create and manage error chains with rich context.
	return "", nil
}

// IsNotFoundError checks if an error is (or wraps) ErrUserNotFound.
func IsNotFoundError(err error) bool {
	// TODO: Implement this function.

	// This function demonstrates how to check for a specific error within a chain.

	// Step 1: Handle the nil case.
	// - If `err` is `nil`, there's no error, so return `false`.

	// Step 2: Use `errors.Is`.
	// - `return errors.Is(err, ErrUserNotFound)`
	// - `errors.Is` traverses the error chain (by calling the `Unwrap()` method on each error).
	// - It returns `true` if any error in the chain matches the target error (`ErrUserNotFound`).
	// - This is the correct way to check for sentinel errors when you are using error wrapping. A simple `err == ErrUserNotFound` would fail if the error has been wrapped.
	return false
}

// GetUserWithFallback attempts to get a user, falling back to "guest" if not found.
func GetUserWithFallback(id int) (string, error) {
	// TODO: Implement this function.

	// This function shows a practical use case for `errors.Is`.

	// Step 1: Call `FindUser(id)`.
	// Step 2: If an error occurs, check if it's a "not found" error.
	// - `if errors.Is(err, ErrUserNotFound) { ... }`
	// - If it is, return the fallback value ("guest") and a `nil` error, because you've successfully handled the error condition.

	// Step 3: If it's a different kind of error, return it to the caller.
	// - `return "", err`

	// Step 4: If there was no error, return the username.
	//
	// Memory & Language Comparison:
	// - Go: `errors.Is` is the modern, idiomatic way to check for error identity.
	// - C++: You would typically use `dynamic_cast` to see if a caught exception matches a specific type.
	// - Java: You would use a `catch (FileNotFoundException e)` block to handle a specific type of exception.
	// - Rust: You would use a `match` statement to handle different variants of an error `enum`.
	return "", nil
}

// ValidationError represents a validation failure with details.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface for ValidationError.
func (e ValidationError) Error() string {
	// TODO: Implement this method.

	// To be a valid error, a type must implement the `error` interface, which has a single method: `Error() string`.
	// This method should return a human-readable string representation of the error.
	// Use `fmt.Sprintf` to create a descriptive message, like "validation error: [Field] [Message]".
	return ""
}

// ValidateUsername checks if a username is valid.
func ValidateUsername(username string) error {
	// TODO: Implement this function.

	// This function shows how to return your custom error type.

	// Step 1: Check for an empty username.
	// - If `username` is empty, return a `ValidationError` with the `Field` as "username" and `Message` as "cannot be empty".

	// Step 2: Check for minimum length.
	// - If the username is too short (e.g., less than 3 characters), return a `ValidationError` with the appropriate message.

	// Step 3: Check for maximum length.
	// - If the username is too long (e.g., more than 20 characters), return a `ValidationError` with the appropriate message.

	// Step 4: If all checks pass, return `nil`.

	// Memory & Language Comparison:
	// - Go: Creating custom error types (structs that implement the `error` interface) is the idiomatic way to provide rich, structured information about an error.
	// - Java/C#: This is equivalent to creating a custom exception class (e.g., `class ValidationException extends Exception { ... }`).
	// - Python: You would create a custom exception class by inheriting from `Exception`.
	// - Rust: You would define a `struct` or `enum` and implement the `std::error::Error` trait for it. Libraries like `thiserror` make this very easy.
	return nil
}

// GetValidationField extracts the field name from a ValidationError.
func GetValidationField(err error) (string, bool) {
	// TODO: Implement this function.

	// This function demonstrates how to extract a specific error type from an error chain.

	// Step 1: Handle the nil case.
	// - If `err` is `nil`, return `"", false`.

	// Step 2: Use `errors.As`.
	// - `var ve ValidationError` - Declare a variable of the type you want to find.
	// - `if errors.As(err, &ve) { ... }`
	//   - `errors.As` traverses the error chain.
	//   - If it finds an error that can be assigned to `ValidationError`, it assigns the error to `ve` and returns `true`.
	//   - This is more powerful than a type assertion because it works with wrapped errors.
	// - If it returns `true`, you can now access the fields of `ve`, so return `ve.Field, true`.

	// Step 3: If no `ValidationError` is found, return `"", false`.

	// Memory & Language Comparison:
	// - Go: `errors.As` is the modern, idiomatic way to extract a specific error type from a chain, allowing you to access its fields.
	// - Java: This is similar to a `catch (ValidationException e)` block, where you can then call methods on the caught exception `e`.
	// - C#: `catch (ValidationException e)` or using `is` and `as` keywords.
	// - Python: `except ValidationException as e:`.
	// - Rust: `if let Some(e) = err.downcast_ref::<ValidationError>() { ... }` provides a similar mechanism for trait objects.
	return "", false
}

// DatabaseError wraps another error and adds database context.
type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

func (e DatabaseError) Error() string {
	// TODO: Implement this method.
	// - Return a descriptive string that includes the operation, the table, and the underlying error's message.
	// - Example: "database error: [Operation] on [Table]: [underlying error]"
	return ""
}

func (e DatabaseError) Unwrap() error {
	// TODO: Implement this method.

	// This `Unwrap()` method is what makes `DatabaseError` a "wrapper" error.
	// - It should simply return the underlying error (`e.Err`).
	// - When `errors.Is` or `errors.As` are called on a `DatabaseError`, they will use this method to get the next error in the chain and continue their search.
	// - If you forget to implement this method, error wrapping will not work with your custom type.
	return nil
}

// QueryUser simulates a database query.
func QueryUser(id int) (string, error) {
	// TODO: Implement this function.

	// This function shows how to create and return your custom wrapper error.

	// Step 1: Call `FindUser(id)`.
	// Step 2: If an error occurs, create a `DatabaseError`.
	// - The `Operation` should be "SELECT".
	// - The `Table` should be "users".
	// - The `Err` field should be set to the error you received from `FindUser`.
	// - Return this `DatabaseError`.

	// Step 3: If successful, return the username.
	return "", nil
}

// MultiError holds multiple errors.
type MultiError struct {
	Errors []error
}

func (m MultiError) Error() string {
	// TODO: Implement this method.
	// - It should provide a summary of the errors.
	// - If there are no errors, return "no errors".
	// - If there is one error, return that error's message.
	// - If there are multiple errors, return a summary like "multiple errors: [first error message] (and [N] more)".
	return ""
}

func (m MultiError) Unwrap() []error {
	// TODO: Implement this method.

	// This is a special method introduced in Go 1.20.
	// - If an `Unwrap()` method returns a `[]error`, then `errors.Is` and `errors.As` will check *each* error in that slice.
	// - This allows you to check if a `MultiError` contains a specific underlying error.
	// - It should simply return the `m.Errors` slice.
	return nil
}

// ValidateUsers validates multiple usernames.
func ValidateUsers(usernames []string) error {
	// TODO: Implement this function.

	// This function shows how to collect multiple errors.

	// Step 1: Create a `MultiError` variable.
	// Step 2: Loop through the `usernames`.
	// Step 3: For each username, call `ValidateUsername`.
	// Step 4: If `ValidateUsername` returns an error, append it to the `Errors` slice of your `MultiError`.
	// Step 5: After the loop, if the `Errors` slice is not empty, return the `MultiError`.
	// Step 6: Otherwise, return `nil`.

	// Memory & Language Comparison:
	// - Go: The ability to return multiple errors as a single `error` is a powerful pattern, especially with the `Unwrap() []error` feature. The standard library doesn't have a built-in `MultiError` type, so you often create your own or use a third-party library (like `hashicorp/go-multierror`).
	// - C#: `AggregateException` is a standard library class that serves the exact same purpose.
	// - JavaScript: `AggregateError` is a standard error type for this purpose, often used with `Promise.any()`.
	// - Rust: This is less common. The idiomatic approach is often to fail on the first error, but you can manually collect errors into a `Vec<Error>`.
	return nil
}

// ProcessUser demonstrates guard clauses and error handling patterns.
func ProcessUser(username string) error {
	// TODO: Implement this function.

	// This function demonstrates the "guard clause" or "early return" pattern, which is very common in Go.
	// Instead of nesting `if/else` blocks, you check for errors at the beginning of the function and return immediately if one occurs.
	// This keeps the "happy path" (the successful case) at the lowest level of indentation, making the code easier to read.

	// Step 1 (Guard Clause): Validate the username.
	// - Call `ValidateUsername(username)`.
	// - If it returns an error, wrap it and return immediately.

	// Step 2 (Guard Clause): Check for banned users.
	// - If the username is "banned", return a new error.

	// Step 3 (Guard Clause): Create the user.
	// - Call `CreateUser(username)`.
	// - If it returns an error, wrap it and return immediately.

	// Step 4 (Happy Path): If all checks pass, return `nil`.

	// Memory & Language Comparison:
	// - Go: This pattern is extremely idiomatic in Go due to the explicit `if err != nil` check.
	// - Python: Also very common in Python.
	// - C#: Less common, as `try/catch` blocks are the primary mechanism for error handling.
	// - Rust: The `?` operator provides a very concise way to achieve the same "early return" pattern. `do_something()?` will automatically return if the result is an `Err`.
	return nil
}

// RetryableError indicates an error that can be retried.
type RetryableError struct {
	Err     error
	Retries int
}

func (e RetryableError) Error() string {
	// TODO: Implement this method.
	// - Return a descriptive message that includes the number of retries and the underlying error message.
	// - Example: "retryable error (attempt [N]): [underlying error]"
	return ""
}

func (e RetryableError) Unwrap() error {
	// TODO: Implement this method.
	// - To allow `errors.Is` and `errors.As` to inspect the error chain, this must return the wrapped error (`e.Err`).
	return nil
}

// IsRetryable checks if an error is retryable.
func IsRetryable(err error) bool {
	// TODO: Implement this function.

	// This function combines everything: checking for a specific error type in a chain and inspecting its state.

	// Step 1: Handle the nil case.
	// - If `err` is `nil`, it's not retryable.

	// Step 2: Use `errors.As` to find a `RetryableError`.
	// - `var re RetryableError`
	// - `if errors.As(err, &re) { ... }`

	// Step 3: If one is found, check its state.
	// - The error is retryable only if the number of retries is below a certain threshold (e.g., 3).
	// - `return re.Retries < 3`

	// Step 4: If no `RetryableError` is found, it's not retryable.
	return false
}
