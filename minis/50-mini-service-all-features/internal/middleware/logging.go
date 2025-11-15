package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// Logging logs HTTP requests and responses
func Logging(logger zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := GetRequestID(r.Context())

			// Log request
			logger.Info().
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Msg("request started")

			// Wrap response writer
			rw := NewResponseWriter(w)

			next.ServeHTTP(rw, r)

			// Log response
			duration := time.Since(start)
			logger.Info().
				Str("request_id", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", rw.StatusCode()).
				Int("bytes", rw.BytesWritten()).
				Dur("duration", duration).
				Msg("request completed")
		})
	}
}
