//go:build !solution && !reference

package jwtauthmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *User, secret []byte, expiresIn time.Duration) (string, error) {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func ValidateToken(tokenString string, secret []byte) (*Claims, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	// TODO: Implement this function
	panic("not implemented")
}

func RequireRole(role string) func(http.Handler) http.Handler {
	// TODO: Implement this function
	panic("not implemented")
}

func GetClaims(r *http.Request) (*Claims, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func RefreshToken(refreshTokenString string, secret []byte, newExpiresIn time.Duration) (string, error) {
	// TODO: Implement this function
	panic("not implemented")
}
