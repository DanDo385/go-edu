//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents a user in the system
type User struct {
	ID       int
	Username string
	Password string // In production: bcrypt hash
	Roles    []string
}

// Claims represents the JWT claims
type Claims struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token for the given user.
//
// The token should:
// - Include user ID, username, and roles in claims
// - Set expiration time based on expiresIn parameter
// - Set issued at time to current time
// - Use HMAC-SHA256 signing algorithm
// - Sign with the provided secret
//
// Parameters:
//   - user: User to generate token for
//   - secret: Secret key for signing
//   - expiresIn: Token expiration duration
//
// Returns:
//   - string: Signed JWT token
//   - error: Non-nil if token generation fails
func GenerateToken(user *User, secret []byte, expiresIn time.Duration) (string, error) {
	// TODO: Implement token generation
	// Hint: Create Claims, then jwt.NewWithClaims, then SignedString
	return "", nil
}

// ValidateToken validates a JWT token and returns the claims.
//
// The function should:
// - Parse the token string
// - Verify the signature using the provided secret
// - Check that the signing method is HMAC-SHA256 (prevent algorithm confusion)
// - Verify token hasn't expired
// - Return the parsed claims
//
// Parameters:
//   - tokenString: JWT token to validate
//   - secret: Secret key for verification
//
// Returns:
//   - *Claims: Parsed claims if valid
//   - error: Non-nil if token is invalid or expired
func ValidateToken(tokenString string, secret []byte) (*Claims, error) {
	// TODO: Implement token validation
	// Hint: Use jwt.ParseWithClaims with a key function that validates the signing method
	return nil, nil
}

// AuthMiddleware returns a middleware that validates JWT tokens.
//
// The middleware should:
// - Extract the Authorization header
// - Check for "Bearer <token>" format
// - Validate the token using ValidateToken
// - Add claims to request context with key "claims"
// - Call next handler if valid
// - Return 401 Unauthorized if token is missing or invalid
//
// Parameters:
//   - secret: Secret key for token validation
//
// Returns:
//   - Middleware function that wraps http.Handler
func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement authentication middleware
			// Hint: Get Authorization header, parse "Bearer <token>", validate, add to context
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole returns a middleware that checks if the user has the required role.
//
// The middleware should:
// - Extract claims from request context
// - Check if user has the required role
// - Call next handler if user has the role
// - Return 403 Forbidden if user doesn't have the role
//
// Parameters:
//   - role: Required role (e.g., "admin")
//
// Returns:
//   - Middleware function that wraps http.Handler
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement role-based access control
			// Hint: Get claims from context, check if role is in claims.Roles
			next.ServeHTTP(w, r)
		})
	}
}

// GetClaims extracts claims from the request context.
//
// Helper function to extract claims added by AuthMiddleware.
//
// Parameters:
//   - r: HTTP request
//
// Returns:
//   - *Claims: Claims from context
//   - error: Non-nil if no claims in context
func GetClaims(r *http.Request) (*Claims, error) {
	// TODO: Implement claims extraction
	// Hint: Use r.Context().Value("claims")
	return nil, nil
}

// RefreshToken generates a new access token using a valid refresh token.
//
// Stretch goal: Implement refresh token pattern
//
// Parameters:
//   - refreshTokenString: Valid refresh token
//   - secret: Secret key
//   - newExpiresIn: Expiration duration for new access token
//
// Returns:
//   - string: New access token
//   - error: Non-nil if refresh token is invalid
func RefreshToken(refreshTokenString string, secret []byte, newExpiresIn time.Duration) (string, error) {
	// TODO: Stretch goal - implement refresh token logic
	return "", nil
}
