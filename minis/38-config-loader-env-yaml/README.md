# 38-config-loader-env-yaml

**Configuration Loading**

Load config from environment variables and YAML files.

## What You'll Learn

- YAML parsing
- Environment variable binding
- Configuration precedence
- Struct tags

## Functions to Implement

| Function | Description |
|----------|-------------|
| Load configuration | From env + YAML |

## Project Structure

```
38-config-loader-env-yaml/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── config/
│   ├── default.yaml     # Default config
│   └── production.yaml  # Production overrides
├── internal/configloaderenvyaml/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/38-config-loader-env-yaml

# Load config from file
go run ./cmd/app/main.go --config config/default.yaml

# Override with env var
PORT=9000 go run ./cmd/app/main.go --config config/default.yaml

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Load from file
go run ./cmd/app/main.go --config config/default.yaml

# With environment overrides
DATABASE_URL=postgres://localhost/db go run ./cmd/app/main.go

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **YAML Tags**: `yaml:"field_name"`
2. **Env Tags**: `env:"FIELD_NAME"`
3. **Precedence**: env > file > defaults
4. **Validation**: Required fields check

## Next Steps

After completing this exercise, proceed to `minis/39-sha256-hasher`.
