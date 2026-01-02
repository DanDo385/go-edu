# 35-jwt-auth-middleware

**JWT Authentication**

Build JWT authentication middleware.

## What You'll Learn

- JWT structure (header.payload.signature)
- Token signing and verification
- HTTP middleware pattern
- Claims extraction

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement JWT middleware | Token verification |

## Project Structure

```
35-jwt-auth-middleware/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/jwtauthmiddleware/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/35-jwt-auth-middleware

# Start server
go run ./cmd/app/main.go --secret mysecret

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Start server with JWT auth
go run ./cmd/app/main.go --secret mysecret

# Get a token (login endpoint)
curl -X POST http://localhost:8080/login -d '{"user":"test"}'

# Use token
curl -H "Authorization: Bearer <token>" http://localhost:8080/protected

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **JWT**: JSON Web Token
2. **HMAC Signing**: Secret-based signatures
3. **Claims**: Payload data (sub, exp, etc.)
4. **Middleware**: Extract and verify token

## Next Steps

After completing this exercise, proceed to `minis/36-caching-reverse-proxy`.
