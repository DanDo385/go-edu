//go:build !solution && !reference

package configloaderenvyaml

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
func LoadConfig(filename string) (*Config, error) {
	// TODO: Implement the full configuration loading pipeline.

	// Step 1: Read the raw YAML file into a byte slice.
	// - Use `os.ReadFile(filename)`.
	// - Wrap any error with `fmt.Errorf` to add context.

	// Step 2: Substitute environment variables in the raw YAML content.
	// - Convert the byte slice to a string and pass it to your `substituteEnvVars` helper function.

	// Step 3: Parse the substituted YAML into the `Config` struct.
	// - `var config Config`
	// - `err := yaml.Unmarshal([]byte(substitutedYAML), &config)`
	// - The `yaml.v3` library will parse the YAML and populate the fields of your struct based on the `yaml:` tags.

	// Step 4: Apply default values for any fields that were not present in the YAML file.
	// - `config.ApplyDefaults()`

	// Step 5: Validate the final configuration to ensure all required fields are present and have sane values.
	// - `if err := config.Validate(); err != nil { ... }`

	// Step 6: If all steps succeed, return the populated and validated `*Config`.
	return nil, nil
}

// substituteEnvVars replaces ${VAR} and ${VAR:-default} patterns with environment variable values.
func substituteEnvVars(input string) string {
	// TODO: Implement environment variable substitution.

	// Step 1: Define the regular expression to find all occurrences of `${...}`.
	// - The pattern `\$\{([A-Z_][A-Z0-9_]*)(:-([^}]*))?\}` works well.
	//   - `\$\{ ... \}`: Matches the literal `${` and `}`.
	//   - `([A-Z_][A-Z0-9_]*)`: The first capture group. Matches a valid environment variable name (starts with letter or underscore, followed by letters, numbers, or underscores).
	//   - `(:-([^}]*))?`: An optional second capture group for the default value. It matches `:-` followed by any characters that are not `}`.

	// Step 2: Use `re.ReplaceAllStringFunc` to process each match.
	// - This function takes a callback that receives the matched string (e.g., "${DB_HOST:-localhost}").
	// - Inside the callback, you need to parse the matched string to get the variable name and the default value. `re.FindStringSubmatch` is perfect for this.

	// Step 3: Implement the logic inside the callback.
	// - Get the variable name from the first capture group.
	// - Get the default value from the (optional) third capture group.
	// - Call `os.Getenv()` with the variable name.
	// - If the environment variable exists and is not empty, return its value.
	// - Otherwise, if a default value was provided, return the default value.
	// - Otherwise, return the original matched string (e.g., "${DB_HOST}") so the user knows a variable was missing.
	return input
}

// ApplyDefaults sets default values for any zero-value fields.
func (c *Config) ApplyDefaults() {
	// TODO: Implement this method.

	// For each field in the `Config` struct (and its nested structs), check if it has its "zero value".
	// If it does, it means the value was not provided in the YAML file, so you should set it to a sensible default.

	// Example for Server.Host:
	// - `if c.Server.Host == "" { c.Server.Host = "0.0.0.0" }`

	// Example for Server.Port:
	// - `if c.Server.Port == 0 { c.Server.Port = 8080 }`

	// Do this for all the fields listed in the exercise description.
}

// Validate checks that the configuration is valid and returns an error if not.
func (c *Config) Validate() error {
	// TODO: Implement this method.

	// Step 1: Create a slice to hold any validation error messages.
	// - `var errors []string`

	// Step 2: Perform all the validation checks listed in the exercise.
	// - For required fields: `if c.Database.Host == "" { errors = append(errors, "database.host is required") }`
	// - For port ranges: `if c.Server.Port < 1 || c.Server.Port > 65535 { ... }`
	// - For enum-like fields (like `logging.level`), it's common to use a map for a clean lookup:
	//   ```
	//   validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	//   if !validLogLevels[c.Logging.Level] { ... }
	//   ```

	// Step 3: After all checks, see if you found any errors.
	// - `if len(errors) > 0`:
	//   - Join the errors together into a single, nicely formatted string (e.g., separated by newlines).
	//   - Return a new error containing this string: `return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, ", "))`

	// Step 4: If there were no errors, return `nil`.
	return nil
}
