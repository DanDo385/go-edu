//go:build solution
// +build solution

package exercise

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
	now := time.Now()

	// Create claims with user information
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "jwt-auth-server",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
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
	// Parse token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method to prevent algorithm confusion attack
		// This is critical! Without this check, an attacker could change the algorithm
		// from RS256 to HS256 and sign with the public key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
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
			// Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Parse "Bearer <token>" format
			// Expected format: "Bearer eyJhbGciOi..."
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization format, expected 'Bearer {token}'", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := ValidateToken(tokenString, secret)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized)
				return
			}

			// Add claims to request context for use by downstream handlers
			ctx := context.WithValue(r.Context(), "claims", claims)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
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
			// Extract claims from context (should be added by AuthMiddleware)
			claims, err := GetClaims(r)
			if err != nil {
				http.Error(w, "unauthorized: no claims in context", http.StatusUnauthorized)
				return
			}

			// Check if user has required role
			hasRole := false
			for _, userRole := range claims.Roles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, fmt.Sprintf("forbidden: requires '%s' role", role), http.StatusForbidden)
				return
			}

			// User has required role, proceed
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
	// Extract value from context
	value := r.Context().Value("claims")
	if value == nil {
		return nil, fmt.Errorf("no claims in context")
	}

	// Type assert to *Claims
	claims, ok := value.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type in context")
	}

	return claims, nil
}

// RefreshToken generates a new access token using a valid refresh token.
//
// Stretch goal: Implement refresh token pattern
//
// This is a simplified implementation. In production:
// - Refresh tokens should be stored in a database with user ID
// - Should support revocation (blacklist or token versioning)
// - Should rotate refresh tokens (issue new refresh token on use)
// - Should have longer expiration than access tokens
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
	// Validate the refresh token
	claims, err := ValidateToken(refreshTokenString, secret)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// In production, you would:
	// 1. Check if refresh token is in database (not revoked)
	// 2. Verify the token type (should have a "type": "refresh" claim)
	// 3. Update last used timestamp in database
	// 4. Optionally rotate the refresh token

	// Create a new user object from claims
	user := &User{
		ID:       claims.UserID,
		Username: claims.Username,
		Roles:    claims.Roles,
	}

	// Generate new access token with shorter expiration
	newAccessToken, err := GenerateToken(user, secret, newExpiresIn)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newAccessToken, nil
}
