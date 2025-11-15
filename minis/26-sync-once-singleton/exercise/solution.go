//go:build solution
// +build solution

// Package exercise contains solutions for sync.Once and singleton exercises.

package exercise

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ============================================================================
// EXERCISE 1: Basic sync.Once
// ============================================================================

// Counter with exactly-once initialization.
type Counter struct {
	value int
	once  sync.Once
}

// Initialize sets the counter value exactly once.
func (c *Counter) Initialize(initialValue int) {
	c.once.Do(func() {
		c.value = initialValue
	})
}

// GetValue returns the current counter value.
func (c *Counter) GetValue() int {
	return c.value
}

// ============================================================================
// EXERCISE 2: Configuration Singleton
// ============================================================================

// ConfigManager manages a singleton configuration.
type ConfigManager struct {
	cfg  *Config
	once sync.Once
}

// Global config manager instance
var configManager = &ConfigManager{}

// GetConfig returns the singleton configuration.
func GetConfig() *Config {
	configManager.once.Do(func() {
		configManager.cfg = &Config{
			DatabaseURL: "postgres://localhost:5432/app",
			APIKey:      "secret-key-123",
			Port:        8080,
		}
	})
	return configManager.cfg
}

// ============================================================================
// EXERCISE 3: Database Singleton with Error Handling
// ============================================================================

// DatabaseManager manages a singleton database connection.
type DatabaseManager struct {
	db      *Database
	once    sync.Once
	initErr error
}

var dbManager = &DatabaseManager{}

// GetDatabase returns the singleton database connection or an error.
func GetDatabase(connectionURL string) (*Database, error) {
	dbManager.once.Do(func() {
		if connectionURL == "" {
			dbManager.initErr = errors.New("connection URL cannot be empty")
			return
		}

		dbManager.db = &Database{
			URL:       connectionURL,
			Connected: true,
		}
	})

	return dbManager.db, dbManager.initErr
}

// ============================================================================
// EXERCISE 4: Logger Singleton
// ============================================================================

var (
	logger     *Logger
	loggerOnce sync.Once
)

// GetLogger returns the singleton logger instance.
func GetLogger() *Logger {
	loggerOnce.Do(func() {
		logger = &Logger{
			Name:   "AppLogger",
			Output: []string{},
		}
	})
	return logger
}

// ============================================================================
// EXERCISE 5: Cache Singleton
// ============================================================================

var (
	cache     *Cache
	cacheOnce sync.Once
)

// GetCache returns the singleton cache instance.
func GetCache() *Cache {
	cacheOnce.Do(func() {
		cache = NewCache()
	})
	return cache
}

// ============================================================================
// EXERCISE 6: Multiple Independent Singletons
// ============================================================================

// MetricsManager manages singleton metrics.
type MetricsManager struct {
	metrics *Metrics
	once    sync.Once
}

var metricsManager = &MetricsManager{}

// GetMetrics returns the singleton metrics instance.
func GetMetrics() *Metrics {
	metricsManager.once.Do(func() {
		metricsManager.metrics = &Metrics{
			RequestCount: 0,
			ErrorCount:   0,
		}
	})
	return metricsManager.metrics
}

// ============================================================================
// EXERCISE 7: Lazy Field Initialization
// ============================================================================

// Application with lazily-initialized components.
type Application struct {
	dbOnce sync.Once
	db     *Database

	loggerOnce sync.Once
	logger     *Logger

	cacheOnce sync.Once
	cache     *Cache
}

// GetDB returns the database, initializing it lazily.
func (app *Application) GetDB() *Database {
	app.dbOnce.Do(func() {
		app.db = &Database{
			URL:       "localhost:5432",
			Connected: true,
		}
	})
	return app.db
}

// GetAppLogger returns the logger, initializing it lazily.
func (app *Application) GetAppLogger() *Logger {
	app.loggerOnce.Do(func() {
		app.logger = &Logger{
			Name:   "AppLogger",
			Output: []string{},
		}
	})
	return app.logger
}

// GetAppCache returns the cache, initializing it lazily.
func (app *Application) GetAppCache() *Cache {
	app.cacheOnce.Do(func() {
		app.cache = NewCache()
	})
	return app.cache
}

// ============================================================================
// EXERCISE 8: Idempotent Initialization
// ============================================================================

// InitOnce is a wrapper for idempotent initialization.
type InitOnce struct {
	once        sync.Once
	initialized uint32 // atomic flag
}

// Do runs f exactly once and marks initialization as complete.
func (io *InitOnce) Do(f func()) {
	io.once.Do(func() {
		defer atomic.StoreUint32(&io.initialized, 1)
		f()
	})
}

// IsInitialized returns whether initialization has completed.
func (io *InitOnce) IsInitialized() bool {
	return atomic.LoadUint32(&io.initialized) == 1
}

// ============================================================================
// EXERCISE 9: Resettable Once (for testing)
// ============================================================================

// ResettableOnce is like sync.Once but can be reset for testing.
//
// WARNING: This is ONLY for tests. Not thread-safe during reset!
type ResettableOnce struct {
	once sync.Once
}

// Do runs f exactly once (until Reset is called).
func (ro *ResettableOnce) Do(f func()) {
	ro.once.Do(f)
}

// Reset resets the once so Do() will run again.
//
// WARNING: Not thread-safe! Only call in tests when no goroutines are running.
func (ro *ResettableOnce) Reset() {
	ro.once = sync.Once{}
}

// ============================================================================
// EXERCISE 10: Factory Singleton
// ============================================================================

// FactorySingleton manages a singleton created by a factory function.
type FactorySingleton struct {
	instance interface{}
	once     sync.Once
}

// GetOrCreate returns the singleton, creating it with factory if needed.
func (fs *FactorySingleton) GetOrCreate(factory func() interface{}) interface{} {
	fs.once.Do(func() {
		fs.instance = factory()
	})
	return fs.instance
}
