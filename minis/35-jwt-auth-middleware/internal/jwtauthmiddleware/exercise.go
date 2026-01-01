//go:build !solution && !reference


package jwtauthmiddleware

// TODO: Import required packages
// You'll need:
// - "context" for passing data through request context
// - "fmt" for error formatting
// - "net/http" for HTTP handlers and middleware
// - "strings" for parsing "Bearer token" format
// - "time" for token expiration
// - "github.com/golang-jwt/jwt/v5" for JWT operations
//
// import (
//     "context"
//     "fmt"
//     "net/http"
//     "strings"
//     "time"
//     "github.com/golang-jwt/jwt/v5"
// )

// ============================================================================
// JWT AUTHENTICATION MIDDLEWARE: Stateless Authentication
// ============================================================================
//
// What is JWT (JSON Web Token)?
// - Self-contained token that carries user information
// - Signed using a secret key (HMAC) or public/private key pair (RSA)
// - Three parts: Header.Payload.Signature (base64 encoded)
//
// Example JWT:
//   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.       ← Header
//   eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFsaWNlIn0. ← Payload (claims)
//   SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c    ← Signature
//
// Why JWT for authentication?
// - Stateless: Server doesn't need to store sessions in database
// - Scalable: Works across multiple servers (no shared session store)
// - Standard: Widely supported by libraries and services
// - Self-contained: Token includes all user information needed
//
// Security considerations:
// - NEVER store sensitive data in JWT (it's base64, not encrypted)
// - Always verify signature (prevents tampering)
// - Set short expiration times (15 minutes for access tokens)
// - Use HTTPS (prevents token interception)
// - Protect against algorithm confusion attacks (verify signing method)
//
// JWT structure:
// - Header: {"alg": "HS256", "typ": "JWT"}
// - Payload: {"user_id": 1, "exp": 1234567890}
// - Signature: HMACSHA256(base64(header) + "." + base64(payload), secret)
//
// Memory management:
// - Tokens are strings (immutable, stored on heap if > 32 bytes)
// - Claims are structs (stack allocated if small, heap if large)
// - Context values are interface{} (always heap allocated)
//
// ============================================================================

// ============================================================================
// Exercise 1: Define User and Claims Types
// ============================================================================

// User represents a user in the system.
// TODO: This is already defined for you.
//
// In production, this would come from a database.
// For this exercise, we use in-memory users for simplicity.
//
// Password field:
// - In production: Use bcrypt.GenerateFromPassword()
// - Never store plain text passwords!
// - bcrypt automatically handles salting and multiple rounds
//
// Roles field:
// - Used for authorization (who can do what)
// - Example roles: "admin", "user", "moderator"
// - Checked by RequireRole middleware
//
type User struct {
	ID       int      // Unique user ID
	Username string   // Username for login
	Password string   // In production: bcrypt hash, not plain text!
	Roles    []string // User roles for authorization
}

// Claims represents the JWT claims.
// TODO: This is already defined for you.
//
// Embeds jwt.RegisteredClaims which provides:
// - ExpiresAt: When token expires (time.Time)
// - IssuedAt: When token was created
// - NotBefore: Token not valid before this time
// - Issuer: Who issued the token
// - Subject: Who the token is about
//
// Custom fields:
// - UserID: Our custom field (not in JWT standard)
// - Username: Our custom field
// - Roles: Our custom field (for authorization)
//
// JSON tags:
// - Tells JWT library how to marshal/unmarshal
// - "user_id" becomes the JSON key in the token payload
//
// Memory:
// - UserID: 8 bytes (int on 64-bit system)
// - Username: 16 bytes (string header) + string length
// - Roles: 24 bytes (slice header) + slice contents
// - RegisteredClaims: ~80 bytes
// - Total: ~128 bytes + variable data
//
type Claims struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// ============================================================================
// Exercise 2: Implement Token Generation
// ============================================================================

// GenerateToken creates a signed JWT token for the given user.
// TODO: Implement token generation with proper claims and signing.
//
// The token should:
// - Include user ID, username, and roles in claims
// - Set expiration time based on expiresIn parameter
// - Set issued at time to current time
// - Use HMAC-SHA256 signing algorithm (jwt.SigningMethodHS256)
// - Sign with the provided secret
//
// Algorithm steps:
// 1. Create Claims struct with user data and timestamps
// 2. Create token with jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 3. Sign token with token.SignedString(secret)
// 4. Return signed token string
//
// Parameters:
// - user: User to generate token for
// - secret: Secret key for signing (keep this safe!)
// - expiresIn: Token expiration duration (e.g., 15*time.Minute)
//
// Returns:
// - string: Signed JWT token (three parts separated by dots)
// - error: Non-nil if token generation fails
//
// Example usage:
//   user := &User{ID: 1, Username: "alice", Roles: []string{"user"}}
//   token, err := GenerateToken(user, []byte("secret"), 15*time.Minute)
//
// Security notes:
// - Secret must be strong (at least 32 bytes random data)
// - Use environment variables for secrets (never hardcode)
// - Rotate secrets periodically
//
// func GenerateToken(user *User, secret []byte, expiresIn time.Duration) (string, error) {
//     now := time.Now()
//
//     // Create claims with user information
//     claims := &Claims{
//         UserID:   user.ID,
//         Username: user.Username,
//         Roles:    user.Roles,
//         RegisteredClaims: jwt.RegisteredClaims{
//             ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
//             IssuedAt:  jwt.NewNumericDate(now),
//             NotBefore: jwt.NewNumericDate(now),
//             Issuer:    "jwt-auth-server",  // Who issued this token
//             Subject:   fmt.Sprintf("%d", user.ID),  // Who token is about
//         },
//     }
//
//     // Create token with HS256 signing method
//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
//     // Sign and get the complete encoded token as a string
//     tokenString, err := token.SignedString(secret)
//     if err != nil {
//         return "", fmt.Errorf("failed to sign token: %w", err)
//     }
//
//     return tokenString, nil
// }

// ============================================================================
// Exercise 3: Implement Token Validation
// ============================================================================

// ValidateToken validates a JWT token and returns the claims.
// TODO: Implement token validation with signature verification.
//
// The function should:
// - Parse the token string
// - Verify the signature using the provided secret
// - Check that the signing method is HMAC-SHA256 (prevent algorithm confusion)
// - Verify token hasn't expired
// - Return the parsed claims
//
// Algorithm steps:
// 1. Call jwt.ParseWithClaims with a key function
// 2. In key function: Verify signing method is HMAC (critical security check!)
// 3. Return secret for signature verification
// 4. Extract and return claims if token is valid
//
// Parameters:
// - tokenString: JWT token to validate
// - secret: Secret key for verification (must match signing secret)
//
// Returns:
// - *Claims: Parsed claims if valid
// - error: Non-nil if token is invalid, expired, or signature doesn't match
//
// Algorithm confusion attack:
// - Attacker changes "alg" from "HS256" to "none"
// - Token validation accepts unsigned tokens
// - Prevention: Always verify signing method!
//
// Example attack (without method check):
//   Header: {"alg": "none", "typ": "JWT"}
//   Payload: {"user_id": 1, "roles": ["admin"]}
//   Signature: (empty)
//   → Server accepts it as valid! (VERY BAD)
//
// With method check:
//   → Server rejects because method is "none", not "HS256"
//
// func ValidateToken(tokenString string, secret []byte) (*Claims, error) {
//     // Parse token with claims
//     token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//         // CRITICAL SECURITY CHECK: Verify signing method
//         // This prevents algorithm confusion attacks
//         // Without this check, an attacker could change the algorithm
//         // from RS256 to HS256 and sign with the public key
//         if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//             return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//         }
//         return secret, nil
//     })
//
//     if err != nil {
//         return nil, fmt.Errorf("failed to parse token: %w", err)
//     }
//
//     // Extract and validate claims
//     if claims, ok := token.Claims.(*Claims); ok && token.Valid {
//         return claims, nil
//     }
//
//     return nil, fmt.Errorf("invalid token")
// }

// ============================================================================
// Exercise 4: Implement Authentication Middleware
// ============================================================================

// AuthMiddleware returns a middleware that validates JWT tokens.
// TODO: Implement authentication middleware.
//
// The middleware should:
// - Extract the Authorization header
// - Check for "Bearer <token>" format
// - Validate the token using ValidateToken
// - Add claims to request context with key "claims"
// - Call next handler if valid
// - Return 401 Unauthorized if token is missing or invalid
//
// Algorithm:
// 1. Get Authorization header from request
// 2. Check if header exists (if not, return 401)
// 3. Parse "Bearer <token>" format (split by space)
// 4. Validate format (must be 2 parts, first must be "Bearer")
// 5. Extract token (second part)
// 6. Validate token with ValidateToken
// 7. If valid: Add claims to context and call next handler
// 8. If invalid: Return 401 with error message
//
// Parameters:
// - secret: Secret key for token validation
//
// Returns:
// - Middleware function that wraps http.Handler
//
// HTTP status codes:
// - 401 Unauthorized: Missing or invalid token
// - 403 Forbidden: Valid token but insufficient permissions (see RequireRole)
//
// Context usage:
// - context.WithValue creates a new context with added value
// - Downstream handlers can access claims with GetClaims(r)
// - Context values are type interface{} (need type assertion)
//
// Example request:
//   GET /api/users
//   Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
//
// func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             // Extract Authorization header
//             authHeader := r.Header.Get("Authorization")
//             if authHeader == "" {
//                 http.Error(w, "missing authorization header", http.StatusUnauthorized)
//                 return
//             }
//
//             // Parse "Bearer <token>" format
//             // Expected format: "Bearer eyJhbGciOi..."
//             parts := strings.Split(authHeader, " ")
//             if len(parts) != 2 || parts[0] != "Bearer" {
//                 http.Error(w, "invalid authorization format, expected 'Bearer {token}'", http.StatusUnauthorized)
//                 return
//             }
//
//             tokenString := parts[1]
//
//             // Validate token
//             claims, err := ValidateToken(tokenString, secret)
//             if err != nil {
//                 http.Error(w, fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized)
//                 return
//             }
//
//             // Add claims to request context for use by downstream handlers
//             ctx := context.WithValue(r.Context(), "claims", claims)
//
//             // Call next handler with updated context
//             next.ServeHTTP(w, r.WithContext(ctx))
//         })
//     }
// }

// ============================================================================
// Exercise 5: Implement Role-Based Access Control
// ============================================================================

// RequireRole returns a middleware that checks if the user has the required role.
// TODO: Implement role-based access control middleware.
//
// The middleware should:
// - Extract claims from request context (added by AuthMiddleware)
// - Check if user has the required role
// - Call next handler if user has the role
// - Return 403 Forbidden if user doesn't have the role
//
// Algorithm:
// 1. Get claims from context using GetClaims(r)
// 2. If no claims found, return 401 (shouldn't happen if AuthMiddleware ran first)
// 3. Iterate through user's roles
// 4. If required role found, call next handler
// 5. If not found, return 403 Forbidden
//
// Parameters:
// - role: Required role (e.g., "admin", "moderator")
//
// Returns:
// - Middleware function that wraps http.Handler
//
// HTTP status codes:
// - 401 Unauthorized: No claims in context (AuthMiddleware didn't run)
// - 403 Forbidden: Valid token but user doesn't have required role
//
// Example usage:
//   Chain handlers:
//   1. AuthMiddleware (validates token, adds claims to context)
//   2. RequireRole("admin") (checks if user is admin)
//   3. Handler (executes if user is authenticated AND is admin)
//
// Example:
//   http.Handle("/admin", AuthMiddleware(secret)(RequireRole("admin")(adminHandler)))
//
// func RequireRole(role string) func(http.Handler) http.Handler {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             // Extract claims from context (should be added by AuthMiddleware)
//             claims, err := GetClaims(r)
//             if err != nil {
//                 http.Error(w, "unauthorized: no claims in context", http.StatusUnauthorized)
//                 return
//             }
//
//             // Check if user has required role
//             hasRole := false
//             for _, userRole := range claims.Roles {
//                 if userRole == role {
//                     hasRole = true
//                     break
//                 }
//             }
//
//             if !hasRole {
//                 http.Error(w, fmt.Sprintf("forbidden: requires '%s' role", role), http.StatusForbidden)
//                 return
//             }
//
//             // User has required role, proceed
//             next.ServeHTTP(w, r)
//         })
//     }
// }

// ============================================================================
// Exercise 6: Implement Claims Extraction Helper
// ============================================================================

// GetClaims extracts claims from the request context.
// TODO: Implement claims extraction from context.
//
// Helper function to extract claims added by AuthMiddleware.
//
// Algorithm:
// 1. Get value from context using r.Context().Value("claims")
// 2. Check if value is nil (no claims in context)
// 3. Type assert to *Claims
// 4. Return claims or error
//
// Parameters:
// - r: HTTP request (contains context)
//
// Returns:
// - *Claims: Claims from context
// - error: Non-nil if no claims in context or wrong type
//
// Type assertion:
// - value.(type) checks if value is of given type
// - Returns (value, true) if correct type, (zero, false) if not
// - Use two-value form to avoid panic
//
// Example usage in handler:
//   func myHandler(w http.ResponseWriter, r *http.Request) {
//       claims, err := GetClaims(r)
//       if err != nil {
//           // Shouldn't happen if AuthMiddleware ran
//           http.Error(w, "unauthorized", http.StatusUnauthorized)
//           return
//       }
//       fmt.Fprintf(w, "Hello, %s (ID: %d)!", claims.Username, claims.UserID)
//   }
//
// func GetClaims(r *http.Request) (*Claims, error) {
//     // Extract value from context
//     value := r.Context().Value("claims")
//     if value == nil {
//         return nil, fmt.Errorf("no claims in context")
//     }
//
//     // Type assert to *Claims
//     claims, ok := value.(*Claims)
//     if !ok {
//         return nil, fmt.Errorf("invalid claims type in context")
//     }
//
//     return claims, nil
// }

// ============================================================================
// Exercise 7: Implement Token Refresh (STRETCH GOAL)
// ============================================================================

// RefreshToken generates a new access token using a valid refresh token.
// TODO: Implement refresh token pattern (STRETCH GOAL).
//
// Refresh token pattern:
// - Access tokens: Short-lived (15 minutes)
// - Refresh tokens: Long-lived (7 days)
// - When access token expires, use refresh token to get new access token
// - Don't need to re-enter credentials
//
// This is a simplified implementation. In production:
// - Refresh tokens should be stored in a database with user ID
// - Should support revocation (blacklist or token versioning)
// - Should rotate refresh tokens (issue new refresh token on use)
// - Should have longer expiration than access tokens
// - Should be stored in httpOnly cookies (more secure than localStorage)
//
// Algorithm:
// 1. Validate the refresh token (same as ValidateToken)
// 2. Extract user info from refresh token claims
// 3. Create new User object from claims
// 4. Generate new access token with GenerateToken
// 5. Return new access token
//
// Parameters:
// - refreshTokenString: Valid refresh token
// - secret: Secret key
// - newExpiresIn: Expiration duration for new access token (e.g., 15 minutes)
//
// Returns:
// - string: New access token
// - error: Non-nil if refresh token is invalid
//
// Production considerations:
// - Store refresh tokens in database with:
//   * user_id
//   * token_id (unique identifier)
//   * created_at
//   * expires_at
//   * revoked (boolean)
// - Check if token is revoked before issuing new access token
// - Implement token rotation (new refresh token on each use)
//
// func RefreshToken(refreshTokenString string, secret []byte, newExpiresIn time.Duration) (string, error) {
//     // Validate the refresh token
//     claims, err := ValidateToken(refreshTokenString, secret)
//     if err != nil {
//         return "", fmt.Errorf("invalid refresh token: %w", err)
//     }
//
//     // In production, you would:
//     // 1. Check if refresh token is in database (not revoked)
//     // 2. Verify the token type (should have a "type": "refresh" claim)
//     // 3. Update last used timestamp in database
//     // 4. Optionally rotate the refresh token
//
//     // Create a new user object from claims
//     user := &User{
//         ID:       claims.UserID,
//         Username: claims.Username,
//         Roles:    claims.Roles,
//     }
//
//     // Generate new access token with shorter expiration
//     newAccessToken, err := GenerateToken(user, secret, newExpiresIn)
//     if err != nil {
//         return "", fmt.Errorf("failed to generate new access token: %w", err)
//     }
//
//     return newAccessToken, nil
// }

// ============================================================================
// Testing and Running
// ============================================================================
//
// After implementing all functions:
// - Run: go test -v ./...
// - Check: All tests should pass
//
// To test manually with curl:
//   # Generate token
//   curl -X POST http://localhost:8080/login \
//     -H "Content-Type: application/json" \
//     -d '{"username":"alice","password":"password"}'
//
//   # Use token to access protected endpoint
//   curl http://localhost:8080/protected \
//     -H "Authorization: Bearer <token>"
//
//   # Try to access admin endpoint (should fail if not admin)
//   curl http://localhost:8080/admin \
//     -H "Authorization: Bearer <token>"
//
// Common errors:
// - 401 Unauthorized: Missing or invalid token
// - 403 Forbidden: Valid token but insufficient permissions
// - Invalid signature: Token was tampered with or wrong secret
// - Token expired: Generate a new token
//
// Key concepts learned:
// 1. JWT structure (header, payload, signature)
// 2. HMAC-SHA256 signing
// 3. Algorithm confusion attack prevention
// 4. Stateless authentication
// 5. Context values for passing data through middleware
// 6. Role-based access control (RBAC)
// 7. Token refresh pattern
//
// Compare with solution.go to see detailed implementations and explanations.
