# Project 31: Static File Server

## 1. What Is This About?

### Real-World Scenario

You're building a web application. Users need to download images, videos, CSS, JavaScript files.

**Without optimization:**
1. Every request downloads the entire file (even if cached)
2. Video downloads can't resume if interrupted
3. Browsers re-download unchanged files constantly
4. Server sends 5MB image even though browser has it
5. Slow page loads, high bandwidth costs

**With proper static file server:**
1. ETags: Browser asks "has this file changed?" (304 Not Modified)
2. Range requests: Download resumes from byte 2,458,932
3. Cache headers: Browser caches for 1 year, no network requests
4. Proper MIME types: Browser knows image/jpeg vs text/html
5. Fast loads, reduced bandwidth (90% savings possible)

This project teaches you how to build **production-grade static file servers** with:
- **ETags**: Efficient cache validation without re-downloading
- **Range requests**: Partial content, resumable downloads
- **Caching headers**: Cache-Control, Expires, Last-Modified
- **MIME types**: Correct Content-Type detection
- **Security**: Path traversal prevention

### What You'll Learn

1. **HTTP file serving**: ServeFile, FileServer, custom handlers
2. **ETags**: Content-based cache validation (MD5/SHA256)
3. **Range requests**: HTTP 206 Partial Content
4. **Cache headers**: Cache-Control, Last-Modified, Expires
5. **MIME detection**: mime.TypeByExtension, content sniffing
6. **Security**: filepath.Clean, path traversal attacks

### The Challenge

Build a static file server with:
- ETag generation and validation (If-None-Match)
- Range request support (Accept-Ranges, Content-Range)
- Cache headers (Cache-Control, Last-Modified, If-Modified-Since)
- MIME type detection with fallback to application/octet-stream
- Directory listing (optional, with security considerations)
- Path traversal protection

---

## 2. First Principles: HTTP File Serving

### What is a Static File Server?

A **static file server** serves files from disk over HTTP.

**Basic flow**:
```
Client                Server
  |                      |
  |--- GET /image.jpg --|
  |                      |
  |                   Read file
  |                   Set headers
  |                      |
  |--- 200 OK + data ---|
  |    Content-Type: image/jpeg
```

**In Go**:
```go
http.Handle("/static/", http.FileServer(http.Dir("./public")))
```

This serves files from `./public` directory.

### What are ETags?

**ETags** (Entity Tags) are unique identifiers for file versions.

**Purpose**: Avoid re-downloading unchanged files.

**Flow without ETag**:
```
Client: GET /style.css
Server: 200 OK + 50KB CSS file
(1 hour later)
Client: GET /style.css
Server: 200 OK + 50KB CSS file (same file, wasted bandwidth)
```

**Flow with ETag**:
```
Server: 200 OK + ETag: "abc123" + 50KB CSS
(1 hour later)
Client: GET /style.css + If-None-Match: "abc123"
Server: File hasn't changed
Server: 304 Not Modified (no body, just 100 bytes of headers)
```

**ETag generation**:
```go
// Option 1: MD5 hash of content
hash := md5.Sum(fileContent)
etag := fmt.Sprintf("\"%x\"", hash)

// Option 2: ModTime + Size (faster, less accurate)
etag := fmt.Sprintf("\"%x-%x\"", modTime.Unix(), size)
```

**Why quotes?**
HTTP spec requires ETags to be quoted: `"abc123"` not `abc123`

### What are Range Requests?

**Range requests** allow downloading part of a file.

**Use cases**:
- **Resume downloads**: Continue from byte 5,000,000
- **Video streaming**: Download chunks as needed
- **PDF preview**: Download first 100KB to show first page

**Flow**:
```
Client: GET /video.mp4
        Range: bytes=0-1023
Server: 206 Partial Content
        Content-Range: bytes 0-1023/5242880
        (1KB of 5MB file)
```

**Headers**:
```
Request:  Range: bytes=0-1023
Response: Content-Range: bytes 0-1023/5242880
          Content-Length: 1024
          Accept-Ranges: bytes
```

**Range syntax**:
```
bytes=0-1023      → First 1024 bytes
bytes=1024-2047   → Second 1024 bytes
bytes=-1024       → Last 1024 bytes
bytes=1024-       → From byte 1024 to end
```

### What are Cache Headers?

**Cache headers** tell browsers how long to cache files.

**Cache-Control**:
```
Cache-Control: max-age=31536000    → Cache for 1 year
Cache-Control: no-cache            → Revalidate every time
Cache-Control: no-store            → Never cache
Cache-Control: public, max-age=86400 → Cache for 1 day
```

**Last-Modified / If-Modified-Since**:
```
Server: Last-Modified: Wed, 21 Oct 2023 07:28:00 GMT
(Later)
Client: If-Modified-Since: Wed, 21 Oct 2023 07:28:00 GMT
Server: 304 Not Modified (file unchanged)
```

**Expires**:
```
Expires: Thu, 01 Dec 2024 16:00:00 GMT
```

**Best practices**:
```go
// Static assets (CSS, JS, images with hash in name)
w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

// HTML pages
w.Header().Set("Cache-Control", "no-cache, must-revalidate")

// User-specific data
w.Header().Set("Cache-Control", "private, max-age=3600")
```

### What are MIME Types?

**MIME types** (Media Types) tell browsers how to handle files.

**Format**: `type/subtype`

**Common types**:
```
text/html           → HTML page
text/css            → CSS stylesheet
application/javascript → JavaScript
image/jpeg          → JPEG image
image/png           → PNG image
video/mp4           → MP4 video
application/pdf     → PDF document
application/json    → JSON data
application/octet-stream → Unknown binary
```

**In Go**:
```go
import "mime"

// By extension
mime.TypeByExtension(".jpg")  // "image/jpeg"
mime.TypeByExtension(".css")  // "text/css"

// Add custom types
mime.AddExtensionType(".webp", "image/webp")
```

**Why it matters**:
```
Content-Type: text/html     → Browser renders as HTML
Content-Type: text/plain    → Browser shows as text
Content-Type: image/jpeg    → Browser displays image
Content-Type: application/octet-stream → Browser downloads
```

---

## 3. Breaking Down the Solution

### Step 1: Basic File Server

```go
func serveFile(w http.ResponseWriter, r *http.Request, path string) {
    // Open file
    file, err := os.Open(path)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }
    defer file.Close()

    // Get file info
    stat, err := file.Stat()
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }

    // Set Content-Type
    contentType := mime.TypeByExtension(filepath.Ext(path))
    if contentType == "" {
        contentType = "application/octet-stream"
    }
    w.Header().Set("Content-Type", contentType)

    // Copy file to response
    http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}
```

**Why `http.ServeContent`?**
- Handles Range requests automatically
- Sets Last-Modified automatically
- Handles If-Modified-Since automatically
- More efficient than manual copy

### Step 2: ETag Generation

```go
func generateETag(stat os.FileInfo) string {
    // Simple ETag: modtime + size
    return fmt.Sprintf("\"%x-%x\"", stat.ModTime().Unix(), stat.Size())
}

func checkETag(w http.ResponseWriter, r *http.Request, etag string) bool {
    // Check If-None-Match header
    if match := r.Header.Get("If-None-Match"); match != "" {
        if match == etag {
            w.WriteHeader(http.StatusNotModified)
            return true // File not modified
        }
    }
    return false // File modified or no ETag
}
```

**Usage**:
```go
etag := generateETag(stat)
w.Header().Set("ETag", etag)

if checkETag(w, r, etag) {
    return // Already sent 304
}

// Serve file
http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
```

### Step 3: Range Request Handling

```go
func parseRange(r *http.Request, size int64) (start, end int64, err error) {
    rangeHeader := r.Header.Get("Range")
    if rangeHeader == "" {
        return 0, size - 1, nil // Full file
    }

    // Parse "bytes=0-1023"
    if !strings.HasPrefix(rangeHeader, "bytes=") {
        return 0, 0, fmt.Errorf("invalid range")
    }

    ranges := strings.TrimPrefix(rangeHeader, "bytes=")
    parts := strings.Split(ranges, "-")
    if len(parts) != 2 {
        return 0, 0, fmt.Errorf("invalid range")
    }

    // Parse start
    if parts[0] != "" {
        start, err = strconv.ParseInt(parts[0], 10, 64)
        if err != nil {
            return 0, 0, err
        }
    }

    // Parse end
    if parts[1] != "" {
        end, err = strconv.ParseInt(parts[1], 10, 64)
        if err != nil {
            return 0, 0, err
        }
    } else {
        end = size - 1
    }

    // Validate
    if start > end || start >= size {
        return 0, 0, fmt.Errorf("invalid range")
    }
    if end >= size {
        end = size - 1
    }

    return start, end, nil
}
```

**Serving range**:
```go
start, end, err := parseRange(r, stat.Size())
if err != nil {
    w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", stat.Size()))
    http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
    return
}

// Set headers
w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, stat.Size()))
w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
w.WriteHeader(http.StatusPartialContent)

// Seek and copy
file.Seek(start, io.SeekStart)
io.CopyN(w, file, end-start+1)
```

**Note**: `http.ServeContent` handles this automatically!

### Step 4: Cache Headers

```go
func setCacheHeaders(w http.ResponseWriter, path string, stat os.FileInfo) {
    // Last-Modified
    w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))

    // Cache-Control based on file type
    ext := filepath.Ext(path)
    switch ext {
    case ".html", ".htm":
        // HTML: revalidate every time
        w.Header().Set("Cache-Control", "no-cache, must-revalidate")
    case ".css", ".js":
        // Static assets: cache for 1 year (assumes versioned URLs)
        w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
    case ".jpg", ".jpeg", ".png", ".gif", ".webp":
        // Images: cache for 30 days
        w.Header().Set("Cache-Control", "public, max-age=2592000")
    default:
        // Default: cache for 1 hour
        w.Header().Set("Cache-Control", "public, max-age=3600")
    }
}
```

**If-Modified-Since check**:
```go
func checkModifiedSince(r *http.Request, modTime time.Time) bool {
    if ims := r.Header.Get("If-Modified-Since"); ims != "" {
        t, err := http.ParseTime(ims)
        if err == nil {
            // Truncate to seconds (HTTP date precision)
            modTime = modTime.Truncate(time.Second)
            if !modTime.After(t) {
                return false // Not modified
            }
        }
    }
    return true // Modified or no If-Modified-Since
}
```

### Step 5: MIME Type Detection

```go
func detectContentType(path string, file *os.File) string {
    // Try by extension first
    contentType := mime.TypeByExtension(filepath.Ext(path))
    if contentType != "" {
        return contentType
    }

    // Fallback: detect from content
    buffer := make([]byte, 512)
    n, err := file.Read(buffer)
    if err != nil && err != io.EOF {
        return "application/octet-stream"
    }

    // Reset file position
    file.Seek(0, io.SeekStart)

    // Detect
    contentType = http.DetectContentType(buffer[:n])
    return contentType
}
```

**How `http.DetectContentType` works**:
- Reads first 512 bytes
- Looks for magic numbers (file signatures)
- Returns best guess or "application/octet-stream"

**Examples**:
```
[0xFF, 0xD8, 0xFF] → image/jpeg
[0x89, 0x50, 0x4E, 0x47] → image/png
<!DOCTYPE html> → text/html
```

### Step 6: Security - Path Traversal

```go
func securePath(root, requestPath string) (string, error) {
    // Clean path (resolve .., ., etc.)
    requestPath = filepath.Clean(requestPath)

    // Remove leading slash
    requestPath = strings.TrimPrefix(requestPath, "/")

    // Join with root
    fullPath := filepath.Join(root, requestPath)

    // Ensure result is within root
    if !strings.HasPrefix(fullPath, filepath.Clean(root)) {
        return "", fmt.Errorf("path traversal attempt")
    }

    return fullPath, nil
}
```

**Why this is needed**:
```
Malicious request: GET /../../../etc/passwd
Without cleaning: /var/www/static/../../../etc/passwd → /etc/passwd
With cleaning:    Rejected (not within /var/www/static)
```

---

## 4. Complete Solution Walkthrough

### FileServer Structure

```go
type FileServer struct {
    root       string           // Root directory
    enableETag bool             // Enable ETag support
    cache      map[string][]byte // Optional: in-memory cache
}

func NewFileServer(root string) *FileServer {
    return &FileServer{
        root:       root,
        enableETag: true,
        cache:      make(map[string][]byte),
    }
}
```

### Main Handler

```go
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Only GET and HEAD
    if r.Method != http.MethodGet && r.Method != http.MethodHead {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Secure path
    path, err := securePath(fs.root, r.URL.Path)
    if err != nil {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Open file
    file, err := os.Open(path)
    if err != nil {
        if os.IsNotExist(err) {
            http.NotFound(w, r)
        } else {
            http.Error(w, "Internal error", http.StatusInternalServerError)
        }
        return
    }
    defer file.Close()

    // Get file info
    stat, err := file.Stat()
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }

    // Directory listing (if enabled)
    if stat.IsDir() {
        fs.serveDirectory(w, r, file)
        return
    }

    // Serve file
    fs.serveFile(w, r, file, stat)
}
```

### Serve File with All Features

```go
func (fs *FileServer) serveFile(w http.ResponseWriter, r *http.Request, file *os.File, stat os.FileInfo) {
    // Detect Content-Type
    contentType := detectContentType(file.Name(), file)
    w.Header().Set("Content-Type", contentType)

    // Set cache headers
    setCacheHeaders(w, file.Name(), stat)

    // Generate and set ETag
    if fs.enableETag {
        etag := generateETag(stat)
        w.Header().Set("ETag", etag)

        // Check If-None-Match
        if checkETag(w, r, etag) {
            return // 304 already sent
        }
    }

    // Check If-Modified-Since
    if !checkModifiedSince(r, stat.ModTime()) {
        w.WriteHeader(http.StatusNotModified)
        return
    }

    // Set Accept-Ranges
    w.Header().Set("Accept-Ranges", "bytes")

    // Serve with range support
    http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}
```

**What `http.ServeContent` does for us**:
1. Handles Range requests (206 Partial Content)
2. Sets Content-Length automatically
3. Handles If-Range header
4. Efficient copying with io.CopyN
5. Sets Last-Modified if not already set

### Directory Listing

```go
func (fs *FileServer) serveDirectory(w http.ResponseWriter, r *http.Request, dir *os.File) {
    entries, err := dir.Readdir(-1)
    if err != nil {
        http.Error(w, "Cannot read directory", http.StatusInternalServerError)
        return
    }

    // Sort by name
    sort.Slice(entries, func(i, j int) bool {
        return entries[i].Name() < entries[j].Name()
    })

    // Render HTML
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, "<html><body><h1>Directory: %s</h1><ul>\n", r.URL.Path)

    // Parent directory link
    if r.URL.Path != "/" {
        fmt.Fprintf(w, "<li><a href=\"..\">..</a></li>\n")
    }

    // List entries
    for _, entry := range entries {
        name := entry.Name()
        if entry.IsDir() {
            name += "/"
        }
        fmt.Fprintf(w, "<li><a href=\"%s\">%s</a> (%d bytes)</li>\n",
            url.PathEscape(name), html.EscapeString(name), entry.Size())
    }

    fmt.Fprintf(w, "</ul></body></html>\n")
}
```

---

## 5. Key Concepts Explained

### Concept 1: Why http.ServeContent is Powerful

```go
// Manual approach (DON'T DO THIS)
func badServe(w http.ResponseWriter, r *http.Request, file *os.File) {
    io.Copy(w, file)
}

// Good approach
func goodServe(w http.ResponseWriter, r *http.Request, file *os.File, modTime time.Time) {
    http.ServeContent(w, r, "file.txt", modTime, file)
}
```

**What ServeContent handles**:
- ✅ Range requests (bytes=0-1023)
- ✅ If-Modified-Since validation
- ✅ If-Range validation
- ✅ Content-Length calculation
- ✅ Efficient streaming
- ✅ Multiple range handling (multipart/byteranges)

### Concept 2: ETag vs Last-Modified

**Last-Modified**:
- Based on file modification time
- Precision: 1 second (HTTP date format)
- Problem: File content might change, time might not

**ETag**:
- Based on file content (hash) or metadata
- Precision: exact
- More reliable for cache validation

**Both together**:
```go
w.Header().Set("ETag", etag)
w.Header().Set("Last-Modified", modTime.Format(http.TimeFormat))

// Client can send both
// If-None-Match: "abc123"
// If-Modified-Since: Wed, 21 Oct 2023 07:28:00 GMT
```

**Precedence**: If-None-Match (ETag) takes priority over If-Modified-Since

### Concept 3: Cache Busting Strategies

**Problem**: How to update cached files?

**Strategy 1: Query strings**
```html
<script src="/app.js?v=1.2.3"></script>
```

**Strategy 2: Filename hashing**
```html
<script src="/app.abc123.js"></script>
```

**Strategy 3: Path versioning**
```html
<script src="/v1.2.3/app.js"></script>
```

**Best practice**: Filename hashing + long max-age
```go
// app.abc123.js → cache forever
w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

// index.html → never cache
w.Header().Set("Cache-Control", "no-cache")
```

### Concept 4: Content-Type Security

**Problem**: Wrong Content-Type can cause security issues

```go
// User uploads "image.jpg" with content:
// <script>alert('XSS')</script>

// If served as text/html → script executes
w.Header().Set("Content-Type", "text/html")

// If served as image/jpeg → browser treats as image
w.Header().Set("Content-Type", "image/jpeg")
```

**Defense**:
```go
// 1. Validate file extensions
allowedExts := map[string]bool{
    ".jpg": true, ".png": true, ".pdf": true,
}

// 2. Detect content type from actual content
contentType := http.DetectContentType(buffer)

// 3. Set X-Content-Type-Options
w.Header().Set("X-Content-Type-Options", "nosniff")
```

### Concept 5: Performance - Sendfile System Call

```go
// http.ServeContent uses sendfile() on Linux
// This is a zero-copy operation

// Without sendfile:
// Disk → Kernel buffer → User space → Kernel buffer → Network
// (4 copies, 2 context switches)

// With sendfile:
// Disk → Kernel buffer → Network
// (2 copies, 0 context switches to user space)
```

**Result**: Much faster, lower CPU usage

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Conditional Request Handler

```go
func ConditionalHandler(etag string, modTime time.Time, handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check ETag
        if match := r.Header.Get("If-None-Match"); match == etag {
            w.WriteHeader(http.StatusNotModified)
            return
        }

        // Check Modified-Since
        if ims := r.Header.Get("If-Modified-Since"); ims != "" {
            t, err := http.ParseTime(ims)
            if err == nil && !modTime.After(t) {
                w.WriteHeader(http.StatusNotModified)
                return
            }
        }

        // Set headers
        w.Header().Set("ETag", etag)
        w.Header().Set("Last-Modified", modTime.Format(http.TimeFormat))

        handler(w, r)
    }
}
```

### Pattern 2: In-Memory File Cache

```go
type FileCache struct {
    mu    sync.RWMutex
    files map[string]*CachedFile
}

type CachedFile struct {
    content     []byte
    contentType string
    etag        string
    modTime     time.Time
}

func (fc *FileCache) Get(path string) (*CachedFile, bool) {
    fc.mu.RLock()
    defer fc.mu.RUnlock()
    file, ok := fc.files[path]
    return file, ok
}

func (fc *FileCache) Set(path string, file *CachedFile) {
    fc.mu.Lock()
    defer fc.mu.Unlock()
    fc.files[path] = file
}
```

### Pattern 3: Compression Middleware

```go
func GzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if client accepts gzip
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        // Wrap response writer
        gzw := gzip.NewWriter(w)
        defer gzw.Close()

        w.Header().Set("Content-Encoding", "gzip")
        gzipWriter := &gzipResponseWriter{ResponseWriter: w, Writer: gzw}

        next.ServeHTTP(gzipWriter, r)
    })
}
```

### Pattern 4: SPA Fallback

```go
func SPAHandler(staticPath, indexPath string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        path := filepath.Join(staticPath, r.URL.Path)

        // Check if file exists
        if _, err := os.Stat(path); os.IsNotExist(err) {
            // Serve index.html for client-side routing
            http.ServeFile(w, r, filepath.Join(staticPath, indexPath))
            return
        }

        // Serve static file
        http.ServeFile(w, r, path)
    }
}
```

### Pattern 5: CORS for Static Assets

```go
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")

            for _, allowed := range allowedOrigins {
                if origin == allowed || allowed == "*" {
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
                    w.Header().Set("Access-Control-Max-Age", "86400")
                    break
                }
            }

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 7. Real-World Applications

### CDN Origin Server

Static file servers serve as CDN origin.

```go
// CloudFlare, Fastly, Akamai pull from origin
server := &FileServer{
    root: "/var/www/static",
    enableETag: true,
}

// Cache-Control tells CDN how long to cache
w.Header().Set("Cache-Control", "public, max-age=31536000")
```

Companies: Netflix, Spotify, YouTube (all use CDNs with origin servers)

### Static Website Hosting

Serve entire websites (Jekyll, Hugo, Next.js static export).

```go
server := SPAHandler("/var/www/site", "index.html")
http.ListenAndServe(":8080", server)
```

Companies: Netlify, Vercel, GitHub Pages

### Asset Server for Web Apps

Serve CSS, JS, images for web applications.

```go
// Versioned assets
http.Handle("/assets/", http.StripPrefix("/assets/", FileServer("/var/www/assets")))

// HTML pages (no cache)
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "no-cache")
    http.ServeFile(w, r, "index.html")
})
```

### Download Server

Serve large files with resume support.

```go
// Range requests allow resume
server := &FileServer{root: "/var/downloads"}

// Client can resume:
// Range: bytes=5000000-
```

Companies: SourceForge, FileZilla, Download managers

### Video Streaming

Serve video chunks with range requests.

```go
// HLS/DASH streaming
// Client requests: video-segment-001.ts
// Server: 206 Partial Content with video chunk
```

Companies: Twitch, Vimeo, AWS Media Services

---

## 8. Common Mistakes to Avoid

### Mistake 1: Not Validating Paths

**❌ Wrong**:
```go
path := filepath.Join(root, r.URL.Path)
http.ServeFile(w, r, path)
// Allows: /../../../etc/passwd
```

**✅ Correct**:
```go
path, err := securePath(root, r.URL.Path)
if err != nil {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
}
```

### Mistake 2: Wrong ETag Quotes

**❌ Wrong**:
```go
w.Header().Set("ETag", "abc123")  // Missing quotes
```

**✅ Correct**:
```go
w.Header().Set("ETag", "\"abc123\"")  // Quoted
```

### Mistake 3: Not Setting Accept-Ranges

**❌ Wrong**:
```go
// Client doesn't know if ranges are supported
http.ServeFile(w, r, path)
```

**✅ Correct**:
```go
w.Header().Set("Accept-Ranges", "bytes")
http.ServeContent(w, r, name, modTime, file)
```

### Mistake 4: Wrong Cache-Control for HTML

**❌ Wrong**:
```go
// HTML cached for 1 year → users never see updates
w.Header().Set("Cache-Control", "max-age=31536000")
```

**✅ Correct**:
```go
// HTML: no cache or short cache with revalidation
w.Header().Set("Cache-Control", "no-cache, must-revalidate")
```

### Mistake 5: Not Handling HEAD Requests

**❌ Wrong**:
```go
if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
}
```

**✅ Correct**:
```go
if r.Method != http.MethodGet && r.Method != http.MethodHead {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
}
// http.ServeContent handles HEAD automatically
```

### Mistake 6: Ignoring If-Range

**❌ Wrong**:
```go
// Always serve range, even if file changed
if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
    serveRange(w, r, file)
}
```

**✅ Correct**:
```go
// http.ServeContent checks If-Range
// If file changed, serves full file instead of range
http.ServeContent(w, r, name, modTime, file)
```

---

## 9. Stretch Goals

### Goal 1: Implement Compression ⭐

Add gzip compression for text files.

**Hint**:
```go
if strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".js") {
    gzw := gzip.NewWriter(w)
    defer gzw.Close()
    w.Header().Set("Content-Encoding", "gzip")
    // Serve to gzw
}
```

### Goal 2: Add Directory Listing JSON API ⭐⭐

Return directory contents as JSON.

**Hint**:
```go
type FileInfo struct {
    Name    string `json:"name"`
    Size    int64  `json:"size"`
    ModTime string `json:"modTime"`
    IsDir   bool   `json:"isDir"`
}

// GET /api/list?path=/images
json.NewEncoder(w).Encode(files)
```

### Goal 3: Implement Content-MD5 ⭐⭐

Add MD5 checksum header for integrity verification.

**Hint**:
```go
import "crypto/md5"

hash := md5.Sum(content)
w.Header().Set("Content-MD5", base64.StdEncoding.EncodeToString(hash[:]))
```

### Goal 4: Add Upload Support ⭐⭐⭐

Allow file uploads with multipart/form-data.

**Hint**:
```go
func handleUpload(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(10 << 20) // 10 MB
    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    dst, err := os.Create(filepath.Join(uploadDir, handler.Filename))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    io.Copy(dst, file)
    w.WriteHeader(http.StatusCreated)
}
```

### Goal 5: Implement Byte Range Serving from Memory ⭐⭐⭐

Serve ranges from in-memory cache.

**Hint**:
```go
type ByteRangeReader struct {
    data []byte
    pos  int64
}

func (r *ByteRangeReader) Read(p []byte) (n int, err error) {
    // Read from r.data[r.pos:]
}

func (r *ByteRangeReader) Seek(offset int64, whence int) (int64, error) {
    // Update r.pos
}
```

---

## How to Run

```bash
# Run the file server
go run ./minis/31-static-file-server/cmd/file-server

# Test with curl
curl -v http://localhost:8080/README.md

# Test ETag
curl -H "If-None-Match: \"abc123\"" http://localhost:8080/README.md

# Test Range request
curl -H "Range: bytes=0-99" http://localhost:8080/README.md

# Test with browser (open multiple times to see caching)
open http://localhost:8080/
```

---

## Summary

**What you learned**:
- ✅ Static file serving with http.ServeContent
- ✅ ETag generation and validation for efficient caching
- ✅ Range request handling for resumable downloads
- ✅ Cache headers (Cache-Control, Last-Modified, Expires)
- ✅ MIME type detection and security
- ✅ Path traversal prevention

**Why this matters**:
Every web application serves static files. Proper caching can reduce bandwidth by 90% and improve load times by 10x. Range requests enable video streaming and resumable downloads.

**Key takeaway**:
ETags + Cache-Control + Range requests = Fast, efficient file serving

**Next steps**:
- Project 32: Learn WebSocket servers for real-time bidirectional communication
- Project 33: Learn HTTP/2 server push for optimized asset delivery

Build efficiently!
