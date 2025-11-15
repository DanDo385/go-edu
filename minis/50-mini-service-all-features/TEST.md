# Testing Guide

This guide shows how to test the microservice with curl commands.

## Prerequisites

Make sure the service is running:
```bash
go run ./cmd/service
```

The service will start on `http://localhost:8080`.

## Health Checks

### Liveness Probe
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"healthy"}
```

### Readiness Probe
```bash
curl http://localhost:8080/ready
```

Expected response:
```json
{"status":"ready"}
```

## Authentication

### Login (Get JWT Token)

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'
```

Expected response:
```json
{
  "token": "eyJhbGc...very-long-token...",
  "expires_at": 1705334400
}
```

Save the token for authenticated requests:
```bash
export TOKEN="<paste-token-here>"
```

Or use this one-liner:
```bash
export TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}' \
  | jq -r '.token')
```

## Protected Endpoints

### List All Users

```bash
curl http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN"
```

Expected response:
```json
[
  {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "created_at": "2024-12-15T10:30:00Z"
  },
  {
    "id": 2,
    "username": "bob",
    "email": "bob@example.com",
    "created_at": "2024-12-30T14:20:00Z"
  }
]
```

### Get Single User

```bash
curl http://localhost:8080/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

Expected response:
```json
{
  "id": 1,
  "username": "alice",
  "email": "alice@example.com",
  "created_at": "2024-12-15T10:30:00Z"
}
```

## Error Cases

### Missing Authentication

```bash
curl http://localhost:8080/users
```

Expected response (401):
```
Unauthorized
```

### Invalid Token

```bash
curl http://localhost:8080/users \
  -H "Authorization: Bearer invalid-token"
```

Expected response (401):
```
Invalid token: invalid token format
```

### User Not Found

```bash
curl http://localhost:8080/users/999 \
  -H "Authorization: Bearer $TOKEN"
```

Expected response (404):
```json
{
  "code": "USER_NOT_FOUND",
  "message": "User not found"
}
```

### Invalid Credentials

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"wrong"}'
```

Expected response (401):
```json
{
  "code": "INVALID_CREDENTIALS",
  "message": "Invalid username or password"
}
```

## Metrics

### Prometheus Metrics

```bash
curl http://localhost:8080/metrics
```

Expected response (sample):
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/users",status="200"} 5

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/users",status="200",le="0.005"} 2
http_request_duration_seconds_sum{method="GET",path="/users",status="200"} 0.123
http_request_duration_seconds_count{method="GET",path="/users",status="200"} 5

# HELP http_active_requests Number of active HTTP requests
# TYPE http_active_requests gauge
http_active_requests 0
```

## CORS Testing

```bash
curl -X OPTIONS http://localhost:8080/users \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

Check for CORS headers in response:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## Rate Limiting

Send multiple requests rapidly to trigger rate limiting:

```bash
for i in {1..150}; do
  curl -s http://localhost:8080/health > /dev/null
  echo "Request $i sent"
done
```

You should see some requests fail with 429 Too Many Requests.

## Load Testing

Use a tool like `ab` (Apache Bench) or `hey`:

```bash
# Install hey: go install github.com/rakyll/hey@latest

# 1000 requests, 50 concurrent
hey -n 1000 -c 50 http://localhost:8080/health
```

## Graceful Shutdown

1. Start the service
2. Send Ctrl+C (SIGINT) or `kill -TERM <pid>`
3. Observe logs showing graceful shutdown:
   ```
   {"level":"info","time":"...","message":"Shutting down server..."}
   {"level":"info","time":"...","message":"Server stopped gracefully"}
   ```

## Request Tracing

All requests include a `X-Request-ID` header in the response:

```bash
curl -v http://localhost:8080/health 2>&1 | grep X-Request-ID
```

Output:
```
< X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

You can also provide your own request ID:

```bash
curl http://localhost:8080/health \
  -H "X-Request-ID: my-custom-id" \
  -v
```

## Available Test Users

| Username | Password     | User ID |
|----------|--------------|---------|
| alice    | password123  | 1       |
| bob      | password123  | 2       |
| charlie  | password123  | 3       |

## Complete Workflow Example

```bash
# 1. Health check
curl http://localhost:8080/health

# 2. Login
export TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}' \
  | jq -r '.token')

# 3. List users
curl http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN" | jq

# 4. Get specific user
curl http://localhost:8080/users/1 \
  -H "Authorization: Bearer $TOKEN" | jq

# 5. Check metrics
curl http://localhost:8080/metrics | grep http_requests_total
```
