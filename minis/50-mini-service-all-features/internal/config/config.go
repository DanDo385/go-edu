package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Logging   LoggingConfig   `yaml:"logging"`
	CORS      CORSConfig      `yaml:"cors"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	JWT       JWTConfig       `yaml:"jwt"`
}

type ServerConfig struct {
	Addr            string        `yaml:"addr"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

type RateLimitConfig struct {
	RequestsPerSecond float64 `yaml:"requests_per_second"`
	Burst             int     `yaml:"burst"`
}

type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
}

// Load reads config from file and env vars
func Load(configPath string) (*Config, error) {
	// Read YAML file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Override with environment variables
	if addr := os.Getenv("SERVER_ADDR"); addr != "" {
		cfg.Server.Addr = addr
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.JWT.Secret = jwtSecret
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Logging.Level = logLevel
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Addr == "" {
		return fmt.Errorf("server.addr is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if c.JWT.Secret == "your-secret-key-change-in-production" {
		fmt.Println("WARNING: Using default JWT secret. Change this in production!")
	}
	return nil
}
