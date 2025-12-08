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
	if id <= 0 {
		return "", ErrInvalidUserID
	}
	if id > 1000 {
		return "", ErrUserNotFound
	}
	return fmt.Sprintf("user_%d", id), nil
}

// CreateUser simulates creating a new user.
func CreateUser(username string) error {
	if username == "" {
		return ErrInvalidUserID
	}
	if username == "admin" || username == "root" {
		return ErrUserExists
	}
	return nil
}

// ReadConfig simulates reading a configuration file.
func ReadConfig(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		return "", fmt.Errorf("read config for user %d: %w", id, err)
	}
	return username, nil
}

// LoadUserData simulates loading user data from multiple sources.
func LoadUserData(id int) (string, error) {
	username, err := ReadConfig(id)
	if err != nil {
		return "", fmt.Errorf("load user data: %w", err)
	}
	return username, nil
}

// IsNotFoundError checks if an error is (or wraps) ErrUserNotFound.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrUserNotFound)
}

// GetUserWithFallback attempts to get a user, falling back to "guest" if not found.
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

// ValidationError represents a validation failure with details.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface for ValidationError.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Message)
}

// ValidateUsername checks if a username is valid.
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

// GetValidationField extracts the field name from a ValidationError.
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

// DatabaseError wraps another error and adds database context.
type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error: %s on %s: %v", e.Operation, e.Table, e.Err)
}

func (e DatabaseError) Unwrap() error {
	return e.Err
}

// QueryUser simulates a database query.
func QueryUser(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		return "", DatabaseError{
			Operation: "SELECT",
			Table:     "users",
			Err:       err,
		}
	}
	return username, nil
}

// MultiError holds multiple errors.
type MultiError struct {
	Errors []error
}

func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	return fmt.Sprintf("multiple errors: %v (and %d more)", m.Errors[0], len(m.Errors)-1)
}

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

// ProcessUser demonstrates guard clauses and error handling patterns.
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

// RetryableError indicates an error that can be retried.
type RetryableError struct {
	Err     error
	Retries int
}

func (e RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempt %d): %v", e.Retries, e.Err)
}

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
