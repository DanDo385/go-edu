# Quick Start Guide

Get up and running in 5 minutes!

## Step 1: Start the Service

```bash
cd /home/user/go-edu/minis/50-mini-service-all-features
go run ./cmd/service
```

You should see:
```
{"level":"info","time":"...","message":"Starting microservice..."}
{"level":"info","time":"...","message":"Database initialized"}
{"level":"info","time":"...","message":"Server starting on :8080"}
```

## Step 2: Test Health Check

Open a new terminal:
```bash
curl http://localhost:8080/health
```

Response:
```json
{"status":"healthy"}
```

## Step 3: Login and Get Token

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'
```

Copy the token from the response.

## Step 4: Access Protected Endpoint

```bash
export TOKEN="<paste-your-token-here>"

curl http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN"
```

## Step 5: View Metrics

```bash
curl http://localhost:8080/metrics
```

## What's Next?

- **Read the full README.md** for architecture details
- **Check TEST.md** for more API examples
- **Try the exercises** in `exercise/exercise.go`
- **Review DEPLOYMENT.md** for production deployment

## Available Users

| Username | Password    |
|----------|-------------|
| alice    | password123 |
| bob      | password123 |
| charlie  | password123 |

## Key Endpoints

- `GET /health` - Health check (no auth)
- `GET /ready` - Readiness check (no auth)
- `GET /metrics` - Prometheus metrics (no auth)
- `POST /login` - Get JWT token (no auth)
- `GET /users` - List all users (auth required)
- `GET /users/{id}` - Get user by ID (auth required)

## Graceful Shutdown

Press `Ctrl+C` to trigger graceful shutdown. The service will:
1. Stop accepting new requests
2. Wait for in-flight requests to complete (up to 30s)
3. Close database connections
4. Exit cleanly

## Run Tests

```bash
# All tests
go test ./...

# Just exercises
go test ./exercise

# With race detector
go test -race ./...
```

## Configuration

Edit `config.yaml` or use environment variables:

```bash
export SERVER_ADDR=":9090"
export JWT_SECRET="my-secret"
export LOG_LEVEL="debug"
```

## Common Issues

**Port already in use?**
```bash
export SERVER_ADDR=":9090"
go run ./cmd/service
```

**Want JSON logs in console format?**
Edit `config.yaml`:
```yaml
logging:
  format: console  # instead of json
```

**Need debug logs?**
```bash
export LOG_LEVEL="debug"
```

## Features Demonstrated

✅ Structured logging (zerolog)
✅ Prometheus metrics
✅ JWT authentication
✅ Rate limiting
✅ CORS handling
✅ Request tracing
✅ Graceful shutdown
✅ Health checks
✅ Middleware chain
✅ Error handling

Enjoy exploring the microservice! 🚀
