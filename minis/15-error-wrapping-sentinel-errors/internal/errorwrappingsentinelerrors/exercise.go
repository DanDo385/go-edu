//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package errorwrappingsentinelerrors

import "errors"

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
// TODO: implement FindUser.
func FindUser(id int) (string, error) { panic("TODO: implement") }
// TODO: implement CreateUser.
func CreateUser(username string) error { panic("TODO: implement") }
// TODO: implement ReadConfig.
func ReadConfig(id int) (string, error) { panic("TODO: implement") }
// TODO: implement LoadUserData.
func LoadUserData(id int) (string, error) { panic("TODO: implement") }
// TODO: implement IsNotFoundError.
func IsNotFoundError(err error) bool { panic("TODO: implement") }
// TODO: implement GetUserWithFallback.
func GetUserWithFallback(id int) (string, error) { panic("TODO: implement") }
// TODO: implement Error.
func (e ValidationError) Error() string { panic("TODO: implement") }
// TODO: implement ValidateUsername.
func ValidateUsername(username string) error { panic("TODO: implement") }
// TODO: implement GetValidationField.
func GetValidationField(err error) (string, bool) { panic("TODO: implement") }
// TODO: implement Error.
func (e DatabaseError) Error() string { panic("TODO: implement") }
// TODO: implement Unwrap.
func (e DatabaseError) Unwrap() error { panic("TODO: implement") }
// TODO: implement QueryUser.
func QueryUser(id int) (string, error) { panic("TODO: implement") }
// TODO: implement Error.
func (m MultiError) Error() string { panic("TODO: implement") }
// TODO: implement Unwrap.
func (m MultiError) Unwrap() []error { panic("TODO: implement") }
// TODO: implement ValidateUsers.
func ValidateUsers(usernames []string) error { panic("TODO: implement") }
// TODO: implement ProcessUser.
func ProcessUser(username string) error { panic("TODO: implement") }
// TODO: implement Error.
func (e RetryableError) Error() string { panic("TODO: implement") }
// TODO: implement Unwrap.
func (e RetryableError) Unwrap() error { panic("TODO: implement") }
// TODO: implement IsRetryable.
func IsRetryable(err error) bool { panic("TODO: implement") }
