//go:build !solution && !reference

package staticfileserver

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

// NewFileServer implements the exercise.
//
// TODO: Implement this function
func NewFileServer(config FileServerConfig) (*FileServer, error) {
	// TODO: Implement
	return nil, nil
}

// ServeHTTP implements the exercise.
//
// TODO: Implement this function
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

// serveFile implements the exercise.
//
// TODO: Implement this function
func (fs *FileServer) serveFile(w http.ResponseWriter, r *http.Request, path string) {
	// TODO: Implement
}

// serveDirectory implements the exercise.
//
// TODO: Implement this function
func (fs *FileServer) serveDirectory(w http.ResponseWriter, r *http.Request, path string) {
	// TODO: Implement
}
