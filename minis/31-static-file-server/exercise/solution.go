//go:build solution
// +build solution

/*
Problem: Build a static file server with ETags, Range requests, and caching

Requirements:
1. Serve files from a root directory
2. ETag generation and validation
3. Range request support for resumable downloads
4. Proper caching headers (Cache-Control, Last-Modified)
5. MIME type detection
6. Path traversal protection
7. Directory listing (optional)

Why Go is well-suited:
- http.ServeContent: Built-in range request and conditional request handling
- mime package: MIME type detection
- filepath package: Secure path handling
- Zero-copy sendfile() for efficient file serving

Compared to other languages:
- Node.js: express.static similar, but requires library
- Python: SimpleHTTPServer/http.server built-in but less efficient
- Rust: Very fast with actix-files, but more complex
*/

package exercise

import (
	"crypto/md5"
	"fmt"
	"html"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
func NewFileServer(config FileServerConfig) (*FileServer, error) {
	// Validate root directory exists
	info, err := os.Stat(config.Root)
	if err != nil {
		return nil, fmt.Errorf("root directory error: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("root path is not a directory: %s", config.Root)
	}

	// Get absolute path
	absRoot, err := filepath.Abs(config.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	config.Root = absRoot

	return &FileServer{
		config: config,
	}, nil
}

// ServeHTTP implements http.Handler interface.
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only handle GET and HEAD methods
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate and secure the requested path
	path, err := securePath(fs.config.Root, r.URL.Path)
	if err != nil {
		log.Printf("Path security error: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			log.Printf("Stat error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Serve directory or file
	if info.IsDir() {
		// Try index.html first
		indexPath := filepath.Join(path, "index.html")
		if indexInfo, err := os.Stat(indexPath); err == nil && !indexInfo.IsDir() {
			fs.serveFile(w, r, indexPath)
			return
		}

		// Serve directory listing
		fs.serveDirectory(w, r, path)
	} else {
		fs.serveFile(w, r, path)
	}
}

// securePath validates the requested path and returns the absolute filesystem path.
func securePath(root, requestPath string) (string, error) {
	// Get absolute root first
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	// Clean the request path (resolves .., ., etc.)
	// Prepend "/" to ensure it's treated as absolute from root
	requestPath = filepath.Clean("/" + requestPath)

	// Remove leading slash for joining
	requestPath = strings.TrimPrefix(requestPath, "/")

	// Join with root
	fullPath := filepath.Join(absRoot, requestPath)

	// Get absolute path and resolve any symlinks
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	// Ensure result is within root directory
	// Use filepath.Clean to normalize both paths
	absRoot = filepath.Clean(absRoot) + string(filepath.Separator)
	absPath = filepath.Clean(absPath) + string(filepath.Separator)

	// Check if path is within root
	if !strings.HasPrefix(absPath, absRoot) {
		return "", fmt.Errorf("path traversal attempt: %s", requestPath)
	}

	// Return the path without trailing separator
	return strings.TrimSuffix(absPath, string(filepath.Separator)), nil
}

// generateETag generates an ETag for the given file.
func generateETag(stat os.FileInfo) string {
	// Simple ETag based on modification time and size
	// For production, consider using content hash for better accuracy
	return fmt.Sprintf("\"%x-%x\"", stat.ModTime().Unix(), stat.Size())
}

// generateETagWithHash generates a hash-based ETag (more accurate but slower).
func generateETagWithHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(data)
	return fmt.Sprintf("\"%x\"", hash), nil
}

// checkETag checks if the ETag matches the If-None-Match header.
func checkETag(r *http.Request, etag string) bool {
	ifNoneMatch := r.Header.Get("If-None-Match")
	if ifNoneMatch == "" {
		return false
	}

	// Check for exact match or "*"
	if ifNoneMatch == etag || ifNoneMatch == "*" {
		return true
	}

	// Check if ETag is in comma-separated list
	for _, tag := range strings.Split(ifNoneMatch, ",") {
		tag = strings.TrimSpace(tag)
		if tag == etag {
			return true
		}
	}

	return false
}

// checkModifiedSince checks if file was modified since If-Modified-Since header.
func checkModifiedSince(r *http.Request, modTime time.Time) bool {
	ims := r.Header.Get("If-Modified-Since")
	if ims == "" {
		return true // No If-Modified-Since header, file is "modified"
	}

	// Parse the time
	t, err := http.ParseTime(ims)
	if err != nil {
		return true // Invalid time, assume modified
	}

	// Truncate to seconds (HTTP date precision)
	modTime = modTime.Truncate(time.Second)

	// If file is not newer, it hasn't been modified
	return modTime.After(t)
}

// setCacheHeaders sets appropriate Cache-Control headers based on file type.
func setCacheHeaders(w http.ResponseWriter, filename string, maxAge int) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".html", ".htm":
		// HTML: revalidate every time (content changes frequently)
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")

	case ".css", ".js":
		// CSS/JS: cache for 1 year (assumes versioned URLs like app.abc123.css)
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".ico":
		// Images: cache for 30 days
		w.Header().Set("Cache-Control", "public, max-age=2592000")

	case ".woff", ".woff2", ".ttf", ".eot":
		// Fonts: cache for 1 year
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	case ".pdf", ".zip", ".tar", ".gz":
		// Documents/Archives: cache for 1 day
		w.Header().Set("Cache-Control", "public, max-age=86400")

	default:
		// Default: use configured max-age
		if maxAge > 0 {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		} else {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
	}
}

// detectContentType detects the MIME type of a file.
func detectContentType(path string) string {
	// Try extension-based detection first
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	if contentType != "" {
		return contentType
	}

	// Fallback to content sniffing
	file, err := os.Open(path)
	if err != nil {
		return "application/octet-stream"
	}
	defer file.Close()

	// Read first 512 bytes for detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return "application/octet-stream"
	}

	// Detect content type
	contentType = http.DetectContentType(buffer[:n])
	return contentType
}

// serveFile serves a single file with all features.
func (fs *FileServer) serveFile(w http.ResponseWriter, r *http.Request, path string) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get file info
	stat, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Detect and set Content-Type
	contentType := detectContentType(path)
	w.Header().Set("Content-Type", contentType)

	// Set cache headers
	setCacheHeaders(w, filepath.Base(path), fs.config.DefaultCacheMaxAge)

	// Set Last-Modified header
	w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))

	// Generate and check ETag if enabled
	if fs.config.EnableETag {
		etag := generateETag(stat)
		w.Header().Set("ETag", etag)

		// Check If-None-Match (ETag validation)
		if checkETag(r, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// Check If-Modified-Since
	if !checkModifiedSince(r, stat.ModTime()) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Set Accept-Ranges header if range support enabled
	if fs.config.EnableRange {
		w.Header().Set("Accept-Ranges", "bytes")
	}

	// Serve content (http.ServeContent handles Range requests automatically)
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}

// serveDirectory serves a directory listing if enabled.
func (fs *FileServer) serveDirectory(w http.ResponseWriter, r *http.Request, path string) {
	if !fs.config.EnableDirListing {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Open directory
	dir, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer dir.Close()

	// Read directory entries
	entries, err := dir.Readdir(-1)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Sort entries by name
	sort.Slice(entries, func(i, j int) bool {
		// Directories first, then alphabetical
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	// Set headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")

	// Generate HTML
	fmt.Fprintf(w, "<!DOCTYPE html>\n")
	fmt.Fprintf(w, "<html>\n<head>\n")
	fmt.Fprintf(w, "<title>Directory: %s</title>\n", html.EscapeString(r.URL.Path))
	fmt.Fprintf(w, "<style>\n")
	fmt.Fprintf(w, "body { font-family: Arial, sans-serif; max-width: 1000px; margin: 0 auto; padding: 20px; }\n")
	fmt.Fprintf(w, "h1 { border-bottom: 2px solid #007bff; padding-bottom: 10px; }\n")
	fmt.Fprintf(w, "table { width: 100%%; border-collapse: collapse; }\n")
	fmt.Fprintf(w, "th, td { text-align: left; padding: 8px; border-bottom: 1px solid #ddd; }\n")
	fmt.Fprintf(w, "th { background-color: #f5f5f5; }\n")
	fmt.Fprintf(w, "tr:hover { background-color: #f9f9f9; }\n")
	fmt.Fprintf(w, "a { color: #007bff; text-decoration: none; }\n")
	fmt.Fprintf(w, "a:hover { text-decoration: underline; }\n")
	fmt.Fprintf(w, ".dir { font-weight: bold; }\n")
	fmt.Fprintf(w, ".size { text-align: right; }\n")
	fmt.Fprintf(w, "</style>\n")
	fmt.Fprintf(w, "</head>\n<body>\n")
	fmt.Fprintf(w, "<h1>Directory: %s</h1>\n", html.EscapeString(r.URL.Path))

	// Table header
	fmt.Fprintf(w, "<table>\n")
	fmt.Fprintf(w, "<tr><th>Name</th><th class=\"size\">Size</th><th>Modified</th></tr>\n")

	// Parent directory link
	if r.URL.Path != "/" {
		fmt.Fprintf(w, "<tr><td class=\"dir\"><a href=\"..\">..</a></td><td></td><td></td></tr>\n")
	}

	// List entries
	for _, entry := range entries {
		name := entry.Name()
		displayName := name
		href := url.PathEscape(name)

		if entry.IsDir() {
			displayName += "/"
			href += "/"
		}

		// Format size
		size := ""
		if !entry.IsDir() {
			size = formatSize(entry.Size())
		}

		// Format modification time
		modTime := entry.ModTime().Format("2006-01-02 15:04:05")

		// CSS class for directories
		class := ""
		if entry.IsDir() {
			class = " class=\"dir\""
		}

		fmt.Fprintf(w, "<tr><td%s><a href=\"%s\">%s</a></td><td class=\"size\">%s</td><td>%s</td></tr>\n",
			class, href, html.EscapeString(displayName), size, modTime)
	}

	fmt.Fprintf(w, "</table>\n")
	fmt.Fprintf(w, "</body>\n</html>\n")
}

// formatSize formats file size in human-readable format.
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
