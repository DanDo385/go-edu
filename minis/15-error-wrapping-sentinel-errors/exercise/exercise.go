//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for understanding error handling in Go.
//
// LEARNING OBJECTIVES:
// - Create and use sentinel errors
// - Wrap errors with context using %w
// - Check error identity with errors.Is
// - Extract error types with errors.As
// - Implement custom error types
// - Handle errors idiomatically

package exercise

import (
	"errors"
	"fmt"
)

// Suppress unused import errors
var _ = errors.New
var _ = fmt.Sprintf

// TODO: Implement these functions according to the specifications in the tests.

// ============================================================================
// EXERCISE 1: Sentinel Errors
// ============================================================================

// Define sentinel errors for a user management system.
// These should be package-level variables created with errors.New.

// TODO: Create ErrUserNotFound - indicates a user doesn't exist
var ErrUserNotFound error

// TODO: Create ErrUserExists - indicates a user already exists
var ErrUserExists error

// TODO: Create ErrInvalidUserID - indicates the user ID is invalid
var ErrInvalidUserID error

// FindUser simulates finding a user by ID.
//
// REQUIREMENTS:
// - If id <= 0, return ErrInvalidUserID
// - If id > 1000, return ErrUserNotFound
// - Otherwise, return a username in the format "user_<id>" and nil error
//
// EXAMPLES:
//   FindUser(-1) → "", ErrInvalidUserID
//   FindUser(5000) → "", ErrUserNotFound
//   FindUser(42) → "user_42", nil
func FindUser(id int) (string, error) {
	// TODO: Implement this
	return "", nil
}

// CreateUser simulates creating a new user.
//
// REQUIREMENTS:
// - If username is empty, return ErrInvalidUserID
// - If username is "admin" or "root", return ErrUserExists
// - Otherwise, return nil (success)
//
// EXAMPLES:
//   CreateUser("") → ErrInvalidUserID
//   CreateUser("admin") → ErrUserExists
//   CreateUser("alice") → nil
func CreateUser(username string) error {
	// TODO: Implement this
	return nil
}

// ============================================================================
// EXERCISE 2: Error Wrapping with %w
// ============================================================================

// ReadConfig simulates reading a configuration file.
//
// REQUIREMENTS:
// - Call FindUser with the provided id
// - If there's an error, wrap it with context using fmt.Errorf and %w
// - The wrapped error should include the message: "read config for user <id>: %w"
// - If successful, return the username
//
// EXAMPLES:
//   ReadConfig(0) → "", wrapped error with message "read config for user 0: <original error>"
//   ReadConfig(42) → "user_42", nil
//
// HINT: Use fmt.Errorf with %w to preserve the error chain
func ReadConfig(id int) (string, error) {
	// TODO: Implement this
	return "", nil
}

// LoadUserData simulates loading user data from multiple sources.
//
// REQUIREMENTS:
// - Call ReadConfig with the provided id
// - If there's an error, wrap it AGAIN with: "load user data: %w"
// - This creates a multi-level error chain
// - If successful, return the username
//
// This demonstrates how errors accumulate context as they bubble up.
func LoadUserData(id int) (string, error) {
	// TODO: Implement this
	return "", nil
}

// ============================================================================
// EXERCISE 3: errors.Is (Checking Error Identity)
// ============================================================================

// IsNotFoundError checks if an error is (or wraps) ErrUserNotFound.
//
// REQUIREMENTS:
// - Use errors.Is to check if err is or contains ErrUserNotFound
// - Return true if it is, false otherwise
// - Handle nil errors (return false)
//
// EXAMPLES:
//   IsNotFoundError(nil) → false
//   IsNotFoundError(ErrUserNotFound) → true
//   IsNotFoundError(fmt.Errorf("wrapped: %w", ErrUserNotFound)) → true
//   IsNotFoundError(ErrInvalidUserID) → false
//
// HINT: errors.Is traverses the error chain for you
func IsNotFoundError(err error) bool {
	// TODO: Implement this
	return false
}

// GetUserWithFallback attempts to get a user, falling back to a default if not found.
//
// REQUIREMENTS:
// - Call FindUser with id
// - If the error is ErrUserNotFound (check with errors.Is), return "guest" and nil
// - If there's any other error, return it
// - If successful, return the username
//
// EXAMPLES:
//   GetUserWithFallback(5000) → "guest", nil (not found, use fallback)
//   GetUserWithFallback(-1) → "", ErrInvalidUserID (different error, return it)
//   GetUserWithFallback(42) → "user_42", nil (found)
func GetUserWithFallback(id int) (string, error) {
	// TODO: Implement this
	return "", nil
}

// ============================================================================
// EXERCISE 4: Custom Error Types
// ============================================================================

// ValidationError represents a validation failure with details.
//
// REQUIREMENTS:
// - Implement the error interface (Error() method)
// - The Error() method should return: "validation error: <Field> <Message>"
//
// EXAMPLE:
//   err := ValidationError{Field: "email", Message: "invalid format"}
//   err.Error() → "validation error: email invalid format"
type ValidationError struct {
	Field   string // The field that failed validation
	Message string // Why it failed
}

// TODO: Implement the Error() method for ValidationError
func (e ValidationError) Error() string {
	// TODO: Implement this
	return ""
}

// ValidateUsername checks if a username is valid.
//
// REQUIREMENTS:
// - If username is empty, return ValidationError{Field: "username", Message: "cannot be empty"}
// - If len(username) < 3, return ValidationError{Field: "username", Message: "too short"}
// - If len(username) > 20, return ValidationError{Field: "username", Message: "too long"}
// - Otherwise, return nil
//
// EXAMPLES:
//   ValidateUsername("") → ValidationError{Field: "username", Message: "cannot be empty"}
//   ValidateUsername("ab") → ValidationError{Field: "username", Message: "too short"}
//   ValidateUsername("alice") → nil
func ValidateUsername(username string) error {
	// TODO: Implement this
	return nil
}

// ============================================================================
// EXERCISE 5: errors.As (Extracting Error Types)
// ============================================================================

// GetValidationField extracts the field name from a ValidationError.
//
// REQUIREMENTS:
// - Use errors.As to check if err is (or wraps) a ValidationError
// - If it is, return the Field value and true
// - If it's not, return "" and false
// - Handle nil errors (return "", false)
//
// EXAMPLES:
//   GetValidationField(nil) → "", false
//   GetValidationField(ValidationError{Field: "email"}) → "email", true
//   GetValidationField(fmt.Errorf("wrapped: %w", ValidationError{Field: "age"})) → "age", true
//   GetValidationField(errors.New("other")) → "", false
//
// HINT: errors.As takes a pointer to the target type
func GetValidationField(err error) (string, bool) {
	// TODO: Implement this
	return "", false
}

// ============================================================================
// EXERCISE 6: Custom Error with Wrapping
// ============================================================================

// DatabaseError wraps another error and adds database context.
type DatabaseError struct {
	Operation string // The operation that failed (e.g., "SELECT", "INSERT")
	Table     string // The table involved
	Err       error  // The underlying error
}

// TODO: Implement Error() method for DatabaseError
// Format: "database error: <Operation> on <Table>: <Err>"
// Example: "database error: SELECT on users: resource not found"
func (e DatabaseError) Error() string {
	// TODO: Implement this
	return ""
}

// TODO: Implement Unwrap() method for DatabaseError
// This should return the wrapped error (Err field)
// This allows errors.Is and errors.As to traverse the chain
func (e DatabaseError) Unwrap() error {
	// TODO: Implement this
	return nil
}

// QueryUser simulates a database query.
//
// REQUIREMENTS:
// - Call FindUser with id
// - If there's an error, wrap it in a DatabaseError with:
//   - Operation: "SELECT"
//   - Table: "users"
//   - Err: the error from FindUser
// - If successful, return the username
//
// EXAMPLES:
//   QueryUser(5000) → "", DatabaseError{Operation: "SELECT", Table: "users", Err: ErrUserNotFound}
//   QueryUser(42) → "user_42", nil
func QueryUser(id int) (string, error) {
	// TODO: Implement this
	return "", nil
}

// ============================================================================
// EXERCISE 7: Multi-Error Handling
// ============================================================================

// MultiError holds multiple errors.
type MultiError struct {
	Errors []error
}

// TODO: Implement Error() method for MultiError
// REQUIREMENTS:
// - If len(Errors) == 0, return "no errors"
// - If len(Errors) == 1, return Errors[0].Error()
// - Otherwise, return "multiple errors: <first error message> (and N more)"
//   where N = len(Errors) - 1
//
// EXAMPLES:
//   MultiError{Errors: nil}.Error() → "no errors"
//   MultiError{Errors: []error{errors.New("fail")}}.Error() → "fail"
//   MultiError{Errors: []error{errors.New("fail1"), errors.New("fail2")}}.Error()
//     → "multiple errors: fail1 (and 1 more)"
func (m MultiError) Error() string {
	// TODO: Implement this
	return ""
}

// TODO: Implement Unwrap() method for MultiError
// REQUIREMENTS:
// - Return the Errors slice
// - This enables Go 1.20+ multi-error support for errors.Is and errors.As
func (m MultiError) Unwrap() []error {
	// TODO: Implement this
	return nil
}

// ValidateUsers validates multiple usernames.
//
// REQUIREMENTS:
// - For each username, call ValidateUsername
// - Collect all errors in a MultiError
// - If there are any errors, return the MultiError
// - If all validations pass, return nil
//
// EXAMPLES:
//   ValidateUsers([]string{"alice", "bob"}) → nil
//   ValidateUsers([]string{"", "ab"}) → MultiError with 2 errors
//   ValidateUsers([]string{"alice", "", "bob"}) → MultiError with 1 error
func ValidateUsers(usernames []string) error {
	// TODO: Implement this
	return nil
}

// ============================================================================
// EXERCISE 8: Error Handling Patterns
// ============================================================================

// ProcessUser demonstrates guard clauses and error handling patterns.
//
// REQUIREMENTS:
// - Validate the username using ValidateUsername
// - If validation fails, wrap the error with "validate username: %w"
// - If username is "banned", return a NEW error: "user is banned"
// - Call CreateUser with the username
// - If creation fails, wrap with "create user: %w"
// - If everything succeeds, return nil
//
// This exercises:
// - Guard clauses (early returns)
// - Error wrapping
// - Creating new errors
//
// EXAMPLES:
//   ProcessUser("") → wrapped ValidationError
//   ProcessUser("banned") → error "user is banned"
//   ProcessUser("admin") → wrapped ErrUserExists
//   ProcessUser("alice") → nil
func ProcessUser(username string) error {
	// TODO: Implement this
	return nil
}

// ============================================================================
// EXERCISE 9: Optional - Advanced Error Chain
// ============================================================================

// RetryableError indicates an error that can be retried.
type RetryableError struct {
	Err     error // The underlying error
	Retries int   // How many retries have been attempted
}

// TODO: Implement Error() method for RetryableError
// Format: "retryable error (attempt <Retries>): <Err>"
func (e RetryableError) Error() string {
	// TODO: Implement this
	return ""
}

// TODO: Implement Unwrap() method for RetryableError
func (e RetryableError) Unwrap() error {
	// TODO: Implement this
	return nil
}

// IsRetryable checks if an error is retryable.
//
// REQUIREMENTS:
// - Use errors.As to check if err contains a RetryableError
// - If it does AND Retries < 3, return true
// - Otherwise, return false
func IsRetryable(err error) bool {
	// TODO: Implement this
	return false
}
