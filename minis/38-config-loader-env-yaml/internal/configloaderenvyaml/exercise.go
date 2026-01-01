//go:build !solution && !reference

package configloaderenvyaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
	"time"
)

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
	Level  string `yaml:"level"`  // debug, info, warn, error
	Format string `yaml:"format"` // json, text
	Output string `yaml:"output"` // stdout, stderr, file path
}

// LoadConfig implements the exercise.
//
// TODO: Implement this function
func LoadConfig(filename string) (*Config, error) {
	// TODO: Implement
	return nil, nil
}

// ApplyDefaults implements the exercise.
//
// TODO: Implement this function
func (c *Config) ApplyDefaults() {
	// TODO: Implement
}

// Validate implements the exercise.
//
// TODO: Implement this function
func (c *Config) Validate() error {
	// TODO: Implement
	return nil
}
