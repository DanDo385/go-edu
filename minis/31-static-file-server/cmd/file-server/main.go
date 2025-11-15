package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/go-10x-minis/minis/31-static-file-server/exercise"
)

func main() {
	// Parse command-line flags
	var (
		addr        = flag.String("addr", ":8080", "HTTP server address")
		root        = flag.String("root", "./testdata", "Root directory to serve")
		enableETag  = flag.Bool("etag", true, "Enable ETag support")
		enableRange = flag.Bool("range", true, "Enable Range request support")
		listDir     = flag.Bool("list", true, "Enable directory listing")
	)
	flag.Parse()

	// Create server configuration
	config := exercise.FileServerConfig{
		Root:              *root,
		EnableETag:        *enableETag,
		EnableRange:       *enableRange,
		EnableDirListing:  *listDir,
		DefaultCacheMaxAge: 3600, // 1 hour
	}

	// Create file server
	fileServer, err := exercise.NewFileServer(config)
	if err != nil {
		log.Fatalf("Failed to create file server: %v", err)
	}

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         *addr,
		Handler:      fileServer,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		log.Printf("Starting file server on %s", *addr)
		log.Printf("Serving files from: %s", *root)
		log.Printf("Features: ETag=%v Range=%v DirListing=%v", *enableETag, *enableRange, *listDir)
		log.Printf("Access at: http://localhost%s", *addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutdown signal received...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}

func init() {
	// Create testdata directory with sample files if it doesn't exist
	if _, err := os.Stat("./testdata"); os.IsNotExist(err) {
		if err := setupTestData(); err != nil {
			log.Printf("Warning: Could not setup test data: %v", err)
		}
	}
}

func setupTestData() error {
	// Create testdata directory
	if err := os.MkdirAll("./testdata", 0755); err != nil {
		return err
	}

	// Create subdirectory
	if err := os.MkdirAll("./testdata/images", 0755); err != nil {
		return err
	}

	// Create sample files
	files := map[string]string{
		"./testdata/index.html": `<!DOCTYPE html>
<html>
<head>
    <title>Static File Server Demo</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <h1>Welcome to the Static File Server</h1>
    <p>This server demonstrates ETags, Range requests, and caching.</p>
    <img src="images/sample.txt" alt="Sample">
    <script src="app.js"></script>
</body>
</html>`,
		"./testdata/styles.css": `body {
    font-family: Arial, sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    background-color: #f5f5f5;
}

h1 {
    color: #333;
    border-bottom: 2px solid #007bff;
    padding-bottom: 10px;
}

p {
    line-height: 1.6;
    color: #666;
}`,
		"./testdata/app.js": `console.log('Static file server loaded!');

document.addEventListener('DOMContentLoaded', function() {
    const heading = document.querySelector('h1');
    if (heading) {
        heading.style.color = '#007bff';
    }
});`,
		"./testdata/README.md": `# Static File Server Test Data

This directory contains sample files for testing the static file server.

## Features Demonstrated

1. **ETags**: Try refreshing the page - the server will send 304 Not Modified
2. **Range Requests**: Use curl with Range header to download partial files
3. **Caching**: Different cache policies for HTML vs CSS/JS
4. **MIME Types**: Proper Content-Type for each file

## Examples

### Test ETag
curl -v http://localhost:8080/README.md
# Note the ETag header, then:
curl -H "If-None-Match: \"<etag>\"" http://localhost:8080/README.md

### Test Range Request
curl -H "Range: bytes=0-99" http://localhost:8080/README.md

### Test Last-Modified
curl -v http://localhost:8080/styles.css
# Note the Last-Modified header, then:
curl -H "If-Modified-Since: <date>" http://localhost:8080/styles.css
`,
		"./testdata/images/sample.txt": `This is a sample image placeholder.
In a real application, this would be a PNG or JPEG file.
The server will detect the MIME type correctly.`,
		"./testdata/large.txt": generateLargeFile(1024 * 100), // 100KB file for range testing
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	log.Println("Created test data in ./testdata directory")
	return nil
}

func generateLargeFile(size int) string {
	content := ""
	line := "This is line %d of a large file for testing range requests and streaming.\n"
	lineNum := 1
	for len(content) < size {
		content += fmt.Sprintf(line, lineNum)
		lineNum++
	}
	return content[:size]
}
