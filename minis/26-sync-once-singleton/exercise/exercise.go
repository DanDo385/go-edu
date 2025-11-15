//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for sync.Once and singletons.
//
// LEARNING OBJECTIVES:
// - Use sync.Once for exactly-once initialization
// - Implement the singleton pattern correctly
// - Handle initialization errors properly
// - Create lazy initialization patterns
// - Avoid common concurrency pitfalls

package exercise

import "sync"

// ============================================================================
// EXERCISE 1: Basic sync.Once
// ============================================================================

// Counter with exactly-once initialization.
type Counter struct {
	value int
	once  sync.Once
}

// Initialize sets the counter value exactly once.
//
// REQUIREMENTS:
// - Use sync.Once to ensure initialization runs exactly once
// - Set c.value to the provided initialValue
// - Subsequent calls should be no-ops (value doesn't change)
//
// HINT: Use c.once.Do(func() { ... })
func (c *Counter) Initialize(initialValue int) {
	// TODO: Implement this method
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
	// TODO: Add necessary fields
	// You'll need:
	// - A *Config field to store the configuration
	// - A sync.Once to ensure one-time initialization
}

// Global config manager instance
var configManager = &ConfigManager{}

// GetConfig returns the singleton configuration.
//
// REQUIREMENTS:
// - Initialize the config exactly once using sync.Once
// - Return a *Config with:
//   - DatabaseURL: "postgres://localhost:5432/app"
//   - APIKey: "secret-key-123"
//   - Port: 8080
// - All calls should return the same *Config instance
//
// HINT: Use sync.Once.Do() to initialize the config field
func GetConfig() *Config {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 3: Database Singleton with Error Handling
// ============================================================================

// DatabaseManager manages a singleton database connection.
type DatabaseManager struct {
	// TODO: Add necessary fields
	// You'll need:
	// - A *Database field
	// - A sync.Once
	// - An error field to store initialization errors
}

var dbManager = &DatabaseManager{}

// GetDatabase returns the singleton database connection or an error.
//
// REQUIREMENTS:
// - Initialize the database exactly once
// - If connectionURL is empty, set initErr to an error and return it
// - If connectionURL is valid, create a Database with:
//   - URL: connectionURL
//   - Connected: true
// - Return the same error on subsequent calls if initialization failed
//
// HINT: Store the initialization error in a field and return it
func GetDatabase(connectionURL string) (*Database, error) {
	// TODO: Implement this function
	return nil, nil
}

// ============================================================================
// EXERCISE 4: Logger Singleton
// ============================================================================

var (
	// TODO: Add global variables for logger singleton
	// You'll need:
	// - logger *Logger
	// - loggerOnce sync.Once
)

// GetLogger returns the singleton logger instance.
//
// REQUIREMENTS:
// - Initialize the logger exactly once with name "AppLogger"
// - Create a new Logger: &Logger{Name: name, Output: []string{}}
// - Return the same instance on all calls
//
// HINT: Use package-level variables
func GetLogger() *Logger {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 5: Cache Singleton
// ============================================================================

var (
	// TODO: Add global variables for cache singleton
)

// GetCache returns the singleton cache instance.
//
// REQUIREMENTS:
// - Initialize the cache exactly once using NewCache()
// - Return the same instance on all calls
//
// HINT: Similar to logger, but call NewCache() for initialization
func GetCache() *Cache {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 6: Multiple Independent Singletons
// ============================================================================

// MetricsManager manages singleton metrics.
type MetricsManager struct {
	// TODO: Add fields for metrics singleton
}

var metricsManager = &MetricsManager{}

// GetMetrics returns the singleton metrics instance.
//
// REQUIREMENTS:
// - Initialize metrics exactly once with &Metrics{RequestCount: 0, ErrorCount: 0}
// - Return the same instance on all calls
//
// HINT: Same pattern as GetConfig
func GetMetrics() *Metrics {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 7: Lazy Field Initialization
// ============================================================================

// Application with lazily-initialized components.
type Application struct {
	// TODO: Add fields for lazy initialization
	// Each component needs:
	// - Its own sync.Once
	// - Storage for the component instance
	//
	// Components:
	// - database (*Database)
	// - logger (*Logger)
	// - cache (*Cache)
}

// GetDB returns the database, initializing it lazily.
//
// REQUIREMENTS:
// - Initialize database on first call with &Database{URL: "localhost:5432", Connected: true}
// - Return the same instance on subsequent calls
// - Use app's database-specific sync.Once
//
// HINT: Each component has its own Once field
func (app *Application) GetDB() *Database {
	// TODO: Implement this method
	return nil
}

// GetAppLogger returns the logger, initializing it lazily.
//
// REQUIREMENTS:
// - Initialize logger on first call with &Logger{Name: "AppLogger", Output: []string{}}
// - Return the same instance on subsequent calls
// - Use app's logger-specific sync.Once
func (app *Application) GetAppLogger() *Logger {
	// TODO: Implement this method
	return nil
}

// GetAppCache returns the cache, initializing it lazily.
//
// REQUIREMENTS:
// - Initialize cache on first call using NewCache()
// - Return the same instance on subsequent calls
// - Use app's cache-specific sync.Once
func (app *Application) GetAppCache() *Cache {
	// TODO: Implement this method
	return nil
}

// ============================================================================
// EXERCISE 8: Idempotent Initialization
// ============================================================================

// InitOnce is a wrapper for idempotent initialization.
type InitOnce struct {
	// TODO: Add fields for tracking initialization
	// You'll need:
	// - sync.Once
	// - A field to track if initialization completed
}

// Do runs f exactly once and marks initialization as complete.
//
// REQUIREMENTS:
// - Run f exactly once using sync.Once
// - After f completes (even if it panics), mark as initialized
//
// HINT: Use defer to ensure marking happens even on panic
func (io *InitOnce) Do(f func()) {
	// TODO: Implement this method
}

// IsInitialized returns whether initialization has completed.
//
// REQUIREMENTS:
// - Return true if Do() has been called and completed
// - Return false otherwise
//
// HINT: Check the initialized field you added
func (io *InitOnce) IsInitialized() bool {
	// TODO: Implement this method
	return false
}

// ============================================================================
// EXERCISE 9: Resettable Once (for testing)
// ============================================================================

// ResettableOnce is like sync.Once but can be reset for testing.
//
// WARNING: This is ONLY for tests. Not thread-safe during reset!
type ResettableOnce struct {
	// TODO: Add fields
	// You'll need:
	// - sync.Once
}

// Do runs f exactly once (until Reset is called).
func (ro *ResettableOnce) Do(f func()) {
	// TODO: Implement this method
}

// Reset resets the once so Do() will run again.
//
// WARNING: Not thread-safe! Only call in tests when no goroutines are running.
func (ro *ResettableOnce) Reset() {
	// TODO: Implement this method
	// HINT: Create a new sync.Once
}

// ============================================================================
// EXERCISE 10: Factory Singleton
// ============================================================================

// FactorySingleton manages a singleton created by a factory function.
type FactorySingleton struct {
	// TODO: Add fields
	// You'll need:
	// - instance interface{} (to store any type)
	// - once sync.Once
}

// GetOrCreate returns the singleton, creating it with factory if needed.
//
// REQUIREMENTS:
// - On first call, run factory() and store the result
// - On subsequent calls, return the stored instance
// - Use sync.Once to ensure factory runs exactly once
//
// HINT: Store the result of factory() in the instance field
func (fs *FactorySingleton) GetOrCreate(factory func() interface{}) interface{} {
	// TODO: Implement this method
	return nil
}
