// Package exercise contains types and interfaces for sync.Once exercises.

package exercise

import "sync"

// Stringer is the standard string representation interface.
type Stringer interface {
	String() string
}

// Config represents application configuration.
type Config struct {
	DatabaseURL string
	APIKey      string
	Port        int
}

// Database represents a database connection.
type Database struct {
	URL       string
	Connected bool
}

// Logger represents a logging system.
type Logger struct {
	Name   string
	Output []string
}

// Write writes a log message.
func (l *Logger) Write(msg string) {
	l.Output = append(l.Output, msg)
}

// Cache represents an in-memory cache.
type Cache struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewCache creates a new cache.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

// Get retrieves a value from the cache.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

// Set stores a value in the cache.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Metrics represents a metrics collector.
type Metrics struct {
	RequestCount uint64
	ErrorCount   uint64
}
