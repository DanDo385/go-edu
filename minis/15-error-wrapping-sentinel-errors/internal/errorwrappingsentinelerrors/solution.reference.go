//go:build reference

package errorwrappingsentinelerrors

/*
Reference Solution - Error Wrapping, Sentinel Errors, Custom Error Types
======================================================================

This file demonstrates Go's error handling idioms: sentinel errors, error
wrapping with %w, errors.Is/errors.As for unwrapping, and custom error types
that implement Unwrap for error chains.

Key concepts:
- Sentinel errors: package-level vars (ErrUserNotFound) for comparison with errors.Is
- fmt.Errorf("...: %w", err): wraps err so errors.Is/As can unwrap
- Unwrap(): return inner error; errors.Is walks the chain; errors.As extracts type
- MultiError: Unwrap() []error for multiple wrapped errors (Go 1.20+)
*/

import (
	"errors"
	"fmt"
)

// Sentinel errors: compare with errors.Is(err, ErrUserNotFound)
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrInvalidUserID = errors.New("invalid user id")
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

// FindUser - Returns sentinel errors for validation/lookup failures.
func FindUser(id int) (string, error) {
	if id <= 0 {
		return "", ErrInvalidUserID
	}
	if id > 1000 {
		return "", ErrUserNotFound
	}
	return fmt.Sprintf("user_%d", id), nil
}

// CreateUser - Returns sentinel errors; note reused ErrInvalidUserID for empty username.
func CreateUser(username string) error {
	if username == "" {
		return ErrInvalidUserID
	}
	if username == "admin" || username == "root" {
		return ErrUserExists
	}
	return nil
}

// ReadConfig - Wraps FindUser error with %w so errors.Is still finds ErrUserNotFound.
func ReadConfig(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		return "", fmt.Errorf("read config for user %d: %w", id, err)
	}
	return username, nil
}

// LoadUserData - Wraps again; chain: LoadUserData -> ReadConfig -> FindUser.
func LoadUserData(id int) (string, error) {
	username, err := ReadConfig(id)
	if err != nil {
		return "", fmt.Errorf("load user data: %w", err)
	}
	return username, nil
}

// IsNotFoundError - errors.Is walks unwrap chain; finds ErrUserNotFound even when wrapped.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrUserNotFound)
}

// GetUserWithFallback - If ErrUserNotFound, return "guest" instead of error.
func GetUserWithFallback(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "guest", nil
		}
		return "", err
	}
	return username, nil
}

// Error - ValidationError has no Unwrap; errors.As won't find wrapped sentinels through it.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// ValidateUsername - Returns ValidationError (custom type) for different validation failures.
func ValidateUsername(username string) error {
	if username == "" {
		return ValidationError{Field: "username", Message: "cannot be empty"}
	}
	if len(username) < 3 {
		return ValidationError{Field: "username", Message: "too short"}
	}
	if len(username) > 20 {
		return ValidationError{Field: "username", Message: "too long"}
	}
	return nil
}

// GetValidationField - errors.As extracts ValidationError from chain; returns Field.
func GetValidationField(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	var ve ValidationError
	if errors.As(err, &ve) {
		return ve.Field, true
	}
	return "", false
}

// Error implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error: %s on %s: %v", e.Operation, e.Table, e.Err)
}

// Unwrap - DatabaseError wraps inner error; errors.Is(err, ErrUserNotFound) works through it.
func (e DatabaseError) Unwrap() error {
	return e.Err
}

// QueryUser - Wraps FindUser in DatabaseError; chain preserves ErrUserNotFound.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func QueryUser(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		return "", DatabaseError{Operation: "SELECT", Table: "users", Err: err}
	}
	return username, nil
}

// Error - MultiError: format first error + count for readability.
func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	return fmt.Sprintf("multiple errors: %v (and %d more)", m.Errors[0], len(m.Errors)-1)
}

// Unwrap - MultiError implements Unwrap() []error (Go 1.20+); errors.Join uses this.
func (m MultiError) Unwrap() []error {
	return m.Errors
}

// ValidateUsers - Collects all validation errors; returns MultiError if any.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func ValidateUsers(usernames []string) error {
	var me MultiError
	for _, username := range usernames {
		if err := ValidateUsername(username); err != nil {
			me.Errors = append(me.Errors, err)
		}
	}
	if len(me.Errors) > 0 {
		return me
	}
	return nil
}

// ProcessUser - Chains ValidateUsername and CreateUser with %w; full trace preserved.
func ProcessUser(username string) error {
	if err := ValidateUsername(username); err != nil {
		return fmt.Errorf("validate username: %w", err)
	}
	if username == "banned" {
		return errors.New("user is banned")
	}
	if err := CreateUser(username); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// Error - RetryableError carries attempt count for logging.
func (e RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempt %d): %v", e.Retries, e.Err)
}

// Unwrap - RetryableError wraps inner error.
func (e RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable - errors.As extracts RetryableError; retries < 3 means still retryable.
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
