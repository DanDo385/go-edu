//go:build solution
// +build solution

package exercise

import (
	"errors"
	"fmt"
)

// ============================================================================
// TYPE DEFINITIONS (shared with exercise.go)
// ============================================================================

// ValidationError represents a validation failure with details.
type ValidationError struct {
	Field   string // The field that failed validation
	Message string // Why it failed
}

// DatabaseError wraps another error and adds database context.
type DatabaseError struct {
	Operation string // The operation that failed (e.g., "SELECT", "INSERT")
	Table     string // The table involved
	Err       error  // The underlying error
}

// MultiError holds multiple errors.
type MultiError struct {
	Errors []error
}

// RetryableError indicates an error that can be retried.
type RetryableError struct {
	Err     error // The underlying error
	Retries int   // How many retries have been attempted
}

// ============================================================================
// EXERCISE 1: Sentinel Errors
// ============================================================================

// Sentinel errors for user management
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrInvalidUserID = errors.New("invalid user ID")
)

// FindUser simulates finding a user by ID.
func FindUser(id int) (string, error) {
	// Check for invalid ID
	if id <= 0 {
		return "", ErrInvalidUserID
	}

	// Check if user doesn't exist
	if id > 1000 {
		return "", ErrUserNotFound
	}

	// Return the user
	return fmt.Sprintf("user_%d", id), nil
}

// CreateUser simulates creating a new user.
func CreateUser(username string) error {
	// Check for empty username
	if username == "" {
		return ErrInvalidUserID
	}

	// Check if user already exists
	if username == "admin" || username == "root" {
		return ErrUserExists
	}

	// Success
	return nil
}

// ============================================================================
// EXERCISE 2: Error Wrapping with %w
// ============================================================================

// ReadConfig simulates reading a configuration file.
func ReadConfig(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		// Wrap the error with context using %w
		return "", fmt.Errorf("read config for user %d: %w", id, err)
	}
	return username, nil
}

// LoadUserData simulates loading user data from multiple sources.
func LoadUserData(id int) (string, error) {
	username, err := ReadConfig(id)
	if err != nil {
		// Wrap again, creating a multi-level chain
		return "", fmt.Errorf("load user data: %w", err)
	}
	return username, nil
}

// ============================================================================
// EXERCISE 3: errors.Is (Checking Error Identity)
// ============================================================================

// IsNotFoundError checks if an error is (or wraps) ErrUserNotFound.
func IsNotFoundError(err error) bool {
	// Handle nil
	if err == nil {
		return false
	}

	// Use errors.Is to traverse the error chain
	return errors.Is(err, ErrUserNotFound)
}

// GetUserWithFallback attempts to get a user, falling back to a default if not found.
func GetUserWithFallback(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		// Check if it's specifically ErrUserNotFound
		if errors.Is(err, ErrUserNotFound) {
			// Use fallback
			return "guest", nil
		}
		// Different error, return it
		return "", err
	}

	// Success
	return username, nil
}

// ============================================================================
// EXERCISE 4: Custom Error Types
// ============================================================================

// Error implements the error interface for ValidationError.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Message)
}

// ValidateUsername checks if a username is valid.
func ValidateUsername(username string) error {
	if username == "" {
		return ValidationError{
			Field:   "username",
			Message: "cannot be empty",
		}
	}

	if len(username) < 3 {
		return ValidationError{
			Field:   "username",
			Message: "too short",
		}
	}

	if len(username) > 20 {
		return ValidationError{
			Field:   "username",
			Message: "too long",
		}
	}

	return nil
}

// ============================================================================
// EXERCISE 5: errors.As (Extracting Error Types)
// ============================================================================

// GetValidationField extracts the field name from a ValidationError.
func GetValidationField(err error) (string, bool) {
	// Handle nil
	if err == nil {
		return "", false
	}

	// Use errors.As to extract the ValidationError
	var ve ValidationError
	if errors.As(err, &ve) {
		return ve.Field, true
	}

	return "", false
}

// ============================================================================
// EXERCISE 6: Custom Error with Wrapping
// ============================================================================

// Error implements the error interface for DatabaseError.
func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error: %s on %s: %v", e.Operation, e.Table, e.Err)
}

// Unwrap returns the wrapped error.
func (e DatabaseError) Unwrap() error {
	return e.Err
}

// QueryUser simulates a database query.
func QueryUser(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		// Wrap in DatabaseError
		return "", DatabaseError{
			Operation: "SELECT",
			Table:     "users",
			Err:       err,
		}
	}
	return username, nil
}

// ============================================================================
// EXERCISE 7: Multi-Error Handling
// ============================================================================

// Error implements the error interface for MultiError.
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

// Unwrap returns all errors (Go 1.20+ multi-error support).
func (m MultiError) Unwrap() []error {
	return m.Errors
}

// ValidateUsers validates multiple usernames.
func ValidateUsers(usernames []string) error {
	var multi MultiError

	for _, username := range usernames {
		if err := ValidateUsername(username); err != nil {
			multi.Errors = append(multi.Errors, err)
		}
	}

	if len(multi.Errors) > 0 {
		return multi
	}

	return nil
}

// ============================================================================
// EXERCISE 8: Error Handling Patterns
// ============================================================================

// ProcessUser demonstrates guard clauses and error handling patterns.
func ProcessUser(username string) error {
	// Guard clause 1: Validate username
	if err := ValidateUsername(username); err != nil {
		return fmt.Errorf("validate username: %w", err)
	}

	// Guard clause 2: Check for banned users
	if username == "banned" {
		return errors.New("user is banned")
	}

	// Guard clause 3: Create user
	if err := CreateUser(username); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// ============================================================================
// EXERCISE 9: Optional - Advanced Error Chain
// ============================================================================

// Error implements the error interface for RetryableError.
func (e RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempt %d): %v", e.Retries, e.Err)
}

// Unwrap returns the wrapped error.
func (e RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable checks if an error is retryable.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var re RetryableError
	if errors.As(err, &re) {
		return re.Retries < 3
	}

	return false
}
