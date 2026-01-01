package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/example/go-10x-minis/minis/38-config-loader-env-yaml/exercise"
)

func main() {
	fmt.Println("=== Configuration Loader Demo ===\n")

	// Determine the testdata directory
	// When run from project root: minis/38-.../testdata
	testdataDir := filepath.Join("minis", "38-config-loader-env-yaml", "testdata")

	// Demo 1: Load basic configuration
	fmt.Println("1. Loading basic configuration from YAML...")
	fmt.Println("   File: testdata/basic.yaml")
	basicConfig, err := exercise.LoadConfig(filepath.Join(testdataDir, "basic.yaml"))
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Loaded successfully!\n")
		fmt.Printf("   Server: %s:%d\n", basicConfig.Server.Host, basicConfig.Server.Port)
		fmt.Printf("   Database: %s:%d (max_conns: %d)\n",
			basicConfig.Database.Host, basicConfig.Database.Port, basicConfig.Database.MaxConns)
		fmt.Printf("   Logging: level=%s, format=%s\n\n",
			basicConfig.Logging.Level, basicConfig.Logging.Format)
	}

	// Demo 2: Environment variable substitution
	fmt.Println("2. Testing environment variable substitution...")
	fmt.Println("   File: testdata/with-env-vars.yaml")
	fmt.Println("   Setting environment variables:")

	// Set some environment variables
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "9000")
	os.Setenv("DB_HOST", "postgres.example.com")
	os.Setenv("DB_PASSWORD", "super_secret_password")

	fmt.Println("   - APP_HOST=0.0.0.0")
	fmt.Println("   - APP_PORT=9000")
	fmt.Println("   - DB_HOST=postgres.example.com")
	fmt.Println("   - DB_PASSWORD=super_secret_password")

	envConfig, err := exercise.LoadConfig(filepath.Join(testdataDir, "with-env-vars.yaml"))
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Loaded with environment substitution!\n")
		fmt.Printf("   Server: %s:%d\n", envConfig.Server.Host, envConfig.Server.Port)
		fmt.Printf("   Database: %s (password: %s)\n\n",
			envConfig.Database.Host, maskPassword(envConfig.Database.Password))
	}

	// Demo 3: Default values
	fmt.Println("3. Testing default values...")
	fmt.Println("   File: testdata/minimal.yaml (only required fields)")

	minimalConfig, err := exercise.LoadConfig(filepath.Join(testdataDir, "minimal.yaml"))
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Defaults applied successfully!\n")
		fmt.Printf("   Server host: %s (default)\n", minimalConfig.Server.Host)
		fmt.Printf("   Server port: %d (default)\n", minimalConfig.Server.Port)
		fmt.Printf("   Read timeout: %v (default)\n", minimalConfig.Server.ReadTimeout)
		fmt.Printf("   Write timeout: %v (default)\n", minimalConfig.Server.WriteTimeout)
		fmt.Printf("   Database max connections: %d (default)\n", minimalConfig.Database.MaxConns)
		fmt.Printf("   Logging level: %s (default)\n\n", minimalConfig.Logging.Level)
	}

	// Demo 4: Validation errors
	fmt.Println("4. Testing validation (should fail)...")
	fmt.Println("   File: testdata/invalid.yaml")

	invalidConfig, err := exercise.LoadConfig(filepath.Join(testdataDir, "invalid.yaml"))
	if err != nil {
		fmt.Printf("   ✓ Validation caught errors as expected:\n")
		fmt.Printf("   %v\n\n", err)
	} else {
		fmt.Printf("   ✗ Unexpected: config should have failed validation\n")
		fmt.Printf("   Config: %+v\n\n", invalidConfig)
	}

	// Demo 5: Environment variable with defaults
	fmt.Println("5. Testing environment variable defaults (${VAR:-default})...")
	fmt.Println("   File: testdata/env-with-defaults.yaml")

	// Unset a variable to test default
	os.Unsetenv("UNDEFINED_VAR")
	fmt.Println("   UNDEFINED_VAR is not set, should use default value")

	defaultEnvConfig, err := exercise.LoadConfig(filepath.Join(testdataDir, "env-with-defaults.yaml"))
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Default values from env var syntax applied!\n")
		fmt.Printf("   Server port: %d (from ${UNDEFINED_PORT:-8080})\n", defaultEnvConfig.Server.Port)
		fmt.Printf("   Log level: %s (from ${UNDEFINED_LOG_LEVEL:-info})\n\n", defaultEnvConfig.Logging.Level)
	}

	// Demo 6: Complete configuration with all fields
	fmt.Println("6. Loading complete production-like configuration...")
	fmt.Println("   File: testdata/production.yaml")

	// Set production environment variables
	os.Setenv("PROD_DB_HOST", "prod-db.example.com")
	os.Setenv("PROD_DB_PASSWORD", "prod_secure_password_123")
	os.Setenv("API_KEY", "sk_live_abc123def456")

	prodConfig, err := exercise.LoadConfig(filepath.Join(testdataDir, "production.yaml"))
	if err != nil {
		log.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Production config loaded!\n")
		fmt.Printf("   Server: %s:%d\n", prodConfig.Server.Host, prodConfig.Server.Port)
		fmt.Printf("   Timeouts: read=%v, write=%v\n",
			prodConfig.Server.ReadTimeout, prodConfig.Server.WriteTimeout)
		fmt.Printf("   Database: %s:%d/%s (user: %s)\n",
			prodConfig.Database.Host, prodConfig.Database.Port,
			prodConfig.Database.Database, prodConfig.Database.Username)
		fmt.Printf("   Max connections: %d\n", prodConfig.Database.MaxConns)
		fmt.Printf("   Logging: level=%s, format=%s, output=%s\n\n",
			prodConfig.Logging.Level, prodConfig.Logging.Format, prodConfig.Logging.Output)
	}

	// Summary
	fmt.Println("=== Demo Complete ===")
	fmt.Println("\nKey takeaways:")
	fmt.Println("  • YAML files provide structured configuration")
	fmt.Println("  • Environment variables allow runtime customization")
	fmt.Println("  • ${VAR} syntax enables env var substitution in YAML")
	fmt.Println("  • ${VAR:-default} provides fallback values")
	fmt.Println("  • Validation catches configuration errors early")
	fmt.Println("  • Defaults make configuration less verbose")
	fmt.Println("\nTry modifying the YAML files in testdata/ and re-running!")
}

// maskPassword replaces all but the first 3 characters with asterisks
func maskPassword(password string) string {
	if len(password) <= 3 {
		return "***"
	}
	return password[:3] + "***"
}
