package exercise

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var testSecret = []byte("test-secret-key")

func TestGenerateToken_Success(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "alice",
		Roles:    []string{"user"},
	}

	token, err := GenerateToken(user, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Token should have 3 parts (header.payload.signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Expected 3 parts in JWT, got %d", len(parts))
	}
}

func TestGenerateToken_ContainsClaims(t *testing.T) {
	user := &User{
		ID:       42,
		Username: "testuser",
		Roles:    []string{"admin", "user"},
	}

	token, err := GenerateToken(user, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Parse token to verify claims
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return testSecret, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse generated token: %v", err)
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID=%d, got %d", user.ID, claims.UserID)
	}
	if claims.Username != user.Username {
		t.Errorf("Expected Username=%s, got %s", user.Username, claims.Username)
	}
	if len(claims.Roles) != len(user.Roles) {
		t.Errorf("Expected %d roles, got %d", len(user.Roles), len(claims.Roles))
	}
}

func TestGenerateToken_Expiration(t *testing.T) {
	user := &User{ID: 1, Username: "alice"}
	expiresIn := time.Second

	token, err := GenerateToken(user, testSecret, expiresIn)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Token should be valid immediately
	_, err = ValidateToken(token, testSecret)
	if err != nil {
		t.Errorf("Token should be valid immediately: %v", err)
	}

	// Wait for expiration
	time.Sleep(expiresIn + 100*time.Millisecond)

	// Token should be expired
	_, err = ValidateToken(token, testSecret)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestValidateToken_Success(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "alice",
		Roles:    []string{"user"},
	}

	token, err := GenerateToken(user, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID=%d, got %d", user.ID, claims.UserID)
	}
	if claims.Username != user.Username {
		t.Errorf("Expected Username=%s, got %s", user.Username, claims.Username)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	invalidTokens := []struct {
		name  string
		token string
	}{
		{"empty", ""},
		{"malformed", "not.a.valid.jwt"},
		{"random", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"},
	}

	for _, tt := range invalidTokens {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateToken(tt.token, testSecret)
			if err == nil {
				t.Error("Expected error for invalid token")
			}
		})
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	user := &User{ID: 1, Username: "alice"}

	token, err := GenerateToken(user, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	wrongSecret := []byte("wrong-secret")
	_, err = ValidateToken(token, wrongSecret)
	if err == nil {
		t.Error("Expected error when validating with wrong secret")
	}
}

func TestValidateToken_AlgorithmConfusion(t *testing.T) {
	// Create a token with a different algorithm
	claims := &Claims{
		UserID:   1,
		Username: "alice",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	// Try to create token with HS512 instead of HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(testSecret)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// ValidateToken should reject tokens with wrong algorithm
	_, err = ValidateToken(tokenString, testSecret)
	if err == nil {
		t.Error("Expected error for token with wrong signing algorithm")
	}
}

func TestAuthMiddleware_Success(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "alice",
		Roles:    []string{"user"},
	}

	token, err := GenerateToken(user, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Create a test handler that checks for claims in context
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetClaims(r)
		if err != nil {
			t.Errorf("Expected claims in context, got error: %v", err)
		}
		if claims.UserID != user.ID {
			t.Errorf("Expected UserID=%d, got %d", user.ID, claims.UserID)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "success")
	})

	// Wrap with middleware
	wrappedHandler := AuthMiddleware(testSecret)(handler)

	// Create request with valid token
	req := httptest.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without token")
	})

	wrappedHandler := AuthMiddleware(testSecret)(handler)

	req := httptest.NewRequest("GET", "/api/profile", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	testCases := []struct {
		name   string
		header string
	}{
		{"no bearer prefix", "token123"},
		{"wrong prefix", "Basic token123"},
		{"missing token", "Bearer"},
		{"extra spaces", "Bearer  token1  token2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Handler should not be called with invalid auth format")
			})

			wrappedHandler := AuthMiddleware(testSecret)(handler)

			req := httptest.NewRequest("GET", "/api/profile", nil)
			req.Header.Set("Authorization", tc.header)

			rec := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("Expected status 401, got %d", rec.Code)
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
	})

	wrappedHandler := AuthMiddleware(testSecret)(handler)

	req := httptest.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	user := &User{ID: 1, Username: "alice"}

	// Create token that expires immediately
	token, err := GenerateToken(user, testSecret, time.Millisecond)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with expired token")
	})

	wrappedHandler := AuthMiddleware(testSecret)(handler)

	req := httptest.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestRequireRole_Success(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "alice",
		Roles:    []string{"user", "admin"},
	}

	token, _ := GenerateToken(user, testSecret, time.Hour)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "admin access granted")
	})

	// Chain auth middleware and role middleware
	wrappedHandler := AuthMiddleware(testSecret)(RequireRole("admin")(handler))

	req := httptest.NewRequest("GET", "/api/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestRequireRole_MissingRole(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "alice",
		Roles:    []string{"user"}, // No admin role
	}

	token, _ := GenerateToken(user, testSecret, time.Hour)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without required role")
	})

	wrappedHandler := AuthMiddleware(testSecret)(RequireRole("admin")(handler))

	req := httptest.NewRequest("GET", "/api/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}

func TestGetClaims_Success(t *testing.T) {
	claims := &Claims{
		UserID:   42,
		Username: "testuser",
	}

	ctx := context.WithValue(context.Background(), "claims", claims)
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

	extractedClaims, err := GetClaims(req)
	if err != nil {
		t.Errorf("GetClaims failed: %v", err)
	}

	if extractedClaims.UserID != claims.UserID {
		t.Errorf("Expected UserID=%d, got %d", claims.UserID, extractedClaims.UserID)
	}
}

func TestGetClaims_NoClaims(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	_, err := GetClaims(req)
	if err == nil {
		t.Error("Expected error when no claims in context")
	}
}

// Benchmark tests
func BenchmarkGenerateToken(b *testing.B) {
	user := &User{
		ID:       1,
		Username: "alice",
		Roles:    []string{"user"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateToken(user, testSecret, time.Hour)
	}
}

func BenchmarkValidateToken(b *testing.B) {
	user := &User{ID: 1, Username: "alice"}
	token, _ := GenerateToken(user, testSecret, time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidateToken(token, testSecret)
	}
}

func BenchmarkAuthMiddleware(b *testing.B) {
	user := &User{ID: 1, Username: "alice"}
	token, _ := GenerateToken(user, testSecret, time.Hour)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := AuthMiddleware(testSecret)(handler)

	req := httptest.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
	}
}
