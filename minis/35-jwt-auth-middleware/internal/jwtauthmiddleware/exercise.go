//go:build !solution && !reference

package jwtauthmiddleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
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

// GenerateToken implements the exercise.
//
// TODO: Implement this function
func GenerateToken(user *User, secret []byte, expiresIn time.Duration) (string, error) {
	// TODO: Implement
	return "", nil
}

// ValidateToken implements the exercise.
//
// TODO: Implement this function
func ValidateToken(tokenString string, secret []byte) (*Claims, error) {
	// TODO: Implement
	return nil, nil
}

// AuthMiddleware implements the exercise.
//
// TODO: Implement this function
func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	// TODO: Implement
	return nil
}

// RequireRole implements the exercise.
//
// TODO: Implement this function
func RequireRole(role string) func(http.Handler) http.Handler {
	// TODO: Implement
	return nil
}

// GetClaims implements the exercise.
//
// TODO: Implement this function
func GetClaims(r *http.Request) (*Claims, error) {
	// TODO: Implement
	return nil, nil
}

// RefreshToken implements the exercise.
//
// TODO: Implement this function
func RefreshToken(refreshTokenString string, secret []byte, newExpiresIn time.Duration) (string, error) {
	// TODO: Implement
	return "", nil
}
