package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Chain applies middleware in order
// First middleware in the list wraps all others
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Apply in reverse so first middleware wraps all others
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ResponseWriter wraps http.ResponseWriter to capture status code and bytes written
type ResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default 200
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *ResponseWriter) BytesWritten() int {
	return rw.bytesWritten
}
