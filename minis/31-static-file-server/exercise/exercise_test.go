package exercise

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// setupTestFiles creates temporary test files for testing.
func setupTestFiles(t *testing.T) string {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.html": "<html><body>Home</body></html>",
		"style.css":  "body { color: red; }",
		"app.js":     "console.log('test');",
		"image.png":  "fake-png-data",
		"doc.pdf":    "fake-pdf-data",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	// Create subdirectory
	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(subdir, "test.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("Failed to create subdir file: %v", err)
	}

	return tmpDir
}

func TestNewFileServer(t *testing.T) {
	tmpDir := setupTestFiles(t)

	tests := []struct {
		name      string
		config    FileServerConfig
		expectErr bool
	}{
		{
			name: "valid config",
			config: FileServerConfig{
				Root:       tmpDir,
				EnableETag: true,
			},
			expectErr: false,
		},
		{
			name: "non-existent root",
			config: FileServerConfig{
				Root: "/nonexistent/path/12345",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, err := NewFileServer(tt.config)
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if fs == nil {
					t.Error("Expected non-nil FileServer")
				}
			}
		})
	}
}

func TestServeHTTP_BasicFileServing(t *testing.T) {
	tmpDir := setupTestFiles(t)

	config := FileServerConfig{
		Root:       tmpDir,
		EnableETag: true,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		checkContentType string
	}{
		{
			name:           "serve HTML file",
			path:           "/index.html",
			expectedStatus: http.StatusOK,
			expectedBody:   "<html><body>Home</body></html>",
			checkContentType: "text/html",
		},
		{
			name:           "serve CSS file",
			path:           "/style.css",
			expectedStatus: http.StatusOK,
			expectedBody:   "body { color: red; }",
			checkContentType: "text/css",
		},
		{
			name:           "serve JS file",
			path:           "/app.js",
			expectedStatus: http.StatusOK,
			expectedBody:   "console.log('test');",
			checkContentType: "javascript",
		},
		{
			name:           "not found",
			path:           "/nonexistent.txt",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "subdirectory file",
			path:           "/subdir/test.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			fs.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" {
				body := w.Body.String()
				if body != tt.expectedBody {
					t.Errorf("Expected body %q, got %q", tt.expectedBody, body)
				}
			}

			if tt.checkContentType != "" {
				contentType := w.Header().Get("Content-Type")
				if !strings.Contains(contentType, tt.checkContentType) {
					t.Errorf("Expected Content-Type to contain %q, got %q", tt.checkContentType, contentType)
				}
			}
		})
	}
}

func TestServeHTTP_MethodNotAllowed(t *testing.T) {
	tmpDir := setupTestFiles(t)

	config := FileServerConfig{
		Root: tmpDir,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/index.html", nil)
			w := httptest.NewRecorder()

			fs.ServeHTTP(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for %s, got %d", http.StatusMethodNotAllowed, method, w.Code)
			}
		})
	}
}

func TestServeHTTP_ETag(t *testing.T) {
	tmpDir := setupTestFiles(t)

	config := FileServerConfig{
		Root:       tmpDir,
		EnableETag: true,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	// First request - get ETag
	req1 := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	w1 := httptest.NewRecorder()
	fs.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w1.Code)
	}

	etag := w1.Header().Get("ETag")
	if etag == "" {
		t.Fatal("Expected ETag header")
	}

	// Second request with If-None-Match
	req2 := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	req2.Header.Set("If-None-Match", etag)
	w2 := httptest.NewRecorder()
	fs.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotModified {
		t.Errorf("Expected status 304, got %d", w2.Code)
	}

	if w2.Body.Len() > 0 {
		t.Error("Expected empty body for 304 response")
	}
}

func TestServeHTTP_LastModified(t *testing.T) {
	tmpDir := setupTestFiles(t)

	config := FileServerConfig{
		Root: tmpDir,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	// First request - get Last-Modified
	req1 := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	w1 := httptest.NewRecorder()
	fs.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w1.Code)
	}

	lastModified := w1.Header().Get("Last-Modified")
	if lastModified == "" {
		t.Fatal("Expected Last-Modified header")
	}

	// Second request with If-Modified-Since (same time)
	req2 := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	req2.Header.Set("If-Modified-Since", lastModified)
	w2 := httptest.NewRecorder()
	fs.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotModified {
		t.Errorf("Expected status 304, got %d", w2.Code)
	}

	// Third request with If-Modified-Since (future time)
	futureTime := time.Now().Add(24 * time.Hour).UTC().Format(http.TimeFormat)
	req3 := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	req3.Header.Set("If-Modified-Since", futureTime)
	w3 := httptest.NewRecorder()
	fs.ServeHTTP(w3, req3)

	if w3.Code != http.StatusNotModified {
		t.Errorf("Expected status 304 for future time, got %d", w3.Code)
	}
}

func TestServeHTTP_CacheControl(t *testing.T) {
	tmpDir := setupTestFiles(t)

	config := FileServerConfig{
		Root:               tmpDir,
		DefaultCacheMaxAge: 3600,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		contains string
	}{
		{
			name:     "HTML no-cache",
			path:     "/index.html",
			contains: "no-cache",
		},
		{
			name:     "CSS long cache",
			path:     "/style.css",
			contains: "max-age=31536000",
		},
		{
			name:     "JS long cache",
			path:     "/app.js",
			contains: "max-age=31536000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			fs.ServeHTTP(w, req)

			cacheControl := w.Header().Get("Cache-Control")
			if !strings.Contains(cacheControl, tt.contains) {
				t.Errorf("Expected Cache-Control to contain %q, got %q", tt.contains, cacheControl)
			}
		})
	}
}

func TestServeHTTP_RangeRequest(t *testing.T) {
	tmpDir := setupTestFiles(t)

	// Create a larger file for range testing
	largeContent := strings.Repeat("0123456789", 100) // 1000 bytes
	largePath := filepath.Join(tmpDir, "large.txt")
	if err := os.WriteFile(largePath, []byte(largeContent), 0644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	config := FileServerConfig{
		Root:        tmpDir,
		EnableRange: true,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	// Check Accept-Ranges header
	req1 := httptest.NewRequest(http.MethodGet, "/large.txt", nil)
	w1 := httptest.NewRecorder()
	fs.ServeHTTP(w1, req1)

	if w1.Header().Get("Accept-Ranges") != "bytes" {
		t.Error("Expected Accept-Ranges: bytes header")
	}

	// Test range request
	req2 := httptest.NewRequest(http.MethodGet, "/large.txt", nil)
	req2.Header.Set("Range", "bytes=0-99")
	w2 := httptest.NewRecorder()
	fs.ServeHTTP(w2, req2)

	if w2.Code != http.StatusPartialContent {
		t.Errorf("Expected status 206, got %d", w2.Code)
	}

	body := w2.Body.String()
	if len(body) != 100 {
		t.Errorf("Expected 100 bytes, got %d", len(body))
	}

	expectedBody := largeContent[:100]
	if body != expectedBody {
		t.Errorf("Range content mismatch")
	}

	contentRange := w2.Header().Get("Content-Range")
	if !strings.Contains(contentRange, "bytes 0-99/1000") {
		t.Errorf("Expected Content-Range bytes 0-99/1000, got %q", contentRange)
	}
}

func TestServeHTTP_DirectoryListing(t *testing.T) {
	tmpDir := setupTestFiles(t)

	config := FileServerConfig{
		Root:             tmpDir,
		EnableDirListing: true,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/subdir/", nil)
	w := httptest.NewRecorder()
	fs.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "test.txt") {
		t.Error("Expected directory listing to contain test.txt")
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected Content-Type text/html, got %q", contentType)
	}
}

func TestServeHTTP_DirectoryListingDisabled(t *testing.T) {
	tmpDir := setupTestFiles(t)

	config := FileServerConfig{
		Root:             tmpDir,
		EnableDirListing: false,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/subdir/", nil)
	w := httptest.NewRecorder()
	fs.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestSecurePath_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		requestPath string
		expectError bool
	}{
		{
			name:        "normal path",
			requestPath: "/index.html",
			expectError: false,
		},
		{
			name:        "subdirectory",
			requestPath: "/subdir/file.txt",
			expectError: false,
		},
		{
			name:        "dot dot in path gets cleaned",
			requestPath: "/subdir/../index.html",
			expectError: false, // Gets cleaned to /index.html which is safe
		},
		{
			name:        "multiple dots get cleaned",
			requestPath: "/././././index.html",
			expectError: false, // Gets cleaned to /index.html
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := securePath(tmpDir, tt.requestPath)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error for path traversal attempt")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Verify path is within tmpDir
				if !strings.HasPrefix(path, tmpDir) {
					t.Errorf("Path %q is not within tmpDir %q", path, tmpDir)
				}
			}
		})
	}
}

// TestSecurePath_ActualSecurity tests that even tricky paths stay within root
func TestSecurePath_ActualSecurity(t *testing.T) {
	tmpDir := t.TempDir()

	trickyPaths := []string{
		"/../../../etc/passwd",
		"/subdir/../../etc/passwd",
		"/./../../etc/passwd",
		"/..",
		"/../",
		"/../../",
	}

	for _, requestPath := range trickyPaths {
		t.Run(requestPath, func(t *testing.T) {
			path, err := securePath(tmpDir, requestPath)
			// These paths should either error OR resolve to something within tmpDir
			if err == nil {
				// If no error, path must be within tmpDir
				if !strings.HasPrefix(path, tmpDir) {
					t.Errorf("SECURITY: Path %q escaped tmpDir %q for request %q",
						path, tmpDir, requestPath)
				}
			}
			// If there's an error, that's also fine - we rejected it
		})
	}
}

func TestGenerateETag(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	stat, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	etag1 := generateETag(stat)
	if etag1 == "" {
		t.Error("Expected non-empty ETag")
	}

	// Check format (should be quoted)
	if !strings.HasPrefix(etag1, "\"") || !strings.HasSuffix(etag1, "\"") {
		t.Errorf("ETag should be quoted, got %q", etag1)
	}

	// Same file should generate same ETag
	etag2 := generateETag(stat)
	if etag1 != etag2 {
		t.Error("Same file should generate same ETag")
	}
}

func TestCheckETag(t *testing.T) {
	tests := []struct {
		name        string
		etag        string
		ifNoneMatch string
		expected    bool
	}{
		{
			name:        "matching ETag",
			etag:        "\"abc123\"",
			ifNoneMatch: "\"abc123\"",
			expected:    true,
		},
		{
			name:        "non-matching ETag",
			etag:        "\"abc123\"",
			ifNoneMatch: "\"xyz789\"",
			expected:    false,
		},
		{
			name:        "no If-None-Match",
			etag:        "\"abc123\"",
			ifNoneMatch: "",
			expected:    false,
		},
		{
			name:        "wildcard",
			etag:        "\"abc123\"",
			ifNoneMatch: "*",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.ifNoneMatch != "" {
				req.Header.Set("If-None-Match", tt.ifNoneMatch)
			}

			result := checkETag(req, tt.etag)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCheckModifiedSince(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	past := now.Add(-1 * time.Hour)

	tests := []struct {
		name        string
		modTime     time.Time
		ims         string
		expected    bool
		description string
	}{
		{
			name:        "modified after IMS",
			modTime:     now,
			ims:         past.Format(http.TimeFormat),
			expected:    true,
			description: "file modified, should serve",
		},
		{
			name:        "not modified since IMS",
			modTime:     past,
			ims:         now.Format(http.TimeFormat),
			expected:    false,
			description: "file not modified, should 304",
		},
		{
			name:        "same time",
			modTime:     now,
			ims:         now.Format(http.TimeFormat),
			expected:    false,
			description: "same time, should 304",
		},
		{
			name:        "no IMS header",
			modTime:     now,
			ims:         "",
			expected:    true,
			description: "no IMS header, should serve",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.ims != "" {
				req.Header.Set("If-Modified-Since", tt.ims)
			}

			result := checkModifiedSince(req, tt.modTime)
			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.description, tt.expected, result)
			}
		})
	}
}

func TestDetectContentType(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		filename    string
		content     string
		expectedType string
	}{
		{
			name:        "HTML file",
			filename:    "index.html",
			content:     "<html></html>",
			expectedType: "text/html",
		},
		{
			name:        "CSS file",
			filename:    "style.css",
			content:     "body { }",
			expectedType: "text/css",
		},
		{
			name:        "JavaScript file",
			filename:    "app.js",
			content:     "console.log('test');",
			expectedType: "javascript",
		},
		{
			name:        "JSON file",
			filename:    "data.json",
			content:     "{}",
			expectedType: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			contentType := detectContentType(path)
			if !strings.Contains(strings.ToLower(contentType), tt.expectedType) {
				t.Errorf("Expected content type to contain %q, got %q", tt.expectedType, contentType)
			}
		})
	}
}

func TestHEADRequest(t *testing.T) {
	tmpDir := setupTestFiles(t)

	config := FileServerConfig{
		Root:       tmpDir,
		EnableETag: true,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatalf("Failed to create file server: %v", err)
	}

	// HEAD request
	req := httptest.NewRequest(http.MethodHead, "/index.html", nil)
	w := httptest.NewRecorder()
	fs.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// HEAD should have headers but no body
	if w.Body.Len() > 0 {
		t.Error("HEAD request should not have body")
	}

	// Should have Content-Type and other headers
	if w.Header().Get("Content-Type") == "" {
		t.Error("Expected Content-Type header")
	}

	if w.Header().Get("ETag") == "" {
		t.Error("Expected ETag header")
	}
}

// Benchmark file serving
func BenchmarkServeFile(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := strings.Repeat("Hello, World!", 1000)

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	config := FileServerConfig{
		Root:       tmpDir,
		EnableETag: true,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		b.Fatalf("Failed to create file server: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test.txt", nil)
		w := httptest.NewRecorder()
		fs.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Errorf("Expected status 200, got %d", w.Code)
		}
	}
}

// Benchmark ETag validation (304 responses)
func BenchmarkETagValidation(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	config := FileServerConfig{
		Root:       tmpDir,
		EnableETag: true,
	}

	fs, err := NewFileServer(config)
	if err != nil {
		b.Fatalf("Failed to create file server: %v", err)
	}

	// Get ETag
	req1 := httptest.NewRequest(http.MethodGet, "/test.txt", nil)
	w1 := httptest.NewRecorder()
	fs.ServeHTTP(w1, req1)
	etag := w1.Header().Get("ETag")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test.txt", nil)
		req.Header.Set("If-None-Match", etag)
		w := httptest.NewRecorder()
		fs.ServeHTTP(w, req)

		if w.Code != http.StatusNotModified {
			b.Errorf("Expected status 304, got %d", w.Code)
		}
	}
}
