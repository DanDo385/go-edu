# Project 35: JWT Authentication Middleware

## 1. What Is This About?

### Real-World Scenario

You're building a REST API for a web application. Users need to log in and access protected resources:

**❌ Naive approach:** Send username/password with every request
```go
// Every API request includes credentials
GET /api/profile?username=alice&password=secret123
```

**Problems**:
- Credentials exposed in every request (logs, network sniffing)
- Database lookup on every request (slow)
- Hard to revoke access
- Doesn't scale across multiple services

**✅ Modern approach:** JWT (JSON Web Token) authentication
```go
// 1. Login once, get a token
POST /login → {"token": "eyJhbGciOi..."}

// 2. Include token in subsequent requests
GET /api/profile
Authorization: Bearer eyJhbGciOi...
```

**Benefits**:
- ✅ Credentials sent only once during login
- ✅ Stateless (no server-side session storage)
- ✅ Works across multiple services (microservices)
- ✅ Contains user info (claims) - no database lookup
- ✅ Cryptographically signed (tamper-proof)

This project teaches you how to build **production-grade JWT authentication** that is:
- **Secure**: Uses proper signing algorithms (HMAC-SHA256, RSA)
- **Stateless**: No server-side session storage required
- **Flexible**: Supports custom claims, expiration, refresh tokens
- **Middleware-based**: Reusable across multiple endpoints

### What You'll Learn

1. **JWT structure**: Header, payload, signature
2. **Token generation**: Creating signed tokens with claims
3. **Token validation**: Verifying signatures and expiration
4. **Middleware pattern**: Authentication middleware for HTTP handlers
5. **Security best practices**: Secret management, token expiration, HTTPS
6. **Common attacks**: How to prevent token theft, replay attacks, etc.

### The Challenge

Build a JWT authentication system that:
- Generates tokens on successful login
- Validates tokens on protected endpoints
- Supports token expiration and refresh
- Implements middleware for easy authentication
- Handles errors gracefully
- Follows security best practices

---

## 2. First Principles: What is JWT?

### JWT Structure

A JWT is a string with three parts separated by dots (`.`):

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

**Decoded**:
```
HEADER.PAYLOAD.SIGNATURE
```

#### Part 1: Header (Algorithm + Token Type)

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

- `alg`: Signing algorithm (HS256 = HMAC-SHA256, RS256 = RSA with SHA256)
- `typ`: Token type (always "JWT")

**Base64URL encoded**: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`

#### Part 2: Payload (Claims)

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true,
  "iat": 1516239022,
  "exp": 1516242622
}
```

**Standard claims** (optional but recommended):
- `sub` (subject): User ID
- `iat` (issued at): Timestamp when token was created
- `exp` (expiration): Timestamp when token expires
- `nbf` (not before): Token not valid before this time
- `iss` (issuer): Who issued the token
- `aud` (audience): Who the token is intended for

**Custom claims**: Any data you want (user role, permissions, etc.)

**Base64URL encoded**: `eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ`

#### Part 3: Signature (Cryptographic Proof)

```
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret
)
```

**Purpose**: Proves that:
1. Token hasn't been tampered with
2. Token was issued by someone with the secret key

**How it works**:
1. Take `header + "." + payload`
2. Sign with secret key using algorithm from header
3. Base64URL encode the signature

**Result**: `SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`

### Why Base64URL (not Base64)?

**Base64** uses `+`, `/`, and `=` which are problematic in URLs and headers.

**Base64URL** uses `-`, `_`, and no padding (`=`) → safe for URLs and headers.

```
Base64:    "Hello+World/=="
Base64URL: "Hello-World_"
```

### Symmetric vs Asymmetric Signing

#### HMAC (HS256) - Symmetric

**Same secret** for signing and verification.

```go
// Server A signs
signature := HMAC-SHA256(data, "shared-secret")

// Server B verifies
valid := HMAC-SHA256(data, "shared-secret") == signature
```

**Pros**:
- Fast
- Simple

**Cons**:
- Both issuer and verifier need the same secret
- If secret leaks, anyone can create valid tokens

**Use case**: Single service or trusted microservices

#### RSA (RS256) - Asymmetric

**Private key** for signing, **public key** for verification.

```go
// Auth service signs with private key
signature := RSA-SHA256(data, privateKey)

// Other services verify with public key
valid := RSA-Verify(data, signature, publicKey)
```

**Pros**:
- Separate signing and verification
- Public key can be distributed safely
- Even if public key leaks, can't create tokens

**Cons**:
- Slower than HMAC
- More complex setup

**Use case**: Microservices where different services verify tokens

---

## 3. Breaking Down the Solution

### Step 1: Define User and Claims

```go
// User represents a user in the system
type User struct {
    ID       int
    Username string
    Password string  // In production: hashed with bcrypt
}

// Claims represents the JWT claims
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

**Why embed `jwt.RegisteredClaims`?**
- Provides standard fields: `ExpiresAt`, `IssuedAt`, `NotBefore`, `Issuer`, `Subject`, `Audience`
- Handles expiration checking automatically

### Step 2: Generate Token on Login

```go
func generateToken(user *User, secret []byte, expiration time.Duration) (string, error) {
    // Create claims
    claims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign token
    return token.SignedString(secret)
}
```

**Step-by-step**:
1. Create claims with user info + expiration
2. Create unsigned token with claims
3. Sign with secret key → returns JWT string

### Step 3: Validate Token

```go
func validateToken(tokenString string, secret []byte) (*Claims, error) {
    // Parse and validate
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secret, nil
    })

    if err != nil {
        return nil, err
    }

    // Extract claims
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}
```

**What happens during validation**:
1. **Parse**: Decode Base64URL parts
2. **Verify signature**: Recalculate signature with secret, compare
3. **Check expiration**: Ensure `exp` > current time
4. **Check signing method**: Prevent algorithm confusion attacks
5. **Return claims**: If all checks pass

### Step 4: Middleware for Protected Routes

```go
func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing authorization header", http.StatusUnauthorized)
                return
            }

            // Format: "Bearer <token>"
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "invalid authorization format", http.StatusUnauthorized)
                return
            }

            // Validate token
            claims, err := validateToken(parts[1], secret)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }

            // Add claims to request context
            ctx := context.WithValue(r.Context(), "claims", claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**Middleware flow**:
```
Request → Extract token → Validate → Add to context → Call next handler
            ↓               ↓            ↓
        Missing?        Invalid?    Store user info
        Return 401      Return 401   for later use
```

### Step 5: Login Handler

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    // Parse credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    // Find user (in production: query database)
    user := findUser(creds.Username)
    if user == nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    // Verify password (in production: use bcrypt.CompareHashAndPassword)
    if user.Password != creds.Password {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    // Generate token
    token, err := generateToken(user, []byte(secret), 24*time.Hour)
    if err != nil {
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    // Return token
    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
    })
}
```

### Step 6: Protected Handler

```go
func profileHandler(w http.ResponseWriter, r *http.Request) {
    // Extract claims from context
    claims, ok := r.Context().Value("claims").(*Claims)
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // Use claims
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user_id":  claims.UserID,
        "username": claims.Username,
    })
}
```

---

## 4. Complete Solution Walkthrough

### Full Flow Example

**1. User logs in**:
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123"}'
```

**Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFsaWNlIiwiZXhwIjoxNzAwMDAwMDAwfQ.xyz"
}
```

**2. User accesses protected resource**:
```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer eyJhbGciOi..."
```

**Response**:
```json
{
  "user_id": 1,
  "username": "alice"
}
```

**3. Invalid token**:
```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer invalid-token"
```

**Response**:
```
HTTP/1.1 401 Unauthorized
invalid token
```

### Request Flow Diagram

```
┌─────────┐                                  ┌─────────┐
│ Client  │                                  │ Server  │
└────┬────┘                                  └────┬────┘
     │                                            │
     │ POST /login {username, password}           │
     │───────────────────────────────────────────>│
     │                                            │
     │                                  Verify credentials
     │                                  Generate JWT
     │                                  Sign with secret
     │                                            │
     │  {"token": "eyJhbGci..."}                  │
     │<───────────────────────────────────────────│
     │                                            │
     │ GET /api/profile                           │
     │ Authorization: Bearer eyJhbGci...          │
     │───────────────────────────────────────────>│
     │                                            │
     │                                  Extract token
     │                                  Validate signature
     │                                  Check expiration
     │                                  Extract claims
     │                                            │
     │  {"user_id": 1, "username": "alice"}       │
     │<───────────────────────────────────────────│
     │                                            │
```

---

## 5. Key Concepts Explained

### Concept 1: Token Expiration

**Why tokens should expire**:
- Limits damage from token theft
- Forces re-authentication
- Allows revoking access

**Access Token vs Refresh Token**:

**Access Token** (short-lived: 15 minutes - 1 hour):
- Used for API requests
- Short expiration
- No revocation (stateless)

**Refresh Token** (long-lived: 7 days - 30 days):
- Used to get new access tokens
- Stored securely (httpOnly cookie)
- Can be revoked (stored in database)

**Flow**:
```
1. Login → Get access token (1 hour) + refresh token (7 days)
2. Use access token for API requests
3. When access token expires → Use refresh token to get new access token
4. When refresh token expires → Login again
```

### Concept 2: Token Storage

**Where to store JWTs**:

**❌ localStorage**:
- Vulnerable to XSS (JavaScript can access)
- Any script can steal token

**❌ sessionStorage**:
- Same XSS vulnerability as localStorage

**✅ httpOnly Cookie**:
- Not accessible via JavaScript
- Protected from XSS
- Automatically sent with requests

**✅ In-memory (JavaScript variable)**:
- Safe from XSS (if not exposed globally)
- Lost on page refresh (need refresh token)

**Best practice**:
- Access token: In-memory or httpOnly cookie
- Refresh token: httpOnly cookie with Secure and SameSite flags

### Concept 3: Security Best Practices

#### 1. Use Strong Secrets

**❌ Weak**:
```go
secret := []byte("secret123")
```

**✅ Strong**:
```go
// Generate random 256-bit key
secret := make([]byte, 32)
rand.Read(secret)

// Or use environment variable
secret := []byte(os.Getenv("JWT_SECRET"))
```

#### 2. Validate Signing Method (Algorithm Confusion Attack)

**Attack**: Attacker changes algorithm from RS256 to HS256
- Server expects RS256 (asymmetric)
- Attacker signs with public key using HS256 (symmetric)
- Server verifies with public key → succeeds! (wrong)

**Prevention**:
```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Explicitly check signing method
    if token.Method.Alg() != "HS256" {
        return nil, fmt.Errorf("unexpected signing method")
    }
    return secret, nil
})
```

#### 3. Always Use HTTPS

**Why**: JWTs are not encrypted, only signed.

Without HTTPS:
- Token visible in network traffic
- Man-in-the-middle can steal token
- Attacker can replay token

With HTTPS:
- Token encrypted in transit
- Protected from network sniffing

#### 4. Include Audience and Issuer Claims

**Why**: Prevents token reuse across different services.

```go
claims := &Claims{
    RegisteredClaims: jwt.RegisteredClaims{
        Issuer:    "auth.example.com",
        Audience:  jwt.ClaimStrings{"api.example.com"},
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
    },
}
```

#### 5. Implement Token Revocation (For Critical Operations)

JWTs are stateless → can't revoke them directly.

**Solutions**:

**Blacklist approach**:
```go
// Store revoked tokens in Redis with expiration
func revokeToken(tokenID string, expiresAt time.Time) {
    redis.Set("revoked:"+tokenID, "1", time.Until(expiresAt))
}

func isRevoked(tokenID string) bool {
    return redis.Exists("revoked:" + tokenID)
}
```

**Version approach**:
```go
// Add version to claims
type Claims struct {
    UserID  int `json:"user_id"`
    Version int `json:"version"`  // Increment on logout/password change
    jwt.RegisteredClaims
}

// Check version during validation
if claims.Version != user.TokenVersion {
    return errors.New("token revoked")
}
```

### Concept 4: Refresh Token Pattern

```go
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

func generateTokenPair(user *User) (*TokenPair, error) {
    // Short-lived access token (15 minutes)
    accessToken, err := generateToken(user, secret, 15*time.Minute)
    if err != nil {
        return nil, err
    }

    // Long-lived refresh token (7 days)
    refreshToken, err := generateToken(user, secret, 7*24*time.Hour)
    if err != nil {
        return nil, err
    }

    // Store refresh token in database for revocation
    storeRefreshToken(user.ID, refreshToken)

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
    // Extract refresh token
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    // Validate refresh token
    claims, err := validateToken(req.RefreshToken, secret)
    if err != nil {
        http.Error(w, "invalid refresh token", http.StatusUnauthorized)
        return
    }

    // Check if revoked
    if isRefreshTokenRevoked(claims.UserID, req.RefreshToken) {
        http.Error(w, "token revoked", http.StatusUnauthorized)
        return
    }

    // Generate new access token
    user := findUserByID(claims.UserID)
    newAccessToken, _ := generateToken(user, secret, 15*time.Minute)

    json.NewEncoder(w).Encode(map[string]string{
        "access_token": newAccessToken,
    })
}
```

### Concept 5: Role-Based Access Control (RBAC)

**Add roles to claims**:
```go
type Claims struct {
    UserID   int      `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}
```

**Middleware to check roles**:
```go
func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := r.Context().Value("claims").(*Claims)

            // Check if user has required role
            hasRole := false
            for _, userRole := range claims.Roles {
                if userRole == role {
                    hasRole = true
                    break
                }
            }

            if !hasRole {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

**Usage**:
```go
// Public route
http.Handle("/", publicHandler)

// Authenticated route
http.Handle("/api/profile", AuthMiddleware(secret)(profileHandler))

// Admin-only route
http.Handle("/api/admin",
    AuthMiddleware(secret)(
        RequireRole("admin")(adminHandler),
    ),
)
```

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Middleware Chaining

```go
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
    return func(final http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            final = middlewares[i](final)
        }
        return final
    }
}

// Usage
http.Handle("/api/admin",
    Chain(
        LoggingMiddleware,
        AuthMiddleware(secret),
        RequireRole("admin"),
    )(adminHandler),
)
```

### Pattern 2: Extract Claims Helper

```go
func GetClaims(r *http.Request) (*Claims, error) {
    claims, ok := r.Context().Value("claims").(*Claims)
    if !ok {
        return nil, errors.New("no claims in context")
    }
    return claims, nil
}

// Usage in handlers
func profileHandler(w http.ResponseWriter, r *http.Request) {
    claims, err := GetClaims(r)
    if err != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    // Use claims...
}
```

### Pattern 3: Token Response Format

```go
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
}

func respondWithToken(w http.ResponseWriter, accessToken, refreshToken string, expiresIn time.Duration) {
    response := TokenResponse{
        AccessToken:  accessToken,
        TokenType:    "Bearer",
        ExpiresIn:    int(expiresIn.Seconds()),
        RefreshToken: refreshToken,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### Pattern 4: Error Response Format

```go
type ErrorResponse struct {
    Error       string `json:"error"`
    Description string `json:"error_description,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, err, desc string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error:       err,
        Description: desc,
    })
}

// Usage
respondWithError(w, http.StatusUnauthorized, "invalid_token", "The access token is expired")
```

### Pattern 5: CORS Middleware (For SPAs)

```go
func CORSMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            w.Header().Set("Access-Control-Allow-Credentials", "true")

            // Handle preflight
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 7. Real-World Applications

### Single Page Applications (SPAs)

**Use case**: React/Vue/Angular app + Go API

```go
// CORS for SPA
http.Handle("/", CORSMiddleware("http://localhost:3000")(router))

// Login returns JWT
http.HandleFunc("/api/login", loginHandler)

// Protected API routes
http.Handle("/api/profile", AuthMiddleware(secret)(profileHandler))
http.Handle("/api/todos", AuthMiddleware(secret)(todosHandler))
```

**Frontend (JavaScript)**:
```javascript
// Login
const response = await fetch('http://localhost:8080/api/login', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({username: 'alice', password: 'secret123'})
});
const {token} = await response.json();

// Store in memory (or httpOnly cookie)
let accessToken = token;

// Use for API calls
const profile = await fetch('http://localhost:8080/api/profile', {
    headers: {'Authorization': `Bearer ${accessToken}`}
});
```

Companies using this: Almost every modern web app

### Microservices Authentication

**Use case**: Multiple services need to verify tokens

```go
// Shared public key (for RS256)
publicKey := loadPublicKey()

// Each service can verify tokens
func verifyTokenRS256(tokenString string, publicKey *rsa.PublicKey) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected method: %v", token.Header["alg"])
        }
        return publicKey, nil
    })
    // ...
}
```

Companies: Netflix, Uber, Amazon (microservices with JWT)

### Mobile API Authentication

**Use case**: iOS/Android app + Go backend

```swift
// iOS (Swift)
struct LoginResponse: Codable {
    let token: String
}

func login(username: String, password: String) async throws -> String {
    let url = URL(string: "https://api.example.com/login")!
    var request = URLRequest(url: url)
    request.httpMethod = "POST"
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")

    let body = ["username": username, "password": password]
    request.httpBody = try JSONEncoder().encode(body)

    let (data, _) = try await URLSession.shared.data(for: request)
    let response = try JSONDecoder().decode(LoginResponse.self, from: data)

    // Store in Keychain (secure storage)
    Keychain.save(response.token, forKey: "access_token")

    return response.token
}

func fetchProfile() async throws -> Profile {
    let url = URL(string: "https://api.example.com/api/profile")!
    var request = URLRequest(url: url)

    // Get token from Keychain
    let token = Keychain.load(forKey: "access_token")
    request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

    let (data, _) = try await URLSession.shared.data(for: request)
    return try JSONDecoder().decode(Profile.self, from: data)
}
```

Companies: Instagram, Twitter, Spotify (mobile apps with JWT)

### Third-Party API Integration

**Use case**: Provide API for external developers

```go
// Developer generates API key → Gets JWT for their apps

type APICredentials struct {
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
    var creds APICredentials
    json.NewDecoder(r.Body).Decode(&creds)

    // Verify credentials
    client := findClient(creds.ClientID)
    if client == nil || client.Secret != creds.ClientSecret {
        respondWithError(w, http.StatusUnauthorized, "invalid_client", "")
        return
    }

    // Generate token with API scopes
    claims := &Claims{
        UserID: client.ID,
        Scopes: client.Scopes,  // e.g., ["read:users", "write:posts"]
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            Audience:  jwt.ClaimStrings{"api.example.com"},
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString(secret)

    respondWithToken(w, tokenString, "", 24*time.Hour)
}
```

Companies: GitHub, Stripe, Twilio (provide APIs with JWT/OAuth)

---

## 8. Common Mistakes to Avoid

### Mistake 1: Storing Secrets in Code

**❌ Wrong**:
```go
var secret = []byte("my-secret-key-123")
```

**Problem**: Exposed in version control, can't rotate easily.

**✅ Correct**:
```go
secret := []byte(os.Getenv("JWT_SECRET"))
if len(secret) == 0 {
    log.Fatal("JWT_SECRET environment variable not set")
}
```

### Mistake 2: No Expiration

**❌ Wrong**:
```go
claims := &Claims{
    UserID: user.ID,
    // No expiration!
}
```

**Problem**: Token valid forever, can't revoke.

**✅ Correct**:
```go
claims := &Claims{
    UserID: user.ID,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    },
}
```

### Mistake 3: Not Validating Signing Method

**❌ Wrong**:
```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return secret, nil  // No algorithm check!
})
```

**Problem**: Algorithm confusion attack.

**✅ Correct**:
```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method")
    }
    return secret, nil
})
```

### Mistake 4: Storing Sensitive Data in Claims

**❌ Wrong**:
```go
claims := &Claims{
    UserID:   user.ID,
    Password: user.Password,  // Never!
    SSN:      user.SSN,       // Never!
}
```

**Problem**: JWTs are **not encrypted**, only **encoded**. Anyone can decode and read claims.

**✅ Correct**:
```go
claims := &Claims{
    UserID:   user.ID,
    Username: user.Username,
    Roles:    user.Roles,
}
```

### Mistake 5: Not Using HTTPS

**❌ Wrong**:
```go
http.ListenAndServe(":8080", handler)
```

**Problem**: Tokens transmitted in plain text → easily stolen.

**✅ Correct**:
```go
http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", handler)
```

### Mistake 6: Accepting Tokens from Query Parameters

**❌ Wrong**:
```
GET /api/profile?token=eyJhbGci...
```

**Problem**: Tokens logged in server logs, browser history, referrer headers.

**✅ Correct**:
```
GET /api/profile
Authorization: Bearer eyJhbGci...
```

### Mistake 7: Long-Lived Tokens Without Refresh

**❌ Wrong**:
```go
// 30-day access token, no refresh token
claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour))
```

**Problem**: If stolen, attacker has access for 30 days.

**✅ Correct**:
```go
// Short-lived access token (15 min) + refresh token (7 days)
accessToken := generateToken(user, 15*time.Minute)
refreshToken := generateRefreshToken(user, 7*24*time.Hour)
```

---

## 9. Stretch Goals

### Goal 1: Implement Refresh Token Pattern ⭐⭐

Add refresh token functionality with token rotation.

**Hint**: Store refresh tokens in database, rotate on use.

### Goal 2: Add Role-Based Access Control (RBAC) ⭐⭐

Implement middleware to check user roles/permissions.

**Hint**: Add `roles` field to claims, create `RequireRole` middleware.

### Goal 3: Implement Token Blacklist ⭐⭐⭐

Add ability to revoke tokens before expiration.

**Hint**: Use Redis with token ID as key, expiration time as TTL.

### Goal 4: Support RSA (RS256) Signing ⭐⭐⭐

Implement asymmetric key signing for microservices.

**Hint**: Generate RSA key pair, use `jwt.SigningMethodRS256`.

### Goal 5: Add Rate Limiting to Login ⭐⭐

Prevent brute force attacks on login endpoint.

**Hint**: Track login attempts per IP, use sliding window.

### Goal 6: Implement OAuth2 Flow ⭐⭐⭐⭐

Add OAuth2 authorization code flow for third-party integrations.

**Hint**: Implement authorization endpoint, token endpoint, consent screen.

---

## How to Run

```bash
# Get dependencies
go get github.com/golang-jwt/jwt/v5

# Run the server
go run ./minis/35-jwt-auth-middleware/cmd/jwt-server/main.go

# Test login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123"}'

# Test protected endpoint
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer <token-from-login>"

# Run tests
go test ./minis/35-jwt-auth-middleware/...

# Run with race detector
go test -race ./minis/35-jwt-auth-middleware/...
```

---

## Summary

**What you learned**:
- ✅ JWT structure (header, payload, signature)
- ✅ Token generation and validation
- ✅ Middleware pattern for authentication
- ✅ Security best practices (expiration, HTTPS, signing method validation)
- ✅ Refresh token pattern
- ✅ Role-based access control

**Why this matters**:
JWT is the industry standard for stateless authentication in modern web applications. Understanding JWT deeply allows you to build secure, scalable authentication systems that work across microservices, mobile apps, and SPAs.

**Key formulas**:
- **Token structure**: `Base64URL(header).Base64URL(payload).Base64URL(signature)`
- **Signature**: `HMAC-SHA256(header.payload, secret)`
- **Validation**: Verify signature + check expiration + validate claims

**Next steps**:
- Project 36: Build a caching reverse proxy
- Project 37: Learn advanced middleware patterns
- OAuth2 and OpenID Connect for enterprise authentication

Stay secure! 🔐
