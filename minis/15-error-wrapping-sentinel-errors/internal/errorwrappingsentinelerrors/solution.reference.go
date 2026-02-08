//go:build reference

package errorwrappingsentinelerrors

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

import (
	"errors"
	"fmt"
)

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

// FindUser implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func FindUser(id int) (string, error) {
	if id <= 0 {
		return "", ErrInvalidUserID
	}
	if id > 1000 {
		return "", ErrUserNotFound
	}
	return fmt.Sprintf("user_%d", id), nil
}

// CreateUser implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func CreateUser(username string) error {
	if username == "" {
		return ErrInvalidUserID
	}
	if username == "admin" || username == "root" {
		return ErrUserExists
	}
	return nil
}

// ReadConfig implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func ReadConfig(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		return "", fmt.Errorf("read config for user %d: %w", id, err)
	}
	return username, nil
}

// LoadUserData implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func LoadUserData(id int) (string, error) {
	username, err := ReadConfig(id)
	if err != nil {
		return "", fmt.Errorf("load user data: %w", err)
	}
	return username, nil
}

// IsNotFoundError implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrUserNotFound)
}

// GetUserWithFallback implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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

// Error implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// ValidateUsername implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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

// GetValidationField implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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

// Unwrap implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (e DatabaseError) Unwrap() error {
	return e.Err
}

// QueryUser implements the reference behavior for this exercise.
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

// Error implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return "no errors"
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	return fmt.Sprintf("multiple errors: %v (and %d more)", m.Errors[0], len(m.Errors)-1)
}

// Unwrap implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (m MultiError) Unwrap() []error {
	return m.Errors
}

// ValidateUsers implements the reference behavior for this exercise.
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

// ProcessUser implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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

// Error implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (e RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempt %d): %v", e.Retries, e.Err)
}

// Unwrap implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (e RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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
