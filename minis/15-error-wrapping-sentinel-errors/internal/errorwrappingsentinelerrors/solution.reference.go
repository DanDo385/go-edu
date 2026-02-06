//go:build reference

package errorwrappingsentinelerrors

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

func FindUser(id int) (string, error) {
	if id <= 0 {
		return "", ErrInvalidUserID
	}
	if id > 1000 {
		return "", ErrUserNotFound
	}
	return fmt.Sprintf("user_%d", id), nil
}

func CreateUser(username string) error {
	if username == "" {
		return ErrInvalidUserID
	}
	if username == "admin" || username == "root" {
		return ErrUserExists
	}
	return nil
}

func ReadConfig(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		return "", fmt.Errorf("read config for user %d: %w", id, err)
	}
	return username, nil
}

func LoadUserData(id int) (string, error) {
	username, err := ReadConfig(id)
	if err != nil {
		return "", fmt.Errorf("load user data: %w", err)
	}
	return username, nil
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrUserNotFound)
}

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

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

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

func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error: %s on %s: %v", e.Operation, e.Table, e.Err)
}

func (e DatabaseError) Unwrap() error {
	return e.Err
}

func QueryUser(id int) (string, error) {
	username, err := FindUser(id)
	if err != nil {
		return "", DatabaseError{Operation: "SELECT", Table: "users", Err: err}
	}
	return username, nil
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

func (e RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempt %d): %v", e.Retries, e.Err)
}

func (e RetryableError) Unwrap() error {
	return e.Err
}

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
