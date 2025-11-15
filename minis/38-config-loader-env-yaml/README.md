# Project 38: Config Loader (Environment Variables & YAML)

## 1. What Is This About?

### Real-World Scenario

You're deploying a web application to multiple environments:
- **Development**: Local laptop, debug mode enabled, uses local database
- **Staging**: Cloud server, logs to file, uses test database
- **Production**: Multiple cloud servers, strict security, uses production database cluster

**❌ Bad approach:** Hard-code everything in source code
```go
const dbHost = "localhost"           // Won't work in production!
const debugMode = true               // Security risk in production!
const apiKey = "sk_test_123"         // API keys in source code = bad!
```

**Problems:**
- Need to rebuild app for each environment
- Secrets exposed in version control
- Can't change settings without redeploying
- Different team members need different settings

**✅ Better approach:** External configuration
```yaml
# config.yaml
database:
  host: ${DB_HOST}
  port: 5432
  timeout: 30s

server:
  port: 8080
  debug: false

api:
  key: ${API_KEY}  # Load from environment variable
  timeout: 10s
```

This project teaches you **configuration management** in Go:
1. **Environment variables**: OS-level config (12-factor apps)
2. **YAML parsing**: Structured configuration files
3. **Config validation**: Ensure required fields are present
4. **Default values**: Sensible fallbacks
5. **Type safety**: Parse strings to proper types (duration, int, bool)
6. **Env var substitution**: Replace `${VAR}` with environment values

### What You'll Learn

1. **Configuration hierarchy**: Environment variables override YAML defaults
2. **YAML parsing**: Using `gopkg.in/yaml.v3` for structured config
3. **Validation**: Required fields, valid ranges, format checking
4. **Type conversion**: Parse durations, parse ports, validate URLs
5. **Error handling**: Clear error messages for config problems
6. **Best practices**: The 12-factor app methodology

### The Challenge

Build a configuration loader that:
- Reads YAML configuration files
- Substitutes environment variables (e.g., `${DB_HOST}`)
- Validates required fields and types
- Provides sensible defaults
- Returns typed, validated configuration structs
- Gives clear error messages for config problems

---

## 2. First Principles: Why Configuration Management Matters

### What is Configuration?

**Configuration** = Settings that change how your program behaves without changing code.

**Analogy**: Configuration is like the settings on your phone:
- You don't recompile iOS to change WiFi settings
- You don't rebuild Android to enable dark mode
- Settings are **external** to the app itself

### The 12-Factor App Methodology

The [12-Factor App](https://12factor.net/) is a widely-adopted methodology for building modern applications. Factor III states:

> **Store config in the environment**

**Why?**
1. **Separation of concerns**: Code (logic) vs config (settings)
2. **Security**: Secrets never committed to version control
3. **Flexibility**: Same binary runs in any environment
4. **Scalability**: Each instance can have different config

**Example of 12-factor config:**
```bash
# Development
export DB_HOST=localhost
export DB_PORT=5432
export DEBUG=true
./myapp

# Production
export DB_HOST=prod-db.example.com
export DB_PORT=5432
export DEBUG=false
./myapp  # Same binary!
```

### Configuration Hierarchy

Most applications use a **layered approach** where each layer can override the previous:

```
1. Hardcoded defaults (in code)
   ↓ overridden by ↓
2. Configuration file (e.g., config.yaml)
   ↓ overridden by ↓
3. Environment variables
   ↓ overridden by ↓
4. Command-line flags
```

**Example:**
```go
// 1. Code default
port := 8080

// 2. YAML file might say: port: 3000
if config.Port != 0 {
    port = config.Port  // Now port = 3000
}

// 3. Environment variable: PORT=9000
if envPort := os.Getenv("PORT"); envPort != "" {
    port, _ = strconv.Atoi(envPort)  // Now port = 9000
}

// Final port = 9000 (env var wins)
```

### What is YAML?

**YAML** (YAML Ain't Markup Language) is a human-friendly data serialization format.

**Comparison with JSON:**

```json
// JSON: Great for machines, verbose for humans
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "credentials": {
      "username": "admin",
      "password": "secret"
    }
  }
}
```

```yaml
# YAML: Great for humans, easier to read/write
database:
  host: localhost
  port: 5432
  credentials:
    username: admin
    password: secret
```

**Key features:**
- Indentation-based structure (like Python)
- Comments with `#`
- No quotes needed for simple strings
- Lists, maps, and nested structures
- Type inference (numbers, bools, strings)

**Common YAML types:**
```yaml
# Strings
name: myapp
description: "Quoted strings can have: special chars"

# Numbers
port: 8080
pi: 3.14159

# Booleans
debug: true
enabled: false

# Lists (two syntaxes)
servers:
  - web1.example.com
  - web2.example.com
  - web3.example.com

tags: [production, critical, database]

# Maps (dictionaries)
database:
  host: localhost
  port: 5432

# Null
nothing: null
nothing_implicit:  # Also null
```

### Environment Variable Substitution

A powerful pattern is to **embed environment variable references** in your YAML:

```yaml
database:
  host: ${DB_HOST}
  password: ${DB_PASSWORD}

server:
  port: ${PORT:-8080}  # Use PORT env var, or default to 8080
```

**Why this is useful:**
1. YAML file can be committed to git (no secrets!)
2. Secrets come from environment (via env vars or secret management systems)
3. Same config file works in all environments
4. Easy to override specific values without changing entire file

**Implementation strategy:**
1. Read YAML file into string
2. Find all `${VAR}` or `${VAR:-default}` patterns
3. Replace with actual environment values
4. Parse the substituted YAML into structs

---

## 3. Breaking Down the Solution

### Step 1: Define the Configuration Structure

First, we need to define **what** we're configuring. Use Go structs with YAML tags:

```go
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Logging  LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
    Host         string        `yaml:"host"`
    Port         int           `yaml:"port"`
    ReadTimeout  time.Duration `yaml:"read_timeout"`
    WriteTimeout time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Database string `yaml:"database"`
    MaxConns int    `yaml:"max_connections"`
}

type LoggingConfig struct {
    Level  string `yaml:"level"`
    Format string `yaml:"format"`
    Output string `yaml:"output"`
}
```

**YAML tags** tell the YAML parser which struct field maps to which YAML key.

### Step 2: Parse YAML Files

Using `gopkg.in/yaml.v3`:

```go
func LoadYAML(filename string) (*Config, error) {
    // Read file
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("reading config file: %w", err)
    }

    // Parse YAML
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("parsing YAML: %w", err)
    }

    return &config, nil
}
```

**How YAML unmarshaling works:**
1. YAML parser reads the file bytes
2. Parses YAML structure (maps, lists, scalars)
3. Matches YAML keys to struct field tags
4. Converts YAML values to Go types (string → string, 123 → int, true → bool)
5. Populates the struct

### Step 3: Environment Variable Substitution

We need to replace `${VAR}` patterns with actual environment values:

```go
func substituteEnvVars(input string) string {
    // Regex pattern: ${VAR} or ${VAR:-default}
    re := regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)(:-([^}]+))?\}`)

    return re.ReplaceAllStringFunc(input, func(match string) string {
        // Extract variable name and default value
        parts := re.FindStringSubmatch(match)
        varName := parts[1]
        defaultVal := parts[3]

        // Look up environment variable
        if value := os.Getenv(varName); value != "" {
            return value
        }

        // Use default if provided
        if defaultVal != "" {
            return defaultVal
        }

        // No value found, keep original
        return match
    })
}
```

**Pattern breakdown:**
- `\$\{` - Literal `${`
- `([A-Z_][A-Z0-9_]*)` - Variable name (capture group 1)
  - Must start with letter or underscore
  - Can contain letters, numbers, underscores
- `(:-([^}]+))?` - Optional default value (capture groups 2 and 3)
  - `:-` separator
  - `([^}]+)` - Any characters except `}`
- `\}` - Literal `}`

**Examples:**
- `${DB_HOST}` → `parts[1]="DB_HOST", parts[3]=""`
- `${PORT:-8080}` → `parts[1]="PORT", parts[3]="8080"`

### Step 4: Validate Configuration

After loading, we must **validate** that the config is sensible:

```go
func (c *Config) Validate() error {
    var errors []string

    // Check required fields
    if c.Database.Host == "" {
        errors = append(errors, "database.host is required")
    }
    if c.Database.Database == "" {
        errors = append(errors, "database.database is required")
    }

    // Check valid ranges
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        errors = append(errors, "server.port must be between 1 and 65535")
    }
    if c.Database.MaxConns < 1 {
        errors = append(errors, "database.max_connections must be at least 1")
    }

    // Check valid values
    validLogLevels := map[string]bool{
        "debug": true, "info": true, "warn": true, "error": true,
    }
    if !validLogLevels[c.Logging.Level] {
        errors = append(errors,
            "logging.level must be one of: debug, info, warn, error")
    }

    if len(errors) > 0 {
        return fmt.Errorf("configuration validation failed:\n  - %s",
            strings.Join(errors, "\n  - "))
    }

    return nil
}
```

**Validation categories:**
1. **Required fields**: Must not be empty
2. **Range checks**: Port numbers, connection pools, timeouts
3. **Enum values**: Log levels, formats, modes
4. **Format checks**: URLs, file paths, regex patterns

### Step 5: Apply Defaults

Some fields should have **sensible defaults** if not specified:

```go
func (c *Config) ApplyDefaults() {
    // Server defaults
    if c.Server.Host == "" {
        c.Server.Host = "0.0.0.0"
    }
    if c.Server.Port == 0 {
        c.Server.Port = 8080
    }
    if c.Server.ReadTimeout == 0 {
        c.Server.ReadTimeout = 30 * time.Second
    }
    if c.Server.WriteTimeout == 0 {
        c.Server.WriteTimeout = 30 * time.Second
    }

    // Database defaults
    if c.Database.Port == 0 {
        c.Database.Port = 5432  // PostgreSQL default
    }
    if c.Database.MaxConns == 0 {
        c.Database.MaxConns = 10
    }

    // Logging defaults
    if c.Logging.Level == "" {
        c.Logging.Level = "info"
    }
    if c.Logging.Format == "" {
        c.Logging.Format = "json"
    }
    if c.Logging.Output == "" {
        c.Logging.Output = "stdout"
    }
}
```

**When to use defaults vs required validation:**
- **Defaults**: Reasonable fallback exists, optional setting
- **Required**: No sensible default, application can't work without it

### Step 6: Parse Special Types

YAML gives us strings, but we often need special types:

**Duration parsing:**
```yaml
timeout: 30s  # String "30s" in YAML
```

```go
type ServerConfig struct {
    Timeout time.Duration `yaml:"timeout"`
}

// Custom unmarshaler
func (s *ServerConfig) UnmarshalYAML(value *yaml.Node) error {
    // Define a temporary struct with string field
    type raw struct {
        Timeout string `yaml:"timeout"`
    }

    var r raw
    if err := value.Decode(&r); err != nil {
        return err
    }

    // Parse duration
    if r.Timeout != "" {
        d, err := time.ParseDuration(r.Timeout)
        if err != nil {
            return fmt.Errorf("invalid timeout: %w", err)
        }
        s.Timeout = d
    }

    return nil
}
```

**Valid duration formats:**
- `"300ms"` = 300 milliseconds
- `"1.5s"` = 1.5 seconds
- `"2m"` = 2 minutes
- `"1h"` = 1 hour
- `"1h30m"` = 1 hour 30 minutes

---

## 4. Complete Solution Walkthrough

Let's build the complete config loader step by step.

### The LoadConfig Function

This is the main entry point:

```go
func LoadConfig(filename string) (*Config, error) {
    // Step 1: Read file
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("reading config file %s: %w", filename, err)
    }

    // Step 2: Substitute environment variables
    substituted := substituteEnvVars(string(data))

    // Step 3: Parse YAML
    var config Config
    if err := yaml.Unmarshal([]byte(substituted), &config); err != nil {
        return nil, fmt.Errorf("parsing YAML: %w", err)
    }

    // Step 4: Apply defaults
    config.ApplyDefaults()

    // Step 5: Validate
    if err := config.Validate(); err != nil {
        return nil, err
    }

    return &config, nil
}
```

**Why this order?**
1. Read first (need data to process)
2. Substitute before parsing (YAML parser needs final values)
3. Parse into struct
4. Apply defaults (before validation, so defaults are validated too)
5. Validate (catch problems before using config)

### The substituteEnvVars Function

```go
func substituteEnvVars(input string) string {
    // Pattern: ${VAR} or ${VAR:-default}
    re := regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)(:-([^}]*))?\}`)

    return re.ReplaceAllStringFunc(input, func(match string) string {
        matches := re.FindStringSubmatch(match)
        varName := matches[1]
        defaultValue := ""
        if len(matches) > 3 {
            defaultValue = matches[3]
        }

        // Try environment variable first
        if value := os.Getenv(varName); value != "" {
            return value
        }

        // Fall back to default
        if defaultValue != "" {
            return defaultValue
        }

        // No replacement available - keep original
        return match
    })
}
```

**Edge cases:**
- `${UNDEFINED}` → stays as `${UNDEFINED}` (no default, no env var)
- `${UNDEFINED:-fallback}` → becomes `fallback`
- `${DEFINED}` where `DEFINED=value` → becomes `value`
- `${DEFINED:-fallback}` where `DEFINED=value` → becomes `value` (env var wins)

### Environment Variable Override Pattern

Sometimes you want **direct environment variable overrides** without YAML:

```go
func LoadConfigWithEnvOverrides(filename string) (*Config, error) {
    config, err := LoadConfig(filename)
    if err != nil {
        return nil, err
    }

    // Override with direct env vars
    if port := os.Getenv("SERVER_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            config.Server.Port = p
        }
    }

    if host := os.Getenv("SERVER_HOST"); host != "" {
        config.Server.Host = host
    }

    if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
        config.Database.Host = dbHost
    }

    // Re-validate after overrides
    if err := config.Validate(); err != nil {
        return nil, err
    }

    return config, nil
}
```

This allows both patterns:
- YAML with substitution: `host: ${DB_HOST}`
- Direct env override: `DB_HOST=localhost` (works even if YAML says something else)

---

## 5. Key Concepts Explained

### Concept 1: Struct Tags

**Struct tags** are metadata attached to struct fields:

```go
type Config struct {
    Port    int    `yaml:"port" json:"port" env:"PORT"`
    Timeout string `yaml:"timeout,omitempty"`
}
```

**Tag format:** `` `key1:"value1" key2:"value2"` ``

**Common tag keys:**
- `yaml:"field_name"` - YAML field mapping
- `json:"field_name"` - JSON field mapping
- `env:"VAR_NAME"` - Environment variable mapping (custom)
- `default:"value"` - Default value (custom)
- `validate:"required"` - Validation rule (using validator libraries)

**Special modifiers:**
- `yaml:"field,omitempty"` - Omit field if zero value when marshaling
- `yaml:"-"` - Never marshal/unmarshal this field
- `yaml:",inline"` - Embed struct fields at same level

**Example:**
```go
type Server struct {
    Host string `yaml:"host" json:"host"`
    Port int    `yaml:"port" json:"port"`

    // Not in YAML/JSON output
    runtime string `yaml:"-" json:"-"`
}
```

### Concept 2: YAML Anchors and Aliases

YAML supports **DRY (Don't Repeat Yourself)** with anchors:

```yaml
# Define an anchor with &
defaults: &defaults
  timeout: 30s
  retries: 3

# Reference with *
service_a:
  <<: *defaults  # Merge defaults
  host: service-a.example.com

service_b:
  <<: *defaults
  host: service-b.example.com
  timeout: 60s  # Override specific value
```

**Result after parsing:**
```yaml
service_a:
  timeout: 30s
  retries: 3
  host: service-a.example.com

service_b:
  timeout: 60s    # Overridden
  retries: 3
  host: service-b.example.com
```

### Concept 3: Configuration Hot Reload

Advanced pattern: Reload config without restarting the app:

```go
type App struct {
    config atomic.Value // Thread-safe config storage
}

func (a *App) ReloadConfig(filename string) error {
    newConfig, err := LoadConfig(filename)
    if err != nil {
        return err
    }

    // Atomic swap
    a.config.Store(newConfig)
    return nil
}

func (a *App) GetConfig() *Config {
    return a.config.Load().(*Config)
}

// Watch file for changes
func (a *App) WatchConfigFile(filename string) {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add(filename)

    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Println("Config file changed, reloading...")
                if err := a.ReloadConfig(filename); err != nil {
                    log.Printf("Failed to reload: %v", err)
                } else {
                    log.Println("Config reloaded successfully")
                }
            }
        }
    }
}
```

### Concept 4: Secret Management

**Never** commit secrets to version control. Use environment variables or secret management systems:

**Bad:**
```yaml
database:
  password: super_secret_password  # In git = bad!
```

**Good:**
```yaml
database:
  password: ${DB_PASSWORD}  # Loaded from environment
```

**Better (for production):**
Use dedicated secret managers:
- **AWS Secrets Manager**
- **HashiCorp Vault**
- **Kubernetes Secrets**
- **Azure Key Vault**

**Example with AWS Secrets Manager:**
```go
func loadSecret(secretName string) (string, error) {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess)

    result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    })
    if err != nil {
        return "", err
    }

    return *result.SecretString, nil
}

// In config loading:
if strings.HasPrefix(c.Database.Password, "secret://") {
    secretName := strings.TrimPrefix(c.Database.Password, "secret://")
    password, err := loadSecret(secretName)
    if err != nil {
        return nil, fmt.Errorf("loading secret: %w", err)
    }
    c.Database.Password = password
}
```

### Concept 5: Configuration Profiles

Support multiple environments with profiles:

```yaml
# config.yaml
common: &common
  server:
    timeout: 30s

development:
  <<: *common
  server:
    host: localhost
    port: 8080
  database:
    host: localhost

production:
  <<: *common
  server:
    host: 0.0.0.0
    port: 80
  database:
    host: prod-db.example.com
```

```go
func LoadConfigWithProfile(filename, profile string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // Parse entire file
    var allConfigs map[string]*Config
    if err := yaml.Unmarshal(data, &allConfigs); err != nil {
        return nil, err
    }

    // Extract specific profile
    config, exists := allConfigs[profile]
    if !exists {
        return nil, fmt.Errorf("profile %q not found", profile)
    }

    return config, nil
}

// Usage:
// APP_ENV=production ./myapp
profile := os.Getenv("APP_ENV")
if profile == "" {
    profile = "development"
}
config, _ := LoadConfigWithProfile("config.yaml", profile)
```

---

## 6. Real-World Applications

### Web Applications

**Use case:** Configure HTTP server, database, caching, logging

```yaml
server:
  host: ${HOST:-0.0.0.0}
  port: ${PORT:-8080}
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s

database:
  host: ${DB_HOST}
  port: ${DB_PORT:-5432}
  username: ${DB_USER}
  password: ${DB_PASSWORD}
  database: ${DB_NAME}
  max_connections: ${DB_MAX_CONNS:-20}

redis:
  host: ${REDIS_HOST:-localhost}
  port: ${REDIS_PORT:-6379}

logging:
  level: ${LOG_LEVEL:-info}
  format: json
```

Companies using this: Every SaaS company (Stripe, GitHub, Shopify, etc.)

### Microservices

**Use case:** Each service has its own config, but shares common settings

```yaml
# common.yaml
common:
  tracing:
    endpoint: ${JAEGER_ENDPOINT}
  metrics:
    endpoint: ${PROMETHEUS_ENDPOINT}

# user-service.yaml
service:
  name: user-service
  port: 8001

database:
  host: ${USER_DB_HOST}

common: !include common.yaml
```

Companies: Netflix, Uber, Airbnb (microservices architecture)

### CLI Tools

**Use case:** User-configurable defaults for command-line tools

```yaml
# ~/.mytool/config.yaml
defaults:
  output_format: json
  timeout: 30s
  retries: 3

profiles:
  production:
    api_endpoint: https://api.example.com

  staging:
    api_endpoint: https://staging-api.example.com
```

Examples: `kubectl`, `aws-cli`, `terraform`

### Deployment Automation

**Use case:** Infrastructure as Code with configurable parameters

```yaml
# deploy.yaml
infrastructure:
  region: ${AWS_REGION:-us-east-1}
  instance_type: ${INSTANCE_TYPE:-t3.medium}

  autoscaling:
    min_instances: ${MIN_INSTANCES:-2}
    max_instances: ${MAX_INSTANCES:-10}
    target_cpu: 70

  database:
    instance_class: ${DB_INSTANCE:-db.t3.small}
    allocated_storage: ${DB_STORAGE:-100}
```

Tools: Terraform, Ansible, Kubernetes manifests

### Batch Processing Jobs

**Use case:** Configurable data pipelines and ETL jobs

```yaml
job:
  name: ${JOB_NAME}
  schedule: ${CRON_SCHEDULE:-0 0 * * *}  # Daily at midnight

input:
  type: ${INPUT_TYPE:-s3}
  bucket: ${INPUT_BUCKET}
  prefix: ${INPUT_PREFIX}

processing:
  batch_size: ${BATCH_SIZE:-1000}
  workers: ${WORKERS:-10}
  timeout: ${PROCESSING_TIMEOUT:-5m}

output:
  type: ${OUTPUT_TYPE:-database}
  connection: ${OUTPUT_DB}
```

Companies: Data platforms (Snowflake, Databricks, Airflow)

---

## 7. Common Mistakes to Avoid

### Mistake 1: Secrets in Version Control

**❌ Wrong:**
```yaml
database:
  password: mysecretpassword123  # Committed to git!
```

**✅ Correct:**
```yaml
database:
  password: ${DB_PASSWORD}  # Loaded from environment
```

### Mistake 2: No Validation

**❌ Wrong:**
```go
config, _ := LoadConfig("config.yaml")
// Use config without validation
db.Connect(config.Database.Host)  // Might be empty!
```

**✅ Correct:**
```go
config, err := LoadConfig("config.yaml")
if err != nil {
    log.Fatalf("Config error: %v", err)
}
if err := config.Validate(); err != nil {
    log.Fatalf("Invalid config: %v", err)
}
```

### Mistake 3: Ignoring Defaults

**❌ Wrong:**
```go
// No default for timeout
if config.Timeout == 0 {
    // Oops, zero timeout = infinite wait or immediate timeout!
}
```

**✅ Correct:**
```go
func (c *Config) ApplyDefaults() {
    if c.Timeout == 0 {
        c.Timeout = 30 * time.Second
    }
}
```

### Mistake 4: Type Conversion Errors

**❌ Wrong:**
```yaml
port: "8080"  # String instead of number
```

```go
type Config struct {
    Port int `yaml:"port"`
}
// Parsing fails silently, Port = 0
```

**✅ Correct:**
```yaml
port: 8080  # Correct type
```

Or handle string → int conversion explicitly:
```go
type rawConfig struct {
    Port interface{} `yaml:"port"`
}

// Convert to proper type with validation
```

### Mistake 5: Nested Environment Variables

**❌ Wrong:**
```go
// Trying to use nested env vars
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
// How to map these to nested struct?
```

**✅ Correct:**
Use a naming convention and explicit mapping:
```go
func LoadFromEnv(c *Config) {
    if host := os.Getenv("DATABASE_HOST"); host != "" {
        c.Database.Host = host
    }
    if port := os.Getenv("DATABASE_PORT"); port != "" {
        c.Database.Port, _ = strconv.Atoi(port)
    }
}
```

Or use a library like `envconfig` or `viper`.

### Mistake 6: Mutable Configuration

**❌ Wrong:**
```go
var GlobalConfig *Config  // Global mutable state

func GetConfig() *Config {
    return GlobalConfig  // Caller can modify!
}
```

**✅ Correct:**
```go
type App struct {
    config *Config  // Private
}

func (a *App) GetConfig() Config {
    return *a.config  // Return copy
}

// Or use interfaces for read-only access
type ConfigReader interface {
    GetDatabaseHost() string
    GetServerPort() int
}
```

### Mistake 7: Not Testing Configuration

**❌ Wrong:**
```go
// No tests for config loading
```

**✅ Correct:**
```go
func TestLoadConfig(t *testing.T) {
    tests := []struct {
        name    string
        yaml    string
        env     map[string]string
        wantErr bool
    }{
        {
            name: "valid config",
            yaml: `server: { port: 8080 }`,
            wantErr: false,
        },
        {
            name: "invalid port",
            yaml: `server: { port: 99999 }`,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test config loading
        })
    }
}
```

---

## 8. Stretch Goals

### Goal 1: Support Multiple Config Formats ⭐

Support YAML, JSON, and TOML:

```go
func LoadConfigAuto(filename string) (*Config, error) {
    switch filepath.Ext(filename) {
    case ".yaml", ".yml":
        return LoadConfigYAML(filename)
    case ".json":
        return LoadConfigJSON(filename)
    case ".toml":
        return LoadConfigTOML(filename)
    default:
        return nil, fmt.Errorf("unsupported format: %s", filepath.Ext(filename))
    }
}
```

### Goal 2: Configuration Merge ⭐⭐

Merge multiple config files:

```go
// Load base config
base, _ := LoadConfig("base.yaml")

// Load environment-specific overrides
override, _ := LoadConfig("production.yaml")

// Merge (override takes precedence)
final := MergeConfigs(base, override)
```

### Goal 3: Config Schema Validation ⭐⭐

Use JSON Schema to validate configuration:

```go
// schema.json
{
  "type": "object",
  "properties": {
    "server": {
      "type": "object",
      "properties": {
        "port": { "type": "integer", "minimum": 1, "maximum": 65535 }
      },
      "required": ["port"]
    }
  }
}

func ValidateAgainstSchema(config *Config, schemaFile string) error {
    // Load schema
    // Validate config
}
```

### Goal 4: Config Diff and Reload ⭐⭐⭐

Show what changed between configs and support hot reload:

```go
func DiffConfigs(old, new *Config) []Change {
    var changes []Change
    if old.Server.Port != new.Server.Port {
        changes = append(changes, Change{
            Field: "server.port",
            Old:   old.Server.Port,
            New:   new.Server.Port,
        })
    }
    return changes
}
```

### Goal 5: Remote Configuration ⭐⭐⭐

Load config from remote sources (etcd, Consul, S3):

```go
func LoadConfigFromS3(bucket, key string) (*Config, error) {
    sess := session.Must(session.NewSession())
    svc := s3.New(sess)

    result, err := svc.GetObject(&s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return nil, err
    }

    data, _ := io.ReadAll(result.Body)
    // Parse and return config
}
```

---

## How to Run

```bash
# Run the demo program
cd /home/user/go-edu
make run P=38-config-loader-env-yaml

# Run tests
go test ./minis/38-config-loader-env-yaml/...

# Run specific test
go test ./minis/38-config-loader-env-yaml/exercise -run TestLoadConfig

# Test with environment variables
DB_HOST=localhost PORT=9000 go test ./minis/38-config-loader-env-yaml/...

# Verbose output
go test -v ./minis/38-config-loader-env-yaml/...
```

---

## Summary

**What you learned:**
- ✅ Configuration management best practices (12-factor app)
- ✅ YAML parsing with struct tags
- ✅ Environment variable substitution
- ✅ Configuration validation and defaults
- ✅ Type-safe configuration with custom unmarshaling
- ✅ Error handling for configuration problems

**Why this matters:**
Every production application needs configuration management. The patterns you learned here (YAML + env vars, validation, defaults) are used in virtually all modern cloud applications, from startups to large enterprises.

**Next steps:**
- Use these patterns in your own projects
- Explore configuration libraries: `viper`, `envconfig`
- Learn about secret management: Vault, AWS Secrets Manager
- Study infrastructure as code: Terraform, Kubernetes

**Key takeaway:** Good configuration management makes your application flexible, secure, and easy to deploy across multiple environments.
