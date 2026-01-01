//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package configloaderenvyaml

import "time"

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
// TODO: implement LoadConfig.
func LoadConfig(filename string) (*Config, error) { panic("TODO: implement") }
// TODO: implement substituteEnvVars.
func substituteEnvVars(input string) string { panic("TODO: implement") }
// TODO: implement ApplyDefaults.
func (c *Config) ApplyDefaults() { panic("TODO: implement") }
// TODO: implement Validate.
func (c *Config) Validate() error { panic("TODO: implement") }
