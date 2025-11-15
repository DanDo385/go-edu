//go:build solution
// +build solution

package exercise

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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

// substituteEnvVars replaces ${VAR} and ${VAR:-default} patterns with environment variable values.
//
// Patterns:
//   ${VAR}          - Replace with env var VAR, or leave as-is if not set
//   ${VAR:-default} - Replace with env var VAR, or use "default" if not set
func substituteEnvVars(input string) string {
	// Pattern explanation:
	// \$\{              - Literal ${
	// ([A-Z_][A-Z0-9_]*) - Variable name (capture group 1)
	//                     Must start with letter or underscore
	//                     Can contain letters, numbers, underscores
	// (:-([^}]*))?      - Optional default value (capture groups 2 and 3)
	//                     :- separator
	//                     ([^}]*) - Any characters except }
	// \}                - Literal }
	re := regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)(:-([^}]*))?\}`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract parts from the regex match
		matches := re.FindStringSubmatch(match)
		varName := matches[1]
		defaultValue := ""
		if len(matches) > 3 {
			defaultValue = matches[3]
		}

		// Try to get value from environment
		if value := os.Getenv(varName); value != "" {
			return value
		}

		// Use default if provided
		if defaultValue != "" {
			return defaultValue
		}

		// No replacement available - keep original
		return match
	})
}

// ApplyDefaults sets default values for any zero-value fields.
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
		c.Database.Port = 5432 // PostgreSQL default
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

// Validate checks that the configuration is valid and returns an error if not.
func (c *Config) Validate() error {
	var errors []string

	// Required fields
	if c.Database.Host == "" {
		errors = append(errors, "database.host is required")
	}
	if c.Database.Database == "" {
		errors = append(errors, "database.database is required")
	}

	// Port ranges
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errors = append(errors, "server.port must be between 1 and 65535")
	}
	if c.Database.Port < 1 || c.Database.Port > 65535 {
		errors = append(errors, "database.port must be between 1 and 65535")
	}

	// Connection pool
	if c.Database.MaxConns < 1 {
		errors = append(errors, "database.max_connections must be at least 1")
	}

	// Log level validation
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		errors = append(errors,
			"logging.level must be one of: debug, info, warn, error")
	}

	// Return combined errors
	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n  - %s",
			strings.Join(errors, "\n  - "))
	}

	return nil
}
