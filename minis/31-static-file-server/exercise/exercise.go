//go:build !solution
// +build !solution

package exercise

import (
	"net/http"
	"os"
	"time"
)

// FileServerConfig holds configuration for the file server.
type FileServerConfig struct {
	Root               string // Root directory to serve files from
	EnableETag         bool   // Enable ETag generation and validation
	EnableRange        bool   // Enable HTTP Range request support
	EnableDirListing   bool   // Enable directory listing
	DefaultCacheMaxAge int    // Default Cache-Control max-age in seconds
}

// FileServer serves static files with ETags, Range requests, and caching.
type FileServer struct {
	config FileServerConfig
}

// NewFileServer creates a new file server with the given configuration.
// It should validate that the root directory exists.
func NewFileServer(config FileServerConfig) (*FileServer, error) {
	// TODO: implement
	// - Validate root directory exists
	// - Create and return FileServer
	return nil, nil
}

// ServeHTTP implements http.Handler interface.
// It should:
// 1. Only handle GET and HEAD methods
// 2. Validate and secure the requested path
// 3. Determine if request is for a file or directory
// 4. Serve accordingly with proper headers
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// securePath validates the requested path and returns the absolute filesystem path.
// It must prevent path traversal attacks (e.g., /../../../etc/passwd).
// Returns error if path is invalid or outside root directory.
func securePath(root, requestPath string) (string, error) {
	// TODO: implement
	// - Clean the request path
	// - Join with root
	// - Ensure result is within root directory
	return "", nil
}

// generateETag generates an ETag for the given file.
// You can use different strategies:
// - Simple: fmt.Sprintf("\"%x-%x\"", modTime.Unix(), size)
// - Hash-based: MD5/SHA256 of content
func generateETag(stat os.FileInfo) string {
	// TODO: implement
	return ""
}

// checkETag checks if the ETag matches the If-None-Match header.
// Returns true if matched (should send 304 Not Modified).
func checkETag(r *http.Request, etag string) bool {
	// TODO: implement
	return false
}

// checkModifiedSince checks if file was modified since If-Modified-Since header.
// Returns true if modified (should send file).
// Returns false if not modified (should send 304).
func checkModifiedSince(r *http.Request, modTime time.Time) bool {
	// TODO: implement
	return true
}

// setCacheHeaders sets appropriate Cache-Control headers based on file type.
// Different file types should have different cache policies:
// - HTML: no-cache or short cache with revalidation
// - CSS/JS: long cache (assumes versioned URLs)
// - Images: medium cache
// - Others: default cache
func setCacheHeaders(w http.ResponseWriter, filename string, maxAge int) {
	// TODO: implement
}

// detectContentType detects the MIME type of a file.
// Try extension-based detection first, fallback to content sniffing.
func detectContentType(path string) string {
	// TODO: implement
	return "application/octet-stream"
}

// serveFile serves a single file with all features:
// - ETag generation and validation
// - Last-Modified and If-Modified-Since handling
// - Content-Type detection
// - Cache headers
// - Range request support (via http.ServeContent)
func (fs *FileServer) serveFile(w http.ResponseWriter, r *http.Request, path string) {
	// TODO: implement
	// 1. Open file
	// 2. Get file info (stat)
	// 3. Generate ETag if enabled
	// 4. Check If-None-Match (ETag validation)
	// 5. Check If-Modified-Since
	// 6. Set headers (Content-Type, Cache-Control, ETag, Last-Modified)
	// 7. Set Accept-Ranges if range support enabled
	// 8. Use http.ServeContent for actual serving (handles ranges automatically)
}

// serveDirectory serves a directory listing if enabled.
// Otherwise returns 403 Forbidden.
func (fs *FileServer) serveDirectory(w http.ResponseWriter, r *http.Request, path string) {
	// TODO: implement
	// If directory listing disabled: return 403
	// Otherwise: read directory, generate HTML listing
}
