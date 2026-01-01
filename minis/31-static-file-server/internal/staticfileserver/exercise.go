//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package staticfileserver

import (
	"net/http"

	"os"

	"time"
)

type FileServerConfig struct {
	Root               string // Root directory to serve files from
	EnableETag         bool   // Enable ETag generation and validation
	EnableRange        bool   // Enable HTTP Range request support
	EnableDirListing   bool   // Enable directory listing
	DefaultCacheMaxAge int    // Default Cache-Control max-age in seconds
}

type FileServer struct {
	config FileServerConfig
}
// TODO: implement NewFileServer.
func NewFileServer(config FileServerConfig) (*FileServer, error) { panic("TODO: implement") }
// TODO: implement ServeHTTP.
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) { panic("TODO: implement") }
// TODO: implement securePath.
func securePath(root, requestPath string) (string, error) { panic("TODO: implement") }
// TODO: implement generateETag.
func generateETag(stat os.FileInfo) string { panic("TODO: implement") }
// TODO: implement generateETagWithHash.
func generateETagWithHash(path string) (string, error) { panic("TODO: implement") }
// TODO: implement checkETag.
func checkETag(r *http.Request, etag string) bool { panic("TODO: implement") }
// TODO: implement checkModifiedSince.
func checkModifiedSince(r *http.Request, modTime time.Time) bool { panic("TODO: implement") }
// TODO: implement setCacheHeaders.
func setCacheHeaders(w http.ResponseWriter, filename string, maxAge int) { panic("TODO: implement") }
// TODO: implement detectContentType.
func detectContentType(path string) string { panic("TODO: implement") }
// TODO: implement serveFile.
func (fs *FileServer) serveFile(w http.ResponseWriter, r *http.Request, path string) {
	panic("TODO: implement")
}
// TODO: implement serveDirectory.
func (fs *FileServer) serveDirectory(w http.ResponseWriter, r *http.Request, path string) {
	panic("TODO: implement")
}
// TODO: implement formatSize.
func formatSize(size int64) string { panic("TODO: implement") }
