//go:build solution
// +build solution

/*
Problem: Understanding Go's error handling patterns

Requirements:
1. Use sentinel errors for simple conditions
2. Wrap errors with context using %w
3. Check error identity with errors.Is
4. Extract error types with errors.As
5. Handle multiple errors

Data Structure:
- error: Interface with Error() string method
- Sentinel error: Pre-declared error value
- Wrapped error: Error chain with context
- Custom error: Struct implementing error interface

Algorithm: Error Chain Traversal
- errors.Is: Walk chain, check identity
- errors.As: Walk chain, extract type
- Unwrap(): Return next error in chain

Why error handling is critical:
- Explicit error returns (no exceptions)
- Error chains preserve context
- Type-safe error inspection
- Composable error handling
*/

package exercise

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

type MultiError struct {
	Errors []error
}

type RetryableError struct {
	Err     error
	Retries int
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrInvalidUserID = errors.New("invalid user ID")
)

// FindUser simulates finding a user.
// BREAKPOINT: Set breakpoint here to trace sentinel errors
// DEBUG: Watch 'id' validation
// DEBUG: Watch sentinel error returns
func FindUser(id int) (string, error) {
	// BREAKPOINT: Check invalid ID
	// DEBUG: Watch 'id <= 0' condition
	if id <= 0 {
		return "", ErrInvalidUserID
	}

	// BREAKPOINT: Check not found
	// DEBUG: Watch 'id > 1000' condition
	if id > 1000 {
		return "", ErrUserNotFound
	}

	// DEBUG: Success case
	return fmt.Sprintf("user_%d", id), nil
}

// CreateUser simulates creating a user.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch validation and conflict checks
func CreateUser(username string) error {
	// BREAKPOINT: Validate username
	if username == "" {
		return ErrInvalidUserID
	}

	// BREAKPOINT: Check conflicts
	// DEBUG: Watch reserved username check
	if username == "admin" || username == "root" {
		return ErrUserExists
	}

	return nil
}

// ReadConfig wraps FindUser error.
// BREAKPOINT: Set breakpoint here to trace error wrapping
// DEBUG: Watch error chain creation
func ReadConfig(id int) (string, error) {
	// BREAKPOINT: Call FindUser
	username, err := FindUser(id)
	if err != nil {
		// BREAKPOINT: Wrap error with %w
		// DEBUG: Watch fmt.Errorf with %w create chain
		return "", fmt.Errorf("read config for user %d: %w", id, err)
	}
	return username, nil
}

// LoadUserData creates multi-level error chain.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch error chain grow
func LoadUserData(id int) (string, error) {
	username, err := ReadConfig(id)
	if err != nil {
		// BREAKPOINT: Wrap again
		// DEBUG: Watch second level of wrapping
		return "", fmt.Errorf("load user data: %w", err)
	}
	return username, nil
}

// IsNotFoundError checks error identity.
// BREAKPOINT: Set breakpoint here to trace errors.Is
// DEBUG: Watch error chain traversal
func IsNotFoundError(err error) bool {
	// BREAKPOINT: Check nil
	if err == nil {
		return false
	}

	// BREAKPOINT: Use errors.Is
	// DEBUG: Watch errors.Is traverse chain
	// DEBUG: Returns true if any error in chain is ErrUserNotFound
	return errors.Is(err, ErrUserNotFound)
}

// GetUserWithFallback handles not found gracefully.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch error handling with fallback
func GetUserWithFallback(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		// BREAKPOINT: Check if not found
		// DEBUG: Watch errors.Is check specific error
		if errors.Is(err, ErrUserNotFound) {
			// DEBUG: Return fallback value
			return "guest", nil
		}
		return "", err
	}
	return username, nil
}

// Error implements error interface.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch error message formatting
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Message)
}

// ValidateUsername validates username.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch custom error creation
func ValidateUsername(username string) error {
	// BREAKPOINT: Check empty
	if username == "" {
		return ValidationError{
			Field:   "username",
			Message: "cannot be empty",
		}
	}

	// BREAKPOINT: Check min length
	if len(username) < 3 {
		return ValidationError{
			Field:   "username",
			Message: "too short",
		}
	}

	// BREAKPOINT: Check max length
	if len(username) > 20 {
		return ValidationError{
			Field:   "username",
			Message: "too long",
		}
	}

	return nil
}

// GetValidationField extracts field from error.
// BREAKPOINT: Set breakpoint here to trace errors.As
// DEBUG: Watch type extraction from chain
func GetValidationField(err error) (string, bool) {
	// BREAKPOINT: Check nil
	if err == nil {
		return "", false
	}

	// BREAKPOINT: Use errors.As
	// DEBUG: Watch errors.As extract ValidationError
	var ve ValidationError
	if errors.As(err, &ve) {
		// DEBUG: Watch ve populated with error data
		return ve.Field, true
	}

	return "", false
}

// Error implements error interface.
// BREAKPOINT: Set breakpoint here
func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error: %s on %s: %v", e.Operation, e.Table, e.Err)
}

// Unwrap returns wrapped error.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch error chain unwrapping
func (e DatabaseError) Unwrap() error {
	// DEBUG: Return next error in chain
	return e.Err
}

// QueryUser wraps database operation.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch custom wrapper error creation
func QueryUser(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		// BREAKPOINT: Wrap in DatabaseError
		// DEBUG: Watch custom error wrap sentinel error
		return "", DatabaseError{
			Operation: "SELECT",
			Table:     "users",
			Err:       err,
		}
	}
	return username, nil
}

// Error implements error interface.
// BREAKPOINT: Set breakpoint here
func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	return fmt.Sprintf("multiple errors: %v (and %d more)",
		m.Errors[0], len(m.Errors)-1)
}

// Unwrap returns all errors.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch multi-error unwrap
func (m MultiError) Unwrap() []error {
	// DEBUG: Return slice of errors (Go 1.20+)
	return m.Errors
}

// ValidateUsers validates multiple usernames.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch error collection
func ValidateUsers(usernames []string) error {
	var multi MultiError

	// BREAKPOINT: Loop through usernames
	for _, username := range usernames {
		if err := ValidateUsername(username); err != nil {
			// BREAKPOINT: Collect error
			// DEBUG: Watch multi.Errors grow
			multi.Errors = append(multi.Errors, err)
		}
	}

	// BREAKPOINT: Return multi-error
	if len(multi.Errors) > 0 {
		return multi
	}

	return nil
}

// ProcessUser demonstrates guard clauses.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch early returns (guard clauses)
func ProcessUser(username string) error {
	// BREAKPOINT: Guard clause 1
	// DEBUG: Watch validation guard
	if err := ValidateUsername(username); err != nil {
		return fmt.Errorf("validate username: %w", err)
	}

	// BREAKPOINT: Guard clause 2
	// DEBUG: Watch banned user guard
	if username == "banned" {
		return errors.New("user is banned")
	}

	// BREAKPOINT: Guard clause 3
	// DEBUG: Watch create user guard
	if err := CreateUser(username); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// Error implements error interface.
// BREAKPOINT: Set breakpoint here
func (e RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempt %d): %v", e.Retries, e.Err)
}

// Unwrap returns wrapped error.
// BREAKPOINT: Set breakpoint here
func (e RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable checks if error is retryable.
// BREAKPOINT: Set breakpoint here
// DEBUG: Watch retry logic
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// BREAKPOINT: Extract RetryableError
	var re RetryableError
	if errors.As(err, &re) {
		// DEBUG: Watch retry count check
		return re.Retries < 3
	}

	return false
}
