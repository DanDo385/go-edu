package exercise

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSubstituteEnvVars tests environment variable substitution
func TestSubstituteEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "simple variable",
			input:    "host: ${DB_HOST}",
			envVars:  map[string]string{"DB_HOST": "localhost"},
			expected: "host: localhost",
		},
		{
			name:     "variable with default (var not set)",
			input:    "port: ${PORT:-8080}",
			envVars:  map[string]string{},
			expected: "port: 8080",
		},
		{
			name:     "variable with default (var set)",
			input:    "port: ${PORT:-8080}",
			envVars:  map[string]string{"PORT": "9000"},
			expected: "port: 9000",
		},
		{
			name:     "multiple variables",
			input:    "url: ${PROTOCOL}://${HOST}:${PORT}",
			envVars:  map[string]string{"PROTOCOL": "https", "HOST": "example.com", "PORT": "443"},
			expected: "url: https://example.com:443",
		},
		{
			name:     "undefined variable without default",
			input:    "value: ${UNDEFINED_VAR}",
			envVars:  map[string]string{},
			expected: "value: ${UNDEFINED_VAR}", // Stays as-is
		},
		{
			name:     "empty default value",
			input:    "value: ${VAR:-}",
			envVars:  map[string]string{},
			expected: "value: ",
		},
		{
			name:     "default with special chars",
			input:    "dsn: ${DSN:-user:pass@localhost/db}",
			envVars:  map[string]string{},
			expected: "dsn: user:pass@localhost/db",
		},
		{
			name:     "no substitution needed",
			input:    "plain: text",
			envVars:  map[string]string{},
			expected: "plain: text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Clear any variables we're testing as undefined
			if !strings.Contains(tt.name, "var set") {
				os.Unsetenv("UNDEFINED_VAR")
				os.Unsetenv("VAR")
			}

			result := substituteEnvVars(tt.input)
			if result != tt.expected {
				t.Errorf("substituteEnvVars(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestApplyDefaults tests that default values are correctly applied
func TestApplyDefaults(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		check  func(*testing.T, *Config)
	}{
		{
			name:   "empty config gets all defaults",
			config: Config{},
			check: func(t *testing.T, c *Config) {
				if c.Server.Host != "0.0.0.0" {
					t.Errorf("Server.Host = %q, want %q", c.Server.Host, "0.0.0.0")
				}
				if c.Server.Port != 8080 {
					t.Errorf("Server.Port = %d, want %d", c.Server.Port, 8080)
				}
				if c.Server.ReadTimeout != 30*time.Second {
					t.Errorf("Server.ReadTimeout = %v, want %v", c.Server.ReadTimeout, 30*time.Second)
				}
				if c.Server.WriteTimeout != 30*time.Second {
					t.Errorf("Server.WriteTimeout = %v, want %v", c.Server.WriteTimeout, 30*time.Second)
				}
				if c.Database.Port != 5432 {
					t.Errorf("Database.Port = %d, want %d", c.Database.Port, 5432)
				}
				if c.Database.MaxConns != 10 {
					t.Errorf("Database.MaxConns = %d, want %d", c.Database.MaxConns, 10)
				}
				if c.Logging.Level != "info" {
					t.Errorf("Logging.Level = %q, want %q", c.Logging.Level, "info")
				}
				if c.Logging.Format != "json" {
					t.Errorf("Logging.Format = %q, want %q", c.Logging.Format, "json")
				}
				if c.Logging.Output != "stdout" {
					t.Errorf("Logging.Output = %q, want %q", c.Logging.Output, "stdout")
				}
			},
		},
		{
			name: "existing values not overwritten",
			config: Config{
				Server: ServerConfig{
					Host:         "localhost",
					Port:         9000,
					ReadTimeout:  60 * time.Second,
					WriteTimeout: 60 * time.Second,
				},
				Database: DatabaseConfig{
					Port:     3306,
					MaxConns: 50,
				},
				Logging: LoggingConfig{
					Level:  "debug",
					Format: "text",
					Output: "stderr",
				},
			},
			check: func(t *testing.T, c *Config) {
				if c.Server.Host != "localhost" {
					t.Errorf("Server.Host = %q, want %q", c.Server.Host, "localhost")
				}
				if c.Server.Port != 9000 {
					t.Errorf("Server.Port = %d, want %d", c.Server.Port, 9000)
				}
				if c.Database.Port != 3306 {
					t.Errorf("Database.Port = %d, want %d", c.Database.Port, 3306)
				}
				if c.Logging.Level != "debug" {
					t.Errorf("Logging.Level = %q, want %q", c.Logging.Level, "debug")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.config
			config.ApplyDefaults()
			tt.check(t, &config)
		})
	}
}

// TestValidate tests configuration validation
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errText string // Part of error message to check for
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "mydb",
					MaxConns: 10,
				},
				Logging: LoggingConfig{
					Level: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "missing database host",
			config: Config{
				Server: ServerConfig{Port: 8080},
				Database: DatabaseConfig{
					Database: "mydb",
					Port:     5432,
					MaxConns: 10,
				},
				Logging: LoggingConfig{Level: "info"},
			},
			wantErr: true,
			errText: "host",
		},
		{
			name: "missing database name",
			config: Config{
				Server: ServerConfig{Port: 8080},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					MaxConns: 10,
				},
				Logging: LoggingConfig{Level: "info"},
			},
			wantErr: true,
			errText: "database",
		},
		{
			name: "invalid server port (too low)",
			config: Config{
				Server: ServerConfig{Port: 0},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "mydb",
					MaxConns: 10,
				},
				Logging: LoggingConfig{Level: "info"},
			},
			wantErr: true,
			errText: "port",
		},
		{
			name: "invalid server port (too high)",
			config: Config{
				Server: ServerConfig{Port: 99999},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "mydb",
					MaxConns: 10,
				},
				Logging: LoggingConfig{Level: "info"},
			},
			wantErr: true,
			errText: "port",
		},
		{
			name: "invalid database port",
			config: Config{
				Server: ServerConfig{Port: 8080},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     100000,
					Database: "mydb",
					MaxConns: 10,
				},
				Logging: LoggingConfig{Level: "info"},
			},
			wantErr: true,
			errText: "port",
		},
		{
			name: "invalid max connections",
			config: Config{
				Server: ServerConfig{Port: 8080},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "mydb",
					MaxConns: 0,
				},
				Logging: LoggingConfig{Level: "info"},
			},
			wantErr: true,
			errText: "max_connections",
		},
		{
			name: "invalid log level",
			config: Config{
				Server: ServerConfig{Port: 8080},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Database: "mydb",
					MaxConns: 10,
				},
				Logging: LoggingConfig{Level: "invalid"},
			},
			wantErr: true,
			errText: "level",
		},
		{
			name: "multiple validation errors",
			config: Config{
				Server: ServerConfig{Port: 99999},
				Database: DatabaseConfig{
					// Missing host and database
					Port:     5432,
					MaxConns: 0,
				},
				Logging: LoggingConfig{Level: "invalid"},
			},
			wantErr: true,
			errText: "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error containing %q, got nil", tt.errText)
				} else if !strings.Contains(err.Error(), tt.errText) {
					t.Errorf("Validate() error = %q, want error containing %q", err.Error(), tt.errText)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestLoadConfig tests loading configuration from YAML files
func TestLoadConfig(t *testing.T) {
	// Create temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		yaml     string
		envVars  map[string]string
		wantErr  bool
		validate func(*testing.T, *Config)
	}{
		{
			name: "basic valid config",
			yaml: `
server:
  host: localhost
  port: 8080
database:
  host: localhost
  port: 5432
  database: testdb
  username: user
  password: pass
  max_connections: 20
logging:
  level: debug
`,
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				if c.Server.Host != "localhost" {
					t.Errorf("Server.Host = %q, want %q", c.Server.Host, "localhost")
				}
				if c.Server.Port != 8080 {
					t.Errorf("Server.Port = %d, want %d", c.Server.Port, 8080)
				}
				if c.Database.Database != "testdb" {
					t.Errorf("Database.Database = %q, want %q", c.Database.Database, "testdb")
				}
				if c.Logging.Level != "debug" {
					t.Errorf("Logging.Level = %q, want %q", c.Logging.Level, "debug")
				}
			},
		},
		{
			name: "config with env var substitution",
			yaml: `
server:
  host: ${SERVER_HOST}
  port: ${SERVER_PORT}
database:
  host: ${DB_HOST}
  database: ${DB_NAME}
  password: ${DB_PASSWORD}
  max_connections: 15
logging:
  level: info
`,
			envVars: map[string]string{
				"SERVER_HOST": "0.0.0.0",
				"SERVER_PORT": "9000",
				"DB_HOST":     "postgres.example.com",
				"DB_NAME":     "production",
				"DB_PASSWORD": "secret123",
			},
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				if c.Server.Host != "0.0.0.0" {
					t.Errorf("Server.Host = %q, want %q", c.Server.Host, "0.0.0.0")
				}
				if c.Server.Port != 9000 {
					t.Errorf("Server.Port = %d, want %d", c.Server.Port, 9000)
				}
				if c.Database.Host != "postgres.example.com" {
					t.Errorf("Database.Host = %q, want %q", c.Database.Host, "postgres.example.com")
				}
				if c.Database.Password != "secret123" {
					t.Errorf("Database.Password = %q, want %q", c.Database.Password, "secret123")
				}
			},
		},
		{
			name: "config with defaults applied",
			yaml: `
server:
  port: 3000
database:
  host: localhost
  database: mydb
logging:
  level: warn
`,
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				// Check defaults were applied
				if c.Server.Host != "0.0.0.0" {
					t.Errorf("Server.Host default = %q, want %q", c.Server.Host, "0.0.0.0")
				}
				if c.Server.ReadTimeout != 30*time.Second {
					t.Errorf("Server.ReadTimeout default = %v, want %v", c.Server.ReadTimeout, 30*time.Second)
				}
				if c.Database.Port != 5432 {
					t.Errorf("Database.Port default = %d, want %d", c.Database.Port, 5432)
				}
				if c.Database.MaxConns != 10 {
					t.Errorf("Database.MaxConns default = %d, want %d", c.Database.MaxConns, 10)
				}
				// Check specified values not overwritten
				if c.Server.Port != 3000 {
					t.Errorf("Server.Port = %d, want %d", c.Server.Port, 3000)
				}
			},
		},
		{
			name: "config with env var defaults",
			yaml: `
server:
  host: ${HOST:-0.0.0.0}
  port: ${PORT:-8080}
database:
  host: localhost
  database: mydb
  password: ${DB_PASSWORD:-defaultpass}
logging:
  level: ${LOG_LEVEL:-info}
`,
			envVars: map[string]string{}, // No env vars set
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				if c.Server.Host != "0.0.0.0" {
					t.Errorf("Server.Host = %q, want %q", c.Server.Host, "0.0.0.0")
				}
				if c.Server.Port != 8080 {
					t.Errorf("Server.Port = %d, want %d", c.Server.Port, 8080)
				}
				if c.Database.Password != "defaultpass" {
					t.Errorf("Database.Password = %q, want %q", c.Database.Password, "defaultpass")
				}
				if c.Logging.Level != "info" {
					t.Errorf("Logging.Level = %q, want %q", c.Logging.Level, "info")
				}
			},
		},
		{
			name: "config with duration parsing",
			yaml: `
server:
  port: 8080
  read_timeout: 45s
  write_timeout: 1m30s
database:
  host: localhost
  database: mydb
logging:
  level: info
`,
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				if c.Server.ReadTimeout != 45*time.Second {
					t.Errorf("Server.ReadTimeout = %v, want %v", c.Server.ReadTimeout, 45*time.Second)
				}
				if c.Server.WriteTimeout != 90*time.Second {
					t.Errorf("Server.WriteTimeout = %v, want %v", c.Server.WriteTimeout, 90*time.Second)
				}
			},
		},
		{
			name: "invalid config - missing required fields",
			yaml: `
server:
  port: 8080
database:
  port: 5432
logging:
  level: info
`,
			wantErr: true,
		},
		{
			name: "invalid config - bad port number",
			yaml: `
server:
  port: 999999
database:
  host: localhost
  database: mydb
logging:
  level: info
`,
			wantErr: true,
		},
		{
			name: "invalid config - bad log level",
			yaml: `
server:
  port: 8080
database:
  host: localhost
  database: mydb
logging:
  level: invalid_level
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Clear potentially interfering env vars
			os.Unsetenv("HOST")
			os.Unsetenv("PORT")
			os.Unsetenv("LOG_LEVEL")

			// Write test file
			filename := filepath.Join(tmpDir, tt.name+".yaml")
			if err := os.WriteFile(filename, []byte(tt.yaml), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Load config
			config, err := LoadConfig(filename)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadConfig() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadConfig() unexpected error: %v", err)
				} else if tt.validate != nil {
					tt.validate(t, config)
				}
			}
		})
	}
}

// TestLoadConfigFileErrors tests error handling for file operations
func TestLoadConfigFileErrors(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "file does not exist",
			filename: "/nonexistent/config.yaml",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			filename: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadConfig(tt.filename)
			if tt.wantErr && err == nil {
				t.Error("LoadConfig() expected error for invalid file, got nil")
			}
		})
	}
}
