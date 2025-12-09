//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "fmt" for error formatting and string building
// - "html" for escaping HTML content (XSS prevention)
// - "mime" for MIME type detection
// - "net/http" for HTTP server types and utilities
// - "net/url" for URL path escaping
// - "os" for file system operations
// - "path/filepath" for secure path manipulation
// - "sort" for sorting directory listings
// - "strings" for string manipulation
// - "time" for timestamp handling
//
// import (
//     "fmt"
//     "html"
//     "mime"
//     "net/http"
//     "net/url"
//     "os"
//     "path/filepath"
//     "sort"
//     "strings"
//     "time"
// )

// FileServerConfig holds configuration for the file server.
// This is a STRUCT TYPE (value type, not reference type)
//
// Key Go concepts for this struct:
// - Exported fields (capitalized) can be accessed from other packages
// - bool fields default to false, int fields default to 0
// - When passed to functions, the entire struct is copied (unless passed as pointer)
type FileServerConfig struct {
	Root               string // Root directory to serve files from
	EnableETag         bool   // Enable ETag generation and validation
	EnableRange        bool   // Enable HTTP Range request support
	EnableDirListing   bool   // Enable directory listing
	DefaultCacheMaxAge int    // Default Cache-Control max-age in seconds
}

// FileServer serves static files with ETags, Range requests, and caching.
// This struct holds state (configuration) and implements http.Handler interface.
//
// TODO: Define the FileServer struct
// type FileServer struct {
//     config FileServerConfig  // Embedded config (stored by value, copied during initialization)
// }
//
// Key Go concepts:
// - To implement http.Handler, we need a ServeHTTP(w http.ResponseWriter, r *http.Request) method
// - Methods with pointer receiver (*FileServer) can modify the struct
// - Methods with value receiver (FileServer) get a copy and cannot modify original

// NewFileServer creates a new file server with the given configuration.
//
// TODO: Implement NewFileServer function
// Function signature: func NewFileServer(config FileServerConfig) (*FileServer, error)
//
// Steps to implement:
//
// 1. Validate that root directory exists
//    - Use: info, err := os.Stat(config.Root)
//    - Check if err != nil (directory doesn't exist or permission error)
//    - Check if !info.IsDir() (path exists but is a file, not directory)
//    - Return descriptive errors: fmt.Errorf("root directory error: %w", err)
//
// 2. Get absolute path of root directory
//    - Use: absRoot, err := filepath.Abs(config.Root)
//    - Why? Relative paths like "." or "../data" are resolved relative to current directory
//    - Absolute paths prevent confusion when working directory changes
//    - Update config: config.Root = absRoot
//
// 3. Create and return FileServer
//    - Use: return &FileServer{config: config}, nil
//    - Why return *FileServer? Because http.Handler methods use pointer receiver
//    - The & creates a pointer to a FileServer struct allocated on the heap
//
// Key Go concepts:
// - Returning pointers (*FileServer) vs values (FileServer)
//   * Pointer: Caller gets reference to heap-allocated struct (can be modified)
//   * Value: Caller gets copy of struct (modifications don't affect original)
// - Error wrapping: %w preserves error chain for errors.Is() and errors.As()
// - Multiple return values: (result, error) is idiomatic Go pattern

// TODO: Implement the NewFileServer function below
// func NewFileServer(config FileServerConfig) (*FileServer, error) {
//     return nil, nil
// }

// ServeHTTP implements http.Handler interface.
//
// TODO: Implement ServeHTTP method
// Function signature: func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
//
// Steps to implement:
//
// 1. Only handle GET and HEAD methods
//    - Check: r.Method != http.MethodGet && r.Method != http.MethodHead
//    - Return 405: http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//    - Why? POST/PUT/DELETE don't make sense for static file serving
//
// 2. Validate and secure the requested path
//    - Use: path, err := securePath(fs.config.Root, r.URL.Path)
//    - This prevents path traversal attacks (e.g., /../../../etc/passwd)
//    - If err != nil: return 403 Forbidden
//
// 3. Check if path exists
//    - Use: info, err := os.Stat(path)
//    - If os.IsNotExist(err): return 404 Not Found
//    - If other error: return 500 Internal Server Error
//
// 4. Determine if request is for a file or directory
//    - Use: info.IsDir()
//    - If directory: Try serving index.html first, then directory listing
//    - If file: Serve the file
//
// 5. Serve index.html for directories (if exists)
//    - Use: indexPath := filepath.Join(path, "index.html")
//    - Check if it exists: os.Stat(indexPath)
//    - If exists and is file: fs.serveFile(w, r, indexPath)
//
// 6. Otherwise serve directory listing or file
//    - Directory: fs.serveDirectory(w, r, path)
//    - File: fs.serveFile(w, r, path)
//
// Key Go concepts:
// - Pointer receiver: (fs *FileServer) allows method to access struct fields
//   * fs.config.Root accesses config through the pointer
//   * Even though we don't modify fs, we use pointer receiver for consistency
//     (all http.Handler methods should use pointer receivers)
// - http.ResponseWriter is an INTERFACE (passed by value, but contains pointer internally)
//   * Writing to w sends data to the client
//   * Setting headers: w.Header().Set("Key", "Value")
// - *http.Request is a POINTER to a struct
//   * Contains: Method, URL, Headers, Body
//   * Never modify the request (read-only)

// TODO: Implement the ServeHTTP method below
// func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//     // TODO: implement
// }

// securePath validates the requested path and returns the absolute filesystem path.
// CRITICAL: This function prevents path traversal attacks!
//
// TODO: Implement securePath function
// Function signature: func securePath(root, requestPath string) (string, error)
//
// Steps to implement:
//
// 1. Get absolute root path
//    - Use: absRoot, err := filepath.Abs(root)
//    - Why? Ensure root is absolute for comparison
//
// 2. Clean the request path
//    - Use: requestPath = filepath.Clean("/" + requestPath)
//    - filepath.Clean() resolves ".." and "." components
//    - Example: "/foo/../bar" becomes "/bar"
//    - Prepending "/" ensures path is treated as absolute from root
//
// 3. Remove leading slash for joining
//    - Use: requestPath = strings.TrimPrefix(requestPath, "/")
//    - Why? filepath.Join() treats absolute paths specially
//    - We want to join relative to our root
//
// 4. Join with root directory
//    - Use: fullPath := filepath.Join(absRoot, requestPath)
//    - This creates: /var/www/public + images/cat.jpg = /var/www/public/images/cat.jpg
//
// 5. Get absolute path (resolves symlinks)
//    - Use: absPath, err := filepath.Abs(fullPath)
//    - This resolves symbolic links to their targets
//
// 6. Ensure result is within root directory
//    - Normalize both paths: absRoot = filepath.Clean(absRoot) + string(filepath.Separator)
//    - Add trailing separator to prevent prefix matching false positives
//    - Example: Without separator, "/var/www" would match "/var/www-backup"
//    - Check: !strings.HasPrefix(absPath, absRoot)
//    - If outside root: return error (path traversal attempt!)
//
// 7. Return the secured path
//    - Remove trailing separator: strings.TrimSuffix(...)
//    - Return: (absPath, nil)
//
// Key Go concepts:
// - filepath vs path packages:
//   * filepath: OS-specific (\ on Windows, / on Unix)
//   * path: Always uses / (for URLs)
// - String immutability: filepath.Clean() returns NEW string
// - Security: Always validate user input before file system operations!

// TODO: Implement the securePath function below
// func securePath(root, requestPath string) (string, error) {
//     return "", nil
// }

// generateETag generates an ETag for the given file.
//
// TODO: Implement generateETag function
// Function signature: func generateETag(stat os.FileInfo) string
//
// Steps to implement:
//
// 1. Create ETag from modification time and size
//    - Use: fmt.Sprintf("\"%x-%x\"", stat.ModTime().Unix(), stat.Size())
//    - %x formats as hexadecimal
//    - Quotes around ETag are required by HTTP spec
//    - Example: "507f1f77bcf86cd799439011-1024"
//
// Alternative (stronger but slower):
// - Read file content and compute MD5 hash
// - Trade-off: More accurate (detects changes without time/size change)
//   but requires reading entire file
//
// Key Go concepts:
// - os.FileInfo is an INTERFACE with methods:
//   * Name() string
//   * Size() int64
//   * ModTime() time.Time
//   * IsDir() bool
// - Interfaces passed by value contain: (type pointer, value pointer)
// - Format verbs: %x (hex), %d (decimal), %s (string)

// TODO: Implement the generateETag function below
// func generateETag(stat os.FileInfo) string {
//     return ""
// }

// checkETag checks if the ETag matches the If-None-Match header.
//
// TODO: Implement checkETag function
// Function signature: func checkETag(r *http.Request, etag string) bool
//
// Steps to implement:
//
// 1. Get If-None-Match header
//    - Use: ifNoneMatch := r.Header.Get("If-None-Match")
//    - If empty: return false (no ETag to check)
//
// 2. Check for exact match or wildcard
//    - If ifNoneMatch == etag: return true
//    - If ifNoneMatch == "*": return true
//
// 3. Handle multiple ETags (comma-separated)
//    - Use: strings.Split(ifNoneMatch, ",")
//    - Loop through each tag
//    - Use: strings.TrimSpace(tag) to remove whitespace
//    - If any matches: return true
//
// 4. No match found: return false
//
// Key Go concepts:
// - HTTP headers are case-insensitive (Get() normalizes key)
// - http.Header is map[string][]string (one key, multiple values)
// - Get() returns first value or "" if not present

// TODO: Implement the checkETag function below
// func checkETag(r *http.Request, etag string) bool {
//     return false
// }

// checkModifiedSince checks if file was modified since If-Modified-Since header.
//
// TODO: Implement checkModifiedSince function
// Function signature: func checkModifiedSince(r *http.Request, modTime time.Time) bool
//
// Steps to implement:
//
// 1. Get If-Modified-Since header
//    - Use: ims := r.Header.Get("If-Modified-Since")
//    - If empty: return true (assume modified)
//
// 2. Parse the time
//    - Use: t, err := http.ParseTime(ims)
//    - http.ParseTime() handles multiple HTTP date formats
//    - If error: return true (invalid time, assume modified)
//
// 3. Compare modification times
//    - Truncate to seconds: modTime = modTime.Truncate(time.Second)
//    - Why? HTTP dates only have second precision, not nanoseconds
//    - Return: modTime.After(t)
//    - If file is newer than request time: modified (return true)
//    - If file is same age or older: not modified (return false)
//
// Return value semantics:
// - true: File WAS modified (send 200 with content)
// - false: File was NOT modified (send 304 Not Modified)
//
// Key Go concepts:
// - time.Time is a STRUCT (value type)
//   * Contains: seconds, nanoseconds, location
//   * Passed by value (copied)
// - Truncate() returns NEW time.Time (doesn't modify original)

// TODO: Implement the checkModifiedSince function below
// func checkModifiedSince(r *http.Request, modTime time.Time) bool {
//     return true
// }

// setCacheHeaders sets appropriate Cache-Control headers based on file type.
//
// TODO: Implement setCacheHeaders function
// Function signature: func setCacheHeaders(w http.ResponseWriter, filename string, maxAge int)
//
// Steps to implement:
//
// 1. Get file extension
//    - Use: ext := strings.ToLower(filepath.Ext(filename))
//    - ToLower() for case-insensitive matching (.JPG == .jpg)
//
// 2. Set Cache-Control based on file type
//    - HTML (.html, .htm): "no-cache, must-revalidate"
//      * Why? HTML changes frequently, always revalidate
//    - CSS/JS (.css, .js): "public, max-age=31536000, immutable"
//      * Why? Assume versioned URLs (app.abc123.css), cache forever
//      * immutable = browser won't revalidate even on reload
//    - Images (.jpg, .png, .gif, .webp, .svg): "public, max-age=2592000"
//      * Why? Cache for 30 days (images change less frequently)
//    - Fonts (.woff, .woff2, .ttf): "public, max-age=31536000, immutable"
//      * Why? Fonts rarely change, cache forever
//    - Default: Use maxAge parameter or 3600 (1 hour)
//
// 3. Set the header
//    - Use: w.Header().Set("Cache-Control", value)
//
// Key Go concepts:
// - http.ResponseWriter methods:
//   * Header() returns http.Header (map[string][]string)
//   * Set() overwrites existing value (vs Add() which appends)
// - Cache directives:
//   * public: Can be cached by browser and CDN
//   * private: Only browser cache, not CDN
//   * no-cache: Must revalidate with server before use
//   * max-age: Seconds until cache expires

// TODO: Implement the setCacheHeaders function below
// func setCacheHeaders(w http.ResponseWriter, filename string, maxAge int) {
//     // TODO: implement
// }

// detectContentType detects the MIME type of a file.
//
// TODO: Implement detectContentType function
// Function signature: func detectContentType(path string) string
//
// Steps to implement:
//
// 1. Try extension-based detection first
//    - Use: ext := filepath.Ext(path)
//    - Use: contentType := mime.TypeByExtension(ext)
//    - If contentType != "": return it
//    - Why try extension first? It's faster (no file I/O)
//
// 2. Fallback to content sniffing
//    - Open file: file, err := os.Open(path)
//    - If error: return "application/octet-stream" (generic binary)
//    - Don't forget: defer file.Close()
//
// 3. Read first 512 bytes
//    - Use: buffer := make([]byte, 512)
//    - Use: n, err := file.Read(buffer)
//    - Why 512? http.DetectContentType() uses first 512 bytes
//    - If error: return "application/octet-stream"
//
// 4. Detect content type from bytes
//    - Use: http.DetectContentType(buffer[:n])
//    - This analyzes byte patterns (magic numbers)
//    - Example: [0xFF, 0xD8, 0xFF] = JPEG
//
// Key Go concepts:
// - defer: Schedules function call to run when surrounding function returns
//   * Ensures file.Close() is called even if we return early
//   * Deferred calls run in LIFO order
// - Slice creation: make([]byte, 512) allocates 512-byte slice
// - Slice slicing: buffer[:n] creates view of first n bytes
//   * No copy! Just a new slice header pointing to same array

// TODO: Implement the detectContentType function below
// func detectContentType(path string) string {
//     return "application/octet-stream"
// }

// serveFile serves a single file with all features.
//
// TODO: Implement serveFile method
// Function signature: func (fs *FileServer) serveFile(w http.ResponseWriter, r *http.Request, path string)
//
// Steps to implement:
//
// 1. Open the file
//    - Use: file, err := os.Open(path)
//    - Check error: log and return 500
//    - Use: defer file.Close()
//
// 2. Get file information
//    - Use: stat, err := file.Stat()
//    - Check error: log and return 500
//    - Why? Need modTime and size for ETag and Last-Modified
//
// 3. Detect and set Content-Type
//    - Use: contentType := detectContentType(path)
//    - Set header: w.Header().Set("Content-Type", contentType)
//
// 4. Set cache headers
//    - Use: setCacheHeaders(w, filepath.Base(path), fs.config.DefaultCacheMaxAge)
//    - filepath.Base() gets filename from path
//
// 5. Set Last-Modified header
//    - Use: w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))
//    - UTC() converts to UTC timezone
//    - http.TimeFormat is RFC1123 format required by HTTP
//
// 6. Generate and check ETag (if enabled)
//    - If fs.config.EnableETag:
//      * Generate: etag := generateETag(stat)
//      * Set header: w.Header().Set("ETag", etag)
//      * Check: if checkETag(r, etag) { w.WriteHeader(http.StatusNotModified); return }
//
// 7. Check If-Modified-Since
//    - If !checkModifiedSince(r, stat.ModTime()):
//      * File hasn't been modified
//      * w.WriteHeader(http.StatusNotModified)
//      * return (don't send body)
//
// 8. Set Accept-Ranges header (if range support enabled)
//    - If fs.config.EnableRange:
//      * w.Header().Set("Accept-Ranges", "bytes")
//      * Tells client that Range requests are supported
//
// 9. Serve the content
//    - Use: http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
//    - Why http.ServeContent?
//      * Handles Range requests automatically
//      * Sets Content-Length
//      * Uses efficient sendfile() syscall on Unix
//      * Handles HEAD requests (doesn't send body)
//
// Key Go concepts:
// - Method receiver: (fs *FileServer) accesses struct fields via pointer
// - os.File implements io.ReadSeeker (required by http.ServeContent)
// - http.ServeContent is zero-copy on Unix (sendfile syscall)
// - Headers must be set BEFORE WriteHeader() or Write()
// - 304 Not Modified: Headers sent, no body

// TODO: Implement the serveFile method below
// func (fs *FileServer) serveFile(w http.ResponseWriter, r *http.Request, path string) {
//     // TODO: implement
// }

// serveDirectory serves a directory listing if enabled.
//
// TODO: Implement serveDirectory method
// Function signature: func (fs *FileServer) serveDirectory(w http.ResponseWriter, r *http.Request, path string)
//
// Steps to implement:
//
// 1. Check if directory listing is enabled
//    - If !fs.config.EnableDirListing:
//      * http.Error(w, "Forbidden", http.StatusForbidden)
//      * return
//
// 2. Open the directory
//    - Use: dir, err := os.Open(path)
//    - Check error: log and return 500
//    - Use: defer dir.Close()
//
// 3. Read directory entries
//    - Use: entries, err := dir.Readdir(-1)
//    - -1 means read all entries
//    - Returns []os.FileInfo
//    - Check error: log and return 500
//
// 4. Sort entries
//    - Use: sort.Slice(entries, func(i, j int) bool { ... })
//    - Sort directories first, then alphabetically
//    - Compare: entries[i].IsDir() vs entries[j].IsDir()
//    - Then: entries[i].Name() < entries[j].Name()
//
// 5. Set response headers
//    - Content-Type: "text/html; charset=utf-8"
//    - Cache-Control: "no-cache" (listings change frequently)
//
// 6. Generate HTML listing
//    - Use fmt.Fprintf(w, ...) to write HTML
//    - Include:
//      * DOCTYPE, html, head, title
//      * CSS for styling (table, links, etc.)
//      * Table with columns: Name, Size, Modified
//      * Parent directory link (..) if not root
//      * List all entries with proper escaping
//
// 7. For each entry:
//    - Escape name: html.EscapeString(name) (XSS prevention!)
//    - Escape URL: url.PathEscape(name)
//    - Add "/" suffix for directories
//    - Format size: Use helper function formatSize()
//    - Format time: modTime.Format("2006-01-02 15:04:05")
//
// Key Go concepts:
// - dir.Readdir() returns []os.FileInfo (slice of interface values)
// - sort.Slice() uses custom comparator function (closure)
// - html.EscapeString() prevents XSS: "<script>" → "&lt;script&gt;"
// - url.PathEscape() encodes special chars: "file name.txt" → "file%20name.txt"
// - fmt.Fprintf(w, ...) writes formatted string to ResponseWriter

// TODO: Implement the serveDirectory method below
// func (fs *FileServer) serveDirectory(w http.ResponseWriter, r *http.Request, path string) {
//     // TODO: implement
// }

// formatSize formats file size in human-readable format.
//
// TODO: Implement formatSize helper function
// Function signature: func formatSize(size int64) string
//
// Steps to implement:
//
// 1. Define size constants
//    - KB = 1024
//    - MB = 1024 * KB
//    - GB = 1024 * MB
//
// 2. Use switch/case or if/else to format
//    - If size >= GB: fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
//    - If size >= MB: fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
//    - If size >= KB: fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
//    - Else: fmt.Sprintf("%d B", size)
//
// Key Go concepts:
// - Type conversion: float64(size) converts int64 to float64
// - Format verbs: %.2f (float with 2 decimal places)
// - Constants: const KB = 1024 (untyped constant, can be used with any numeric type)

// TODO: Implement the formatSize function below
// func formatSize(size int64) string {
//     return ""
// }

// After implementing all functions:
// - Run: go test -tags solution ./minis/31-static-file-server/exercise/...
// - Test with: go run ./minis/31-static-file-server/cmd/file-server
// - Try in browser: http://localhost:8080
// - Check headers: curl -I http://localhost:8080/somefile.txt
// - Test ETag: curl -H "If-None-Match: \"abc\"" http://localhost:8080/file.txt
// - Compare with solution.go to see detailed implementation
