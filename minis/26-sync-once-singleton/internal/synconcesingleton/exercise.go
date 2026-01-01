//go:build !solution && !reference

// Package exercise contains hands-on exercises for sync.Once and singletons.
//
// LEARNING OBJECTIVES:
// - Use sync.Once for exactly-once initialization
// - Implement the singleton pattern correctly
// - Handle initialization errors properly
// - Create lazy initialization patterns
// - Avoid common concurrency pitfalls

package synconcesingleton

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
func (c *Counter) Initialize(initialValue int) {
	// TODO: Implement this method.

	// `sync.Once` is a synchronization primitive that will perform an action exactly once.
	// It's commonly used for thread-safe, lazy initialization.

	// Step 1: Use the `Do` method of the `sync.Once` field.
	// - `c.once.Do(func() { ... })`
	// - The function you pass to `Do` will be executed only on the very first call to `Initialize` across all goroutines.
	// - All subsequent calls to `Initialize` will do nothing, but they will wait for the first call's function to complete if it's still running.

	// Step 2: Inside the function, perform the initialization.
	// - `c.value = initialValue`
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
	// TODO: Implement this function.

	// This is the classic singleton pattern in Go. A package-level variable (`configManager`) holds the instance,
	// and a function (`GetConfig`) provides global access to it, ensuring it's initialized only once.

	// Step 1: Use the `Do` method of the global `configManager`'s `sync.Once` field.
	// - `configManager.once.Do(func() { ... })`

	// Step 2: Inside the function, create the configuration.
	// - This is where you would typically read from a file, parse environment variables, etc. For this exercise, you'll just create the struct.
	// - `configManager.cfg = &Config{...}`

	// Step 3: Return the config.
	// - `return configManager.cfg`
	// - Because `Do` blocks until the function completes, you are guaranteed that `configManager.cfg` is not nil when you return it, even on the very first call.
	return nil
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
	// TODO: Implement this function.

	// This pattern extends the singleton to handle initialization that can fail.
	// The key is to perform the initialization once, and if it fails, *always* return the same error on subsequent calls.

	// Step 1: Use `dbManager.once.Do`.
	// Step 2: Inside the `Do` function:
	//   - Check if the `connectionURL` is valid.
	//   - If it's invalid, create a new error and assign it to `dbManager.initErr`. Then `return`.
	//   - If it's valid, create the `*Database` instance and assign it to `dbManager.db`.
	// Step 3: After `Do` has run, return the results.
	//   - `return dbManager.db, dbManager.initErr`
	// - On the first call, `Do` runs and populates either `db` or `initErr`.
	// - On all other calls, `Do` does nothing, and the function simply returns the already-populated fields.
	return nil, nil
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
	// TODO: Implement this function.

	// This is a slightly different style of singleton that uses package-level variables directly instead of a manager struct. It's more direct for simple cases.

	// Step 1: Use the package-level `loggerOnce.Do`.
	// Step 2: Inside the `Do` function, initialize the package-level `logger` variable.
	// Step 3: Return the `logger` variable.
	return nil
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
	// TODO: Implement this function.
	// - This follows the same pattern as `GetLogger`.
	// - Use the package-level `cacheOnce.Do`.
	// - Inside the `Do` function, call `NewCache()` and assign the result to the package-level `cache` variable.
	// - Return the `cache` variable.
	return nil
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
	// TODO: Implement this function.
	// - This follows the same pattern as `GetConfig`.
	// - Use the global `metricsManager.once.Do`.
	// - Initialize `metricsManager.metrics`.
	// - Return `metricsManager.metrics`.
	return nil
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
	// TODO: Implement this method.

	// This pattern is "lazy initialization". The database connection is not created when the `Application` struct is created, but only when it's first requested.
	// This is useful for expensive resources that may not always be needed.

	// Step 1: Use the `sync.Once` specific to the database field.
	// - `app.dbOnce.Do(func() { ... })`

	// Step 2: Inside the `Do` function, initialize `app.db`.
	// Step 3: Return `app.db`.
	return nil
}

// GetAppLogger returns the logger, initializing it lazily.
func (app *Application) GetAppLogger() *Logger {
	// TODO: Implement this method.
	// - Follow the same lazy initialization pattern as `GetDB`, but use `app.loggerOnce` and initialize `app.logger`.
	return nil
}

// GetAppCache returns the cache, initializing it lazily.
func (app *Application) GetAppCache() *Cache {
	// TODO: Implement this method.
	// - Follow the same lazy initialization pattern, but use `app.cacheOnce` and initialize `app.cache` by calling `NewCache()`.
	return nil
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
	// TODO: Implement this method.

	// This extends `sync.Once` to also record that the initialization has finished.

	// Step 1: Use the `Do` method of the embedded `sync.Once`.
	// Step 2: Inside the `Do` function, you need to both run the user's function `f` and set your `initialized` flag.
	// - To ensure the flag is set *after* `f` completes (even if `f` panics), use `defer`.
	// - `defer atomic.StoreUint32(&io.initialized, 1)`
	// - Then, call the function: `f()`.
}

// IsInitialized returns whether initialization has completed.
func (io *InitOnce) IsInitialized() bool {
	// TODO: Implement this method.
	// - Atomically read the `initialized` flag.
	// - Use `atomic.LoadUint32(&io.initialized) == 1`.
	return false
}

// ============================================================================
// EXERCISE 9: Resettable Once (for testing)
// ============================================================================

// ResettableOnce is like sync.Once but can be reset for testing.
type ResettableOnce struct {
	once sync.Once
}

// Do runs f exactly once (until Reset is called).
func (ro *ResettableOnce) Do(f func()) {
	// TODO: Implement this method.
	// - This just calls the `Do` method of the embedded `sync.Once`.
}

// Reset resets the once so Do() will run again.
func (ro *ResettableOnce) Reset() {
	// TODO: Implement this method.
	// - A `sync.Once` cannot be truly "reset". Its internal state is not exposed.
	// - The trick is to replace the `sync.Once` value itself with a new, zero-value `sync.Once`.
	// - `ro.once = sync.Once{}`
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
	// TODO: Implement this method.

	// This pattern is a more generic version of the singleton. Instead of hard-coding the initialization logic,
	// it accepts a "factory" function that knows how to create the instance.

	// Step 1: Use `fs.once.Do`.
	// Step 2: Inside the `Do` function, call the `factory` function and store its result in `fs.instance`.
	// - `fs.instance = factory()`
	// Step 3: Return `fs.instance`.
	return nil
}
