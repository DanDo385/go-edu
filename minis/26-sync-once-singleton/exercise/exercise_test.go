package exercise

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// EXERCISE 1: Basic sync.Once
// ============================================================================

func TestCounter_Initialize(t *testing.T) {
	t.Run("single initialization", func(t *testing.T) {
		counter := &Counter{}

		counter.Initialize(42)
		if got := counter.GetValue(); got != 42 {
			t.Errorf("GetValue() = %d, want 42", got)
		}
	})

	t.Run("multiple calls", func(t *testing.T) {
		counter := &Counter{}

		counter.Initialize(10)
		counter.Initialize(20) // Should be ignored
		counter.Initialize(30) // Should be ignored

		if got := counter.GetValue(); got != 10 {
			t.Errorf("GetValue() = %d, want 10 (first value)", got)
		}
	})

	t.Run("concurrent initialization", func(t *testing.T) {
		counter := &Counter{}
		var wg sync.WaitGroup
		var initCount int32

		// Launch 100 goroutines
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				counter.Initialize(val)
				atomic.AddInt32(&initCount, 1)
			}(i)
		}

		wg.Wait()

		// Value should be from one of the goroutines (exactly once)
		if got := counter.GetValue(); got < 0 || got >= 100 {
			t.Errorf("GetValue() = %d, want value in range [0, 99]", got)
		}

		// All 100 goroutines should have completed
		if initCount != 100 {
			t.Errorf("initCount = %d, want 100", initCount)
		}
	})
}

// ============================================================================
// EXERCISE 2: Configuration Singleton
// ============================================================================

func TestGetConfig(t *testing.T) {
	t.Run("returns configuration", func(t *testing.T) {
		cfg := GetConfig()

		if cfg == nil {
			t.Fatal("GetConfig() returned nil")
		}

		if cfg.DatabaseURL != "postgres://localhost:5432/app" {
			t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "postgres://localhost:5432/app")
		}

		if cfg.APIKey != "secret-key-123" {
			t.Errorf("APIKey = %q, want %q", cfg.APIKey, "secret-key-123")
		}

		if cfg.Port != 8080 {
			t.Errorf("Port = %d, want 8080", cfg.Port)
		}
	})

	t.Run("returns same instance", func(t *testing.T) {
		cfg1 := GetConfig()
		cfg2 := GetConfig()

		if cfg1 != cfg2 {
			t.Error("GetConfig() returned different instances")
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		var wg sync.WaitGroup
		instances := make([]*Config, 100)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				instances[idx] = GetConfig()
			}(i)
		}

		wg.Wait()

		// All should be the same instance
		first := instances[0]
		for i, cfg := range instances {
			if cfg != first {
				t.Errorf("instances[%d] is different from instances[0]", i)
			}
		}
	})
}

// ============================================================================
// EXERCISE 3: Database Singleton with Error Handling
// ============================================================================

func TestGetDatabase(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		// Reset for test (this is a test-only hack)
		dbManager = &DatabaseManager{}

		db, err := GetDatabase("postgres://localhost:5432/testdb")

		if err != nil {
			t.Fatalf("GetDatabase() error = %v, want nil", err)
		}

		if db == nil {
			t.Fatal("GetDatabase() returned nil database")
		}

		if db.URL != "postgres://localhost:5432/testdb" {
			t.Errorf("URL = %q, want %q", db.URL, "postgres://localhost:5432/testdb")
		}

		if !db.Connected {
			t.Error("Connected = false, want true")
		}
	})

	t.Run("empty URL error", func(t *testing.T) {
		// Reset for test
		dbManager = &DatabaseManager{}

		db, err := GetDatabase("")

		if err == nil {
			t.Error("GetDatabase(\"\") error = nil, want error")
		}

		if db != nil {
			t.Errorf("GetDatabase(\"\") returned non-nil database: %v", db)
		}
	})

	t.Run("returns same error on retry", func(t *testing.T) {
		// Reset for test
		dbManager = &DatabaseManager{}

		// First call with empty URL
		_, err1 := GetDatabase("")
		if err1 == nil {
			t.Fatal("First call: expected error, got nil")
		}

		// Second call (even with valid URL, should return same error)
		_, err2 := GetDatabase("valid-url")
		if err2 == nil {
			t.Fatal("Second call: expected error, got nil")
		}

		// Should be the same error
		if err1.Error() != err2.Error() {
			t.Errorf("Different errors: %v vs %v", err1, err2)
		}
	})

	t.Run("concurrent access with valid URL", func(t *testing.T) {
		// Reset for test
		dbManager = &DatabaseManager{}

		var wg sync.WaitGroup
		dbs := make([]*Database, 50)
		errs := make([]error, 50)

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				dbs[idx], errs[idx] = GetDatabase("postgres://localhost:5432/testdb")
			}(i)
		}

		wg.Wait()

		// All should succeed
		for i, err := range errs {
			if err != nil {
				t.Errorf("goroutine %d: unexpected error: %v", i, err)
			}
		}

		// All should be same instance
		first := dbs[0]
		for i, db := range dbs {
			if db != first {
				t.Errorf("dbs[%d] is different from dbs[0]", i)
			}
		}
	})
}

// ============================================================================
// EXERCISE 4: Logger Singleton
// ============================================================================

func TestGetLogger(t *testing.T) {
	t.Run("returns logger", func(t *testing.T) {
		logger := GetLogger()

		if logger == nil {
			t.Fatal("GetLogger() returned nil")
		}

		if logger.Name != "AppLogger" {
			t.Errorf("Name = %q, want %q", logger.Name, "AppLogger")
		}

		if logger.Output == nil {
			t.Error("Output is nil, want empty slice")
		}
	})

	t.Run("returns same instance", func(t *testing.T) {
		log1 := GetLogger()
		log2 := GetLogger()

		if log1 != log2 {
			t.Error("GetLogger() returned different instances")
		}
	})

	t.Run("shared state", func(t *testing.T) {
		log1 := GetLogger()
		log1.Write("message 1")

		log2 := GetLogger()
		log2.Write("message 2")

		// Both should see both messages (same instance)
		if len(log1.Output) != 2 {
			t.Errorf("log1.Output length = %d, want 2", len(log1.Output))
		}

		if len(log2.Output) != 2 {
			t.Errorf("log2.Output length = %d, want 2", len(log2.Output))
		}
	})
}

// ============================================================================
// EXERCISE 5: Cache Singleton
// ============================================================================

func TestGetCache(t *testing.T) {
	t.Run("returns cache", func(t *testing.T) {
		cache := GetCache()

		if cache == nil {
			t.Fatal("GetCache() returned nil")
		}
	})

	t.Run("returns same instance", func(t *testing.T) {
		cache1 := GetCache()
		cache2 := GetCache()

		if cache1 != cache2 {
			t.Error("GetCache() returned different instances")
		}
	})

	t.Run("shared state", func(t *testing.T) {
		cache1 := GetCache()
		cache1.Set("key", "value")

		cache2 := GetCache()
		val, ok := cache2.Get("key")

		if !ok {
			t.Error("cache2.Get(\"key\") not found")
		}

		if val != "value" {
			t.Errorf("cache2.Get(\"key\") = %q, want %q", val, "value")
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		cache := GetCache()
		var wg sync.WaitGroup

		// Multiple writers
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				c := GetCache()
				c.Set(fmt.Sprintf("key%d", idx), fmt.Sprintf("value%d", idx))
			}(i)
		}

		wg.Wait()

		// Verify all writes succeeded (same cache instance)
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("key%d", i)
			val, ok := cache.Get(key)
			if !ok {
				t.Errorf("cache.Get(%q) not found", key)
			}
			expectedVal := fmt.Sprintf("value%d", i)
			if val != expectedVal {
				t.Errorf("cache.Get(%q) = %q, want %q", key, val, expectedVal)
			}
		}
	})
}

// ============================================================================
// EXERCISE 6: Multiple Independent Singletons
// ============================================================================

func TestGetMetrics(t *testing.T) {
	t.Run("returns metrics", func(t *testing.T) {
		metrics := GetMetrics()

		if metrics == nil {
			t.Fatal("GetMetrics() returned nil")
		}

		if metrics.RequestCount != 0 {
			t.Errorf("RequestCount = %d, want 0", metrics.RequestCount)
		}

		if metrics.ErrorCount != 0 {
			t.Errorf("ErrorCount = %d, want 0", metrics.ErrorCount)
		}
	})

	t.Run("returns same instance", func(t *testing.T) {
		m1 := GetMetrics()
		m2 := GetMetrics()

		if m1 != m2 {
			t.Error("GetMetrics() returned different instances")
		}
	})

	t.Run("independent from other singletons", func(t *testing.T) {
		// GetMetrics should not affect GetLogger or GetCache
		metrics := GetMetrics()
		logger := GetLogger()
		cache := GetCache()

		if metrics == nil || logger == nil || cache == nil {
			t.Error("One or more singletons returned nil")
		}

		// They should be different types/instances
		// (This test mainly checks that they don't interfere)
	})
}

// ============================================================================
// EXERCISE 7: Lazy Field Initialization
// ============================================================================

func TestApplication_LazyFields(t *testing.T) {
	t.Run("independent initialization", func(t *testing.T) {
		app := &Application{}

		// Access logger first
		logger := app.GetAppLogger()
		if logger == nil {
			t.Fatal("GetAppLogger() returned nil")
		}

		// Database and cache should still be uninitialized
		// (We can't directly test this without reflection, but we can verify behavior)

		// Access database
		db := app.GetDB()
		if db == nil {
			t.Fatal("GetDB() returned nil")
		}

		// Access cache
		cache := app.GetAppCache()
		if cache == nil {
			t.Fatal("GetAppCache() returned nil")
		}
	})

	t.Run("same instance per field", func(t *testing.T) {
		app := &Application{}

		db1 := app.GetDB()
		db2 := app.GetDB()
		if db1 != db2 {
			t.Error("GetDB() returned different instances")
		}

		log1 := app.GetAppLogger()
		log2 := app.GetAppLogger()
		if log1 != log2 {
			t.Error("GetAppLogger() returned different instances")
		}

		cache1 := app.GetAppCache()
		cache2 := app.GetAppCache()
		if cache1 != cache2 {
			t.Error("GetAppCache() returned different instances")
		}
	})

	t.Run("concurrent access to different fields", func(t *testing.T) {
		app := &Application{}
		var wg sync.WaitGroup

		// Access DB from multiple goroutines
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				db := app.GetDB()
				if db == nil {
					t.Error("GetDB() returned nil")
				}
			}()
		}

		// Access Logger from multiple goroutines
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				logger := app.GetAppLogger()
				if logger == nil {
					t.Error("GetAppLogger() returned nil")
				}
			}()
		}

		wg.Wait()
	})
}

// ============================================================================
// EXERCISE 8: Idempotent Initialization
// ============================================================================

func TestInitOnce(t *testing.T) {
	t.Run("runs exactly once", func(t *testing.T) {
		io := &InitOnce{}
		var count int

		io.Do(func() { count++ })
		io.Do(func() { count++ })
		io.Do(func() { count++ })

		if count != 1 {
			t.Errorf("function ran %d times, want 1", count)
		}
	})

	t.Run("IsInitialized reflects state", func(t *testing.T) {
		io := &InitOnce{}

		if io.IsInitialized() {
			t.Error("IsInitialized() = true before Do(), want false")
		}

		io.Do(func() {})

		if !io.IsInitialized() {
			t.Error("IsInitialized() = false after Do(), want true")
		}
	})

	t.Run("marks initialized even on panic", func(t *testing.T) {
		io := &InitOnce{}

		func() {
			defer func() { recover() }()
			io.Do(func() { panic("test panic") })
		}()

		if !io.IsInitialized() {
			t.Error("IsInitialized() = false after panic, want true")
		}

		// Subsequent calls should not run
		var count int
		io.Do(func() { count++ })
		if count != 0 {
			t.Errorf("function ran after panic, count = %d, want 0", count)
		}
	})
}

// ============================================================================
// EXERCISE 9: Resettable Once
// ============================================================================

func TestResettableOnce(t *testing.T) {
	t.Run("runs exactly once before reset", func(t *testing.T) {
		ro := &ResettableOnce{}
		var count int

		ro.Do(func() { count++ })
		ro.Do(func() { count++ })

		if count != 1 {
			t.Errorf("function ran %d times, want 1", count)
		}
	})

	t.Run("can reset and run again", func(t *testing.T) {
		ro := &ResettableOnce{}
		var count int

		ro.Do(func() { count++ })
		if count != 1 {
			t.Errorf("first run: count = %d, want 1", count)
		}

		ro.Reset()

		ro.Do(func() { count++ })
		if count != 2 {
			t.Errorf("after reset: count = %d, want 2", count)
		}
	})

	t.Run("multiple resets", func(t *testing.T) {
		ro := &ResettableOnce{}
		values := []int{}

		ro.Do(func() { values = append(values, 1) })
		ro.Reset()
		ro.Do(func() { values = append(values, 2) })
		ro.Reset()
		ro.Do(func() { values = append(values, 3) })

		if len(values) != 3 {
			t.Errorf("len(values) = %d, want 3", len(values))
		}

		expected := []int{1, 2, 3}
		for i, v := range values {
			if v != expected[i] {
				t.Errorf("values[%d] = %d, want %d", i, v, expected[i])
			}
		}
	})
}

// ============================================================================
// EXERCISE 10: Factory Singleton
// ============================================================================

func TestFactorySingleton(t *testing.T) {
	t.Run("creates with factory", func(t *testing.T) {
		fs := &FactorySingleton{}

		instance := fs.GetOrCreate(func() interface{} {
			return &Config{Port: 9000}
		})

		cfg, ok := instance.(*Config)
		if !ok {
			t.Fatal("instance is not *Config")
		}

		if cfg.Port != 9000 {
			t.Errorf("Port = %d, want 9000", cfg.Port)
		}
	})

	t.Run("returns same instance", func(t *testing.T) {
		fs := &FactorySingleton{}
		var callCount int

		factory := func() interface{} {
			callCount++
			return &Config{Port: 9000}
		}

		instance1 := fs.GetOrCreate(factory)
		instance2 := fs.GetOrCreate(factory)

		if instance1 != instance2 {
			t.Error("GetOrCreate() returned different instances")
		}

		if callCount != 1 {
			t.Errorf("factory called %d times, want 1", callCount)
		}
	})

	t.Run("concurrent creation", func(t *testing.T) {
		fs := &FactorySingleton{}
		var callCount int32

		factory := func() interface{} {
			atomic.AddInt32(&callCount, 1)
			time.Sleep(10 * time.Millisecond) // Simulate expensive creation
			return &Config{Port: 9000}
		}

		var wg sync.WaitGroup
		instances := make([]interface{}, 100)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				instances[idx] = fs.GetOrCreate(factory)
			}(i)
		}

		wg.Wait()

		// Factory should be called exactly once
		if callCount != 1 {
			t.Errorf("factory called %d times, want 1", callCount)
		}

		// All instances should be the same
		first := instances[0]
		for i, inst := range instances {
			if inst != first {
				t.Errorf("instances[%d] is different from instances[0]", i)
			}
		}
	})

	t.Run("different factories different singletons", func(t *testing.T) {
		fs1 := &FactorySingleton{}
		fs2 := &FactorySingleton{}

		inst1 := fs1.GetOrCreate(func() interface{} { return &Config{Port: 8080} })
		inst2 := fs2.GetOrCreate(func() interface{} { return &Config{Port: 9090} })

		// Different FactorySingleton instances should create different singletons
		if inst1 == inst2 {
			t.Error("Different FactorySingleton instances returned same singleton")
		}
	})
}

// ============================================================================
// BENCHMARK: Performance Comparison
// ============================================================================

func BenchmarkGetConfig(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetConfig()
		}
	})
}

func BenchmarkGetLogger(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetLogger()
		}
	})
}

func BenchmarkGetCache(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetCache()
		}
	})
}

// ============================================================================
// HELPER: Test Error Type
// ============================================================================

var ErrEmptyURL = errors.New("connection URL cannot be empty")
