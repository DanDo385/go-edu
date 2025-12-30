package exercise

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// ============================================================================
// TEST 1: Sentinel Errors
// ============================================================================

func TestSentinelErrors(t *testing.T) {
	t.Run("sentinel errors are defined", func(t *testing.T) {
		if ErrUserNotFound == nil {
			t.Error("ErrUserNotFound should be defined")
		}
		if ErrUserExists == nil {
			t.Error("ErrUserExists should be defined")
		}
		if ErrInvalidUserID == nil {
			t.Error("ErrInvalidUserID should be defined")
		}
	})

	t.Run("sentinel errors have correct messages", func(t *testing.T) {
		if ErrUserNotFound != nil && !strings.Contains(ErrUserNotFound.Error(), "not found") {
			t.Errorf("ErrUserNotFound message should mention 'not found', got: %v", ErrUserNotFound)
		}
		if ErrUserExists != nil && !strings.Contains(ErrUserExists.Error(), "exists") {
			t.Errorf("ErrUserExists message should mention 'exists', got: %v", ErrUserExists)
		}
		if ErrInvalidUserID != nil && !strings.Contains(ErrInvalidUserID.Error(), "invalid") {
			t.Errorf("ErrInvalidUserID message should mention 'invalid', got: %v", ErrInvalidUserID)
		}
	})
}

func TestFindUser(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		wantUser    string
		wantErr     error
		wantErrType bool // true if we expect any error
	}{
		{
			name:     "valid user",
			id:       42,
			wantUser: "user_42",
			wantErr:  nil,
		},
		{
			name:        "invalid ID (negative)",
			id:          -1,
			wantUser:    "",
			wantErr:     ErrInvalidUserID,
			wantErrType: true,
		},
		{
			name:        "invalid ID (zero)",
			id:          0,
			wantUser:    "",
			wantErr:     ErrInvalidUserID,
			wantErrType: true,
		},
		{
			name:        "user not found",
			id:          5000,
			wantUser:    "",
			wantErr:     ErrUserNotFound,
			wantErrType: true,
		},
		{
			name:     "boundary case (1000)",
			id:       1000,
			wantUser: "user_1000",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := FindUser(tt.id)

			if tt.wantErrType {
				if err == nil {
					t.Errorf("FindUser(%d) expected error, got nil", tt.id)
				} else if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
					t.Errorf("FindUser(%d) expected error %v, got %v", tt.id, tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("FindUser(%d) unexpected error: %v", tt.id, err)
				}
			}

			if user != tt.wantUser {
				t.Errorf("FindUser(%d) = %q, want %q", tt.id, user, tt.wantUser)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  error
	}{
		{
			name:     "valid username",
			username: "alice",
			wantErr:  nil,
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  ErrInvalidUserID,
		},
		{
			name:     "admin exists",
			username: "admin",
			wantErr:  ErrUserExists,
		},
		{
			name:     "root exists",
			username: "root",
			wantErr:  ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateUser(tt.username)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("CreateUser(%q) expected error %v, got %v", tt.username, tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("CreateUser(%q) unexpected error: %v", tt.username, err)
				}
			}
		})
	}
}

// ============================================================================
// TEST 2: Error Wrapping
// ============================================================================

func TestReadConfig(t *testing.T) {
	t.Run("successful read", func(t *testing.T) {
		user, err := ReadConfig(42)
		if err != nil {
			t.Errorf("ReadConfig(42) unexpected error: %v", err)
		}
		if user != "user_42" {
			t.Errorf("ReadConfig(42) = %q, want %q", user, "user_42")
		}
	})

	t.Run("wraps error with context", func(t *testing.T) {
		_, err := ReadConfig(5000)
		if err == nil {
			t.Fatal("ReadConfig(5000) expected error, got nil")
		}

		// Check that it wraps ErrUserNotFound
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("ReadConfig should wrap ErrUserNotFound, got: %v", err)
		}

		// Check that the error message includes context
		errMsg := err.Error()
		if !strings.Contains(errMsg, "read config") {
			t.Errorf("Error should contain 'read config', got: %s", errMsg)
		}
	})
}

func TestLoadUserData(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		user, err := LoadUserData(42)
		if err != nil {
			t.Errorf("LoadUserData(42) unexpected error: %v", err)
		}
		if user != "user_42" {
			t.Errorf("LoadUserData(42) = %q, want %q", user, "user_42")
		}
	})

	t.Run("multi-level error wrapping", func(t *testing.T) {
		_, err := LoadUserData(5000)
		if err == nil {
			t.Fatal("LoadUserData(5000) expected error, got nil")
		}

		// Check that the original error is still findable
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("LoadUserData should preserve ErrUserNotFound in chain, got: %v", err)
		}

		// Check that both layers of context are present
		errMsg := err.Error()
		if !strings.Contains(errMsg, "load user data") {
			t.Errorf("Error should contain 'load user data', got: %s", errMsg)
		}
		if !strings.Contains(errMsg, "read config") {
			t.Errorf("Error should contain 'read config', got: %s", errMsg)
		}
	})
}

// ============================================================================
// TEST 3: errors.Is
// ============================================================================

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "direct ErrUserNotFound",
			err:  ErrUserNotFound,
			want: true,
		},
		{
			name: "wrapped ErrUserNotFound",
			err:  fmt.Errorf("wrapped: %w", ErrUserNotFound),
			want: true,
		},
		{
			name: "multi-wrapped ErrUserNotFound",
			err:  fmt.Errorf("level2: %w", fmt.Errorf("level1: %w", ErrUserNotFound)),
			want: true,
		},
		{
			name: "different error",
			err:  ErrInvalidUserID,
			want: false,
		},
		{
			name: "wrapped different error",
			err:  fmt.Errorf("wrapped: %w", ErrInvalidUserID),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotFoundError(tt.err)
			if got != tt.want {
				t.Errorf("IsNotFoundError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestGetUserWithFallback(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		wantUser string
		wantErr  error
	}{
		{
			name:     "user found",
			id:       42,
			wantUser: "user_42",
			wantErr:  nil,
		},
		{
			name:     "user not found, use fallback",
			id:       5000,
			wantUser: "guest",
			wantErr:  nil,
		},
		{
			name:     "invalid ID, return error",
			id:       -1,
			wantUser: "",
			wantErr:  ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := GetUserWithFallback(tt.id)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetUserWithFallback(%d) expected error %v, got %v", tt.id, tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("GetUserWithFallback(%d) unexpected error: %v", tt.id, err)
				}
			}

			if user != tt.wantUser {
				t.Errorf("GetUserWithFallback(%d) = %q, want %q", tt.id, user, tt.wantUser)
			}
		})
	}
}

// ============================================================================
// TEST 4: Custom Error Types
// ============================================================================

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "email",
		Message: "invalid format",
	}

	// Test that it implements error interface
	var _ error = err

	// Test Error() method
	errMsg := err.Error()
	if !strings.Contains(errMsg, "validation error") {
		t.Errorf("Error message should contain 'validation error', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "email") {
		t.Errorf("Error message should contain field name, got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "invalid format") {
		t.Errorf("Error message should contain message, got: %s", errMsg)
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		wantErr      bool
		wantField    string
		wantContains string
	}{
		{
			name:     "valid username",
			username: "alice",
			wantErr:  false,
		},
		{
			name:         "empty username",
			username:     "",
			wantErr:      true,
			wantField:    "username",
			wantContains: "empty",
		},
		{
			name:         "too short",
			username:     "ab",
			wantErr:      true,
			wantField:    "username",
			wantContains: "short",
		},
		{
			name:         "too long",
			username:     "abcdefghijklmnopqrstuvwxyz",
			wantErr:      true,
			wantField:    "username",
			wantContains: "long",
		},
		{
			name:     "minimum length",
			username: "abc",
			wantErr:  false,
		},
		{
			name:     "maximum length",
			username: "12345678901234567890",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateUsername(%q) expected error, got nil", tt.username)
					return
				}

				// Check that it's a ValidationError
				var ve ValidationError
				if !errors.As(err, &ve) {
					t.Errorf("Expected ValidationError, got %T", err)
					return
				}

				if ve.Field != tt.wantField {
					t.Errorf("ValidationError.Field = %q, want %q", ve.Field, tt.wantField)
				}

				if !strings.Contains(ve.Message, tt.wantContains) {
					t.Errorf("ValidationError.Message should contain %q, got %q", tt.wantContains, ve.Message)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateUsername(%q) unexpected error: %v", tt.username, err)
				}
			}
		})
	}
}

// ============================================================================
// TEST 5: errors.As
// ============================================================================

func TestGetValidationField(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantField string
		wantOk    bool
	}{
		{
			name:      "nil error",
			err:       nil,
			wantField: "",
			wantOk:    false,
		},
		{
			name:      "direct ValidationError",
			err:       ValidationError{Field: "email", Message: "invalid"},
			wantField: "email",
			wantOk:    true,
		},
		{
			name:      "wrapped ValidationError",
			err:       fmt.Errorf("wrapped: %w", ValidationError{Field: "age", Message: "invalid"}),
			wantField: "age",
			wantOk:    true,
		},
		{
			name:      "different error type",
			err:       errors.New("other error"),
			wantField: "",
			wantOk:    false,
		},
		{
			name:      "wrapped different error",
			err:       fmt.Errorf("wrapped: %w", ErrUserNotFound),
			wantField: "",
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := GetValidationField(tt.err)

			if ok != tt.wantOk {
				t.Errorf("GetValidationField(%v) ok = %v, want %v", tt.err, ok, tt.wantOk)
			}

			if field != tt.wantField {
				t.Errorf("GetValidationField(%v) field = %q, want %q", tt.err, field, tt.wantField)
			}
		})
	}
}

// ============================================================================
// TEST 6: Custom Error with Wrapping
// ============================================================================

func TestDatabaseError(t *testing.T) {
	baseErr := ErrUserNotFound
	dbErr := DatabaseError{
		Operation: "SELECT",
		Table:     "users",
		Err:       baseErr,
	}

	// Test Error() method
	errMsg := dbErr.Error()
	if !strings.Contains(errMsg, "database error") {
		t.Errorf("Error message should contain 'database error', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "SELECT") {
		t.Errorf("Error message should contain operation, got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "users") {
		t.Errorf("Error message should contain table, got: %s", errMsg)
	}

	// Test Unwrap() method
	unwrapped := dbErr.Unwrap()
	if unwrapped != baseErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, baseErr)
	}

	// Test that errors.Is works through the wrapper
	if !errors.Is(dbErr, ErrUserNotFound) {
		t.Error("errors.Is should find ErrUserNotFound through DatabaseError")
	}
}

func TestQueryUser(t *testing.T) {
	t.Run("successful query", func(t *testing.T) {
		user, err := QueryUser(42)
		if err != nil {
			t.Errorf("QueryUser(42) unexpected error: %v", err)
		}
		if user != "user_42" {
			t.Errorf("QueryUser(42) = %q, want %q", user, "user_42")
		}
	})

	t.Run("wraps error in DatabaseError", func(t *testing.T) {
		_, err := QueryUser(5000)
		if err == nil {
			t.Fatal("QueryUser(5000) expected error, got nil")
		}

		// Check that it's a DatabaseError
		var dbErr DatabaseError
		if !errors.As(err, &dbErr) {
			t.Fatalf("Expected DatabaseError, got %T", err)
		}

		if dbErr.Operation != "SELECT" {
			t.Errorf("DatabaseError.Operation = %q, want %q", dbErr.Operation, "SELECT")
		}
		if dbErr.Table != "users" {
			t.Errorf("DatabaseError.Table = %q, want %q", dbErr.Table, "users")
		}

		// Check that the original error is preserved
		if !errors.Is(err, ErrUserNotFound) {
			t.Error("QueryUser should preserve ErrUserNotFound in chain")
		}
	})
}

// ============================================================================
// TEST 7: Multi-Error Handling
// ============================================================================

func TestMultiError(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		me := MultiError{Errors: nil}
		if me.Error() != "no errors" {
			t.Errorf("MultiError with no errors should return 'no errors', got: %s", me.Error())
		}
	})

	t.Run("single error", func(t *testing.T) {
		me := MultiError{Errors: []error{errors.New("fail")}}
		if me.Error() != "fail" {
			t.Errorf("MultiError with single error should return that error, got: %s", me.Error())
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		me := MultiError{Errors: []error{
			errors.New("fail1"),
			errors.New("fail2"),
			errors.New("fail3"),
		}}
		errMsg := me.Error()
		if !strings.Contains(errMsg, "multiple errors") {
			t.Errorf("Error message should contain 'multiple errors', got: %s", errMsg)
		}
		if !strings.Contains(errMsg, "and 2 more") {
			t.Errorf("Error message should mention additional errors, got: %s", errMsg)
		}
	})

	t.Run("Unwrap returns errors slice", func(t *testing.T) {
		errs := []error{errors.New("fail1"), errors.New("fail2")}
		me := MultiError{Errors: errs}
		unwrapped := me.Unwrap()
		if len(unwrapped) != 2 {
			t.Errorf("Unwrap() should return all errors, got %d", len(unwrapped))
		}
	})
}

func TestValidateUsers(t *testing.T) {
	tests := []struct {
		name      string
		usernames []string
		wantErr   bool
		wantCount int // Number of errors expected
	}{
		{
			name:      "all valid",
			usernames: []string{"alice", "bob", "charlie"},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:      "all invalid",
			usernames: []string{"", "ab"},
			wantErr:   true,
			wantCount: 2,
		},
		{
			name:      "mixed valid and invalid",
			usernames: []string{"alice", "", "bob", "xy"},
			wantErr:   true,
			wantCount: 2,
		},
		{
			name:      "empty slice",
			usernames: []string{},
			wantErr:   false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsers(tt.usernames)

			if tt.wantErr {
				if err == nil {
					t.Fatal("ValidateUsers expected error, got nil")
				}

				var me MultiError
				if !errors.As(err, &me) {
					t.Fatalf("Expected MultiError, got %T", err)
				}

				if len(me.Errors) != tt.wantCount {
					t.Errorf("Expected %d errors, got %d", tt.wantCount, len(me.Errors))
				}
			} else {
				if err != nil {
					t.Errorf("ValidateUsers unexpected error: %v", err)
				}
			}
		})
	}
}

// ============================================================================
// TEST 8: Error Handling Patterns
// ============================================================================

func TestProcessUser(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		wantErr    bool
		checkError func(t *testing.T, err error)
	}{
		{
			name:     "valid user",
			username: "alice",
			wantErr:  false,
		},
		{
			name:     "validation fails",
			username: "",
			wantErr:  true,
			checkError: func(t *testing.T, err error) {
				if !strings.Contains(err.Error(), "validate username") {
					t.Errorf("Error should be wrapped with 'validate username', got: %v", err)
				}
				var ve ValidationError
				if !errors.As(err, &ve) {
					t.Error("Should preserve ValidationError in chain")
				}
			},
		},
		{
			name:     "banned user",
			username: "banned",
			wantErr:  true,
			checkError: func(t *testing.T, err error) {
				if !strings.Contains(err.Error(), "banned") {
					t.Errorf("Error should mention 'banned', got: %v", err)
				}
			},
		},
		{
			name:     "user exists",
			username: "admin",
			wantErr:  true,
			checkError: func(t *testing.T, err error) {
				if !strings.Contains(err.Error(), "create user") {
					t.Errorf("Error should be wrapped with 'create user', got: %v", err)
				}
				if !errors.Is(err, ErrUserExists) {
					t.Error("Should preserve ErrUserExists in chain")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessUser(tt.username)

			if tt.wantErr {
				if err == nil {
					t.Fatal("ProcessUser expected error, got nil")
				}
				if tt.checkError != nil {
					tt.checkError(t, err)
				}
			} else {
				if err != nil {
					t.Errorf("ProcessUser unexpected error: %v", err)
				}
			}
		})
	}
}

// ============================================================================
// TEST 9: Advanced Error Chain
// ============================================================================

func TestRetryableError(t *testing.T) {
	baseErr := errors.New("network timeout")
	retryErr := RetryableError{
		Err:     baseErr,
		Retries: 1,
	}

	// Test Error() method
	errMsg := retryErr.Error()
	if !strings.Contains(errMsg, "retryable") {
		t.Errorf("Error message should contain 'retryable', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "1") {
		t.Errorf("Error message should contain retry count, got: %s", errMsg)
	}

	// Test Unwrap() method
	unwrapped := retryErr.Unwrap()
	if unwrapped != baseErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "retryable with 1 attempt",
			err:  RetryableError{Err: errors.New("fail"), Retries: 1},
			want: true,
		},
		{
			name: "retryable with 2 attempts",
			err:  RetryableError{Err: errors.New("fail"), Retries: 2},
			want: true,
		},
		{
			name: "retryable with 3 attempts (max)",
			err:  RetryableError{Err: errors.New("fail"), Retries: 3},
			want: false,
		},
		{
			name: "wrapped retryable",
			err:  fmt.Errorf("wrapped: %w", RetryableError{Err: errors.New("fail"), Retries: 1}),
			want: true,
		},
		{
			name: "non-retryable error",
			err:  errors.New("permanent failure"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryable(tt.err)
			if got != tt.want {
				t.Errorf("IsRetryable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
