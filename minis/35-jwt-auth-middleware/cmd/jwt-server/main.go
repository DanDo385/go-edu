package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents a user in the system
// In production, this would be fetched from a database
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // In production: hashed with bcrypt
	Roles    []string `json:"roles"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse represents the login response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description,omitempty"`
}

// Mock database of users
var users = map[string]*User{
	"alice": {
		ID:       1,
		Username: "alice",
		Password: "secret123", // In production: bcrypt hash
		Roles:    []string{"user"},
	},
	"bob": {
		ID:       2,
		Username: "bob",
		Password: "password456",
		Roles:    []string{"user", "admin"},
	},
}

// Secret key for signing JWTs
// In production: Load from environment variable
var jwtSecret = []byte("my-secret-key-change-this-in-production")

func main() {
	// Override secret from environment if provided
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
		log.Println("Using JWT_SECRET from environment")
	} else {
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable in production!")
	}

	// Setup routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/", homeHandler)

	// Protected routes
	mux.Handle("/api/profile", authMiddleware(http.HandlerFunc(profileHandler)))
	mux.Handle("/api/admin", authMiddleware(requireRole("admin")(http.HandlerFunc(adminHandler))))

	// Start server
	addr := ":8080"
	log.Printf("Starting JWT authentication server on %s", addr)
	log.Printf("\nTry these commands:")
	log.Printf("  # Login as Alice")
	log.Printf("  curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{\"username\":\"alice\",\"password\":\"secret123\"}'")
	log.Printf("\n  # Login as Bob (admin)")
	log.Printf("  curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{\"username\":\"bob\",\"password\":\"password456\"}'")
	log.Printf("\n  # Access profile (replace TOKEN with actual token from login)")
	log.Printf("  curl http://localhost:8080/api/profile -H 'Authorization: Bearer TOKEN'")
	log.Printf("\n  # Access admin endpoint (requires admin role)")
	log.Printf("  curl http://localhost:8080/api/admin -H 'Authorization: Bearer TOKEN'\n")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// homeHandler handles the root path
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "JWT Authentication API",
		"endpoints": map[string]string{
			"POST /login":      "Login and get JWT token",
			"GET /api/profile": "Get user profile (requires auth)",
			"GET /api/admin":   "Admin endpoint (requires admin role)",
		},
		"example": map[string]string{
			"login": `curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"alice","password":"secret123"}'`,
			"profile": "curl http://localhost:8080/api/profile -H 'Authorization: Bearer <token>'",
		},
	})
}

// loginHandler handles user login and token generation
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only POST is allowed")
		return
	}

	// Parse credentials
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON")
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "invalid_request", "Username and password required")
		return
	}

	// Find user
	user, exists := users[req.Username]
	if !exists {
		respondWithError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid username or password")
		return
	}

	// Verify password
	// In production: use bcrypt.CompareHashAndPassword(user.Password, req.Password)
	if user.Password != req.Password {
		respondWithError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid username or password")
		return
	}

	// Generate token
	expiresIn := 24 * time.Hour
	token, err := generateToken(user, expiresIn)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "server_error", "Failed to generate token")
		return
	}

	// Log successful login
	log.Printf("User '%s' logged in successfully", user.Username)

	// Return token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(expiresIn.Seconds()),
	})
}

// profileHandler returns the authenticated user's profile
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract claims from context (added by authMiddleware)
	claims, err := getClaimsFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}

	// Return user profile
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"roles":    claims.Roles,
		"message":  fmt.Sprintf("Hello, %s!", claims.Username),
	})
}

// adminHandler handles admin-only operations
func adminHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Welcome to admin panel, %s!", claims.Username),
		"stats": map[string]int{
			"total_users": len(users),
		},
	})
}

// generateToken creates a signed JWT token for the given user
func generateToken(user *User, expiresIn time.Duration) (string, error) {
	now := time.Now()

	// Create claims
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

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// validateToken validates a JWT token and returns the claims
func validateToken(tokenString string) (*Claims, error) {
	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method to prevent algorithm confusion attack
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// authMiddleware is a middleware that validates JWT tokens
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "missing_token", "Authorization header required")
			return
		}

		// Parse "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "invalid_token", "Authorization header format must be 'Bearer {token}'")
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			respondWithError(w, http.StatusUnauthorized, "invalid_token", "Token is invalid or expired")
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), "claims", claims)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireRole returns a middleware that checks if user has the required role
func requireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := getClaimsFromContext(r)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "unauthorized", err.Error())
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
				respondWithError(w, http.StatusForbidden, "forbidden", fmt.Sprintf("Requires '%s' role", role))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClaimsFromContext extracts claims from request context
func getClaimsFromContext(r *http.Request) (*Claims, error) {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		return nil, fmt.Errorf("no claims in context")
	}
	return claims, nil
}

// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, code int, err, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:       err,
		Description: desc,
	})
}
