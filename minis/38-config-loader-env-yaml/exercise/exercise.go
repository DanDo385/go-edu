//go:build !solution
// +build !solution

package exercise

import (
	"time"
)

// Config represents the complete application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	MaxConns int    `yaml:"max_connections"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`  // debug, info, warn, error
	Format string `yaml:"format"` // json, text
	Output string `yaml:"output"` // stdout, stderr, file path
}

// LoadConfig loads configuration from a YAML file.
// It performs the following steps:
// 1. Reads the YAML file
// 2. Substitutes environment variables (${VAR} or ${VAR:-default})
// 3. Parses the YAML into a Config struct
// 4. Applies default values for missing fields
// 5. Validates the configuration
//
// Example usage:
//   config, err := LoadConfig("config.yaml")
//   if err != nil {
//       log.Fatal(err)
//   }
func LoadConfig(filename string) (*Config, error) {
	// TODO: Implement
	// Steps:
	// 1. Read file with os.ReadFile
	// 2. Call substituteEnvVars on the file contents
	// 3. Unmarshal YAML into Config struct
	// 4. Call ApplyDefaults
	// 5. Call Validate
	// 6. Return config or error
	return nil, nil
}

// substituteEnvVars replaces ${VAR} and ${VAR:-default} patterns with environment variable values.
//
// Patterns:
//   ${VAR}          - Replace with env var VAR, or leave as-is if not set
//   ${VAR:-default} - Replace with env var VAR, or use "default" if not set
//
// Example:
//   input:  "host: ${DB_HOST:-localhost}"
//   output: "host: localhost" (if DB_HOST not set)
//   output: "host: prod-db" (if DB_HOST=prod-db)
func substituteEnvVars(input string) string {
	// TODO: Implement
	// Hints:
	// 1. Use regexp.MustCompile with pattern: `\$\{([A-Z_][A-Z0-9_]*)(:-([^}]*))?\}`
	// 2. Use ReplaceAllStringFunc to process each match
	// 3. Extract variable name and optional default from match
	// 4. Use os.Getenv to look up variable
	// 5. Return env var value, or default, or original match
	return input
}

// ApplyDefaults sets default values for any zero-value fields.
// This is called after parsing but before validation.
//
// Defaults:
//   Server.Host:         "0.0.0.0"
//   Server.Port:         8080
//   Server.ReadTimeout:  30s
//   Server.WriteTimeout: 30s
//   Database.Port:       5432
//   Database.MaxConns:   10
//   Logging.Level:       "info"
//   Logging.Format:      "json"
//   Logging.Output:      "stdout"
func (c *Config) ApplyDefaults() {
	// TODO: Implement
	// Check each field and set default if zero value
	// Use time.Second constants for durations
}

// Validate checks that the configuration is valid and returns an error if not.
// It checks:
//   - Required fields are not empty
//   - Port numbers are in valid range (1-65535)
//   - MaxConns is positive
//   - Log level is one of: debug, info, warn, error
//
// Returns a detailed error message listing all validation failures.
func (c *Config) Validate() error {
	// TODO: Implement
	// Steps:
	// 1. Create a slice to collect error messages
	// 2. Check all validation rules
	// 3. If errors exist, join them and return as error
	// 4. Otherwise return nil
	//
	// Validation rules:
	// - database.host must not be empty
	// - database.database must not be empty
	// - server.port must be 1-65535
	// - database.port must be 1-65535
	// - database.max_connections must be >= 1
	// - logging.level must be one of: debug, info, warn, error
	return nil
}
