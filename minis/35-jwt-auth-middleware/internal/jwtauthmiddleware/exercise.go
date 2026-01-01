//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package jwtauthmiddleware

import (
	"net/http"

	"time"
)

type User struct {
	ID       int
	Username string
	Password string // In production: bcrypt hash
	Roles    []string
}

type Claims struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}
// TODO: implement GenerateToken.
func GenerateToken(user *User, secret []byte, expiresIn time.Duration) (string, error) {
	panic("TODO: implement")
}
// TODO: implement ValidateToken.
func ValidateToken(tokenString string, secret []byte) (*Claims, error) { panic("TODO: implement") }
// TODO: implement AuthMiddleware.
func AuthMiddleware(secret []byte) func(http.Handler) http.Handler { panic("TODO: implement") }
// TODO: implement RequireRole.
func RequireRole(role string) func(http.Handler) http.Handler { panic("TODO: implement") }
// TODO: implement GetClaims.
func GetClaims(r *http.Request) (*Claims, error) { panic("TODO: implement") }
// TODO: implement RefreshToken.
func RefreshToken(refreshTokenString string, secret []byte, newExpiresIn time.Duration) (string, error) {
	panic("TODO: implement")
}
