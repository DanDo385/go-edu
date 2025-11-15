# Project Structure

Complete directory tree and file descriptions for the microservice capstone project.

## Directory Tree

```
50-mini-service-all-features/
├── README.md                      # Main documentation (1,848 lines)
├── QUICKSTART.md                  # 5-minute getting started guide
├── TEST.md                        # API testing examples
├── DEPLOYMENT.md                  # Production deployment guide
├── PROJECT_SUMMARY.md             # Project overview and statistics
├── STRUCTURE.md                   # This file
├── Makefile                       # Build automation
├── .gitignore                     # Git ignore patterns
├── config.yaml                    # Default configuration
├── go.mod                         # Go module definition
│
├── cmd/
│   └── service/
│       └── main.go                # Application entry point (149 lines)
│
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management (90 lines)
│   │
│   ├── metrics/
│   │   └── metrics.go             # Prometheus metrics (27 lines)
│   │
│   ├── models/
│   │   ├── user.go                # User data model (15 lines)
│   │   └── errors.go              # Error types (11 lines)
│   │
│   ├── database/
│   │   └── database.go            # In-memory database (144 lines)
│   │
│   ├── middleware/
│   │   ├── middleware.go          # Chain builder and ResponseWriter (55 lines)
│   │   ├── recovery.go            # Panic recovery (23 lines)
│   │   ├── requestid.go           # Request ID generation (39 lines)
│   │   ├── logging.go             # Structured logging (40 lines)
│   │   ├── metrics.go             # Metrics collection (34 lines)
│   │   ├── cors.go                # CORS handling (48 lines)
│   │   ├── ratelimit.go           # Rate limiting (24 lines)
│   │   └── auth.go                # JWT authentication (134 lines)
│   │
│   └── handlers/
│       ├── health.go              # Health endpoints (35 lines)
│       ├── auth.go                # Login handler (78 lines)
│       └── users.go               # User handlers (105 lines)
│
└── exercise/
    ├── exercise.go                # Exercise skeletons (216 lines)
    ├── exercise_test.go           # Exercise tests (149 lines)
    └── solution.go                # Complete solutions (377 lines)
```

## File Descriptions

### Root Level

#### Documentation Files

**README.md** (1,848 lines)
- What is this project about?
- Real-world scenario and motivation
- First principles of microservice architecture
- Complete implementation walkthrough
- Key concepts and patterns
- Testing strategies
- Deployment considerations
- Performance benchmarks
- Stretch goals

**QUICKSTART.md**
- 5-minute quick start guide
- Basic commands to get started
- Available endpoints
- Test users
- Common issues

**TEST.md** (293 lines)
- Complete API testing guide
- curl examples for all endpoints
- Authentication flow
- Error handling examples
- Metrics collection
- Load testing
- Request tracing

**DEPLOYMENT.md** (581 lines)
- Local development setup
- Docker deployment
- Kubernetes manifests
- Production checklist
- Monitoring setup
- Troubleshooting guide

**PROJECT_SUMMARY.md**
- Project overview
- What was created
- Statistics
- Architecture patterns
- Key takeaways

**STRUCTURE.md**
- This file
- Complete directory tree
- File descriptions

#### Configuration Files

**config.yaml** (21 lines)
- Server configuration (address, timeouts)
- Logging configuration (level, format)
- CORS configuration
- Rate limiting configuration
- JWT configuration

**go.mod** (18 lines)
- Module definition
- Dependencies:
  - github.com/google/uuid
  - github.com/prometheus/client_golang
  - github.com/rs/zerolog
  - golang.org/x/time
  - gopkg.in/yaml.v3

**Makefile** (30 lines)
- build: Build the binary
- run: Run the service
- test: Run all tests
- test-exercise: Run exercise tests
- test-race: Run with race detector
- clean: Clean artifacts
- lint: Run linter
- fmt: Format code

**.gitignore**
- Binaries
- Test outputs
- Dependencies
- IDE files
- OS files

### cmd/service/

**main.go** (149 lines)

Main application entry point with:
- Configuration loading
- Logger setup
- Metrics initialization
- Database initialization
- Router setup with middleware
- HTTP server creation
- Graceful shutdown handling

Key functions:
- `main()` - Entry point
- `setupLogger()` - Logger configuration
- `setupRouter()` - Route and middleware setup

### internal/config/

**config.go** (90 lines)

Configuration management with:
- YAML file parsing
- Environment variable overrides
- Validation
- Type-safe configuration structs

Types:
- `Config` - Main config
- `ServerConfig` - Server settings
- `LoggingConfig` - Log settings
- `CORSConfig` - CORS settings
- `RateLimitConfig` - Rate limit settings
- `JWTConfig` - JWT settings

Functions:
- `Load()` - Load and validate config
- `Validate()` - Validate configuration

### internal/metrics/

**metrics.go** (27 lines)

Prometheus metrics with:
- HTTP request counter
- HTTP request duration histogram
- HTTP active requests gauge

Type:
- `Metrics` - Metrics collection

Function:
- `New()` - Create metrics registry

### internal/models/

**user.go** (15 lines)

Data models:
- `User` - User entity
- `LoginRequest` - Login payload
- `LoginResponse` - Login response with token

**errors.go** (11 lines)

Error types:
- `ErrorResponse` - Structured error
- `NewErrorResponse()` - Error constructor

### internal/database/

**database.go** (144 lines)

In-memory database with:
- Sample user data
- User credentials (username/password)
- Thread-safe operations

Type:
- `DB` - Database wrapper

Functions:
- `New()` - Create database with sample data
- `Health()` - Health check
- `Close()` - Close connections
- `GetUser()` - Get user by ID
- `ListUsers()` - List all users
- `Authenticate()` - Verify credentials

### internal/middleware/

**middleware.go** (55 lines)

Middleware infrastructure:
- Middleware type definition
- Chain builder
- ResponseWriter wrapper

Types:
- `Middleware` - Middleware function type
- `ResponseWriter` - Response wrapper

Functions:
- `Chain()` - Apply middleware chain
- `NewResponseWriter()` - Create wrapper
- `WriteHeader()` - Capture status code
- `Write()` - Capture bytes written

**recovery.go** (23 lines)

Panic recovery:
- Catch panics
- Log stack trace
- Return 500 error

Function:
- `Recovery()` - Recovery middleware

**requestid.go** (39 lines)

Request ID generation:
- Generate unique UUIDs
- Store in context
- Add to response headers

Functions:
- `RequestID()` - Request ID middleware
- `GetRequestID()` - Extract from context

**logging.go** (40 lines)

Structured logging:
- Log requests and responses
- Include request ID
- Measure duration

Function:
- `Logging()` - Logging middleware

**metrics.go** (34 lines)

Metrics collection:
- Track request count
- Measure duration
- Count active requests

Function:
- `Metrics()` - Metrics middleware

**cors.go** (48 lines)

CORS handling:
- Set CORS headers
- Handle preflight requests

Function:
- `CORS()` - CORS middleware

**ratelimit.go** (24 lines)

Rate limiting:
- Token bucket algorithm
- Configurable limits

Function:
- `RateLimit()` - Rate limit middleware

**auth.go** (134 lines)

JWT authentication:
- Token generation (HS256)
- Token validation
- Claims extraction
- Context propagation

Types:
- `Claims` - JWT claims

Functions:
- `GenerateJWT()` - Create token
- `ValidateJWT()` - Verify token
- `Auth()` - Auth middleware
- `GetUser()` - Extract from context

### internal/handlers/

**health.go** (35 lines)

Health check endpoints:
- Liveness probe
- Readiness probe (checks DB)

Functions:
- `Health()` - Health handler
- `Ready()` - Readiness handler

**auth.go** (78 lines)

Authentication handlers:
- Login endpoint
- JWT token generation
- Input validation

Function:
- `Login()` - Login handler

**users.go** (105 lines)

User CRUD handlers:
- List all users
- Get user by ID
- Extract user from context

Functions:
- `ListUsers()` - List handler
- `GetUser()` - Get handler

### exercise/

**exercise.go** (216 lines)

Exercise skeletons for 8 exercises:
1. Cache with TTL
2. Circuit breaker
3. Timeout middleware
4. Worker pool
5. Retry with exponential backoff
6. Per-user rate limiter
7. Structured error handling
8. Event bus

Each exercise includes:
- Type definitions
- Function signatures
- TODO comments

**exercise_test.go** (149 lines)

Comprehensive tests:
- TestCache
- TestCircuitBreaker
- TestRetryWithBackoff
- TestUserRateLimiter
- TestAppError
- TestEventBus
- TestWorkerPool

**solution.go** (377 lines)

Complete implementations:
- `SolutionCache` - Cache with TTL and cleanup
- `SolutionCircuitBreaker` - Circuit breaker with states
- `SolutionTimeoutMiddleware` - Request timeout
- `SolutionWorkerPool` - Concurrent job processing
- `SolutionRetryWithBackoff` - Exponential backoff
- `SolutionUserRateLimiter` - Per-user limits
- `SolutionAppError` - Structured errors
- `SolutionEventBus` - Publish/subscribe

## Import Paths

All internal packages use the import path:
```go
import "github.com/user/go-edu/minis/50-mini-service-all-features/internal/..."
```

Examples:
```go
import (
    "github.com/user/go-edu/minis/50-mini-service-all-features/internal/config"
    "github.com/user/go-edu/minis/50-mini-service-all-features/internal/middleware"
    "github.com/user/go-edu/minis/50-mini-service-all-features/internal/handlers"
)
```

## Package Dependencies

```
main
├── config
├── database
├── metrics
├── middleware
│   ├── config (for CORS, RateLimit)
│   └── metrics
├── handlers
│   ├── config (for JWT)
│   ├── database
│   ├── middleware (for Auth)
│   └── models
└── models
```

## Key Files to Start With

For learning, read files in this order:

1. **README.md** - Understand architecture
2. **QUICKSTART.md** - Get service running
3. **config.yaml** - See configuration options
4. **cmd/service/main.go** - Application entry point
5. **internal/middleware/middleware.go** - Middleware basics
6. **internal/handlers/health.go** - Simple handler example
7. **internal/middleware/auth.go** - JWT implementation
8. **exercise/exercise.go** - Practice exercises
9. **TEST.md** - Test the service
10. **DEPLOYMENT.md** - Production deployment

## Line Count Summary

| Category | Files | Lines |
|----------|-------|-------|
| Documentation | 6 | 2,722+ |
| Source Code | 20 | 1,833 |
| Configuration | 4 | ~100 |
| **Total** | **30** | **4,655+** |

## External Dependencies

| Package | Purpose |
|---------|---------|
| github.com/rs/zerolog | Structured logging |
| github.com/prometheus/client_golang | Metrics |
| github.com/google/uuid | UUID generation |
| golang.org/x/time/rate | Rate limiting |
| gopkg.in/yaml.v3 | YAML parsing |

All dependencies are allowed per project requirements.
